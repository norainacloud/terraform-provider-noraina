package noraina

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/norainacloud/terraform-provider-noraina/go-sdk"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NORAINA_EMAIL", nil),
				Description: descriptions["email"],
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NORAINA_PASSWORD", nil),
				Description: descriptions["password"],
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"noraina_ece":         resourceEce(),
			"noraina_certificate": resourceCertificate(),
		},
		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"email":    "Your Noraina email",
		"password": "Your Noraina password",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return go_sdk.NewNorainaApiClient(
		d.Get("email").(string),
		d.Get("password").(string),
	)
}
