package types

type Bootstrap struct {
	ClusterID string `json:"cluster_id"`
	ProjectID string `json:"project_id"`
	PublicKey string `json:"public_key"`
}
