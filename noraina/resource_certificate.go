package noraina

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/norainacloud/terraform-provider-noraina/go-sdk"
	"io"
	"os"
)

func resourceCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceCertificateCreate,
		Read:   resourceCertificateRead,
		Delete: resourceCertificateDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cert": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				StateFunc: hashCertPart,
			},
			"key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ForceNew:  true,
				StateFunc: hashCertPart,
			},
			"chain": {
				Type:      schema.TypeString,
				Required:  false,
				Optional:  true,
				ForceNew:  true,
				StateFunc: hashCertPart,
			},
			"cert_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*go_sdk.NorainaApiClient)

	iCert := go_sdk.CertificateCreateRequest{
		Name: d.Get("name").(string),
		Cert: d.Get("cert").(string),
		Key:  d.Get("key").(string),
	}

	if chain, ok := d.Get("chain").(string); ok {
		iCert.Chain = chain
	}

	res, err := c.CreateCertificate(iCert)
	if err != nil {
		return err
	}

	d.SetId(res.Data.CertId)

	return resourceCertificateRead(d, meta)
}

func resourceCertificateRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*go_sdk.NorainaApiClient)

	res, err := c.GetCertificate(d.Id())
	if err != nil {
		d.SetId("")
		return err
	}

	d.Set("name", res.Data.Name)
	d.Set("created_date", res.Data.CreatedDate)

	return nil
}

func resourceCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	return meta.(*go_sdk.NorainaApiClient).DeleteCertificate(d.Id())
}

func hashCertPart(cert interface{}) string {
	if cert == nil || cert == (*string)(nil) {
		return ""
	}

	var rawCert string
	switch cert.(type) {
	case string:
		rawCert = cert.(string)
	case *string:
		rawCert = *cert.(*string)
	default:
		return ""
	}

	file, err := os.Open(rawCert)
	if err != nil {
		return ""
	}
	defer file.Close()

	h := md5.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(h.Sum(nil))
}
