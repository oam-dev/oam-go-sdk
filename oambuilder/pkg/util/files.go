package util

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/imports"
)

const (
	ApiPkgImportScaffoldMarker    = "// +kubebuilder:scaffold:imports"
	ApiSchemeScaffoldMarker       = "// +kubebuilder:scaffold:scheme"
	ReconcilerSetupScaffoldMarker = "// +kubebuilder:scaffold:builder"
)

// Append strings after markers.
// This file copy from kubebuilder internal package
func InsertStringsInFile(path string, markerAndValues map[string][]string) error {
	isGoFile := false
	if ext := filepath.Ext(path); ext == ".go" {
		isGoFile = true
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	r, err := insertStrings(f, markerAndValues)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	content, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	formattedContent := content
	if isGoFile {
		formattedContent, err = imports.Process(path, content, nil)
		if err != nil {
			return err
		}
	}

	// use Go import process to format the content
	err = ioutil.WriteFile(path, formattedContent, os.ModePerm)
	if err != nil {
		return err
	}

	return err
}

// insertStrings reads content from given reader and insert string below the
// line containing marker string. So for ex. in insertStrings(r, {'m1':
// [v1], 'm2': [v2]})
// v1 will be inserted below the lines containing m1 string and v2 will be inserted
// below line containing m2 string.
func insertStrings(r io.Reader, markerAndValues map[string][]string) (io.Reader, error) {
	// reader clone is needed since we will be reading twice from the given reader
	buf := new(bytes.Buffer)
	rClone := io.TeeReader(r, buf)

	err := filterExistingValues(rClone, markerAndValues)
	if err != nil {
		return nil, err
	}

	out := new(bytes.Buffer)

	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := scanner.Text()
		_, err := out.WriteString(line + "\n")
		if err != nil {
			return nil, err
		}

		for marker, vals := range markerAndValues {
			if strings.TrimSpace(line) == strings.TrimSpace(marker) {
				for i := len(vals) - 1; i >= 0; i-- {
					_, err := out.WriteString(vals[i])
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// filterExistingValues removes the single-line values that already exists in
// the given reader. Multi-line values are ignore currently simply because we
// don't have a use-case for it.
func filterExistingValues(r io.Reader, markerAndValues map[string][]string) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		for marker, vals := range markerAndValues {
			for i, val := range vals {
				if strings.TrimSpace(line) == strings.TrimSpace(val) {
					markerAndValues[marker] = append(vals[:i], vals[i+1:]...)
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
