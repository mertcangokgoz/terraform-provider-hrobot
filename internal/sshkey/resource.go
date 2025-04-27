package sshkey

import (
	"context"
	"log"

	"github.com/floshodan/hrobot-go/hrobot"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const ResourceType = "hrobot_ssh_key"

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSSHKeyCreate,
		ReadContext:   resourceSSHKeyRead,
		UpdateContext: resourceSSHKeyUpdate,
		DeleteContext: resourceSSHKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSSHKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hrobot.Client)
	data := &hrobot.CreateKeyOpts{
		Name: d.Get("name").(string),
		Data: d.Get("public_key").(string),
	}
	new_key, resp, err := client.SSHKey.Create(ctx, data)
	if resp.StatusCode == 409 {
		return diag.Errorf("Key already exists")
	}

	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("fingerprint", new_key.Fingerprint)
	d.SetId(new_key.Fingerprint)

	return resourceSSHKeyRead(ctx, d, m)
}

func resourceSSHKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hrobot.Client)

	sshKeyID := d.Id()

	sshKey, _, err := client.SSHKey.GetByFingerprint(ctx, sshKeyID)
	if err != nil {
		return diag.FromErr(err)
	}
	if sshKey == nil {
		log.Printf("[WARN] SSH key (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	setSSHKeySchema(d, sshKey)

	return nil
}

func resourceSSHKeyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hrobot.Client)

	sshKeyID := d.Id()

	if d.HasChange("name") {
		name := d.Get("name").(string)
		_, _, err := client.SSHKey.Update(ctx, sshKeyID, &hrobot.UpdateKeyOpts{
			Name: name,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSSHKeyRead(ctx, d, m)
}

func resourceSSHKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hrobot.Client)

	sshKeyID := d.Id()

	_, err := client.SSHKey.Delete(ctx, sshKeyID)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil

}

func setSSHKeySchema(d *schema.ResourceData, s *hrobot.SSHKey) {
	for key, val := range getSSHKeyAttributes(s) {
		if key == "fingerprint" {
			d.SetId(val.(string))
		}
		d.Set(key, val)
	}
}

func getSSHKeyAttributes(s *hrobot.SSHKey) map[string]interface{} {
	return map[string]interface{}{
		"name":        s.Name,
		"fingerprint": s.Fingerprint,
		"type":        s.Type,
		"size":        s.Size,
		"data":        s.Data,
	}
}
