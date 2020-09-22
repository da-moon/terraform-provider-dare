package provider

import (
	"os"
	"strconv"
	"time"

	dare "github.com/da-moon/go-dare"
	model "github.com/da-moon/go-dare/model"
	schema "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	stacktrace "github.com/palantir/stacktrace"
)

func resourceEncryptArtifact() *schema.Resource {
	return &schema.Resource{
		Create: resourceEncryptArtifactCreate,
		Read:   resourceEncryptArtifactRead,
		Update: resourceEncryptArtifactUpdate,
		Delete: resourceEncryptArtifactDelete,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "path of a directory or a single file for encryption",
			},
			"output_dir": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "root directory in which the encrypted files are stored",
			},
			"regex": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "regex for recursive search of files",
			},
			"encrypted_artifacts": {
				Type:        schema.TypeMap,
				Computed:    true,
				ForceNew:    true,
				Description: "encrypted artifacts path and their hash value",
			},
			"random_nonce": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "random nonce generated for encrypting artifacts",
			},
		},
	}
}

func resourceEncryptArtifactCreate(d *schema.ResourceData, m interface{}) error {
	path := d.Get("path").(string)
	output := d.Get("output_dir").(string)
	regex := d.Get("regex").(string)
	key := m.(model.Key)
	err := key.Sanitize()
	if err != nil {
		err = stacktrace.Propagate(err, "could not sanitize key")
		return err
	}
	r, err := key.NewEncryptRequest(path, output, regex)
	if err != nil {
		err = stacktrace.Propagate(err, "could not create encryption request")
		return err
	}
	err = r.Sanitize()
	if err != nil {
		err = stacktrace.Propagate(err, "could not sanitize encryption request")
		return err
	}
	result, err := dare.EncryptFile(r)
	if err != nil {
		err = stacktrace.Propagate(err, "could not encrypt artifact(s).input : '%s' output_dir : '%s' regex : '%s'", path, output, regex)
		return err
	}
	err = result.Sanitize()
	if err != nil {
		err = stacktrace.Propagate(err, "could not sanitize encryption result.input : '%s' output_dir : '%s' regex : '%s'", path, output, regex)
		return err
	}
	err = d.Set("random_nonce", result.RandomNonce)
	if err != nil {
		err = stacktrace.Propagate(err, "could not set random_nonce value .input : '%s' output_dir : '%s' regex : '%s'", path, output, regex)
		return err
	}
	err = d.Set("encrypted_artifacts", result.EncryptedArtifacts)
	if err != nil {
		err = stacktrace.Propagate(err, "could not set encrypted_artifacts value .input : '%s' output_dir : '%s' regex : '%s'", path, output, regex)
		return err
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return nil
}

// resourceEncryptArtifactUpdate deletes all encrypted files and re-encrypts (create) them
func resourceEncryptArtifactUpdate(d *schema.ResourceData, m interface{}) error {
	err := resourceEncryptArtifactDelete(d, m)
	if err != nil {
		stacktrace.Propagate(err, "could not update resource")
	}
	err = resourceEncryptArtifactCreate(d, m)
	if err != nil {
		stacktrace.Propagate(err, "could not create resource")
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return nil
}

//  deletes all encrypted files
func resourceEncryptArtifactDelete(d *schema.ResourceData, m interface{}) error {
	enc := d.Get("encrypted_artifacts").(map[string]model.Hash)
	if enc == nil {
		err := stacktrace.NewError("could not load encrypted artifacts list")
		return err
	}
	for k := range enc {
		// optimistic deletion
		os.Remove(k)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return nil
}

// resourceEncryptArtifactRead doesn't do anything
func resourceEncryptArtifactRead(d *schema.ResourceData, m interface{}) error {
	return nil
}
