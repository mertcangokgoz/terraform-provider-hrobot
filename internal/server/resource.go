package server

import (
	"context"
	"crypto/sha1"
	"fmt"
	"strconv"
	"strings"

	"github.com/floshodan/hrobot-go/hrobot"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceType is the type name of the Hetzner Robot Server resource.
const ResourceType = "hrobot_server"

func Resource() *schema.Resource {
	return &schema.Resource{
		Create: resourceEventCreate,
		Read:   resourceEventRead,
		Update: resourceEventUpdate,
		Delete: resourceEventDelete,

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
				Computed: false,
			},
			"datacenter": {
				Type:     schema.TypeString,
				Computed: false,
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

func getServerAttributes(d *schema.ResourceData, s *hrobot.Server) map[string]interface{} {

	res := map[string]interface{}{
		"id":           s.ServerNumber,
		"name":         s.ServerName,
		"datacenter":   s.Dc,
		"status":       s.Status,
		"server_type":  s.Product,
		"ipv4_address": s.ServerIP,
		"ipv6_address": s.ServerIpv6Net,
		"subnet":       s.Subnet,
	}

	return res
}

func resourceEventCreate(d *schema.ResourceData, m interface{}) error {
	// TODO
	return nil
}

func resourceEventRead(d *schema.ResourceData, m interface{}) error {
	// TODO
	client := m.(*hrobot.Client)
	serverlist, _, err := client.Server.List(context.Background())

	if err != nil {
		return err
	}

	ids := make([]string, len(serverlist))
	tfServers := make([]map[string]interface{}, len(serverlist))
	for i, server := range serverlist {
		ids[i] = strconv.Itoa(server.ServerNumber)
		tfServers[i] = getServerAttributes(d, server)
		d.Set("id", server.ServerNumber)
	}
	d.Set("servers", tfServers)
	d.SetId(fmt.Sprintf("%x", sha1.Sum([]byte(strings.Join(ids, "")))))

	return nil
}

func resourceEventUpdate(d *schema.ResourceData, meta interface{}) error {
	// TODO
	return nil
}

func resourceEventDelete(d *schema.ResourceData, meta interface{}) error {
	// TODO
	return nil
}
