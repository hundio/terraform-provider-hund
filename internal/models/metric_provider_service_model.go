package models

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
)

type MetricProviderServiceModel struct {
	Builtin     *BuiltinServiceModel                   `tfsdk:"builtin"`
	Updown      *UpdownServiceModel                    `tfsdk:"updown"`
	Pingdom     *PingdomMetricProviderServiceModel     `tfsdk:"pingdom"`
	Uptimerobot *UptimerobotMetricProviderServiceModel `tfsdk:"uptimerobot"`
	Webhook     *WebhookMetricProviderServiceModel     `tfsdk:"webhook"`

	NativeIcmp *NativeIcmpServiceModel `tfsdk:"icmp"`
	NativeHttp *NativeHttpServiceModel `tfsdk:"http"`
	NativeDns  *NativeDnsServiceModel  `tfsdk:"dns"`
	NativeTcp  *NativeTcpServiceModel  `tfsdk:"tcp"`
	NativeUdp  *NativeUdpServiceModel  `tfsdk:"udp"`
}

func ToMetricProviderServiceModel(service hundApiV1.ServicesMetricProvider) (MetricProviderServiceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	model := MetricProviderServiceModel{}

	switch service.Discriminator() {
	case "builtin":
		_, err := service.AsServicesMetricProvider0()

		if err != nil {
			diags.Append(ServiceDecodeError(service.Discriminator(), err))

			return model, diags
		}

		model.Builtin = &BuiltinServiceModel{}
	case "updown":
		updown, err := service.AsServicesMetricProvider1()

		if err != nil {
			diags.Append(ServiceDecodeError(service.Discriminator(), err))

			return model, diags
		}

		model.Updown = &UpdownServiceModel{
			MonitorToken: types.StringValue(updown.MonitorToken),
		}
	case "uptimerobot":
		_, err := service.AsServicesMetricProvider2()

		if err != nil {
			diags.Append(ServiceDecodeError(service.Discriminator(), err))

			return model, diags
		}

		model.Uptimerobot = &UptimerobotMetricProviderServiceModel{}
	case "pingdom":
		pingdom, err := service.AsServicesMetricProvider4()

		if err != nil {
			diags.Append(ServiceDecodeError(service.Discriminator(), err))

			return model, diags
		}

		model.Pingdom = &PingdomMetricProviderServiceModel{
			CheckId:   types.StringValue(pingdom.CheckId),
			CheckType: types.StringValue(string(pingdom.CheckType)),
		}
	case "webhook":
		webhook, err := service.AsServicesMetricProvider3()

		if err != nil {
			diags.Append(ServiceDecodeError(service.Discriminator(), err))

			return model, diags
		}

		model.Webhook = &WebhookMetricProviderServiceModel{
			WebhookKey: types.StringPointerValue(webhook.WebhookKey),
		}
	case "native":
		native, err := service.AsServicesMetricProvider6()

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

func (s *MetricProviderServiceModel) ReplaceSensitiveAttributes(orig MetricProviderServiceModel) {
	if s.Pingdom != nil {
		s.Pingdom.ApiToken = orig.Pingdom.ApiToken
	} else if s.Updown != nil {
		s.Updown.MonitorApiKey = orig.Updown.MonitorApiKey
	} else if s.Uptimerobot != nil {
		s.Uptimerobot.MonitorApiKey = orig.Uptimerobot.MonitorApiKey
	}
}

func (s MetricProviderServiceModel) ServiceType() string {
	if s.Builtin != nil {
		return "builtin"
	} else if s.Updown != nil {
		return "updown"
	} else if s.Pingdom != nil {
		return "pingdom"
	} else if s.Uptimerobot != nil {
		return "uptimerobot"
	} else if s.Webhook != nil {
		return "webhook"
	} else if s.NativeIcmp != nil {
		return "icmp"
	} else if s.NativeHttp != nil {
		return "http"
	} else if s.NativeDns != nil {
		return "dns"
	} else if s.NativeTcp != nil {
		return "tcp"
	} else if s.NativeUdp != nil {
		return "udp"
	}

	return "unknown"
}

func (s MetricProviderServiceModel) NativeService() NativeServiceModel {
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

func (s MetricProviderServiceModel) ApiCreateForm() (hundApiV1.FormMetricProviderCreate, error) {
	form := hundApiV1.FormMetricProviderCreate{}

	if s.Builtin != nil {
		err := form.FromFormMetricProviderCreate0(s.Builtin.ApiCreateForm())

		return form, err
	} else if s.Updown != nil {
		err := form.FromFormMetricProviderCreate1(s.Updown.ApiCreateForm())

		return form, err
	} else if s.Uptimerobot != nil {
		err := form.FromFormMetricProviderCreate2(s.Uptimerobot.ApiCreateForm())

		return form, err
	} else if s.Webhook != nil {
		err := form.FromFormMetricProviderCreate3(s.Webhook.ApiCreateForm())

		return form, err
	} else if s.Pingdom != nil {
		err := form.FromFormMetricProviderCreate4(s.Pingdom.ApiCreateForm())

		return form, err
	}

	nativeForm, err := s.NativeApiCreateForm()

	if err != nil {
		return form, err
	}

	err = form.FromFormMetricProviderCreate5(nativeForm)

	return form, err
}

func (s MetricProviderServiceModel) NativeApiCreateForm() (hundApiV1.NativeFormCreate, error) {
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

func (s MetricProviderServiceModel) ApiUpdateForm() (hundApiV1.FormMetricProviderUpdate, error) {
	form := hundApiV1.FormMetricProviderUpdate{}

	if s.Builtin != nil {
		err := form.FromFormMetricProviderUpdate0(s.Builtin.ApiUpdateForm())

		return form, err
	} else if s.Updown != nil {
		err := form.FromFormMetricProviderUpdate1(s.Updown.ApiUpdateForm())

		return form, err
	} else if s.Uptimerobot != nil {
		err := form.FromFormMetricProviderUpdate2(s.Uptimerobot.ApiUpdateForm())

		return form, err
	} else if s.Webhook != nil {
		err := form.FromFormMetricProviderUpdate3(s.Webhook.ApiUpdateForm())

		return form, err
	} else if s.Pingdom != nil {
		err := form.FromFormMetricProviderUpdate4(s.Pingdom.ApiUpdateForm())

		return form, err
	}

	nativeForm, err := s.NativeApiUpdateForm()

	if err != nil {
		return form, err
	}

	err = form.FromFormMetricProviderUpdate5(nativeForm)

	return form, err
}

func (s MetricProviderServiceModel) NativeApiUpdateForm() (hundApiV1.NativeFormUpdate, error) {
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

type BuiltinServiceModel struct{}

func (b BuiltinServiceModel) ApiCreateForm() hundApiV1.BuiltinFormCreate {
	return hundApiV1.BuiltinFormCreate{
		Type: hundApiV1.BuiltinFormCreateTypeBuiltin,
	}
}

func (b BuiltinServiceModel) ApiUpdateForm() hundApiV1.BuiltinFormUpdate {
	return hundApiV1.BuiltinFormUpdate{}
}

type UptimerobotMetricProviderServiceModel struct {
	MonitorApiKey types.String `tfsdk:"monitor_api_key"`
}

func (s UptimerobotMetricProviderServiceModel) ApiCreateForm() hundApiV1.UptimerobotFormCreate {
	return hundApiV1.UptimerobotFormCreate{
		Type:          hundApiV1.UptimerobotFormCreateTypeUptimerobot,
		MonitorApiKey: s.MonitorApiKey.ValueString(),
	}
}

func (s UptimerobotMetricProviderServiceModel) ApiUpdateForm() hundApiV1.UptimerobotFormUpdate {
	return hundApiV1.UptimerobotFormUpdate{
		MonitorApiKey: s.MonitorApiKey.ValueStringPointer(),
	}
}

type PingdomMetricProviderServiceModel struct {
	ApiToken  types.String `tfsdk:"api_token"`
	CheckId   types.String `tfsdk:"check_id"`
	CheckType types.String `tfsdk:"check_type"`
}

func (s PingdomMetricProviderServiceModel) ApiCreateForm() hundApiV1.PingdomFormCreate {
	return hundApiV1.PingdomFormCreate{
		Type:      hundApiV1.PingdomFormCreateTypePingdom,
		ApiToken:  s.ApiToken.ValueString(),
		CheckId:   s.CheckId.ValueString(),
		CheckType: (*hundApiV1.PingdomFormCreateCheckType)(s.CheckType.ValueStringPointer()),
	}
}

func (s PingdomMetricProviderServiceModel) ApiUpdateForm() hundApiV1.PingdomFormUpdate {
	return hundApiV1.PingdomFormUpdate{
		ApiToken:  s.ApiToken.ValueStringPointer(),
		CheckId:   s.CheckId.ValueStringPointer(),
		CheckType: (*hundApiV1.PingdomFormUpdateCheckType)(s.CheckType.ValueStringPointer()),
	}
}

type WebhookMetricProviderServiceModel struct {
	WebhookKey types.String `tfsdk:"webhook_key"`
}

func (s WebhookMetricProviderServiceModel) ApiCreateForm() hundApiV1.WebhookFormCreate {
	return hundApiV1.WebhookFormCreate{
		Type:       hundApiV1.WebhookFormCreateTypeWebhook,
		WebhookKey: s.WebhookKey.ValueStringPointer(),
	}
}

func (s WebhookMetricProviderServiceModel) ApiUpdateForm() hundApiV1.WebhookFormUpdate {
	return hundApiV1.WebhookFormUpdate{
		WebhookKey: s.WebhookKey.ValueStringPointer(),
	}
}
