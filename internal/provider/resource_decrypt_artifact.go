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

func resourceDecryptArtifact() *schema.Resource {
	return &schema.Resource{
		Create: resourceDecryptArtifactCreate,
		Read:   resourceDecryptArtifactRead,
		Update: resourceDecryptArtifactUpdate,
		Delete: resourceDecryptArtifactDelete,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "path of a directory or a single file for decryption",
			},
			"nonce": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "random nonce used for encryption.",
			},
			"output_dir": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "root directory in which the encrypted files are stored",
			},
			"decrypted_artifacts": {
				Type:        schema.TypeMap,
				Computed:    true,
				ForceNew:    true,
				Description: "encrypted artifacts path and their hash value",
			},
		},
	}
}

func resourceDecryptArtifactCreate(d *schema.ResourceData, m interface{}) error {
	path := d.Get("path").(string)
	output := d.Get("output_dir").(string)
	nonce := d.Get("nonce").(string)
	key := m.(model.Key)
	key.Nonce = nonce
	err := key.Sanitize()
	if err != nil {
		err = stacktrace.Propagate(err, "could not sanitize key")
		return err
	}
	r, err := key.NewDecryptRequest(path, output)
	if err != nil {
		err = stacktrace.Propagate(err, "could not create decryption request")
		return err
	}
	err = r.Sanitize()
	if err != nil {
		err = stacktrace.Propagate(err, "could not sanitize decryption request")
		return err
	}
	result, err := dare.DecryptFile(r)
	if err != nil {
		err = stacktrace.Propagate(err, "could not encrypt artifact(s).input : '%s' output_dir : '%s' ", path, output)
		return err
	}
	err = result.Sanitize()
	if err != nil {
		err = stacktrace.Propagate(err, "could not sanitize encryption result.input : '%s' output_dir : '%s' ", path, output)
		return err
	}
	err = d.Set("decrypted_artifacts", result.DecryptedArtifacts)
	if err != nil {
		err = stacktrace.Propagate(err, "could not set decrypted_artifacts value .input : '%s' output_dir : '%s' ", path, output)
		return err
	}
	// removing encrypted artifacts
	for k := range r.Targets {
		// optimistic deletion
		os.Remove(k)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}

// resourceDecryptArtifactUpdate encrypts previously decrypted artifacts ( calling delete method) and then decrypts files with updated values them
func resourceDecryptArtifactUpdate(d *schema.ResourceData, m interface{}) error {
	var err error

	err = resourceDecryptArtifactDelete(d, m)
	if err != nil {
		stacktrace.Propagate(err, "could not update resource")
	}
	err = resourceDecryptArtifactCreate(d, m)
	if err != nil {
		stacktrace.Propagate(err, "could not create resource")
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}

// resourceDecryptArtifactDelete encrypts data
// encrypted payload is stored in same dir as the plaintext form
func resourceDecryptArtifactDelete(d *schema.ResourceData, m interface{}) error {
	var err error
	key := m.(model.Key)
	err = key.Sanitize()
	if err != nil {
		err = stacktrace.Propagate(err, "could not sanitize key")
		return err
	}
	enc := d.Get("decrypted_artifacts").(map[string]model.Hash)
	if enc == nil {
		err := stacktrace.NewError("could not load encrypted artifacts list")
		return err
	}
	for k := range enc {
		request, err := key.NewEncryptRequest(k, "", "")
		if err != nil {
			err = stacktrace.Propagate(err, "could not create encryption request")
			return err
		}
		err = request.Sanitize()
		if err != nil {
			err = stacktrace.Propagate(err, "could not sanitize encryption request")
			return err
		}
		result, err := dare.EncryptFile(request)
		if err != nil {
			err = stacktrace.Propagate(err, "could not encrypt artifact(s).")
		}
		err = result.Sanitize()
		if err != nil {
			err = stacktrace.Propagate(err, "could not sanitize encryption result.")
			return err
		}
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return nil
}

// resourceDecryptArtifactRead doesn't do anything
func resourceDecryptArtifactRead(d *schema.ResourceData, m interface{}) error {
	return nil
}
