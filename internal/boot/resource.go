package boot

import (
	"context"
	"strconv"

	"github.com/floshodan/hrobot-go/hrobot"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const ResourceType = "hrobot_boot_rescue"

func resourceBootRescuev2() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"server_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_ipv6_net": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"os": {
				Type:     schema.TypeString,
				Required: true,
			},
			"arch": {
				Type:     schema.TypeInt,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authorized_key": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"fingerprint": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"size": {
							Type:     schema.TypeFloat,
							Required: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"host_key": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"boot_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBootRescueCreate,
		ReadContext:   resourceBootRescueRead,
		UpdateContext: resourceBootRescueUpdate,
		DeleteContext: resourceBootRescueDelete,
		Schema: map[string]*schema.Schema{
			"rescue": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     resourceBootRescuev2(),
			},
		},
	}
}

func resourceBootRescueCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hrobot.Client)

	serverID := d.Get("server_id").(string)

	rescueOpts := &hrobot.RescueOpts{
		OS: d.Get("os").(string),
	}

	if v, ok := d.GetOk("authorized_key"); ok {
		rescueOpts.Authorized_Key = v.(string)
	}

	if v, ok := d.GetOk("keyboard"); ok {
		rescueOpts.Keyboard = v.(string)
	}

	rescue, _, err := client.Boot.ActivateRescue(ctx, serverID, rescueOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(rescue.ServerNumber))

	return nil
}
func resourceBootRescueRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceBootRescueUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceBootRescueDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func setRescueSchema(d *schema.ResourceData, bootrescue *hrobot.BootRescue) {

}

func setBootListSchema(d *schema.ResourceData, i *hrobot.BootList) {
	if i == nil {
		d.SetId("")
		d.Set("rescue", nil)
		d.Set("linux", nil)
		d.Set("vnc", nil)
		d.Set("windows", nil)
		d.Set("plesk", nil)
		d.Set("cpanel", nil)
		return
	}
	d.SetId(i.Rescue.ServerIP)
	d.Set("rescue", flattenBootRescue(&i.Rescue))
	d.Set("linux", flattenBootLinux(&i.Linux))
}

func flattenBootRescue(rescue *hrobot.BootRescue) []interface{} {
	/*
		var authorizedKeys []map[string]interface{}
		for _, key := range rescue.AuthorizedKey {
			if keyMap, ok := key.(map[string]interface{}); ok {
				if keyVal, ok := keyMap["key"].(map[string]interface{}); ok {
					authorizedKeys = append(authorizedKeys, keyVal)
				}
			}
		}
	*/

	os := []interface{}{}
	if rescue.Active {
		os = []interface{}{rescue.Os.(string)}
	} else {
		os = rescue.Os.([]interface{})
	}
	return []interface{}{
		map[string]interface{}{
			"server_ip":       rescue.ServerIP,
			"server_ipv6_net": rescue.ServerIpv6Net,
			"server_id":       rescue.ServerNumber,
			"active":          rescue.Active,
			"os":              os,
			"password":        rescue.Password,
			"authorized_key":  flattenAuthorizedKeys(rescue.AuthorizedKey),
		},
	}
}

func flattenAuthorizedKeys(keys []interface{}) []interface{} {
	res := make([]interface{}, 0)
	for _, key := range keys {
		if keyMap, ok := key.(map[string]interface{}); ok {
			if keyVal, ok := keyMap["key"].(map[string]interface{}); ok {
				res = append(res, map[string]interface{}{
					"name":        keyVal["name"].(string),
					"fingerprint": keyVal["fingerprint"].(string),
					"type":        keyVal["type"].(string),
					"size":        keyVal["size"].(float64),
					"created_at":  keyVal["created_at"].(string),
				})
			}
		}
	}
	return res
}

func flattenBootLinux(linux *hrobot.BootLinux) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"server_ip":       linux.ServerIP,
			"server_ipv6_net": linux.ServerIpv6Net,
			"server_number":   linux.ServerNumber,
			"active":          linux.Active,
			"dist":            linux.Dist,
		},
	}
}

func flattenBootVNC(vnc hrobot.BootVnc) []interface{} {
	return []interface{}{}
}

func flattenBootWindows(windows hrobot.BootWindows) map[string]interface{} {
	return map[string]interface{}{}
}

func setBootRescueSchema(d *schema.ResourceData, rescue *hrobot.BootRescue) {
}
