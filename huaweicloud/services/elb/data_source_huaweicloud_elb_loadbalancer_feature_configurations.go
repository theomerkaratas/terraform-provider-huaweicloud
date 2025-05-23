// Generated by PMS #566
package elb

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

func DataSourceElbLoadbalancerFeatureConfigurations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceElbLoadbalancerFeatureConfigurationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"loadbalancer_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the load balancer ID.`,
			},
			"features": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Specifies the load balancer feature information list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"feature": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Specifies the feature name.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Specifies the type of the feature configuration value.`,
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Specifies the feature value.`,
						},
					},
				},
			},
		},
	}
}

type LoadbalancerFeatureConfigurationsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newLoadbalancerFeatureConfigurationsDSWrapper(d *schema.ResourceData, meta interface{}) *LoadbalancerFeatureConfigurationsDSWrapper {
	return &LoadbalancerFeatureConfigurationsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceElbLoadbalancerFeatureConfigurationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newLoadbalancerFeatureConfigurationsDSWrapper(d, meta)
	lisLoaFeaRst, err := wrapper.ListLoadbalancerFeature()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listLoadbalancerFeatureToSchema(lisLoaFeaRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API ELB GET /v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/features
func (w *LoadbalancerFeatureConfigurationsDSWrapper) ListLoadbalancerFeature() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "elb")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/features"
	uri = strings.ReplaceAll(uri, "{loadbalancer_id}", w.Get("loadbalancer_id").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

func (w *LoadbalancerFeatureConfigurationsDSWrapper) listLoadbalancerFeatureToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("features", schemas.SliceToList(body.Get("features"),
			func(features gjson.Result) any {
				return map[string]any{
					"feature": features.Get("feature").Value(),
					"type":    features.Get("type").Value(),
					"value":   features.Get("value").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
