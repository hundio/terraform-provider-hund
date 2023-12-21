package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hundio/terraform-provider-hund/internal/models"
)

func WatchdogComputeHighFrequency(ctx context.Context, watchdogPath path.Path, watchdogPlan *models.WatchdogModel, watchdogConfig *models.WatchdogModel, resp *resource.ModifyPlanResponse) {
	highFrequencyPath := watchdogPath.AtName("high_frequency")

	if watchdogPlan.Service.Uptimerobot != nil {
		if watchdogConfig.HighFrequency.ValueBool() {
			resp.Diagnostics.AddAttributeError(
				highFrequencyPath,
				"Cannot set Watchdog High-frequency when using UptimeRobot Service Type",
				"UptimeRobot service does not support high-frequency polling. Either set"+
					" high_frequency to false, or remove the attribute.",
			)
		}

		resp.Plan.SetAttribute(ctx, highFrequencyPath, types.BoolValue(false))
		return
	}

	nativeServicePlan := watchdogPlan.Service.NativeService()

	if nativeServicePlan != nil {
		if !watchdogConfig.HighFrequency.IsNull() {
			resp.Diagnostics.AddAttributeError(
				highFrequencyPath,
				"Cannot set Watchdog High-frequency when using Native Service Types",
				"The value for high_frequency is derived from the frequency attribute"+
					" of your chosen Native service. Please set the frequency instead.",
			)
		}

		highFreq := nativeServicePlan.GetFrequency().ValueInt64() < 60000
		resp.Plan.SetAttribute(ctx, highFrequencyPath, types.BoolValue(highFreq))
	}
}
