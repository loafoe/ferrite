package server

import (
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"strings"
	"time"

	"github.com/philips-labs/ferrite/storer"
	"github.com/philips-labs/ferrite/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ClusterService struct {
	Storer *storer.Ferrite
}

type clusterResponse struct {
	Message string `json:"msg"`
}

func (g *ClusterService) Create(c echo.Context) error {
	var cluster types.Cluster
	if err := c.Bind(&cluster); err != nil {
		return c.JSON(http.StatusBadRequest, clusterResponse{err.Error()})
	}
	if cluster.ID != "" { // Update
		return c.JSON(http.StatusBadRequest, clusterResponse{"updates are not supported"})
	}
	now := time.Now()
	id := strings.Replace(uuid.New().String(), "-", "", -1)
	cluster.ID = id
	cluster.CreatedAt = now
	cluster.UpdatedAt = now

	privateKeyPEM, err := cluster.GeneratePrivateKeyPEM()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, clusterResponse{"error generating private key"})
	}
	cluster.PrivateKey = privateKeyPEM

	createdCluster, err := g.Storer.Cluster.Create(cluster)
	if err != nil {
		return c.JSON(http.StatusBadRequest, clusterResponse{err.Error()})
	}
	publicKeyPEM, _ := createdCluster.PublicKeyPEM()
	createdCluster.PublicKey = string(publicKeyPEM)
	return c.JSON(http.StatusCreated, createdCluster)
}

func (g *ClusterService) Get(c echo.Context) error {
	clusterID := c.Param("cluster")
	foundCluster, err := g.Storer.Cluster.FindByID(clusterID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, clusterResponse{err.Error()})
	}
	block, _ := pem.Decode([]byte(foundCluster.PrivateKey))
	if block == nil {
		return c.JSON(http.StatusInternalServerError, clusterResponse{"error decoding private key"})
	}
	foundCluster.PrivateKey = ""

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return c.JSON(http.StatusBadRequest, clusterResponse{err.Error()})
	}
	// Decode private key
	publicKey := &privateKey.PublicKey
	var publicKeyBytes []byte = x509.MarshalPKCS1PublicKey(publicKey)
	publicKeyBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicKeyPEM := pem.EncodeToMemory(publicKeyBlock)
	foundCluster.PublicKey = string(publicKeyPEM)
	return c.JSON(http.StatusOK, foundCluster)
}
