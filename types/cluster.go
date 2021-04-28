package types

type Cluster struct {
	ID         string `json:"id" gorm:"primaryKey"`
	PrivateKey string `json:"private_key,omitempty"`
	PublicKey  string `json:"public_key,omitempty" gorm:"-"`
}
