package avengers

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAvengers() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceAvengersRead,
		Schema: map[string]*schema.Schema{
			"avengers": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the _id value returned from mongodb",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of avenger",
						},
						"alias": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "any alias of avenger",
						},
						"weapon": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "avenger's special weapon",
						},
					},
				},
			},
		},
	}
}

func resourceAvengersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning resourceAvengersRead", d.Id())
	var diags diag.Diagnostics
	c := m.(*ApiClient)
	// no values to pull from schema
	// no api object
	// call api
	res, err := c.avengersClient.GetAllAvengers()
	if err != nil {
		return diag.FromErr(err)
	}

	// marshall response to schema
	if res == nil {
		return diag.Errorf("no avengers found in database")
	}
	resItems := flattenAvengers(&res)
	if err := d.Set("avengers", resItems); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: resourceAvengersRead finished successfully", d.Id())
	return diags
}
