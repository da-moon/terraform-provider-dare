package dare

import (
	//	"github.com/da-moon/go-dare/model"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceEncryptArtifact() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Read: dataSourceEncryptArtifactRead,
	}
}
func dataSourceEncryptArtifactRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}
