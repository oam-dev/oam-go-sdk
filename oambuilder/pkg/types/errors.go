package types

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	Error_PartialWrite = Error("Failed to write all data to file.")
)

// generator
const (
	Error_NeedExecuteFirst      = Error("You should run generator.Execute first")
	Error_FailedToFindGenerator = Error("Can not detect an generator handler.")
)

// generator.kubebuilder
const (
	Error_OnlySupportScaffoldV2 = Error("We support scaffold v2 for kubebuilder only.")
)

// project
const (
	Error_VDRMismatch = Error("Generated Version/Domain/Repo set is not match with OAM")
)

// generator.kubebuilder.crgen
const (
	Error_NeedOAM              = Error("path for file 'OAM' is required.")
	Error_NeedPath             = Error("path is required.")
	Error_NeedOutput           = Error("output is required.")
	Error_InvalidParameterType = Error("invalid parameter type(boolean, string, number, null)")
)
