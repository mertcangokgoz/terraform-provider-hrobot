package sshkey

import (
	"context"
	"strconv"
	"time"

	"github.com/floshodan/hrobot-go/hrobot"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const DataSourceType = "hrobot_ssh_key"
const DataSourceListType = "hrobot_ssh_keys"

// getCommonDataSchema returns a new common schema used by all ssh key data sources.
func getCommonDataSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"fingerprint": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"size": {
			Type:     schema.TypeInt,
			Computed: true,
		},

		"data": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

// DataSource creates a new Terraform schema for the hcloud_ssh_key data
// source.
func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHrobotSSHKeyRead,
		Schema:      getCommonDataSchema(),
	}
}

func DataSourceList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHrobotSSHKeyReadList,
		Schema: map[string]*schema.Schema{
			"ssh_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: getCommonDataSchema(),
				},
			},
		},
	}
}

func dataSourceHrobotSSHKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*hrobot.Client)

	if fingerprint, ok := d.GetOk("fingerprint"); ok {
		s, _, err := client.SSHKey.GetByFingerprint(ctx, fingerprint.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		if s == nil {
			return diag.Errorf("no sshkey found with fingerprint %v", fingerprint)
		}
		setSSHKeySchema(d, s)
		return nil
	}
	return diag.Errorf("please specifiy a fingerprint to lookup the sshkey")
}

func dataSourceHrobotSSHKeyReadList(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*hrobot.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	allKeys, _, err := client.SSHKey.List(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	tfKeys := make([]map[string]interface{}, len(allKeys))
	for i, key := range allKeys {
		tfKeys[i] = getSSHKeyAttributes(key)
	}

	//set id to time_now (internal terraform id)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	if err := d.Set("ssh_keys", tfKeys); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
