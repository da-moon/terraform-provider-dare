package dare

import (
	"log"
	"os"

	model "github.com/da-moon/go-dare/model"
	logger "github.com/da-moon/go-logger"
	urandom "github.com/da-moon/go-urandom"
	schema "github.com/hashicorp/terraform/helper/schema"
	terraform "github.com/hashicorp/terraform/terraform"
	"github.com/palantir/stacktrace"
)

// Provider ...
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Sensitive:   true,
				Optional:    true,
				Description: "hex encoded 32 byte key string used for encryption/decryption",
			},
			"key_file": {
				Type:        schema.TypeString,
				Sensitive:   false,
				Optional:    true,
				Description: "file containing a hex encoded 32 byte key string used for encryption/decryption",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			// "encrypt_artifact": dataSourceEncryptedFile(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(data *schema.ResourceData) (interface{}, error) {
	key := data.Get("key").(string)
	keyFile := data.Get("key_file").(string)
	sink := logger.NewLevelFilter(
		logger.WithWriter(logger.NewGatedWriter(os.Stderr)),
	)
	l := log.New(sink, "", log.LstdFlags)
	uuid, err := urandom.UUID()
	if err != nil {
		err = stacktrace.Propagate(err, "could not return provider configuration function")
		return nil, err
	}
	k, err := model.NewKey(l, uuid,
		model.WithEncryptionKey(key),
		model.WithKeyFile(keyFile),
	)
	if err != nil {
		err = stacktrace.Propagate(err, "could not return provider configuration function")
		return nil, err
	}
	err = k.Sanitize()
	if err != nil {
		err = stacktrace.Propagate(err, "could not return provider configuration function")
		return nil, err
	}
	return k, nil
}
