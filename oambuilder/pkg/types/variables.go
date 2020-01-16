package types

const (
	LABEL_OAM_UUID       = "OAM_UUID"
	LABEL_OAM_GENERATION = "OAM_GENERATION"
	LABEL_OAM_WORKLOAD   = "OAM_WORKLOAD"
)

type WorkloadUIDGetter interface {
	WorkloadUID() string
}
