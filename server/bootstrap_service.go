package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/philips-labs/ferrite/storer"
	"github.com/philips-labs/ferrite/types"

	"github.com/labstack/echo/v4"
)

type BootstrapService struct {
	Storer *storer.Ferrite
}

type bootstrapResponse struct {
	Message string `json:"message"`
}

func (b *BootstrapService) Bootstrap(c echo.Context) error {
	var bootstrap types.Bootstrap
	cluster, err := b.Storer.Cluster.FindLatest()
	now := time.Now()
	if err != nil { // Assume empty
		newCluster := types.Cluster{
			ID:        strings.Replace(uuid.New().String(), "-", "", -1),
			CreatedAt: now,
			UpdatedAt: now,
		}
		newCluster.PrivateKey, _ = newCluster.GeneratePrivateKeyPEM()
		cluster, err = b.Storer.Cluster.Create(newCluster)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, bootstrapResponse{err.Error()})
		}
	}
	project, err := b.Storer.Project.FindLatest()
	if err != nil { // Assume no project exists
		project, err = b.Storer.Project.Create(types.Project{
			ID:        strings.Replace(uuid.New().String(), "-", "", -1),
			CreatedAt: now,
			UpdatedAt: now,
			Name:      "Bootstrapped project",
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, bootstrapResponse{err.Error()})
		}
	}

	// Return bootstrap data
	publicKeyPEM, err := cluster.PublicKeyPEM()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, bootstrapResponse{err.Error()})
	}
	bootstrap.ProjectID = project.ID
	bootstrap.ClusterID = cluster.ID
	bootstrap.PublicKey = publicKeyPEM
	return c.JSON(http.StatusOK, bootstrap)
}
