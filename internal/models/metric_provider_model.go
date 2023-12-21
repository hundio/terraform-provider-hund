package models

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
)

type MetricProviderModel struct {
	Id        types.String                   `tfsdk:"id"`
	Watchdog  types.String                   `tfsdk:"watchdog"`
	Default   types.Bool                     `tfsdk:"default"`
	Instances map[string]MetricInstanceModel `tfsdk:"instances"`
	Service   *MetricProviderServiceModel    `tfsdk:"service"`
}

func ToMetricProviderModel(mp hundApiV1.MetricProvider) (MetricProviderModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	service, diag0 := ToMetricProviderServiceModel(mp.Service)
	diags.Append(diag0...)

	model := MetricProviderModel{
		Id:        types.StringValue(mp.Id),
		Watchdog:  types.StringValue(mp.Watchdog),
		Default:   types.BoolValue(mp.Default),
		Instances: map[string]MetricInstanceModel{},
		Service:   &service,
	}

	for _, mi := range mp.Instances {
		instance, diag0 := ToMetricInstanceModel(mi)
		diags.Append(diag0...)

		model.Instances[instance.DefinitionSlug.ValueString()] = instance
	}

	return model, diags
}

type MetricInstanceModel struct {
	Id                 types.String  `tfsdk:"id"`
	Slug               types.String  `tfsdk:"slug"`
	DefinitionSlug     types.String  `tfsdk:"definition_slug"`
	Enabled            types.Bool    `tfsdk:"enabled"`
	TopLevelEnabled    types.Bool    `tfsdk:"top_level_enabled"`
	Title              types.String  `tfsdk:"title"`
	TitleTranslations  types.Map     `tfsdk:"title_translations"`
	XTitle             types.String  `tfsdk:"x_title"`
	XTitleTranslations types.Map     `tfsdk:"x_title_translations"`
	YTitle             types.String  `tfsdk:"y_title"`
	YTitleTranslations types.Map     `tfsdk:"y_title_translations"`
	XType              types.String  `tfsdk:"x_type"`
	YType              types.String  `tfsdk:"y_type"`
	YSupremum          types.Float64 `tfsdk:"y_supremum"`
	PlotType           types.String  `tfsdk:"plot_type"`
	Aggregation        types.String  `tfsdk:"aggregation"`
	Interpolation      types.String  `tfsdk:"interpolation"`
}

func MetricInstanceModelUnknown() MetricInstanceModel {
	return MetricInstanceModel{
		Id:                 types.StringUnknown(),
		Slug:               types.StringUnknown(),
		DefinitionSlug:     types.StringUnknown(),
		Enabled:            types.BoolUnknown(),
		TopLevelEnabled:    types.BoolUnknown(),
		Title:              types.StringUnknown(),
		TitleTranslations:  types.MapUnknown(types.StringType),
		XTitle:             types.StringUnknown(),
		XTitleTranslations: types.MapUnknown(types.StringType),
		YTitle:             types.StringUnknown(),
		YTitleTranslations: types.MapUnknown(types.StringType),
		XType:              types.StringUnknown(),
		YType:              types.StringUnknown(),
		YSupremum:          types.Float64Unknown(),
		PlotType:           types.StringUnknown(),
		Aggregation:        types.StringUnknown(),
		Interpolation:      types.StringUnknown(),
	}
}

func ToMetricInstanceModel(mi hundApiV1.MetricInstance) (MetricInstanceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	title, titleTranslations, diag0 := hundApiV1.FromI18nString(mi.Title)
	diags.Append(diag0...)

	xTitle, xTitleTranslations, diag0 := hundApiV1.FromI18nString(mi.Title)
	diags.Append(diag0...)

	yTitle, yTitleTranslations, diag0 := hundApiV1.FromI18nString(mi.Title)
	diags.Append(diag0...)

	if diags.HasError() {
		return MetricInstanceModel{}, diags
	}

	model := MetricInstanceModel{
		Id:                 types.StringValue(mi.Id),
		Slug:               types.StringValue(mi.Slug),
		DefinitionSlug:     types.StringValue(mi.DefinitionSlug),
		Enabled:            types.BoolValue(mi.Enabled),
		TopLevelEnabled:    types.BoolValue(mi.TopLevelEnabled),
		Title:              *title,
		TitleTranslations:  *titleTranslations,
		XTitle:             *xTitle,
		XTitleTranslations: *xTitleTranslations,
		YTitle:             *yTitle,
		YTitleTranslations: *yTitleTranslations,
		XType:              types.StringValue(string(mi.XType)),
		YType:              types.StringValue(string(mi.YType)),
		YSupremum:          types.Float64Value(float64(mi.YSupremum)),
		PlotType:           types.StringValue(string(mi.PlotType)),
		Aggregation:        types.StringValue(string(mi.Aggregation)),
		Interpolation:      types.StringValue(string(mi.Interpolation)),
	}

	return model, diags
}

func (mi MetricInstanceModel) ApiCreateForm(definitionSlug string) (hundApiV1.MetricInstanceFormCreate, diag.Diagnostics) {
	var diags diag.Diagnostics

	title, err := hundApiV1.ToI18nStringPtr(mi.Title, mi.TitleTranslations)
	if err != nil {
		diags.Append(I18nStringError(err))
	}

	x_title, err := hundApiV1.ToI18nStringPtr(mi.XTitle, mi.XTitleTranslations)
	if err != nil {
		diags.Append(I18nStringError(err))
	}

	y_title, err := hundApiV1.ToI18nStringPtr(mi.YTitle, mi.YTitleTranslations)
	if err != nil {
		diags.Append(I18nStringError(err))
	}

	return hundApiV1.MetricInstanceFormCreate{
		Aggregation:     (*hundApiV1.METRICAGGREGATION)(mi.Aggregation.ValueStringPointer()),
		DefinitionSlug:  definitionSlug,
		Enabled:         mi.Enabled.ValueBoolPointer(),
		PlotType:        (*hundApiV1.METRICPLOTTYPE)(mi.PlotType.ValueStringPointer()),
		Title:           (*hundApiV1.MetricInstanceFormCreate_Title)(title),
		TopLevelEnabled: mi.TopLevelEnabled.ValueBoolPointer(),
		XTitle:          x_title,
		XType:           (*hundApiV1.METRICAXISTYPE)(mi.XType.ValueStringPointer()),
		YSupremum:       hundApiV1.ToFloat32Ptr(mi.YSupremum.ValueFloat64Pointer()),
		YTitle:          y_title,
		YType:           (*hundApiV1.METRICAXISTYPE)(mi.YType.ValueStringPointer()),
	}, diags
}

func (mi MetricInstanceModel) ApiUpdateForm(definitionSlug string) (hundApiV1.MetricInstanceFormEmbeddedUpdate, diag.Diagnostics) {
	var diags diag.Diagnostics

	title, err := hundApiV1.ToI18nStringPtr(mi.Title, mi.TitleTranslations)
	if err != nil {
		diags.Append(I18nStringError(err))
	}

	x_title, err := hundApiV1.ToI18nStringPtr(mi.XTitle, mi.XTitleTranslations)
	if err != nil {
		diags.Append(I18nStringError(err))
	}

	y_title, err := hundApiV1.ToI18nStringPtr(mi.YTitle, mi.YTitleTranslations)
	if err != nil {
		diags.Append(I18nStringError(err))
	}

	return hundApiV1.MetricInstanceFormEmbeddedUpdate{
		DefinitionSlug:  definitionSlug,
		Aggregation:     (*hundApiV1.METRICAGGREGATION)(mi.Aggregation.ValueStringPointer()),
		Enabled:         mi.Enabled.ValueBoolPointer(),
		PlotType:        (*hundApiV1.METRICPLOTTYPE)(mi.PlotType.ValueStringPointer()),
		Title:           title,
		TopLevelEnabled: mi.TopLevelEnabled.ValueBoolPointer(),
		XTitle:          x_title,
		XType:           (*hundApiV1.METRICAXISTYPE)(mi.XType.ValueStringPointer()),
		YSupremum:       hundApiV1.ToFloat32Ptr(mi.YSupremum.ValueFloat64Pointer()),
		YTitle:          y_title,
		YType:           (*hundApiV1.METRICAXISTYPE)(mi.YType.ValueStringPointer()),
	}, diags
}

func MetricProviderServiceType(service types.Object) string {
	var serviceType string
	for k, attr := range service.Attributes() {
		if !attr.IsNull() {
			serviceType = k
			break
		}
	}

	return serviceType
}

func MetricProviderServiceInstances() map[string][]string {
	return map[string][]string{
		"builtin":     {"percent_uptime", "incidents_reported"},
		"updown":      {"apdex"},
		"pingdom":     {"res"},
		"uptimerobot": {"res"},
		"webhook":     nil,
		"icmp":        {"res", "icmp.total_addresses", "icmp.passed_addresses"},
		"http": {
			"http.redirect_time",
			"http.name_lookup_time",
			"http.tcp_connection_time",
			"http.tls_handshake_time",
			"http.content_generation_time",
			"http.content_transfer_time",
			"http.total_time",
			"http.time_to_first_byte",
		},
		"dns": {},
		"tcp": {
			"res",
			"tcp.connection_time",
			"tcp.initial_response_time",
			"tcp.initial_response_transfer_time",
			"tcp.data_send_transfer_time",
			"tcp.data_send_response_time",
			"tcp.data_send_response_transfer_time",
			"tcp.disconnection_time",
			"tcp.total_time",
		},
		"udp": {
			"udp.response_time",
			"udp.response_transfer_time",
			"udp.total_time",
		},
	}

}
