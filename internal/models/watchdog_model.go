package models

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
)

type WatchdogModel struct {
	Id            types.String         `tfsdk:"id"`
	HighFrequency types.Bool           `tfsdk:"high_frequency"`
	LatestStatus  types.String         `tfsdk:"latest_status"`
	Service       WatchdogServiceModel `tfsdk:"service"`
}

func ToWatchdogModel(watchdog hundApiV1.Watchdog) (WatchdogModel, diag.Diagnostics) {
	service, diags := ToWatchdogServiceModel(watchdog.Service)

	return WatchdogModel{
		Id:            types.StringValue(watchdog.Id),
		HighFrequency: types.BoolValue(watchdog.HighFrequency),
		LatestStatus:  types.StringPointerValue(watchdog.LatestStatus),
		Service:       service,
	}, diags
}
