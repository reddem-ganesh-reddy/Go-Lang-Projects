package discovery

import (
	"encoding/json"
	"fmt"
	"bytes"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"EMS/config"
)

// Client communicates with Consul.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// Service represents a service returned by Consul.
type Service struct {
	ID      string `json:"ID"`
	Name    string `json:"Service"`
	Address string `json:"Address"`
	Port    int    `json:"Port"`
}

// serviceRegistration represents Consul registration data.
type serviceRegistration struct {
	ID      string   `json:"ID"`
	Name    string   `json:"Name"`
	Address string   `json:"Address"`
	Port    int      `json:"Port"`
	Tags    []string `json:"Tags"`
	Check   struct {
		HTTP                           string `json:"HTTP"`
		Interval                       string `json:"Interval"`
		Timeout                        string `json:"Timeout"`
		DeregisterCriticalServiceAfter string `json:"DeregisterCriticalServiceAfter"`
	} `json:"Check"`
}

// New creates a Consul client.
func New(cfg config.Config) *Client {
	return &Client{
		baseURL: fmt.Sprintf(
			"%s://%s",
			cfg.ConsulScheme,
			cfg.ConsulAddress,
		),
		httpClient: &http.Client{
			Timeout: cfg.HTTPTimeout(),
		},
	}
}

// RegisterService registers a service with Consul.
func (c *Client) RegisterService(
	name string,
	id string,
	address string,
	port int,
) error {

	registration := serviceRegistration{
		ID:      id,
		Name:    name,
		Address: address,
		Port:    port,
		Tags:    []string{"ems", "go", "microservice"},
	}

	registration.Check.HTTP = fmt.Sprintf(
		"http://%s:%d/health",
		address,
		port,
	)
	registration.Check.Interval = "10s"
	registration.Check.Timeout = "5s"
	registration.Check.DeregisterCriticalServiceAfter = "1m"

	body, err := json.Marshal(registration)
	if err != nil {
		return fmt.Errorf("marshal service registration: %w", err)
	}

	request, err := http.NewRequest(
		http.MethodPut,
		c.baseURL+"/v1/agent/service/register",
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("create registration request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("register service with Consul: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf(
			"Consul registration failed with status %s",
			response.Status,
		)
	}

	return nil
}

// DiscoverService finds a healthy service instance by service name.
func (c *Client) DiscoverService(name string) (string, error) {
	serviceURL := fmt.Sprintf(
		"%s/v1/health/service/%s?passing=true",
		c.baseURL,
		url.PathEscape(name),
	)

	request, err := http.NewRequest(http.MethodGet, serviceURL, nil)
	if err != nil {
		return "", fmt.Errorf("create discovery request: %w", err)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("discover service %s: %w", name, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"Consul discovery failed with status %s",
			response.Status,
		)
	}

	var instances []struct {
		Service Service `json:"Service"`
	}

	if err := json.NewDecoder(response.Body).Decode(&instances); err != nil {
		return "", fmt.Errorf("decode Consul response: %w", err)
	}

	if len(instances) == 0 {
		return "", fmt.Errorf("no healthy instances found for service %s", name)
	}

	instance := instances[0].Service

	if instance.Address == "" || instance.Port == 0 {
		return "", fmt.Errorf("invalid service address returned by Consul")
	}

	return "http://" + net.JoinHostPort(
		instance.Address,
		strconv.Itoa(instance.Port),
	), nil
}

// DeregisterService removes a service registration from Consul.
func (c *Client) DeregisterService(id string) error {
	request, err := http.NewRequest(
		http.MethodPut,
		c.baseURL+"/v1/agent/service/deregister/"+url.PathEscape(id),
		nil,
	)
	if err != nil {
		return fmt.Errorf("create deregistration request: %w", err)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("deregister service: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf(
			"Consul deregistration failed with status %s",
			response.Status,
		)
	}

	return nil
}

// Keep time imported for future discovery-related timeout configuration.
var _ = time.Second