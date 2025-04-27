package boot

import (
	"context"

	"github.com/floshodan/hrobot-go/hrobot"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const DataSourceType = "hrobot_boot"
const DataSourceTypeRescue = "hrobot_rescue"
const DataSourceListType = "hrobot_boot_list"

func rescueSchema() *schema.Schema {
	return &schema.Schema{}
}

func DataSourceList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHrobotBootRead,
		Schema: map[string]*schema.Schema{
			"server_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rescue": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     resourceBootRescue(),
			},
			"linux": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     resourceBootLinux(),
			},
			"vnc": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     resourceBootVnc(),
			},
			"windows": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"plesk": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"cpanel": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceBootVnc() *schema.Resource {
	return &schema.Resource{}
}

/*func resourceBootRescue() *schema.Resource {
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
				Optional: true,
			},
			"os": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"arch": {
				Type:     schema.TypeString,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"authorized_key": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
*/

func resourceBootLinux() *schema.Resource {
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
			"server_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"dist": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		}}
}

func dataSourceHrobotBootRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hrobot.Client)

	if id, ok := d.GetOk("server_id"); ok {
		i, _, err := client.Boot.GetBootOptions(ctx, id.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		if i == nil {
			return diag.Errorf("no server found with id %d", id)
		}
		setBootListSchema(d, i)
		return nil
	}
	return nil
}
