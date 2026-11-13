
package utilities

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"EMS/shared/discovery"
	"EMS/shared/httpclient"
)

type DepartmentResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type DepartmentClient struct {
	consul      *discovery.Client
	httpClient  *httpclient.Client
}

func NewDepartmentClient(
	consul *discovery.Client,
	httpClient *httpclient.Client,
) *DepartmentClient {
	return &DepartmentClient{
		consul:     consul,
		httpClient: httpClient,
	}
}

func (c *DepartmentClient) GetByID(
	ctx context.Context,
	departmentID uint,
) (*DepartmentResponse, error) {

	if departmentID == 0 {
		return nil, fmt.Errorf("department ID is required")
	}

	baseURL, err := c.consul.DiscoverService("department-service")
	if err != nil {
		return nil, fmt.Errorf("department service discovery failed: %w", err)
	}

	requestURL := baseURL + "/api/departments/" +
		strconv.FormatUint(uint64(departmentID), 10)

	response, err := c.httpClient.Do(
		ctx,
		http.MethodGet,
		requestURL,
		nil,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("department service request failed: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("department not found")
	}

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf(
			"department service returned status %d: %s",
			response.StatusCode,
			string(body),
		)
	}

	var department DepartmentResponse

	if err := json.NewDecoder(response.Body).Decode(&department); err != nil {
		return nil, fmt.Errorf("failed to decode department response: %w", err)
	}

	return &department, nil
}
