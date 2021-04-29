package server

import (
	"github.com/philips-labs/ferrite/storer"
	"github.com/philips-labs/ferrite/types"
	"net/http"

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
	if err != nil { // Assume empty
		cluster, err = b.Storer.Cluster.Create(types.Cluster{})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, bootstrapResponse{err.Error()})
		}
	}
	project, err := b.Storer.Project.FindLatest()
	if err != nil { // Assume no project exists
		project, err = b.Storer.Project.Create(types.Project{
			Name: "Bootstrapped project",
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
	return c.JSON(http.StatusOK, &bootstrap)
}
