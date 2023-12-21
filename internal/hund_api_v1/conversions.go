package hundApiV1

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func Ptr[P any](obj P) *P {
	return &obj
}

func DblPtr[P any](ptr *P) **P {
	return &ptr
}

func (s ServicesWatchdog) Discriminator() string {
	var discriminator struct {
		Discriminator string `json:"type"`
	}

	json.Unmarshal(s.union, &discriminator)

	return discriminator.Discriminator
}

func (s ServicesMetricProvider) Discriminator() string {
	var discriminator struct {
		Discriminator string `json:"type"`
	}

	json.Unmarshal(s.union, &discriminator)

	return discriminator.Discriminator
}

func (n Native) Discriminator() string {
	var discriminator struct {
		Discriminator string `json:"method"`
	}

	json.Unmarshal(n.union, &discriminator)

	return discriminator.Discriminator
}

func ToFloat32Ptr(fl *float64) *float32 {
	if fl != nil {
		fl32 := float32(*fl)
		return &fl32
	}

	return nil
}

func ToIntTimestamp(ts string) (int64, error) {
	t, err := time.Parse(time.RFC3339, ts)

	if err != nil {
		return 0, err
	}

	return t.Unix(), nil
}

func ToIntTimestampPtr(ts *string) (*int64, error) {
	if ts != nil {
		t, err := ToIntTimestamp(*ts)
		return &t, err
	}

	return nil, nil
}

func ToStringTimestamp(ts int64) string {
	return time.Unix(ts, 0).UTC().Format(time.RFC3339)
}

func ToStringTimestampPtr(ts *int64) *string {
	if ts != nil {
		str := ToStringTimestamp(*ts)

		return &str
	}

	return nil
}

func ToInt64Ptr(i *int) *int64 {
	if i != nil {
		i64 := int64(*i)
		return &i64
	}

	return nil
}

func ToIntPtr(i64 *int64) *int {
	if i64 != nil {
		i := int(*i64)
		return &i
	}

	return nil
}

func ToStringList(list []basetypes.StringValue) []string {
	strList := []string{}

	for _, strVal := range list {
		strList = append(strList, strVal.ValueString())
	}

	return strList
}

func ToI18nStringPtr(orig basetypes.StringValue, i18n basetypes.MapValue) (*I18nString, error) {
	i18nString := I18nString{}

	if i18n.IsNull() || i18n.IsUnknown() {
		str := orig.ValueStringPointer()
		if str != nil && !orig.IsUnknown() {
			err := i18nString.FromI18nString0(*str)

			return &i18nString, err
		}

		return nil, nil
	}

	i18nMap := map[string]string{}

	for k, v := range i18n.Elements() {
		strValue, ok := v.(basetypes.StringValue)

		if !ok {
			return nil, fmt.Errorf("i18n key %q has non-string type: %q", k, v.String())
		}

		i18nMap[k] = strValue.ValueString()
	}

	err := i18nString.FromI18nString1(i18nMap)

	return &i18nString, err
}

func ToI18nString(orig basetypes.StringValue, i18n basetypes.MapValue) (I18nString, error) {
	i18nStr, err := ToI18nStringPtr(orig, i18n)

	if i18nStr != nil {
		return *i18nStr, err
	}

	blank := I18nString{}
	blank.FromI18nString0("")

	return blank, err
}

func FromI18nStringWithoutOriginal(i18nString I18nString) (*basetypes.MapValue, diag.Diagnostics) {
	diag := diag.Diagnostics{}

	i18n, err := i18nString.AsI18nString1()

	if err != nil {
		diag.AddError(
			"Could not parse Hund I18nString",
			"Error: "+err.Error(),
		)

		return nil, diag
	}

	attrMap := make(map[string]attr.Value)

	for lang, str := range i18n {
		attrMap[lang] = types.StringValue(str)
	}

	i18nMap, diag0 := types.MapValue(types.StringType, attrMap)
	diag.Append(diag0...)

	return &i18nMap, diag
}

func FromI18nString(i18nString I18nString) (*basetypes.StringValue, *basetypes.MapValue, diag.Diagnostics) {
	diag := diag.Diagnostics{}

	i18n, err := i18nString.AsI18nString1()

	if err != nil {
		diag.AddError(
			"Could not parse Hund I18nString",
			"Error: "+err.Error(),
		)

		return nil, nil, diag
	}

	attrMap := make(map[string]attr.Value)

	origLang, ok := i18n["original"]
	if !ok {
		diag.AddError(
			"Could not parse Hund I18nString",
			"The `original` field was missing",
		)

		return nil, nil, diag
	}

	originalStr, ok := i18n[origLang]

	var original basetypes.StringValue
	if ok {
		original = types.StringValue(originalStr)
	} else {
		original = types.StringNull()
	}

	for lang, str := range i18n {
		attrMap[lang] = types.StringValue(str)
	}

	i18nMap, diag0 := types.MapValue(types.StringType, attrMap)
	diag.Append(diag0...)

	return &original, &i18nMap, diag
}

func UnwrapUpdateExpansionary(ux UpdateExpansionary) (*Update, error) {
	issueId, err := ux.Issue.AsUpdateExpansionaryIssue0()
	if err != nil {
		return nil, err
	}

	return &Update{
		Id:             ux.Id,
		Issue:          issueId,
		Body:           ux.Body,
		BodyHtml:       ux.BodyHtml,
		CreatedAt:      ux.CreatedAt,
		Effective:      ux.Effective,
		EffectiveAfter: ux.EffectiveAfter,
		Label:          ux.Label,
		Reopening:      ux.Reopening,
		StateOverride:  ux.StateOverride,
		UpdatedAt:      ux.UpdatedAt,
		Type:           UpdateType(ux.Type),
		Template:       ux.Template,
	}, nil
}
