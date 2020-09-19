package provider

import (
	// "io/ioutil"
	// "log"
	// "os"
	// "strings"

	// "github.com/da-moon/go-logger"
	// "github.com/da-moon/terraform-provider-dare/pkg/model"
	schema "github.com/hashicorp/terraform/helper/schema"
	terraform "github.com/hashicorp/terraform/terraform"
)

// Provider ...
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			// "key": {
			// 	Type:      schema.TypeString,
			// 	Sensitive: true,
			// 	Required:  true,
			// },
		},
		DataSourcesMap: map[string]*schema.Resource{
			// "encrypted_file": dataSourceEncryptedFile(),
		},
		// ConfigureFunc: providerConfigure,
	}
}

// func providerConfigure(data *schema.ResourceData) (interface{}, error) {
// 	key := data.Get("key").(string)
// 	keyFile := data.Get("key_file").(string)
// 	sink := logger.NewLevelFilter(
// 		logger.WithMinLevel(strings.ToUpper(logLevel)),
// 		logger.WithWriter(logger.NewGatedWriter(os.Stderr)),
// 	)

// 	l := log.New(sink, "", log.LstdFlags)
// 	log.New(ioutil.Discard)
// 	model.NewKey(l, "",
// 		model.WithEncryptionKey(key),
// 	)
// 	return key, nil
// }
