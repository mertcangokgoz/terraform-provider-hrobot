package server

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

const DataSourceType = "hrobot_server"

// DataSource creates a new Terraform schema for the hrobot_server resource.
func DataSource() *schema.Resource {
	return &schema.Resource{

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"server_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datacenter": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv4_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
