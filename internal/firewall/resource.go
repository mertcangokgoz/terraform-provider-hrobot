package firewall

import (
	"context"
	"strconv"

	"github.com/floshodan/hrobot-go/hrobot"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const ResourceType = "hrobot_firewall"

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFirewallCreate,
		ReadContext:   resourceFirewallRead,
		UpdateContext: resourceFirewallUpdate,
		DeleteContext: resourceFirewallDelete,
		/*
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		*/
		Schema: map[string]*schema.Schema{
			"server_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"server_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"active": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filter_ipv6": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"whitelist_hos": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_version": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dst_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dst_port": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"src_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"src_port": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tcp_flags": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"action": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceFirewallCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hrobot.Client)

	server_ip := d.Get("server_ip").(string)

	firewallRules := make([]hrobot.InputRule, 0)
	for _, ruleMap := range d.Get("rules").([]interface{}) {
		ruleProperties := ruleMap.(map[string]interface{})
		firewallRules = append(firewallRules, hrobot.InputRule{
			IPVersion: ruleProperties["ip_version"].(string),
			Name:      ruleProperties["name"].(string),
			SrcIP:     ruleProperties["src_ip"].(string),
			DstPort:   ruleProperties["dst_port"].(string),
			Protocol:  ruleProperties["protocol"].(string),
			TCPFlags:  ruleProperties["tcp_flags"].(string),
			Action:    ruleProperties["action"].(string),
		})
	}

	fwOps := &hrobot.FirewallOps{
		Status:        "active",
		Whitelist_hos: "true",
		Filter_IPv6:   "true",
		Rules:         firewallRules,
	}
	fw, _, err := client.Firewall.PostFirewallByServernumber(ctx, server_ip, fwOps)

	if err != nil {
		return diag.FromErr(err)
	}

	setFirewallSchema(d, fw)

	return resourceFirewallRead(ctx, d, m)
	/*
		var diags diag.Diagnostics

		return diags
	*/
}

func resourceFirewallRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hrobot.Client)

	if id, ok := d.GetOk("server_ip"); ok {
		i, _, err := client.Firewall.GetFirewallByServernumber(ctx, id.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		if i == nil {
			return diag.Errorf("no firewall found with id %d", id)
		}
		setFirewallSchema(d, i)
		return nil
	}

	return nil
	/*
		client := m.(*hrobot.Client)
		firewall, _, err := client.Firewall.GetFirewallByServernumber(context.Background(), d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		if firewall == nil {
			d.SetId("")
			return nil
		}
		var diags diag.Diagnostics

		return diags
	*/
}

func resourceFirewallUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceFirewallCreate(ctx, d, m)
}

func resourceFirewallDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func setFirewallSchema(d *schema.ResourceData, firewall *hrobot.Firewall) {
	d.SetId(strconv.Itoa(firewall.ServerNumber))
	d.Set("server_ip", firewall.ServerIP)
	d.Set("server_number", firewall.ServerNumber)
	d.Set("active", firewall.Status)
	d.Set("whitelist_hos", firewall.WhitelistHos)
	d.Set("port", firewall.Port)

	inputRules := make([]map[string]interface{}, len(firewall.Rules.Input))
	for i, inputRule := range firewall.Rules.Input {
		inputRules[i] = map[string]interface{}{
			"ip_version": inputRule.IPVersion,
			"name":       inputRule.Name,
			"dst_ip":     inputRule.DstIP,
			"src_ip":     inputRule.SrcIP,
			"dst_port":   inputRule.DstPort,
			"src_port":   inputRule.SrcPort,
			"protocol":   inputRule.Protocol,
			"action":     inputRule.Action,
		}
	}

	d.Set("rule", inputRules)
}
