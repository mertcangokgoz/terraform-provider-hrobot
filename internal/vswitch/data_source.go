package vswitch

import (
	"context"
	"strconv"
	"time"

	"github.com/floshodan/hrobot-go/hrobot"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	DataSourceType     = "hrobot_vswitch"
	DataSourceListType = "hrobot_vswitch_list"
)

func getCommonDataSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vlan": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"cancelled": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	}

}

var serverResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"server_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"server_ipv6_net": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"server_number": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
	},
}

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHrobotVswitchRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vlan": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cancelled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"server": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     serverResource,
			},
			"subnet": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mask": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"gateway": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"cloud_network": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mask": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"gateway": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DataSourceList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHrobotVswitchReadList,

		Schema: map[string]*schema.Schema{
			"vswitches": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: getCommonDataSchema(),
				},
			},
		},
	}
}

func dataSourceHrobotVswitchReadList(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hrobot.Client)

	var diags diag.Diagnostics

	allVswitches, _, err := client.VSwitch.GetVSwitchList(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	tfVswitches := make([]map[string]interface{}, len(allVswitches))
	for i, vswitch := range allVswitches {
		tfVswitches[i] = getVSwitchAttributes(vswitch)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	if err := d.Set("vswitches", tfVswitches); err != nil {
		return diag.FromErr(err)
	}

	return diags

}

func dataSourceHrobotVswitchRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hrobot.Client)

	vswitchID := d.Get("id").(string)

	vswitch, _, err := client.VSwitch.GetVSwitchById(ctx, vswitchID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("id").(string))

	if err := d.Set("id", strconv.Itoa(vswitch.ID)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", vswitch.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("vlan", vswitch.Vlan); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("cancelled", vswitch.Cancelled); err != nil {
		return diag.FromErr(err)
	}

	var serverList []interface{}
	for _, s := range vswitch.Server {
		server := s.(map[string]interface{})
		serverMap := map[string]interface{}{
			"server_ip":       server["server_ip"].(string),
			"server_ipv6_net": server["server_ipv6_net"].(string),
			"server_number":   strconv.Itoa(server["server_number"].(int)),
			"status":          server["status"].(string),
		}
		serverList = append(serverList, serverMap)
	}
	d.Set("server", serverList)

	var subnets []interface{}
	for _, s := range vswitch.Subnet {
		subnet := s.(map[string]interface{})
		subnets = append(subnets, map[string]interface{}{
			"ip":      subnet["ip"].(string),
			"mask":    subnet["mask"].(string),
			"gateway": subnet["gateway"].(string),
		})
	}
	/*
		if err := d.Set("subnets", subnets); err != nil {
			return diag.FromErr(err)
		}

		var cloudNetworks []interface{}
		for _, c := range vswitch.CloudNetwork {
			cloudnetwork := c.(map[string]interface{})
			cloudNetworks = append(cloudNetworks, map[string]interface{}{
				"id":      cloudnetwork["id"].(string),
				"ip":      cloudnetwork["ip"].(string),
				"mask":    cloudnetwork["mask"].(string),
				"gateway": cloudnetwork["gateway"].(string),
			})
		}
		if err := d.Set("cloud_networks", cloudNetworks); err != nil {
			return diag.FromErr(err)
		}
	*/

	return nil
}
