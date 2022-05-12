package avengers

import (
	"context"
	"log"

	"terraform-provider-avengers/avengers/aclient"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAvenger() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAvengerCreate,
		ReadContext:   resourceAvengerRead,
		UpdateContext: resourceAvengerUpdate,
		DeleteContext: resourceAvengerDelete,
		Schema: map[string]*schema.Schema{
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

func resourceAvengerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning resourceAvengerCreate", d.Id())
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

	d.SetId(res.ID)
	resourceAvengerRead(ctx, d, m)
	log.Printf("[DEBUG] %s: resourceAvengerCreate finished successfully", d.Id())
	return diags
}

func resourceAvengerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning resourceAvengerRead", d.Id())
	var diags diag.Diagnostics
	c := m.(*ApiClient)
	// no values to pull from schema
	// no api object
	// call api
	res, err := c.avengersClient.GetAvengerById(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// marshall response to schema
	if res == nil {
		return diag.Errorf("no avengers found in database")
	}
	avengerMap := flattenAvenger(res)

	for k, v := range *avengerMap {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[DEBUG] %s: resourceAvengerRead finished successfully", d.Id())
	return diags
}

func resourceAvengerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	log.Printf("[DEBUG] %s: Beginning resourceAvengerUpdate", d.Id())
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
	log.Printf("[DEBUG] %s: resourceAvengerUpdate finished successfully", d.Id())
	return diags
}

func resourceAvengerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning resourceAvengerDelete", d.Id())
	var diags diag.Diagnostics
	c := m.(*ApiClient)
	name := d.Get("name").(string)
	del, err := c.avengersClient.DeleteAvengerByName(name)
	if err != nil {
		return diag.FromErr(err)
	} else if del.DeletedCount < 1 {
		return diag.Errorf("deleting %s failed.  Avenger by that name may not have been found", name)
	}
	if err := d.Set("deleted_count", del.DeletedCount); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[DEBUG] %s: resourceAvengerDelete finished successfully", d.Id())
	return diags
}

func flattenAvenger(avenger *aclient.Avenger) *map[string]interface{} {
	item := make(map[string]interface{})

	if avenger == nil {
		return &item //TODO: is this right?
	}
	item["_id"] = avenger.ID
	item["name"] = avenger.Name
	item["alias"] = avenger.Alias
	item["weapon"] = avenger.Weapon

	return &item
}

func flattenAvengers(avengersList *[]aclient.Avenger) []interface{} {
	if avengersList == nil {
		return make([]interface{}, 0)
	}
	avengers := make([]interface{}, len(*avengersList))
	for i, avenger := range *avengersList {
		item := flattenAvenger(&avenger)

		avengers[i] = *item
	}
	return avengers
}
