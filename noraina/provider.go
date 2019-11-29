package noraina

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	n "github.com/norainacloud/go-client-noraina"
)

func Provider() terraform.ResourceProvider {
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
	c := n.NewClient(nil)

	err := c.GetAuthToken(context.Background(), &n.AuthRequest{
		Email:    d.Get("email").(string),
		Password: d.Get("password").(string),
	})

	if err != nil {
		return nil, err
	}

	return c, nil
}
