package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/loafoe/ferrite/storer"
	"github.com/loafoe/ferrite/types"

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

// Bootstrap implement the bootstrap call to ourselves
func Bootstrap(baseURL, token string) (*types.Bootstrap, error) {
	bootstrapURL := fmt.Sprintf("%s/bootstrap", baseURL)
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}
	req, _ := http.NewRequest(http.MethodPost, bootstrapURL, nil)
	req.Header.Set("Authorization", "OAuth "+token)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error bootstrapping ferrite: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var bootstrap types.Bootstrap
	if err := json.NewDecoder(resp.Body).Decode(&bootstrap); err != nil {
		return nil, fmt.Errorf("error decoding ferrite bootstrap data: %w", err)
	}
	return &bootstrap, nil
}
