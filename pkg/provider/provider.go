package provider

import (
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
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(data *schema.ResourceData) (interface{}, error) {
	key := data.Get("key").(string)
	return key, nil
}
