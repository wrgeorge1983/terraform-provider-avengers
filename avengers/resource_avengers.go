package avengers

import (
	"context"
	"log"

	aclient "github.com/sourav977/avengers-client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAvengers() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAvengersCreate,
		ReadContext:   resourceAvengersRead,
		UpdateContext: resourceAvengersUpdate,
		DeleteContext: resourceAvengersDelete,
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
							Description: "full name of avenger",
						},
						"alias": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "any alias/nickname of avenger",
						},
						"weapon": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "his/her special weapons",
						},
					},
				},
			},
			"_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the _id value returned from mongodb",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "full name of avenger",
			},
			"alias": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "any alias/nickname of avenger",
			},
			"weapon": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "his/her special weapons",
			},
			"deleted_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "deleted item count",
			},
			"matched_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "total matched item found",
			},
			"modified_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "total item modified",
			},
			"upserted_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "total item upserted",
			},
		},
	}
}

func resourceAvengersCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning resourceAvengersCreate", d.Id())
	var diags diag.Diagnostics
	c := m.(*ApiClient)

	// collect values from schema
	name := d.Get("name").(string)
	alias := d.Get("alias").(string)
	weapon := d.Get("weapon").(string)

	// build API object
	a := aclient.Avenger{
		Name:   name,
		Alias:  alias,
		Weapon: weapon,
	}

	// submit object to API
	res, err := c.avengersClient.CreateAvenger(a)
	if err != nil {
		return diag.FromErr(err)
	}

	// marshall response to schema
	if err := d.Set("_id", res.ID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", res.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("alias", res.Alias); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("weapon", res.Weapon); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(res.ID)
	log.Printf("[DEBUG] %s: resourceAvengersCreate finished successfully", d.Id())
	return diags
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

func resourceAvengersUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	log.Printf("[DEBUG] %s: Beginning resourceAvengersUpdate", d.Id())
	var diags diag.Diagnostics
	c := m.(*ApiClient)

	name := d.Get("name").(string)
	alias := d.Get("alias").(string)
	weapon := d.Get("weapon").(string)

	a := aclient.Avenger{
		Name:   name,
		Alias:  alias,
		Weapon: weapon,
	}
	res, err := c.avengersClient.UpdateAvengerByName(a)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("matched_count", res.MatchedCount); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("modified_count", res.ModifiedCount); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("upserted_count", res.UpsertedCount); err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: resourceAvengersUpdate finished successfully", d.Id())
	return diags
}

func resourceAvengersDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	log.Printf("[DEBUG] %s: Beginning resourceAvengersDelete", d.Id())
	var diags diag.Diagnostics
	c := m.(*ApiClient)
	name := d.Get("name").(string)
	del, err := c.avengersClient.DeleteAvengerByName(name)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("deleted_count", del.DeletedCount); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[DEBUG] %s: resourceAvengersDelete finished successfully", d.Id())
	return diags
}

func flattenAvengers(avengersList *[]aclient.Avenger) []interface{} {
	if avengersList == nil {
		return make([]interface{}, 0)
	}
	avengers := make([]interface{}, len(*avengersList))
	for i, avenger := range *avengersList {
		item := make(map[string]interface{})

		item["_id"] = avenger.ID
		item["name"] = avenger.Name
		item["alias"] = avenger.Alias
		item["weapon"] = avenger.Weapon

		avengers[i] = item
	}
	return avengers
}
