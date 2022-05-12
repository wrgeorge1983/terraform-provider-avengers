package avengers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-provider-avengers/avengers/aclient"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"avengers_resource": resourceAvengers(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"avengers_datasource": dataSourceAvengers(),
		},
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AVENGERS_BACKEND_HOST_URL", "http://localhost:8000"),
			},
		},
		ConfigureContextFunc: providerConfigure,
	}
}

type ApiClient struct {
	data           *schema.ResourceData
	avengersClient *aclient.Client
}

func (a *ApiClient) NewAvengersClient() (*aclient.Client, error) {
	host := a.data.Get("host").(string)
	c, err := aclient.NewClient(&host)
	if err != nil {
		return c, err
	}
	return c, nil
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	c := &ApiClient{data: d}
	client, err := c.NewAvengersClient()
	if err != nil {
		return c, diag.FromErr(err)
	}
	c.avengersClient = client
	return c, nil
}
