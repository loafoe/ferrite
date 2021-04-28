package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"ferrite/storer"
	"ferrite/types"
	"net/http"
	"strings"

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
	id := strings.Replace(uuid.New().String(), "-", "", -1)
	cluster.ID = id

	// Generate key
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, clusterResponse{err.Error()})
	}

	var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privateKeyPEM := pem.EncodeToMemory(privateKeyBlock)
	cluster.PrivateKey = string(privateKeyPEM)

	createdCluster, err := g.Storer.Cluster.Create(cluster)
	if err != nil {
		return c.JSON(http.StatusBadRequest, clusterResponse{err.Error()})
	}
	createdCluster.PrivateKey = ""
	publicKey := &privateKey.PublicKey
	var publicKeyBytes []byte = x509.MarshalPKCS1PublicKey(publicKey)
	publicKeyBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicKeyPEM := pem.EncodeToMemory(publicKeyBlock)
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
