package avengers

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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
