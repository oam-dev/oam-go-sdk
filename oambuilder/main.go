package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/oam-dev/oam-go-sdk/oambuilder/cmd/exchange"
	"github.com/oam-dev/oam-go-sdk/oambuilder/cmd/trait"
	"github.com/oam-dev/oam-go-sdk/oambuilder/cmd/workload"
	"github.com/oam-dev/oam-go-sdk/oambuilder/pkg/generator/kubebuilder/crgen"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	logger = ctrl.Log.WithName("oambuilder.main")
)

var builder = &cobra.Command{
	Use:   "oambuilder",
	Short: "code generator for oam style operator",
	Long:  "code generator for oam style operator",
}

func init() {
	builder.AddCommand(workload.Workload)
	builder.AddCommand(trait.Trait)
	builder.AddCommand(exchange.Exchange)
	builder.AddCommand(crgen.CRGen)
}

func main() {
	ctrl.SetLogger(zap.New(func(o *zap.Options) {
		o.Development = true
	}))
	if err := builder.Execute(); err != nil {
		logger.Error(err, "cmd execute error")
		os.Exit(1)
	}
}
