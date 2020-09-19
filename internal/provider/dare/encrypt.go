package dare

func dataSourceEncryptArtifact() *schema.Resource {
	model.EncryptRequest{

	}
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