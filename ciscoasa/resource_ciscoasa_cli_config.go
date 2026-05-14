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
			"commands_hash": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCiscoASACLIConfigCreate(d *schema.ResourceData, meta interface{}) error {
	return applyCiscoASACLIConfig(d, meta)
}

func resourceCiscoASACLIConfigRead(d *schema.ResourceData, meta interface{}) error {
	expectedHash, err := calculateCiscoASACLIConfigHash(d)
	if err != nil {
		return err
	}

	if d.Id() == "" {
		return nil
	}

	if d.Id() != expectedHash {
		d.SetId("")
		return nil
	}

	if err := d.Set("commands_hash", expectedHash); err != nil {
		return err
	}

	return nil
}

func resourceCiscoASACLIConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("commands") ||
		d.HasChange("save_config") ||
		d.HasChange("configure_terminal") {
		return applyCiscoASACLIConfig(d, meta)
	}

	return resourceCiscoASACLIConfigRead(d, meta)
}

func resourceCiscoASACLIConfigDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func applyCiscoASACLIConfig(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ciscoasa.Client)

	commands, err := buildCiscoASACLICommands(d)
	if err != nil {
		return err
	}

	if err := client.PostCLI(commands); err != nil {
		return err
	}

	hash := hashCiscoASACLICommands(commands)

	d.SetId(hash)

	if err := d.Set("commands_hash", hash); err != nil {
		return err
	}

	return nil
}

func calculateCiscoASACLIConfigHash(d *schema.ResourceData) (string, error) {
	commands, err := buildCiscoASACLICommands(d)
	if err != nil {
		return "", err
	}

	return hashCiscoASACLICommands(commands), nil
}

func buildCiscoASACLICommands(d *schema.ResourceData) ([]string, error) {
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
		return nil, fmt.Errorf("no CLI commands to apply")
	}

	return commands, nil
}

func hashCiscoASACLICommands(commands []string) string {
	hash := sha1.Sum([]byte(strings.Join(commands, "\n")))
	return hex.EncodeToString(hash[:])
}
