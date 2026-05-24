package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/oroodrigo/bufo/internal/store"
)

type ApiClient struct {
	client *http.Client
}

func NewApiClient(socketPath string) *ApiClient {
	return &ApiClient{
		client: &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", socketPath)
				},
			},
		},
	}
}

func (c *ApiClient) ListRoutes() (map[string]store.Route, error) {
	resp, err := c.client.Get("http://unix/routes")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var routes map[string]store.Route
	if err := json.Unmarshal(body, &routes); err != nil {
		return nil, err
	}

	return routes, nil
}

func (c *ApiClient) AddRoute(name string, route store.Route) error {
	jsonRoute, err := json.Marshal(route)
	if err != nil {
		return err
	}

	resp, err := c.client.Post("http://unix/routes/"+name, "application/json", bytes.NewBuffer(jsonRoute))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("erro ao adicionar rota: %s", resp.Status)
	}

	return nil
}

func (c *ApiClient) DeleteRoute(name string) error {
	req, err := http.NewRequest(http.MethodDelete, "http://unix/routes/"+name, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erro ao remover rota: %s", resp.Status)
	}

	return nil
}
