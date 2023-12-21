package models

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
)

type WatchdogServiceModel struct {
	Manual      *ManualServiceModel              `tfsdk:"manual"`
	Updown      *UpdownServiceModel              `tfsdk:"updown"`
	Pingdom     *PingdomWatchdogServiceModel     `tfsdk:"pingdom"`
	Uptimerobot *UptimerobotWatchdogServiceModel `tfsdk:"uptimerobot"`
	Webhook     *WebhookWatchdogServiceModel     `tfsdk:"webhook"`

	NativeIcmp *NativeIcmpServiceModel `tfsdk:"icmp"`
	NativeHttp *NativeHttpServiceModel `tfsdk:"http"`
	NativeDns  *NativeDnsServiceModel  `tfsdk:"dns"`
	NativeTcp  *NativeTcpServiceModel  `tfsdk:"tcp"`
	NativeUdp  *NativeUdpServiceModel  `tfsdk:"udp"`
}

func ToWatchdogServiceModel(service hundApiV1.ServicesWatchdog) (WatchdogServiceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	model := WatchdogServiceModel{}

	switch service.Discriminator() {
	case "manual":
		manual, err := service.AsServicesWatchdog0()

		if err != nil {
			diags.Append(ServiceDecodeError(service.Discriminator(), err))

			return model, diags
		}

		model.Manual = &ManualServiceModel{
			State: types.Int64Value(int64(manual.State)),
		}
	case "updown":
		updown, err := service.AsServicesWatchdog1()

		if err != nil {
			diags.Append(ServiceDecodeError(service.Discriminator(), err))

			return model, diags
		}

		model.Updown = &UpdownServiceModel{
			MonitorToken: types.StringValue(updown.MonitorToken),
		}
	case "uptimerobot":
		uptimerobot, err := service.AsServicesWatchdog2()

		if err != nil {
			diags.Append(ServiceDecodeError(service.Discriminator(), err))

			return model, diags
		}

		model.Uptimerobot = &UptimerobotWatchdogServiceModel{
			UnconfirmedIsDown: types.BoolValue(uptimerobot.UnconfirmedIsDown),
		}
	case "pingdom":
		pingdom, err := service.AsServicesWatchdog4()

		if err != nil {
			diags.Append(ServiceDecodeError(service.Discriminator(), err))

			return model, diags
		}

		model.Pingdom = &PingdomWatchdogServiceModel{
			CheckId:           types.StringValue(pingdom.CheckId),
			CheckType:         types.StringValue(string(pingdom.CheckType)),
			UnconfirmedIsDown: types.BoolValue(pingdom.UnconfirmedIsDown),
		}
	case "webhook":
		webhook, err := service.AsServicesWatchdog3()

		if err != nil {
			diags.Append(ServiceDecodeError(service.Discriminator(), err))

			return model, diags
		}

		model.Webhook = &WebhookWatchdogServiceModel{
			WebhookKey:        types.StringPointerValue(webhook.WebhookKey),
			Deadman:           types.BoolValue(webhook.Deadman),
			ConsecutiveChecks: types.Int64PointerValue(hundApiV1.ToInt64Ptr(webhook.ConsecutiveChecks)),
			ReportingInterval: types.Int64PointerValue(hundApiV1.ToInt64Ptr(webhook.ReportingInterval)),
		}
	case "native":
		native, err := service.AsServicesWatchdog9()

		if err != nil {
			diags.Append(ServiceDecodeError(service.Discriminator(), err))

			return model, diags
		}

		nativeModel, diag0 := ToNativeServiceModel(native)
		diags.Append(diag0...)

		switch m := nativeModel.(type) {
		case *NativeIcmpServiceModel:
			model.NativeIcmp = m
		case *NativeHttpServiceModel:
			model.NativeHttp = m
		case *NativeDnsServiceModel:
			model.NativeDns = m
		case *NativeTcpServiceModel:
			model.NativeTcp = m
		case *NativeUdpServiceModel:
			model.NativeUdp = m
		default:
			return model, diags
		}
	default:
		diags.Append(UnknownServiceError(service.Discriminator()))
	}

	return model, diags
}

func (s *WatchdogServiceModel) ReplaceSensitiveAttributes(orig WatchdogServiceModel) {
	if s.Pingdom != nil {
		s.Pingdom.ApiToken = orig.Pingdom.ApiToken
	} else if s.Updown != nil {
		s.Updown.MonitorApiKey = orig.Updown.MonitorApiKey
	} else if s.Uptimerobot != nil {
		s.Uptimerobot.MonitorApiKey = orig.Uptimerobot.MonitorApiKey
	}
}

func (s WatchdogServiceModel) NativeService() NativeServiceModel {
	if s.NativeIcmp != nil {
		return s.NativeIcmp
	} else if s.NativeHttp != nil {
		return s.NativeHttp
	} else if s.NativeDns != nil {
		return s.NativeDns
	} else if s.NativeTcp != nil {
		return s.NativeTcp
	} else if s.NativeUdp != nil {
		return s.NativeUdp
	}

	return nil
}

func (s WatchdogServiceModel) ApiCreateForm() (hundApiV1.FormWatchdogCreate, error) {
	form := hundApiV1.FormWatchdogCreate{}

	if s.Manual != nil {
		err := form.FromFormWatchdogCreate0(s.Manual.ApiCreateForm())

		return form, err
	} else if s.Updown != nil {
		err := form.FromFormWatchdogCreate1(s.Updown.ApiCreateForm())

		return form, err
	} else if s.Uptimerobot != nil {
		err := form.FromFormWatchdogCreate2(s.Uptimerobot.ApiCreateForm())

		return form, err
	} else if s.Webhook != nil {
		err := form.FromFormWatchdogCreate3(s.Webhook.ApiCreateForm())

		return form, err
	} else if s.Pingdom != nil {
		err := form.FromFormWatchdogCreate4(s.Pingdom.ApiCreateForm())

		return form, err
	}

	nativeForm, err := s.NativeApiCreateForm()

	if err != nil {
		return form, err
	}

	err = form.FromFormWatchdogCreate8(nativeForm)

	return form, err
}

func (s WatchdogServiceModel) NativeApiCreateForm() (hundApiV1.NativeFormCreate, error) {
	form := hundApiV1.NativeFormCreate{}

	if s.NativeIcmp != nil {
		err := form.FromNativeFormCreate0(s.NativeIcmp.ApiCreateForm())

		return form, err
	} else if s.NativeHttp != nil {
		err := form.FromNativeFormCreate1(s.NativeHttp.ApiCreateForm())

		return form, err
	} else if s.NativeDns != nil {
		err := form.FromNativeFormCreate2(s.NativeDns.ApiCreateForm())

		return form, err
	} else if s.NativeTcp != nil {
		err := form.FromNativeFormCreate3(s.NativeTcp.ApiCreateForm())

		return form, err
	} else if s.NativeUdp != nil {
		err := form.FromNativeFormCreate4(s.NativeUdp.ApiCreateForm())

		return form, err
	}

	return form, errors.New("could not find non-nil native model")
}

func (s WatchdogServiceModel) ApiUpdateForm() (hundApiV1.FormWatchdogUpdate, error) {
	form := hundApiV1.FormWatchdogUpdate{}

	if s.Manual != nil {
		err := form.FromFormWatchdogUpdate0(s.Manual.ApiUpdateForm())

		return form, err
	} else if s.Updown != nil {
		err := form.FromFormWatchdogUpdate1(s.Updown.ApiUpdateForm())

		return form, err
	} else if s.Uptimerobot != nil {
		err := form.FromFormWatchdogUpdate2(s.Uptimerobot.ApiUpdateForm())

		return form, err
	} else if s.Webhook != nil {
		err := form.FromFormWatchdogUpdate3(s.Webhook.ApiUpdateForm())

		return form, err
	} else if s.Pingdom != nil {
		err := form.FromFormWatchdogUpdate4(s.Pingdom.ApiUpdateForm())

		return form, err
	}

	nativeForm, err := s.NativeApiUpdateForm()

	if err != nil {
		return form, err
	}

	err = form.FromFormWatchdogUpdate8(nativeForm)

	return form, err
}

func (s WatchdogServiceModel) NativeApiUpdateForm() (hundApiV1.NativeFormUpdate, error) {
	form := hundApiV1.NativeFormUpdate{}

	if s.NativeIcmp != nil {
		err := form.FromNativeFormUpdate0(s.NativeIcmp.ApiUpdateForm())

		return form, err
	} else if s.NativeHttp != nil {
		err := form.FromNativeFormUpdate1(s.NativeHttp.ApiUpdateForm())

		return form, err
	} else if s.NativeDns != nil {
		err := form.FromNativeFormUpdate2(s.NativeDns.ApiUpdateForm())

		return form, err
	} else if s.NativeTcp != nil {
		err := form.FromNativeFormUpdate3(s.NativeTcp.ApiUpdateForm())

		return form, err
	} else if s.NativeUdp != nil {
		err := form.FromNativeFormUpdate4(s.NativeUdp.ApiUpdateForm())

		return form, err
	}

	return form, errors.New("could not find non-nil native model")
}

type ManualServiceModel struct {
	State types.Int64 `tfsdk:"state"`
}

func (s ManualServiceModel) ApiCreateForm() hundApiV1.ManualFormCreate {
	return hundApiV1.ManualFormCreate{
		Type:  hundApiV1.ManualFormCreateTypeManual,
		State: hundApiV1.IntegerState(s.State.ValueInt64()),
	}
}

func (s ManualServiceModel) ApiUpdateForm() hundApiV1.ManualFormUpdate {
	return hundApiV1.ManualFormUpdate{
		State: hundApiV1.MaybeIntegerState(hundApiV1.ToIntPtr(s.State.ValueInt64Pointer())),
	}
}

type UpdownServiceModel struct {
	MonitorToken  types.String `tfsdk:"monitor_token"`
	MonitorApiKey types.String `tfsdk:"monitor_api_key"`
}

func (s UpdownServiceModel) ApiCreateForm() hundApiV1.UpdownFormCreate {
	return hundApiV1.UpdownFormCreate{
		Type:          hundApiV1.UpdownFormCreateTypeUpdown,
		MonitorToken:  s.MonitorToken.ValueString(),
		MonitorApiKey: s.MonitorApiKey.ValueString(),
	}
}

func (s UpdownServiceModel) ApiUpdateForm() hundApiV1.UpdownFormUpdate {
	return hundApiV1.UpdownFormUpdate{
		MonitorToken:  s.MonitorToken.ValueStringPointer(),
		MonitorApiKey: s.MonitorApiKey.ValueStringPointer(),
	}
}

type PingdomWatchdogServiceModel struct {
	ApiToken          types.String `tfsdk:"api_token"`
	CheckId           types.String `tfsdk:"check_id"`
	CheckType         types.String `tfsdk:"check_type"`
	UnconfirmedIsDown types.Bool   `tfsdk:"unconfirmed_is_down"`
}

func (s PingdomWatchdogServiceModel) ApiCreateForm() hundApiV1.PingdomWatchdogFormCreate {
	return hundApiV1.PingdomWatchdogFormCreate{
		Type:              hundApiV1.PingdomWatchdogFormCreateTypePingdom,
		ApiToken:          s.ApiToken.ValueString(),
		CheckId:           s.CheckId.ValueString(),
		CheckType:         (*hundApiV1.PingdomWatchdogFormCreateCheckType)(s.CheckType.ValueStringPointer()),
		UnconfirmedIsDown: s.UnconfirmedIsDown.ValueBoolPointer(),
	}
}

func (s PingdomWatchdogServiceModel) ApiUpdateForm() hundApiV1.PingdomWatchdogFormUpdate {
	return hundApiV1.PingdomWatchdogFormUpdate{
		ApiToken:          s.ApiToken.ValueStringPointer(),
		CheckId:           s.CheckId.ValueStringPointer(),
		CheckType:         (*hundApiV1.PingdomWatchdogFormUpdateCheckType)(s.CheckType.ValueStringPointer()),
		UnconfirmedIsDown: s.UnconfirmedIsDown.ValueBoolPointer(),
	}
}

type UptimerobotWatchdogServiceModel struct {
	MonitorApiKey     types.String `tfsdk:"monitor_api_key"`
	UnconfirmedIsDown types.Bool   `tfsdk:"unconfirmed_is_down"`
}

func (s UptimerobotWatchdogServiceModel) ApiCreateForm() hundApiV1.UptimerobotWatchdogFormCreate {
	return hundApiV1.UptimerobotWatchdogFormCreate{
		Type:              hundApiV1.UptimerobotWatchdogFormCreateTypeUptimerobot,
		MonitorApiKey:     s.MonitorApiKey.ValueString(),
		UnconfirmedIsDown: s.UnconfirmedIsDown.ValueBoolPointer(),
	}
}

func (s UptimerobotWatchdogServiceModel) ApiUpdateForm() hundApiV1.UptimerobotWatchdogFormUpdate {
	return hundApiV1.UptimerobotWatchdogFormUpdate{
		MonitorApiKey:     s.MonitorApiKey.ValueStringPointer(),
		UnconfirmedIsDown: s.UnconfirmedIsDown.ValueBoolPointer(),
	}
}

type WebhookWatchdogServiceModel struct {
	WebhookKey        types.String `tfsdk:"webhook_key"`
	Deadman           types.Bool   `tfsdk:"deadman"`
	ConsecutiveChecks types.Int64  `tfsdk:"consecutive_checks"`
	ReportingInterval types.Int64  `tfsdk:"reporting_interval"`
}

func (s WebhookWatchdogServiceModel) ApiCreateForm() hundApiV1.WebhookWatchdogFormCreate {
	return hundApiV1.WebhookWatchdogFormCreate{
		Type:              hundApiV1.WebhookWatchdogFormCreateTypeWebhook,
		WebhookKey:        s.WebhookKey.ValueStringPointer(),
		Deadman:           s.Deadman.ValueBoolPointer(),
		ConsecutiveChecks: hundApiV1.DblPtr(hundApiV1.ToIntPtr(s.ConsecutiveChecks.ValueInt64Pointer())),
		ReportingInterval: hundApiV1.DblPtr(hundApiV1.ToIntPtr(s.ReportingInterval.ValueInt64Pointer())),
	}
}

func (s WebhookWatchdogServiceModel) ApiUpdateForm() hundApiV1.WebhookWatchdogFormUpdate {
	return hundApiV1.WebhookWatchdogFormUpdate{
		WebhookKey:        s.WebhookKey.ValueStringPointer(),
		Deadman:           s.Deadman.ValueBoolPointer(),
		ConsecutiveChecks: hundApiV1.DblPtr(hundApiV1.ToIntPtr(s.ConsecutiveChecks.ValueInt64Pointer())),
		ReportingInterval: hundApiV1.DblPtr(hundApiV1.ToIntPtr(s.ReportingInterval.ValueInt64Pointer())),
	}
}
