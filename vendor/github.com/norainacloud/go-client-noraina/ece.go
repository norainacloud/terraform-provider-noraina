package noraina

import (
	"context"
	"fmt"
	"net/http"
)

const (
	instanceRoute = "api/instance"
)

type InstanceRequest struct {
	Name     string                   `json:"name"`
	Services []InstanceServiceRequest `json:"services"`
}

type InstanceServiceRequest struct {
	Name             string `json:"name"`
	Fqdn             string `json:"fqdn"`
	OriginHostHeader string `json:"origin_hostheader"`
	OriginBackend    string `json:"origin_backend"`
	ProviderRegion   string `json:"provider_region"`
	ProviderName     string `json:"provider_name"`
	CertId           string `json:"cert_id,omitempty"`
}

type InstanceResponse struct {
	Instance InstanceFields `json:"instance"`
	Servers  []ServerFields `json:"servers"`
}

type InstanceFields struct {
	Id          string          `json:"_id"`
	Name        string          `json:"name"`
	TenantId    string          `json:"tenant_id"`
	CreatedDate string          `json:"created_date"`
	Version     int             `json:"__v"`
	Services    []ServiceFields `json:"services"`
}

type ServiceFields struct {
	Id               string `json:"_id"`
	Name             string `json:"name"`
	Fqdn             string `json:"fqdn"`
	OriginHostHeader string `json:"origin_hostheader"`
	OriginBackend    string `json:"origin_backend"`
	ProviderRegion   string `json:"provider_region"`
	ProviderName     string `json:"provider_name"`
	CreatedDate      string `json:"created_date"`
	CertId           string `json:"cert_id,omitempty"`
}

type ServerFields struct {
	DNSRecordID DNSRecordID `json:"dns_record_id"`
	IPAdress    IPAdress    `json:"ip_address"`
	Status      string      `json:"status"`
	Id          string      `json:"_id"`
	Name        string      `json:"name"`
	ContainerID string      `json:"container_id"`
	HostID      string      `json:"host_id"`
	InstanceID  string      `json:"instance_id"`
	CreatedDate string      `json:"created_date"`
	Version     int         `json:"__v"`
}

type DNSRecordID struct {
	A    string `json:"a"`
	AAAA string `json:"aaaa"`
}

type IPAdress struct {
	IPv4 string `json:"ipv4"`
	IPv6 string `json:"ipv6"`
}

func (c *Client) CreateInstance(ctx context.Context, ireq InstanceRequest) (*InstanceResponse, error) {
	path := instanceRoute

	req, err := c.NewRequest(ctx, http.MethodPost, path, ireq)
	if err != nil {
		return nil, err
	}

	instanceResponse := new(InstanceResponse)
	err = c.Do(ctx, req, instanceResponse)
	if err != nil {
		return nil, err
	}

	return instanceResponse, nil
}

func (c *Client) UpdateInstance(ctx context.Context, id string, ireq InstanceRequest) (*InstanceResponse, error) {
	path := fmt.Sprintf("%s/%s", instanceRoute, id)

	req, err := c.NewRequest(ctx, http.MethodPost, path, ireq)
	if err != nil {
		return nil, err
	}

	instanceResponse := new(InstanceResponse)
	err = c.Do(ctx, req, instanceResponse)
	if err != nil {
		return nil, err
	}

	return instanceResponse, nil
}

func (c *Client) GetInstance(ctx context.Context, id string) (*InstanceResponse, error) {
	path := fmt.Sprintf("%s/%s", instanceRoute, id)
	req, err := c.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	instanceResponse := new(InstanceResponse)
	err = c.Do(ctx, req, instanceResponse)
	if err != nil {
		return nil, err
	}

	return instanceResponse, nil
}

func (c *Client) DeleteInstance(ctx context.Context, id string) error {
	return c.deleteResource(ctx, id, instanceRoute)
}
