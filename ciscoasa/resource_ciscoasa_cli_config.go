package ciscoasa

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/CiscoDevNet/go-ciscoasa/ciscoasa"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCiscoASACLIConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceCiscoASACLIConfigCreate,
		Read:   resourceCiscoASACLIConfigRead,
		Update: resourceCiscoASACLIConfigUpdate,
		Delete: resourceCiscoASACLIConfigDelete,

		Schema: map[string]*schema.Schema{
			"commands": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"save_config": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"configure_terminal": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceCiscoASACLIConfigCreate(d *schema.ResourceData, meta interface{}) error {
	return applyCiscoASACLIConfig(d, meta)
}

func resourceCiscoASACLIConfigRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCiscoASACLIConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	return applyCiscoASACLIConfig(d, meta)
}

func resourceCiscoASACLIConfigDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func applyCiscoASACLIConfig(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ciscoasa.Client)

	rawCommands := d.Get("commands").([]interface{})
	commands := make([]string, 0, len(rawCommands)+2)

	if d.Get("configure_terminal").(bool) {
		commands = append(commands, "configure terminal")
	}

	for _, raw := range rawCommands {
		cmd := strings.TrimSpace(raw.(string))
		if cmd != "" {
			commands = append(commands, cmd)
		}
	}

	if d.Get("save_config").(bool) {
		commands = append(commands, "write memory")
	}

	if len(commands) == 0 {
		return fmt.Errorf("no CLI commands to apply")
	}

	if err := client.PostCLI(commands); err != nil {
		return err
	}

	hash := sha1.Sum([]byte(strings.Join(commands, "\n")))
	d.SetId(hex.EncodeToString(hash[:]))

	return nil
}
