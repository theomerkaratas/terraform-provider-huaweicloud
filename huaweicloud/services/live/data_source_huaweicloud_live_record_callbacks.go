// Generated by PMS #479
package live

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceLiveRecordCallbacks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLiveRecordCallbacksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ingest domain name.`,
			},
			"app_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the application name.`,
			},
			"callbacks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The callback configurations.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sign_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The encryption type.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time in the format of yyyy-mm-ddThh:mm:ssZ (UTC time).`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest modification time in the format of yyyy-mm-ddThh:mm:ssZ (UTC time).`,
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The recording callback ID.`,
						},
						"domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ingest domain name.`,
						},
						"app_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The application name.`,
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The callback URL for sending recording notifications.`,
						},
						"types": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The types of recording notifications.`,
						},
					},
				},
			},
		},
	}
}

type RecordCallbacksDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newRecordCallbacksDSWrapper(d *schema.ResourceData, meta interface{}) *RecordCallbacksDSWrapper {
	return &RecordCallbacksDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceLiveRecordCallbacksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newRecordCallbacksDSWrapper(d, meta)
	lisRecCalConRst, err := wrapper.ListRecordCallbackConfigs()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listRecordCallbackConfigsToSchema(lisRecCalConRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API LIVE GET /v1/{project_id}/record/callbacks
func (w *RecordCallbacksDSWrapper) ListRecordCallbackConfigs() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "live")
	if err != nil {
		return nil, err
	}

	uri := "/v1/{project_id}/record/callbacks"
	params := map[string]any{
		"publish_domain": w.Get("domain_name"),
		"app":            w.Get("app_name"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		OffsetPager("callback_config", "offset", "limit", 0).
		Request().
		Result()
}

func (w *RecordCallbacksDSWrapper) listRecordCallbackConfigsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("callbacks", schemas.SliceToList(body.Get("callback_config"),
			func(callbacks gjson.Result) any {
				return map[string]any{
					"sign_type":   callbacks.Get("sign_type").Value(),
					"created_at":  callbacks.Get("create_time").Value(),
					"updated_at":  callbacks.Get("update_time").Value(),
					"id":          callbacks.Get("id").Value(),
					"domain_name": callbacks.Get("publish_domain").Value(),
					"app_name":    callbacks.Get("app").Value(),
					"url":         callbacks.Get("notify_callback_url").Value(),
					"types":       schemas.SliceToStrList(callbacks.Get("notify_event_subscription")),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
