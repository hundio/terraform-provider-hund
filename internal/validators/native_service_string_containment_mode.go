package validators

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
)

func NativeServiceStringContainmentMode() validator.String {
	return stringvalidator.OneOf(
		string(hundApiV1.NATIVESTRINGCONTAINMENTMODEExact),
		string(hundApiV1.NATIVESTRINGCONTAINMENTMODERegex),
	)
}
