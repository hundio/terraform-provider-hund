package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
)

type NativeServiceModel interface {
	GetFrequency() types.Int64
}

func ToNativeServiceModel(native hundApiV1.Native) (NativeServiceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	switch native.Discriminator() {
	case "icmp":
		icmp, err := native.AsNative0()

		if err != nil {
			diags.Append(NativeServiceDecodeError(native.Discriminator(), err))

			return nil, diags
		}

		return &NativeIcmpServiceModel{
			IpVersion:                 types.StringValue(string(icmp.IpVersion)),
			PercentageFailedThreshold: types.Float64Value(float64(icmp.PercentageFailedThreshold)),

			Target:                            types.StringValue(icmp.Target),
			ConsecutiveCheckDegradedThreshold: types.Int64PointerValue(hundApiV1.ToInt64Ptr(icmp.ConsecutiveCheckDegradedThreshold)),
			ConsecutiveCheckOutageThreshold:   types.Int64Value(int64(icmp.ConsecutiveCheckOutageThreshold)),
			Frequency:                         types.Int64Value(int64(icmp.Frequency)),
			PercentageRegionsFailedThreshold:  types.Float64Value(float64(icmp.PercentageRegionsFailedThreshold)),
			Regions:                           ToNativeRegionModel(icmp.Regions),
			Timeout:                           types.Int64Value(int64(icmp.Timeout)),
		}, diags
	case "http":
		http, err := native.AsNative1()

		if err != nil {
			diags.Append(NativeServiceDecodeError(native.Discriminator(), err))

			return nil, diags
		}

		return &NativeHttpServiceModel{
			Headers:                     ToNativeHttpHeadersModel(http.Headers),
			ResponseBodyMustContain:     types.StringPointerValue(http.ResponseBodyMustContain),
			ResponseBodyMustContainMode: types.StringValue(string(http.ResponseBodyMustContainMode)),
			ResponseCodeMustBe:          types.Int64PointerValue(hundApiV1.ToInt64Ptr(http.ResponseCodeMustBe)),
			SslVerifyPeer:               types.BoolValue(http.SslVerifyPeer),
			FollowRedirects:             types.BoolValue(http.FollowRedirects),
			Username:                    types.StringPointerValue(http.Username),

			Target:                            types.StringValue(http.Target),
			ConsecutiveCheckDegradedThreshold: types.Int64PointerValue(hundApiV1.ToInt64Ptr(http.ConsecutiveCheckDegradedThreshold)),
			ConsecutiveCheckOutageThreshold:   types.Int64Value(int64(http.ConsecutiveCheckOutageThreshold)),
			Frequency:                         types.Int64Value(int64(http.Frequency)),
			PercentageRegionsFailedThreshold:  types.Float64Value(float64(http.PercentageRegionsFailedThreshold)),
			Regions:                           ToNativeRegionModel(http.Regions),
			Timeout:                           types.Int64Value(int64(http.Timeout)),
		}, diags
	case "dns":
		dns, err := native.AsNative2()

		if err != nil {
			diags.Append(NativeServiceDecodeError(native.Discriminator(), err))

			return nil, diags
		}

		nameservers := []attr.Value{}

		for _, v := range dns.Nameservers {
			nameservers = append(nameservers, types.StringValue(v))
		}

		nameserversModel, diag0 := types.ListValue(types.StringType, nameservers)
		diags.Append(diag0...)

		assertions := []attr.Value{}

		for _, v := range dns.ResponsesMustContain {
			assertions = append(assertions, types.StringValue(v))
		}

		assertionsModel, diag0 := types.SetValue(types.StringType, assertions)
		diags.Append(diag0...)

		return &NativeDnsServiceModel{
			RecordType:           types.StringValue(string(dns.RecordType)),
			Nameservers:          nameserversModel,
			ResponseContainment:  types.StringValue(string(dns.ResponseContainment)),
			ResponsesMustContain: assertionsModel,

			Target:                            types.StringValue(dns.Target),
			ConsecutiveCheckDegradedThreshold: types.Int64PointerValue(hundApiV1.ToInt64Ptr(dns.ConsecutiveCheckDegradedThreshold)),
			ConsecutiveCheckOutageThreshold:   types.Int64Value(int64(dns.ConsecutiveCheckOutageThreshold)),
			Frequency:                         types.Int64Value(int64(dns.Frequency)),
			PercentageRegionsFailedThreshold:  types.Float64Value(float64(dns.PercentageRegionsFailedThreshold)),
			Regions:                           ToNativeRegionModel(dns.Regions),
			Timeout:                           types.Int64Value(int64(dns.Timeout)),
		}, diags
	case "tcp":
		tcp, err := native.AsNative3()

		if err != nil {
			diags.Append(NativeServiceDecodeError(native.Discriminator(), err))

			return nil, diags
		}

		return &NativeTcpServiceModel{
			IpVersion:               types.StringValue(string(tcp.IpVersion)),
			Port:                    types.Int64Value(int64(tcp.Port)),
			ResponseMustContain:     types.StringPointerValue(tcp.ResponseMustContain),
			ResponseMustContainMode: types.StringValue(string(tcp.ResponseMustContainMode)),
			SendData:                types.StringPointerValue(tcp.SendData),
			WaitForInitialResponse:  types.BoolValue(tcp.WaitForInitialResponse),

			Target:                            types.StringValue(tcp.Target),
			ConsecutiveCheckDegradedThreshold: types.Int64PointerValue(hundApiV1.ToInt64Ptr(tcp.ConsecutiveCheckDegradedThreshold)),
			ConsecutiveCheckOutageThreshold:   types.Int64Value(int64(tcp.ConsecutiveCheckOutageThreshold)),
			Frequency:                         types.Int64Value(int64(tcp.Frequency)),
			PercentageRegionsFailedThreshold:  types.Float64Value(float64(tcp.PercentageRegionsFailedThreshold)),
			Regions:                           ToNativeRegionModel(tcp.Regions),
			Timeout:                           types.Int64Value(int64(tcp.Timeout)),
		}, diags
	case "udp":
		udp, err := native.AsNative4()

		if err != nil {
			diags.Append(NativeServiceDecodeError(native.Discriminator(), err))

			return nil, diags
		}

		return &NativeUdpServiceModel{
			IpVersion:               types.StringValue(string(udp.IpVersion)),
			Port:                    types.Int64Value(int64(udp.Port)),
			ResponseMustContain:     types.StringPointerValue(udp.ResponseMustContain),
			ResponseMustContainMode: types.StringValue(string(udp.ResponseMustContainMode)),
			SendData:                types.StringValue(udp.SendData),

			Target:                            types.StringValue(udp.Target),
			ConsecutiveCheckDegradedThreshold: types.Int64PointerValue(hundApiV1.ToInt64Ptr(udp.ConsecutiveCheckDegradedThreshold)),
			ConsecutiveCheckOutageThreshold:   types.Int64Value(int64(udp.ConsecutiveCheckOutageThreshold)),
			Frequency:                         types.Int64Value(int64(udp.Frequency)),
			PercentageRegionsFailedThreshold:  types.Float64Value(float64(udp.PercentageRegionsFailedThreshold)),
			Regions:                           ToNativeRegionModel(udp.Regions),
			Timeout:                           types.Int64Value(int64(udp.Timeout)),
		}, diags
	default:
		diags.Append(UnknownNativeServiceError(native.Discriminator()))
	}

	return nil, diags
}

type NativeRegionModel []types.String

func ToNativeRegionModel(regions []hundApiV1.NATIVEREGION) NativeRegionModel {
	model := NativeRegionModel{}

	for _, v := range regions {
		model = append(model, types.StringValue(string(v)))
	}

	return model
}

func (s NativeRegionModel) ApiValue() []hundApiV1.NATIVEREGION {
	model := []hundApiV1.NATIVEREGION{}

	for _, v := range s {
		model = append(model, hundApiV1.NATIVEREGION(v.ValueString()))
	}

	return model
}

type NativeIcmpServiceModel struct {
	IpVersion                 types.String  `tfsdk:"ip_version"`
	PercentageFailedThreshold types.Float64 `tfsdk:"percentage_failed_threshold"`

	Target                            types.String      `tfsdk:"target"`
	ConsecutiveCheckDegradedThreshold types.Int64       `tfsdk:"consecutive_check_degraded_threshold"`
	ConsecutiveCheckOutageThreshold   types.Int64       `tfsdk:"consecutive_check_outage_threshold"`
	Frequency                         types.Int64       `tfsdk:"frequency"`
	PercentageRegionsFailedThreshold  types.Float64     `tfsdk:"percentage_regions_failed_threshold"`
	Regions                           NativeRegionModel `tfsdk:"regions"`
	Timeout                           types.Int64       `tfsdk:"timeout"`
}

func (s NativeIcmpServiceModel) GetFrequency() types.Int64 {
	return s.Frequency
}

func (s NativeIcmpServiceModel) ApiCreateForm() hundApiV1.ICMPFormCreate {
	return hundApiV1.ICMPFormCreate{
		Type:                              hundApiV1.ICMPFormCreateTypeNative,
		Method:                            hundApiV1.ICMPFormCreateMethodIcmp,
		ConsecutiveCheckDegradedThreshold: hundApiV1.DblPtr(hundApiV1.ToIntPtr(s.ConsecutiveCheckDegradedThreshold.ValueInt64Pointer())),
		ConsecutiveCheckOutageThreshold:   hundApiV1.ToIntPtr(s.ConsecutiveCheckOutageThreshold.ValueInt64Pointer()),
		Frequency:                         hundApiV1.ToIntPtr(s.Frequency.ValueInt64Pointer()),
		PercentageRegionsFailedThreshold:  hundApiV1.ToFloat32Ptr(s.PercentageRegionsFailedThreshold.ValueFloat64Pointer()),
		Regions:                           s.Regions.ApiValue(),
		Target:                            s.Target.ValueString(),
		Timeout:                           hundApiV1.ToIntPtr(s.Timeout.ValueInt64Pointer()),

		IpVersion:                 (*hundApiV1.ICMPFormCreateIpVersion)(s.IpVersion.ValueStringPointer()),
		PercentageFailedThreshold: hundApiV1.ToFloat32Ptr(s.PercentageFailedThreshold.ValueFloat64Pointer()),
	}
}

func (s NativeIcmpServiceModel) ApiUpdateForm() hundApiV1.ICMPFormUpdate {
	regions := s.Regions.ApiValue()

	return hundApiV1.ICMPFormUpdate{
		ConsecutiveCheckDegradedThreshold: hundApiV1.DblPtr(hundApiV1.ToIntPtr(s.ConsecutiveCheckDegradedThreshold.ValueInt64Pointer())),
		ConsecutiveCheckOutageThreshold:   hundApiV1.ToIntPtr(s.ConsecutiveCheckOutageThreshold.ValueInt64Pointer()),
		Frequency:                         hundApiV1.ToIntPtr(s.Frequency.ValueInt64Pointer()),
		PercentageRegionsFailedThreshold:  hundApiV1.ToFloat32Ptr(s.PercentageRegionsFailedThreshold.ValueFloat64Pointer()),
		Regions:                           &regions,
		Target:                            s.Target.ValueStringPointer(),
		Timeout:                           hundApiV1.ToIntPtr(s.Timeout.ValueInt64Pointer()),

		IpVersion:                 (*hundApiV1.NATIVEIPVERSION)(s.IpVersion.ValueStringPointer()),
		PercentageFailedThreshold: hundApiV1.ToFloat32Ptr(s.PercentageFailedThreshold.ValueFloat64Pointer()),
	}
}

type NativeHttpHeadersModel map[string]types.String

func ToNativeHttpHeadersModel(headers hundApiV1.HTTPHeaders) NativeHttpHeadersModel {
	model := NativeHttpHeadersModel{}

	for k, v := range headers {
		model[k] = types.StringValue(v)
	}

	return model
}

func (s NativeHttpHeadersModel) ApiValue() hundApiV1.HTTPHeaders {
	model := hundApiV1.HTTPHeaders{}

	for k, v := range s {
		model[k] = v.ValueString()
	}

	return model
}

type NativeHttpServiceModel struct {
	Headers                     NativeHttpHeadersModel `tfsdk:"headers"`
	ResponseBodyMustContain     types.String           `tfsdk:"response_body_must_contain"`
	ResponseBodyMustContainMode types.String           `tfsdk:"response_body_must_contain_mode"`
	ResponseCodeMustBe          types.Int64            `tfsdk:"response_code_must_be"`
	SslVerifyPeer               types.Bool             `tfsdk:"ssl_verify_peer"`
	FollowRedirects             types.Bool             `tfsdk:"follow_redirects"`
	Username                    types.String           `tfsdk:"username"`
	Password                    types.String           `tfsdk:"password"`

	Target                            types.String      `tfsdk:"target"`
	ConsecutiveCheckDegradedThreshold types.Int64       `tfsdk:"consecutive_check_degraded_threshold"`
	ConsecutiveCheckOutageThreshold   types.Int64       `tfsdk:"consecutive_check_outage_threshold"`
	Frequency                         types.Int64       `tfsdk:"frequency"`
	PercentageRegionsFailedThreshold  types.Float64     `tfsdk:"percentage_regions_failed_threshold"`
	Regions                           NativeRegionModel `tfsdk:"regions"`
	Timeout                           types.Int64       `tfsdk:"timeout"`
}

func (s NativeHttpServiceModel) GetFrequency() types.Int64 {
	return s.Frequency
}

func (s NativeHttpServiceModel) ApiCreateForm() hundApiV1.HTTPFormCreate {
	headers := s.Headers.ApiValue()

	model := hundApiV1.HTTPFormCreate{
		Type:                              hundApiV1.HTTPFormCreateTypeNative,
		Method:                            hundApiV1.HTTPFormCreateMethodHttp,
		ConsecutiveCheckDegradedThreshold: hundApiV1.DblPtr(hundApiV1.ToIntPtr(s.ConsecutiveCheckDegradedThreshold.ValueInt64Pointer())),
		ConsecutiveCheckOutageThreshold:   hundApiV1.ToIntPtr(s.ConsecutiveCheckOutageThreshold.ValueInt64Pointer()),
		Frequency:                         hundApiV1.ToIntPtr(s.Frequency.ValueInt64Pointer()),
		PercentageRegionsFailedThreshold:  hundApiV1.ToFloat32Ptr(s.PercentageRegionsFailedThreshold.ValueFloat64Pointer()),
		Regions:                           s.Regions.ApiValue(),
		Target:                            s.Target.ValueString(),
		Timeout:                           hundApiV1.ToIntPtr(s.Timeout.ValueInt64Pointer()),

		Headers:                     &headers,
		FollowRedirects:             s.FollowRedirects.ValueBoolPointer(),
		SslVerifyPeer:               s.SslVerifyPeer.ValueBoolPointer(),
		Username:                    hundApiV1.DblPtr(s.Username.ValueStringPointer()),
		Password:                    hundApiV1.DblPtr(s.Password.ValueStringPointer()),
		ResponseBodyMustContain:     hundApiV1.DblPtr(s.ResponseBodyMustContain.ValueStringPointer()),
		ResponseBodyMustContainMode: (*hundApiV1.HTTPFormCreateResponseBodyMustContainMode)(s.ResponseBodyMustContainMode.ValueStringPointer()),
		ResponseCodeMustBe:          hundApiV1.DblPtr(hundApiV1.ToIntPtr(s.ResponseCodeMustBe.ValueInt64Pointer())),
	}

	return model
}

func (s NativeHttpServiceModel) ApiUpdateForm() hundApiV1.HTTPFormUpdate {
	regions := s.Regions.ApiValue()
	headers := s.Headers.ApiValue()

	model := hundApiV1.HTTPFormUpdate{
		ConsecutiveCheckDegradedThreshold: hundApiV1.DblPtr(hundApiV1.ToIntPtr(s.ConsecutiveCheckDegradedThreshold.ValueInt64Pointer())),
		ConsecutiveCheckOutageThreshold:   hundApiV1.ToIntPtr(s.ConsecutiveCheckOutageThreshold.ValueInt64Pointer()),
		Frequency:                         hundApiV1.ToIntPtr(s.Frequency.ValueInt64Pointer()),
		PercentageRegionsFailedThreshold:  hundApiV1.ToFloat32Ptr(s.PercentageRegionsFailedThreshold.ValueFloat64Pointer()),
		Regions:                           &regions,
		Target:                            s.Target.ValueStringPointer(),
		Timeout:                           hundApiV1.ToIntPtr(s.Timeout.ValueInt64Pointer()),

		Headers:                     &headers,
		FollowRedirects:             s.FollowRedirects.ValueBoolPointer(),
		SslVerifyPeer:               s.SslVerifyPeer.ValueBoolPointer(),
		Username:                    hundApiV1.DblPtr(s.Username.ValueStringPointer()),
		Password:                    hundApiV1.DblPtr(s.Password.ValueStringPointer()),
		ResponseBodyMustContain:     hundApiV1.DblPtr(s.ResponseBodyMustContain.ValueStringPointer()),
		ResponseBodyMustContainMode: (*hundApiV1.NATIVESTRINGCONTAINMENTMODE)(s.ResponseBodyMustContainMode.ValueStringPointer()),
		ResponseCodeMustBe:          hundApiV1.DblPtr(hundApiV1.ToIntPtr(s.ResponseCodeMustBe.ValueInt64Pointer())),
	}

	return model
}

type NativeDnsServiceModel struct {
	RecordType           types.String `tfsdk:"record_type"`
	Nameservers          types.List   `tfsdk:"nameservers"`
	ResponseContainment  types.String `tfsdk:"response_containment"`
	ResponsesMustContain types.Set    `tfsdk:"responses_must_contain"`

	Target                            types.String      `tfsdk:"target"`
	ConsecutiveCheckDegradedThreshold types.Int64       `tfsdk:"consecutive_check_degraded_threshold"`
	ConsecutiveCheckOutageThreshold   types.Int64       `tfsdk:"consecutive_check_outage_threshold"`
	Frequency                         types.Int64       `tfsdk:"frequency"`
	PercentageRegionsFailedThreshold  types.Float64     `tfsdk:"percentage_regions_failed_threshold"`
	Regions                           NativeRegionModel `tfsdk:"regions"`
	Timeout                           types.Int64       `tfsdk:"timeout"`
}

func (s NativeDnsServiceModel) GetFrequency() types.Int64 {
	return s.Frequency
}

func (s NativeDnsServiceModel) ApiCreateForm() hundApiV1.DNSFormCreate {
	model := hundApiV1.DNSFormCreate{
		Type:                              hundApiV1.DNSFormCreateTypeNative,
		Method:                            hundApiV1.DNSFormCreateMethodDns,
		ConsecutiveCheckDegradedThreshold: hundApiV1.DblPtr(hundApiV1.ToIntPtr(s.ConsecutiveCheckDegradedThreshold.ValueInt64Pointer())),
		ConsecutiveCheckOutageThreshold:   hundApiV1.ToIntPtr(s.ConsecutiveCheckOutageThreshold.ValueInt64Pointer()),
		Frequency:                         hundApiV1.ToIntPtr(s.Frequency.ValueInt64Pointer()),
		PercentageRegionsFailedThreshold:  hundApiV1.ToFloat32Ptr(s.PercentageRegionsFailedThreshold.ValueFloat64Pointer()),
		Regions:                           s.Regions.ApiValue(),
		Target:                            s.Target.ValueString(),
		Timeout:                           hundApiV1.ToIntPtr(s.Timeout.ValueInt64Pointer()),

		RecordType:          hundApiV1.NATIVEDNSRECORDTYPE(s.RecordType.ValueString()),
		ResponseContainment: (*hundApiV1.DNSFormCreateResponseContainment)(s.ResponseContainment.ValueStringPointer()),
	}

	if !s.Nameservers.IsNull() && !s.Nameservers.IsUnknown() {
		nameservers := []string{}

		for _, v := range s.Nameservers.Elements() {
			nameservers = append(nameservers, v.(types.String).ValueString())
		}

		model.Nameservers = &nameservers
	}

	if !s.ResponsesMustContain.IsNull() && !s.ResponsesMustContain.IsUnknown() {
		assertions := []string{}

		for _, v := range s.ResponsesMustContain.Elements() {
			assertions = append(assertions, v.(types.String).ValueString())
		}

		model.ResponsesMustContain = &assertions
	}

	return model
}

func (s NativeDnsServiceModel) ApiUpdateForm() hundApiV1.DNSFormUpdate {
	regions := s.Regions.ApiValue()

	model := hundApiV1.DNSFormUpdate{
		ConsecutiveCheckDegradedThreshold: hundApiV1.DblPtr(hundApiV1.ToIntPtr(s.ConsecutiveCheckDegradedThreshold.ValueInt64Pointer())),
		ConsecutiveCheckOutageThreshold:   hundApiV1.ToIntPtr(s.ConsecutiveCheckOutageThreshold.ValueInt64Pointer()),
		Frequency:                         hundApiV1.ToIntPtr(s.Frequency.ValueInt64Pointer()),
		PercentageRegionsFailedThreshold:  hundApiV1.ToFloat32Ptr(s.PercentageRegionsFailedThreshold.ValueFloat64Pointer()),
		Regions:                           &regions,
		Target:                            s.Target.ValueStringPointer(),
		Timeout:                           hundApiV1.ToIntPtr(s.Timeout.ValueInt64Pointer()),

		RecordType:          (*hundApiV1.NATIVEDNSRECORDTYPE)(s.RecordType.ValueStringPointer()),
		ResponseContainment: (*hundApiV1.NATIVEDNSRESPONSECONTAINMENT)(s.ResponseContainment.ValueStringPointer()),
	}

	if !s.Nameservers.IsNull() && !s.Nameservers.IsUnknown() {
		nameservers := []string{}

		for _, v := range s.Nameservers.Elements() {
			nameservers = append(nameservers, v.(types.String).ValueString())
		}

		model.Nameservers = &nameservers
	}

	if !s.ResponsesMustContain.IsNull() && !s.ResponsesMustContain.IsUnknown() {
		assertions := []string{}

		for _, v := range s.ResponsesMustContain.Elements() {
			assertions = append(assertions, v.(types.String).ValueString())
		}

		model.ResponsesMustContain = &assertions
	}

	return model
}

type NativeTcpServiceModel struct {
	IpVersion               types.String `tfsdk:"ip_version"`
	Port                    types.Int64  `tfsdk:"port"`
	ResponseMustContain     types.String `tfsdk:"response_must_contain"`
	ResponseMustContainMode types.String `tfsdk:"response_must_contain_mode"`
	SendData                types.String `tfsdk:"send_data"`
	WaitForInitialResponse  types.Bool   `tfsdk:"wait_for_initial_response"`

	Target                            types.String      `tfsdk:"target"`
	ConsecutiveCheckDegradedThreshold types.Int64       `tfsdk:"consecutive_check_degraded_threshold"`
	ConsecutiveCheckOutageThreshold   types.Int64       `tfsdk:"consecutive_check_outage_threshold"`
	Frequency                         types.Int64       `tfsdk:"frequency"`
	PercentageRegionsFailedThreshold  types.Float64     `tfsdk:"percentage_regions_failed_threshold"`
	Regions                           NativeRegionModel `tfsdk:"regions"`
	Timeout                           types.Int64       `tfsdk:"timeout"`
}

func (s NativeTcpServiceModel) GetFrequency() types.Int64 {
	return s.Frequency
}

func (s NativeTcpServiceModel) ApiCreateForm() hundApiV1.TCPFormCreate {
	return hundApiV1.TCPFormCreate{
		Type:                              hundApiV1.TCPFormCreateTypeNative,
		Method:                            hundApiV1.TCPFormCreateMethodTcp,
		ConsecutiveCheckDegradedThreshold: hundApiV1.DblPtr(hundApiV1.ToIntPtr(s.ConsecutiveCheckDegradedThreshold.ValueInt64Pointer())),
		ConsecutiveCheckOutageThreshold:   hundApiV1.ToIntPtr(s.ConsecutiveCheckOutageThreshold.ValueInt64Pointer()),
		Frequency:                         hundApiV1.ToIntPtr(s.Frequency.ValueInt64Pointer()),
		PercentageRegionsFailedThreshold:  hundApiV1.ToFloat32Ptr(s.PercentageRegionsFailedThreshold.ValueFloat64Pointer()),
		Regions:                           s.Regions.ApiValue(),
		Target:                            s.Target.ValueString(),
		Timeout:                           hundApiV1.ToIntPtr(s.Timeout.ValueInt64Pointer()),

		IpVersion:               (*hundApiV1.TCPFormCreateIpVersion)(s.IpVersion.ValueStringPointer()),
		Port:                    int(s.Port.ValueInt64()),
		ResponseMustContain:     hundApiV1.DblPtr(s.ResponseMustContain.ValueStringPointer()),
		ResponseMustContainMode: (*hundApiV1.TCPFormCreateResponseMustContainMode)(s.ResponseMustContainMode.ValueStringPointer()),
		SendData:                hundApiV1.DblPtr(s.SendData.ValueStringPointer()),
		WaitForInitialResponse:  s.WaitForInitialResponse.ValueBoolPointer(),
	}
}

func (s NativeTcpServiceModel) ApiUpdateForm() hundApiV1.TCPFormUpdate {
	regions := s.Regions.ApiValue()

	return hundApiV1.TCPFormUpdate{
		ConsecutiveCheckDegradedThreshold: hundApiV1.DblPtr(hundApiV1.ToIntPtr(s.ConsecutiveCheckDegradedThreshold.ValueInt64Pointer())),
		ConsecutiveCheckOutageThreshold:   hundApiV1.ToIntPtr(s.ConsecutiveCheckOutageThreshold.ValueInt64Pointer()),
		Frequency:                         hundApiV1.ToIntPtr(s.Frequency.ValueInt64Pointer()),
		PercentageRegionsFailedThreshold:  hundApiV1.ToFloat32Ptr(s.PercentageRegionsFailedThreshold.ValueFloat64Pointer()),
		Regions:                           &regions,
		Target:                            s.Target.ValueStringPointer(),
		Timeout:                           hundApiV1.ToIntPtr(s.Timeout.ValueInt64Pointer()),

		IpVersion:               (*hundApiV1.NATIVEIPVERSION)(s.IpVersion.ValueStringPointer()),
		Port:                    hundApiV1.ToIntPtr(s.Port.ValueInt64Pointer()),
		ResponseMustContain:     hundApiV1.DblPtr(s.ResponseMustContain.ValueStringPointer()),
		ResponseMustContainMode: (*hundApiV1.NATIVESTRINGCONTAINMENTMODE)(s.ResponseMustContainMode.ValueStringPointer()),
		SendData:                hundApiV1.DblPtr(s.SendData.ValueStringPointer()),
		WaitForInitialResponse:  s.WaitForInitialResponse.ValueBoolPointer(),
	}
}

type NativeUdpServiceModel struct {
	IpVersion               types.String `tfsdk:"ip_version"`
	Port                    types.Int64  `tfsdk:"port"`
	ResponseMustContain     types.String `tfsdk:"response_must_contain"`
	ResponseMustContainMode types.String `tfsdk:"response_must_contain_mode"`
	SendData                types.String `tfsdk:"send_data"`

	Target                            types.String      `tfsdk:"target"`
	ConsecutiveCheckDegradedThreshold types.Int64       `tfsdk:"consecutive_check_degraded_threshold"`
	ConsecutiveCheckOutageThreshold   types.Int64       `tfsdk:"consecutive_check_outage_threshold"`
	Frequency                         types.Int64       `tfsdk:"frequency"`
	PercentageRegionsFailedThreshold  types.Float64     `tfsdk:"percentage_regions_failed_threshold"`
	Regions                           NativeRegionModel `tfsdk:"regions"`
	Timeout                           types.Int64       `tfsdk:"timeout"`
}

func (s NativeUdpServiceModel) GetFrequency() types.Int64 {
	return s.Frequency
}

func (s NativeUdpServiceModel) ApiCreateForm() hundApiV1.UDPFormCreate {
	return hundApiV1.UDPFormCreate{
		Type:                              hundApiV1.UDPFormCreateTypeNative,
		Method:                            hundApiV1.UDPFormCreateMethodUdp,
		ConsecutiveCheckDegradedThreshold: hundApiV1.DblPtr(hundApiV1.ToIntPtr(s.ConsecutiveCheckDegradedThreshold.ValueInt64Pointer())),
		ConsecutiveCheckOutageThreshold:   hundApiV1.ToIntPtr(s.ConsecutiveCheckOutageThreshold.ValueInt64Pointer()),
		Frequency:                         hundApiV1.ToIntPtr(s.Frequency.ValueInt64Pointer()),
		PercentageRegionsFailedThreshold:  hundApiV1.ToFloat32Ptr(s.PercentageRegionsFailedThreshold.ValueFloat64Pointer()),
		Regions:                           s.Regions.ApiValue(),
		Target:                            s.Target.ValueString(),
		Timeout:                           hundApiV1.ToIntPtr(s.Timeout.ValueInt64Pointer()),

		IpVersion:               (*hundApiV1.UDPFormCreateIpVersion)(s.IpVersion.ValueStringPointer()),
		Port:                    int(s.Port.ValueInt64()),
		ResponseMustContain:     hundApiV1.DblPtr(s.ResponseMustContain.ValueStringPointer()),
		ResponseMustContainMode: (*hundApiV1.UDPFormCreateResponseMustContainMode)(s.ResponseMustContainMode.ValueStringPointer()),
		SendData:                s.SendData.ValueString(),
	}
}

func (s NativeUdpServiceModel) ApiUpdateForm() hundApiV1.UDPFormUpdate {
	regions := s.Regions.ApiValue()

	return hundApiV1.UDPFormUpdate{
		ConsecutiveCheckDegradedThreshold: hundApiV1.DblPtr(hundApiV1.ToIntPtr(s.ConsecutiveCheckDegradedThreshold.ValueInt64Pointer())),
		ConsecutiveCheckOutageThreshold:   hundApiV1.ToIntPtr(s.ConsecutiveCheckOutageThreshold.ValueInt64Pointer()),
		Frequency:                         hundApiV1.ToIntPtr(s.Frequency.ValueInt64Pointer()),
		PercentageRegionsFailedThreshold:  hundApiV1.ToFloat32Ptr(s.PercentageRegionsFailedThreshold.ValueFloat64Pointer()),
		Regions:                           &regions,
		Target:                            s.Target.ValueStringPointer(),
		Timeout:                           hundApiV1.ToIntPtr(s.Timeout.ValueInt64Pointer()),

		IpVersion:               (*hundApiV1.NATIVEIPVERSION)(s.IpVersion.ValueStringPointer()),
		Port:                    hundApiV1.ToIntPtr(s.Port.ValueInt64Pointer()),
		ResponseMustContain:     hundApiV1.DblPtr(s.ResponseMustContain.ValueStringPointer()),
		ResponseMustContainMode: (*hundApiV1.NATIVESTRINGCONTAINMENTMODE)(s.ResponseMustContainMode.ValueStringPointer()),
		SendData:                s.SendData.ValueStringPointer(),
	}
}
