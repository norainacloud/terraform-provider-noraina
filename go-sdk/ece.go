package go_sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
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

func (c *NorainaApiClient) CreateInstance(ireq InstanceRequest) (*InstanceResponse, error) {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(ireq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, norainaDomain+instanceRoute, b)
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-access-token", c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("[ERROR] Create Instance Error with Status code %v", resp.StatusCode))
	}

	log.Printf("[DEBUG] CREATE -> response Status : %s", resp.Status)
	log.Printf("[DEBUG] CREATE -> response Headers : %s", resp.Header)
	log.Printf("[DEBUG] CREATE -> response Body : %s", resp.Body)

	instanceResponse := &InstanceResponse{}
	err = json.NewDecoder(resp.Body).Decode(instanceResponse)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Instance ID %v", instanceResponse.Instance.Id)
	log.Printf("[DEBUG] Tenant ID %v", instanceResponse.Instance.TenantId)

	return instanceResponse, nil
}

func (c *NorainaApiClient) UpdateInstance(id string, ireq InstanceRequest) (*InstanceResponse, error) {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(ireq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, norainaDomain+instanceRoute+"/"+id, b)
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-access-token", c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("[ERROR] Update Instance Error with Status code %v", resp.StatusCode))
	}

	log.Printf("[DEBUG] CREATE -> response Status : %s", resp.Status)
	log.Printf("[DEBUG] CREATE -> response Headers : %s", resp.Header)
	log.Printf("[DEBUG] CREATE -> response Headers : %v", resp.Body)

	instanceResponse := &InstanceResponse{}
	err = json.NewDecoder(resp.Body).Decode(instanceResponse)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Instance ID %v", instanceResponse.Instance.Id)
	log.Printf("[DEBUG] Tenant ID %v", instanceResponse.Instance.TenantId)

	return instanceResponse, nil
}

func (c *NorainaApiClient) GetInstance(id string) (*InstanceResponse, error) {
	body, err := c.getResource(id, instanceRoute)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	instanceResponse := &InstanceResponse{}
	err = json.NewDecoder(body).Decode(instanceResponse)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Instance ID %v", instanceResponse.Instance.Id)
	log.Printf("[DEBUG] Tenant ID %v", instanceResponse.Instance.TenantId)

	return instanceResponse, nil
}

func (c *NorainaApiClient) DeleteInstance(id string) error {
	return c.deleteResource(id, instanceRoute)
}
