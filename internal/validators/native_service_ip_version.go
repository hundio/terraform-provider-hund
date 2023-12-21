package validators

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
)

func NativeServiceIpVersion() validator.String {
	return stringvalidator.OneOf(
		string(hundApiV1.NATIVEIPVERSIONInet),
		string(hundApiV1.NATIVEIPVERSIONInet6),
	)
}
