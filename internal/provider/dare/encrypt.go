package dare

import (
	dare "github.com/da-moon/go-dare"
	model "github.com/da-moon/go-dare/model"
	stacktrace "github.com/palantir/stacktrace"

	schema "github.com/hashicorp/terraform/helper/schema"
)

func resourceEncryptArtifact() *schema.Resource {
	return &schema.Resource{
		Create: resourceEncryptArtifactCreate,
		Read:   resourceEncryptArtifactRead,
		Update: resourceEncryptArtifactUpdate,
		Delete: resourceEncryptArtifactDelete,
		Schema: map[string]*schema.Schema{
			"path": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "path of a directory or a single file for encryption",
			},
			"output_dir": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "root directory in which the encrypted files are stored",
			},
			"regex": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "regex for recursive search of files",
			},
		},
	}
}

func resourceEncryptArtifactCreate(d *schema.ResourceData, m interface{}) error {
	path := d.Get("path").(string)
	output := d.Get("output_dir").(string)
	regex := d.Get("regex").(string)
	key := m.(model.Key)
	r, err := key.NewEncryptRequest(path, output, regex)
	if err != nil {
		err = stacktrace.Propagate(err, "could not create encryption request")
		return err
	}
	result, err := dare.EncryptFile(request)
	if err != nil {
		err = stacktrace.Propagate(err, "could not encrypt artifact(s).input : '%s' output_dir : '%s' regex : '%s'", input, output, regex)
		return err
	}
	return resourceEncryptArtifactRead(d, m)
}

func resourceEncryptArtifactRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceEncryptArtifactUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceEncryptArtifactDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
