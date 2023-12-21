package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func nativeIcmpServiceDataSourceSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed: true,
		Attributes: nativeServiceDataSourceSchema(map[string]schema.Attribute{
			"ip_version": schema.StringAttribute{
				Computed: true,
			},
			"percentage_failed_threshold": schema.Float64Attribute{
				Computed: true,
			},
		}),
	}
}

func nativeHttpServiceDataSourceSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed: true,
		Attributes: nativeServiceDataSourceSchema(map[string]schema.Attribute{
			"headers": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"response_body_must_contain": schema.StringAttribute{
				Computed: true,
			},
			"response_body_must_contain_mode": schema.StringAttribute{
				Computed: true,
			},
			"response_code_must_be": schema.Int64Attribute{
				Computed: true,
			},
			"ssl_verify_peer": schema.BoolAttribute{
				Computed: true,
			},
			"follow_redirects": schema.BoolAttribute{
				Computed: true,
			},
			"username": schema.StringAttribute{
				Computed: true,
			},
			"password": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
		}),
	}
}

func nativeDnsServiceDataSourceSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed: true,
		Attributes: nativeServiceDataSourceSchema(map[string]schema.Attribute{
			"record_type": schema.StringAttribute{
				Computed: true,
			},
			"nameservers": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"response_containment": schema.StringAttribute{
				Computed: true,
			},
			"responses_must_contain": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
		}),
	}
}

func nativeTcpServiceDataSourceSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed: true,
		Attributes: nativeServiceDataSourceSchema(map[string]schema.Attribute{
			"ip_version": schema.StringAttribute{
				Computed: true,
			},
			"port": schema.Int64Attribute{
				Computed: true,
			},
			"response_must_contain": schema.StringAttribute{
				Computed: true,
			},
			"response_must_contain_mode": schema.StringAttribute{
				Computed: true,
			},
			"send_data": schema.StringAttribute{
				Computed: true,
			},
			"wait_for_initial_response": schema.BoolAttribute{
				Computed: true,
			},
		}),
	}
}

func nativeUdpServiceDataSourceSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed: true,
		Attributes: nativeServiceDataSourceSchema(map[string]schema.Attribute{
			"ip_version": schema.StringAttribute{
				Computed: true,
			},
			"port": schema.Int64Attribute{
				Computed: true,
			},
			"response_must_contain": schema.StringAttribute{
				Computed: true,
			},
			"response_must_contain_mode": schema.StringAttribute{
				Computed: true,
			},
			"send_data": schema.StringAttribute{
				Computed: true,
			},
		}),
	}
}

func nativeServiceDataSourceSchema(extension map[string]schema.Attribute) map[string]schema.Attribute {
	schema := map[string]schema.Attribute{
		"target": schema.StringAttribute{
			Computed: true,
		},
		"consecutive_check_degraded_threshold": schema.Int64Attribute{
			Computed: true,
		},
		"consecutive_check_outage_threshold": schema.Int64Attribute{
			Computed: true,
		},
		"frequency": schema.Int64Attribute{
			Computed: true,
		},
		"percentage_regions_failed_threshold": schema.Float64Attribute{
			Computed: true,
		},
		"regions": schema.SetAttribute{
			Computed:    true,
			ElementType: types.StringType,
		},
		"timeout": schema.Int64Attribute{
			Computed: true,
		},
	}

	for k, a := range extension {
		schema[k] = a
	}

	return schema
}
