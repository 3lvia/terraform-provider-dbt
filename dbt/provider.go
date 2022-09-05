package dbt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"dbt_group": resourceGroup(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
		Schema: map[string]*schema.Schema{
			"service_token": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The service token for api-requests to DBT. See https://docs.getdbt.com/docs/dbt-cloud/access-control/enterprise-permissions for required permission sets",
			},
			"account_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The account id for DBT cloud",
			},
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	serviceToken := d.Get("service_token").(string)
	accountId := d.Get("account_id").(int)

	return &DbtProviderInput{serviceToken, accountId}, nil
}

type DbtProviderInput struct {
	ServiceToken string
	AccountId    int
}
