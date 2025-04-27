package vswitch

import (
	"context"
	"strconv"

	"github.com/floshodan/hrobot-go/hrobot"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const ResourceType = "hrobot_vswitch"

func getVSwitchAttributes(s *hrobot.VSwitch) map[string]interface{} {
	return map[string]interface{}{
		"id":        s.ID,
		"name":      s.Name,
		"vlan":      s.Vlan,
		"cancelled": s.Cancelled,
	}
}

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVSwitchCreate,
		ReadContext:   resourceVSwitchRead,
		UpdateContext: resourceVSwitchUpdate,
		DeleteContext: resourceVSwitchDelete,

		Importer: &schema.ResourceImporter{
			State: resourceServerVSwitchImportState,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceServerVSwitchImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	serverID := d.Get("server_id").(int)
	vSwitchID := d.Get("vswitch_id").(int)
	ID := strconv.Itoa(serverID) + "-" + strconv.Itoa(vSwitchID)
	d.SetId(ID)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceVSwitchCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hrobot.Client)
	name := d.Get("name").(string)
	vlan := d.Get("vlan_id").(int)

	opt := &hrobot.AddvSwitchOps{
		Name:    name,
		Vlan_ID: vlan,
	}

	vSwitch, _, err := client.VSwitch.AddVSwitch(ctx, opt)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(vSwitch.ID))

	return resourceVSwitchRead(ctx, d, m)

}

func resourceVSwitchRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hrobot.Client)

	vSwitchID := d.Id()

	vSwitch, _, err := client.VSwitch.GetVSwitchById(ctx, vSwitchID)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.Set("name", vSwitch.Name)
	d.Set("vlan_id", vSwitch.Vlan)
	//d.Set("cancelled", vSwitch.Cancelled)

	return nil
}

func resourceVSwitchUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hrobot.Client)

	vSwitchID := d.Id()

	opt := &hrobot.AddvSwitchOps{
		Name:    d.Get("name").(string),
		Vlan_ID: d.Get("vlan").(int),
	}

	_, _, err := client.VSwitch.UpdateVSwitchById(ctx, vSwitchID, opt)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceVSwitchRead(ctx, d, m)
}

func resourceVSwitchDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hrobot.Client)

	vSwitchID := d.Id()

	_, err := client.VSwitch.CancelVSwitch(ctx, vSwitchID, "")
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
