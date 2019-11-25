package noraina

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/norainacloud/terraform-provider-noraina/go-sdk"
)

func resourceEce() *schema.Resource {
	return &schema.Resource{
		Create: resourceEceCreate,
		Read:   resourceEceRead,
		Update: resourceEceUpdate,
		Delete: resourceEceDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"service": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},
						"fqdn": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},
						"origin_hostheader": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},
						"origin_backend": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},
						"provider_region": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},
						"provider_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},
						"cert_id": {
							Type:     schema.TypeString,
							Required: false,
							ForceNew: false,
							Optional: true,
						},
						"created_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"server": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"container_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"created_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_record_id": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"a": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"aaaa": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"ip_address": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ipv6": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceEceCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*go_sdk.NorainaApiClient)

	services := d.Get("service").([]interface{})
	if len(services) == 0 {
		return errors.New("[ERROR] you must define at least 1 service in the instance")
	}

	ireq := go_sdk.InstanceRequest{
		Name:     d.Get("name").(string),
		Services: make([]go_sdk.InstanceServiceRequest, 0),
	}

	for _, service := range services {
		s := service.(map[string]interface{})

		isr := go_sdk.InstanceServiceRequest{
			Name:             s["name"].(string),
			Fqdn:             s["fqdn"].(string),
			OriginHostHeader: s["origin_hostheader"].(string),
			OriginBackend:    s["origin_backend"].(string),
			ProviderRegion:   s["provider_region"].(string),
			ProviderName:     s["provider_name"].(string),
		}

		if c, ok := s["cert_id"].(string); ok {
			isr.CertId = c
		}

		ireq.Services = append(ireq.Services, isr)
	}

	res, err := c.CreateInstance(ireq)
	if err != nil {
		return err
	}

	d.SetId(res.Instance.Id)

	return resourceEceRead(d, meta)
}

func resourceEceRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*go_sdk.NorainaApiClient)
	res, err := c.GetInstance(d.Id())
	if err != nil {
		d.SetId("")
		return err
	}

	d.Set("name", res.Instance.Name)
	d.Set("created_date", res.Instance.CreatedDate)
	d.Set("tenant_id", res.Instance.TenantId)
	d.Set("version", res.Instance.Version)

	if err = d.Set("service", flattenServices(res.Instance.Services)); err != nil {
		return fmt.Errorf("[ERROR] Error setting `services`: %+v", err)
	}

	if err = d.Set("server", flattenServers(res.Servers)); err != nil {
		return fmt.Errorf("[ERROR] Error setting `servers`: %+v", err)
	}

	return nil
}

func resourceEceUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*go_sdk.NorainaApiClient)

	services := d.Get("service").([]interface{})
	if len(services) == 0 {
		return errors.New("[ERROR] you must define at least 1 service in the instance")
	}

	ireq := go_sdk.InstanceRequest{
		Name:     d.Get("name").(string),
		Services: make([]go_sdk.InstanceServiceRequest, 0),
	}

	for _, service := range services {
		s := service.(map[string]interface{})

		isr := go_sdk.InstanceServiceRequest{
			Name:             s["name"].(string),
			Fqdn:             s["fqdn"].(string),
			OriginHostHeader: s["origin_hostheader"].(string),
			OriginBackend:    s["origin_backend"].(string),
			ProviderRegion:   s["provider_region"].(string),
			ProviderName:     s["provider_name"].(string),
		}

		if c, ok := s["cert_id"].(string); ok {
			isr.CertId = c
		}

		ireq.Services = append(ireq.Services, isr)
	}

	_, err := c.UpdateInstance(d.Id(), ireq)
	if err != nil {
		d.SetId("")
		return err
	}

	return resourceEceRead(d, meta)
}

func resourceEceDelete(d *schema.ResourceData, meta interface{}) error {
	return meta.(*go_sdk.NorainaApiClient).DeleteInstance(d.Id())
}

func flattenServices(sf []go_sdk.ServiceFields) []map[string]interface{} {
	services := make([]map[string]interface{}, 0, len(sf))

	for _, so := range sf {
		s := make(map[string]interface{})
		s["id"] = so.Id
		s["created_date"] = so.CreatedDate
		s["fqdn"] = so.Fqdn
		s["name"] = so.Name
		s["origin_backend"] = so.OriginBackend
		s["origin_hostheader"] = so.OriginHostHeader
		s["provider_name"] = so.ProviderName
		s["provider_region"] = so.ProviderRegion
		s["cert_id"] = so.CertId

		services = append(services, s)
	}

	return services
}

func flattenServers(sf []go_sdk.ServerFields) []map[string]interface{} {
	servers := make([]map[string]interface{}, 0, len(sf))

	for _, so := range sf {
		s := make(map[string]interface{})

		s["id"] = so.Id
		s["name"] = so.Name
		s["status"] = so.Status
		s["container_id"] = so.ContainerID
		s["instance_id"] = so.InstanceID
		s["host_id"] = so.HostID
		s["created_date"] = so.CreatedDate
		s["version"] = so.Version

		dnsl := make([]map[string]interface{}, 0, 1)
		dns := make(map[string]interface{})
		dns["a"] = so.DNSRecordID.A
		dns["aaaa"] = so.DNSRecordID.AAAA
		dnsl = append(dnsl, dns)
		s["dns_record_id"] = dnsl

		ipl := make([]map[string]interface{}, 0, 1)
		ip := make(map[string]interface{})
		ip["ipv4"] = so.IPAdress.IPv4
		ip["ipv6"] = so.IPAdress.IPv6
		ipl = append(ipl, ip)
		s["ip_address"] = ipl

		servers = append(servers, s)
	}

	return servers
}
