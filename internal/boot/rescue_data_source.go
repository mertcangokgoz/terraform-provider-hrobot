package boot

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func resourceBootRescue() *schema.Resource {
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
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
				Optional: true,
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
