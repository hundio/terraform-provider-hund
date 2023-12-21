package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
	"github.com/hundio/terraform-provider-hund/internal/validators"
)

func issueTemplateApplicationSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "An application of an IssueTemplate, which contains a copy of the template fields of the associated IssueTemplate, as well as an object of user-defined variables that parameterize the template. \n\n-> Alterations to this field do not affect the associated `issue_template_id`, and will update the Issue/Update's content accordingly. Conversely, modification/deletion of the associated IssueTemplate do not affect the attributes of this field.",
		Optional:            true,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: idFieldMarkdownDescription("IssueTemplateApplication"),
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"issue_template_id": schema.StringAttribute{
				MarkdownDescription: "The ObjectId of an IssueTemplate to use as the basis of this Application, which will inform the values for `body`, `label`, and `title` (when `kind = \"issue\"`). If this value is changed, then the application will be re-created according to the values of the given IssueTemplate.",
				Required:            true,
			},
			"body": schema.StringAttribute{
				MarkdownDescription: translationOriginalFieldMarkdownDescription("The [Liquid](https://shopify.github.io/liquid/) template for the `body` of the applied Issue/Update."),
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"body_translations": schema.MapAttribute{
				MarkdownDescription: translationFieldMarkdownDescription("The [Liquid](https://shopify.github.io/liquid/) template for the `body` of the applied Issue/Update."),
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "The template for the `label` of the applied Issue/Update.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"schema": schema.MapNestedAttribute{
				MarkdownDescription: "An object defining a set of typed variables that can be provided in `variables`. The variables can be accessed from any field in the IssueTemplate supporting Liquid.\n\n-> This field is normally copied from the underlying `issue_template`, but can be overridden here as necessary. In any case, `variables` must adhere to `schema`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
				NestedObject: issueTemplateVariablesSchema(),
			},
			"variables": schema.MapNestedAttribute{
				MarkdownDescription: "An object of variable assignments used to parameterize the associated IssueTemplate. If the associated IssueTemplate marks a variable as `required`, then it must appear here with an appropriate value. The type of each variable must match the type set in the template's schema.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"string":      schema.StringAttribute{Optional: true},
						"number":      schema.NumberAttribute{Optional: true},
						"i18n_string": schema.MapAttribute{Optional: true, ElementType: types.StringType},
						"datetime":    schema.StringAttribute{Optional: true},
					},
					Validators: []validator.Object{
						validators.ExactlyOneNonNullAttribute(),
					},
				},
			},
		},
	}
}

func issueTemplateApplicationIssueSchema() schema.SingleNestedAttribute {
	templateSchema := issueTemplateApplicationSchema()

	templateSchema.Attributes["title"] = schema.StringAttribute{
		MarkdownDescription: translationOriginalFieldMarkdownDescription("The [Liquid](https://shopify.github.io/liquid/) template for the `title` of the applied Issue."),
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}

	templateSchema.Attributes["title_translations"] = schema.MapAttribute{
		MarkdownDescription: translationFieldMarkdownDescription("The [Liquid](https://shopify.github.io/liquid/) template for the `title` of the applied Issue."),
		Optional:            true,
		Computed:            true,
		ElementType:         types.StringType,
		PlanModifiers: []planmodifier.Map{
			mapplanmodifier.UseStateForUnknown(),
		},
	}

	return templateSchema
}

func issueTemplateVariablesSchema() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				MarkdownDescription: "The expected type of this variable. One of `datetime`, `i18n-string`, `number`, or `string`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("string"),
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(hundApiV1.ISSUETEMPLATEVARIABLETYPEDatetime),
						string(hundApiV1.ISSUETEMPLATEVARIABLETYPEI18nString),
						string(hundApiV1.ISSUETEMPLATEVARIABLETYPENumber),
						string(hundApiV1.ISSUETEMPLATEVARIABLETYPEString),
					),
				},
			},
			"required": schema.BoolAttribute{
				MarkdownDescription: "Whether this variable is required when applying the template to an Issue/Update.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
	}
}

func nativeIcmpServiceSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "A Hund Native Monitoring ICMP Check.",
		Optional:            true,
		Attributes: nativeServiceSchema(map[string]schema.Attribute{
			"ip_version": schema.StringAttribute{
				MarkdownDescription: "The IP version to use when pinging.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("inet"),
				Validators: []validator.String{
					validators.NativeServiceIpVersion(),
				},
			},
			"percentage_failed_threshold": schema.Float64Attribute{
				MarkdownDescription: "The percentage of addresses at the given target that must fail for a region to be counted as failed. This option only matters when there are multiple IP addresses behind the target when the target is a domain.",
				Optional:            true,
				Computed:            true,
				Default:             float64default.StaticFloat64(0.5),
			},
		}),
	}
}

func nativeHttpServiceSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "A Hund Native Monitoring HTTP Check.",
		Optional:            true,
		Attributes: nativeServiceSchema(map[string]schema.Attribute{
			"headers": schema.MapAttribute{
				MarkdownDescription: "A list of additional HTTP headers to send to the target. The following list of\nheader names are reserved and cannot be set by a check:\n\n```\nAccept-Charset\nAccept-Encoding\nAuthentication\nConnection\nContent-Length\nDate\nHost\nKeep-Alive\nOrigin\nProxy-.*\nSec-.*\nReferer\nTE\nTrailer\nTransfer-Encoding\nUser-Agent\nVia\n```\n",
				Optional:            true,
				Computed:            true,
				Default:             mapdefault.StaticValue(types.MapValueMust(types.StringType, map[string]attr.Value{})),
				ElementType:         types.StringType,
			},
			"response_body_must_contain": schema.StringAttribute{
				MarkdownDescription: "This field supports two different matching modes (given by\n`response_body_must_contain_mode`):\n\n  `exact`: If the requested page does not contain this exact (case-sensitive)\nstring, then the check will fail.\n\n  `regex`: If the requested page does not match against the given regex, then\nthe check will fail. [Click here](https://hund.io/help/documentation/regular-expressions) for\nmore information on the use and supported syntax of Hund regexes.\n",
				Optional:            true,
			},
			"response_body_must_contain_mode": schema.StringAttribute{
				MarkdownDescription: "The response containment mode; either `exact` or `regex`. The modes are discussed\nunder `response_body_must_contain`.\n",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("exact"),
				Validators: []validator.String{
					validators.NativeServiceStringContainmentMode(),
				},
			},
			"response_code_must_be": schema.Int64Attribute{
				MarkdownDescription: "If the requested page does not return this response code, then the check will\nfail.\n",
				Optional:            true,
			},
			"ssl_verify_peer": schema.BoolAttribute{
				MarkdownDescription: "Require the target's TLS certificate to be valid.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"follow_redirects": schema.BoolAttribute{
				MarkdownDescription: "Follow any HTTP redirects given by the requested target. Please note that this check will only follow up to 9 redirects.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "An optional HTTP Basic Authentication username.",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "An optional HTTP Basic Authentication password.",
				Optional:            true,
				Sensitive:           true,
			},
		}),
	}
}

func nativeDnsServiceSchema() schema.SingleNestedAttribute {
	dns := nativeServiceSchema(map[string]schema.Attribute{
		"record_type": schema.StringAttribute{
			MarkdownDescription: "The type of DNS record to query for on the target.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					string(hundApiV1.A),
					string(hundApiV1.AAAA),
					string(hundApiV1.CNAME),
					string(hundApiV1.MX),
					string(hundApiV1.NS),
					string(hundApiV1.PTR),
					string(hundApiV1.SOA),
					string(hundApiV1.SRV),
					string(hundApiV1.TXT),
				),
			},
		},
		"nameservers": schema.ListAttribute{
			MarkdownDescription: "An optional list of nameservers to make DNS queries with. This field is\nignored by SOA queries since they use the nameservers yielded by querying NS\non the target.",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"response_containment": schema.StringAttribute{
			MarkdownDescription: "Whether `all` of the assertions in `responses_must_contain` must match the DNS response,\nor rather just `any` of them (i.e. at least one).",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					string(hundApiV1.NATIVEDNSRESPONSECONTAINMENTAll),
					string(hundApiV1.NATIVEDNSRESPONSECONTAINMENTAny),
				),
			},
		},
		"responses_must_contain": schema.SetAttribute{
			MarkdownDescription: "A set of assertions to make against the records yielded by the query. The\nformat of these assertions is *similar* to DNS record syntax, but is\nslightly simplified and allows for only asserting parts of a record's RDATA,\nrather than the entire thing. The check will fail depending on the value of\n`response_containment`.\n\n  This field is ignored by the SOA check, as it does not use assertions to\ndetermine the validity of SOA records. Instead, we ensure that every\nnameserver reported by querying NS on the target reports the same SOA serial.\nIf your target's nameservers report conflicting SOA serials, we consider the\ncheck failed.\n\n  **Example Assertions (for MX record type):**\n```json\n[\n  \"10 mail.example.com\",\n  \"spool.example.com\",\n  \"mail2.example.com\"\n]\n```\n\n  Note above how we can assert both the priority and domain (*without* the\nterminating period required by canonical DNS) of an MX record, or instead\nsimply the domain.\n",
			Optional:            true,
			ElementType:         types.StringType,
		},
	})

	dns["target"] = schema.StringAttribute{
		MarkdownDescription: "The domain/IP address that will be queried. IP addresses do not need to be\nconverted to the `z.y.x.w.in-addr.arpa` format, as this will be done\nautomatically; however, both formats are accepted.\n",
		Required:            true,
	}

	return schema.SingleNestedAttribute{
		MarkdownDescription: "A Hund Native Monitoring DNS Check.",
		Optional:            true,
		Attributes:          dns,
	}
}

func nativeTcpServiceSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "A Hund Native Monitoring TCP Check.",
		Optional:            true,
		Attributes: nativeServiceSchema(map[string]schema.Attribute{
			"ip_version": schema.StringAttribute{
				MarkdownDescription: "The IP version to use when calling the target.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("inet"),
				Validators: []validator.String{
					validators.NativeServiceIpVersion(),
				},
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "The port at the target to connect to.",
				Required:            true,
			},
			"response_must_contain": schema.StringAttribute{
				MarkdownDescription: "This field supports two different matching modes (given by `response_must_contain_mode`):\n\n  `exact`: Text that the response from the target must contain exactly\n(case-sensitive). In exact match mode, this field supports\n[escape codes](https://hund.io/help/documentation/text-field-escape-codes).\n\n  `regex`: A regex that the response from the target must match against.\n[Click here](https://hund.io/help/documentation/regular-expressions) for more information on\nthe use and supported syntax of Hund regexes.\n\n  If you send data and expect the target to reply, you must populate this field.\nLeaving this field blank will prevent the check from receiving data from the\ntarget unless forced to wait for an initial response.\n\n  The \"response\" from the target that this text is asserted against will be the\nresponse from the target *after* sending data. If data is not sent to the\ntarget, this text is asserted against the *initial* response.\n",
				Optional:            true,
			},
			"response_must_contain_mode": schema.StringAttribute{
				MarkdownDescription: "The response containment mode; either `exact` or `regex`. The modes are discussed\nunder `response_must_contain`.\n",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("exact"),
				Validators: []validator.String{
					validators.NativeServiceStringContainmentMode(),
				},
			},
			"send_data": schema.StringAttribute{
				MarkdownDescription: "Optional data to send to the target after connecting. If this field is left\nblank, nothing is sent to the target after connecting. This field supports [escape codes](https://hund.io/help/documentation/text-field-escape-codes).\n",
				Optional:            true,
			},
			"wait_for_initial_response": schema.BoolAttribute{
				MarkdownDescription: "Whether or not to wait for an initial response from the target before sending\ndata or closing the connection.\n",
				Optional:            true,
			},
		}),
	}
}

func nativeUdpServiceSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "A Hund Native Monitoring UDP Check.",
		Optional:            true,
		Attributes: nativeServiceSchema(map[string]schema.Attribute{
			"ip_version": schema.StringAttribute{
				MarkdownDescription: "The IP version to use when calling the target.\n",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("inet"),
				Validators: []validator.String{
					validators.NativeServiceIpVersion(),
				},
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "The port at the target to connect to.",
				Required:            true,
			},
			"response_must_contain": schema.StringAttribute{
				MarkdownDescription: "This field supports two different matching modes (given by `response_must_contain_mode`):\n\n  `exact`: Text that the response from the target must contain exactly\n(case-sensitive). In exact match mode, this field supports\n[escape codes](https://hund.io/help/documentation/text-field-escape-codes).\n\n  `regex`: A regex that the response from the target must match against.\n[Click here](https://hund.io/help/documentation/regular-expressions) for more information on\nthe use and supported syntax of Hund regexes.\n\n  Leaving this field blank will still cause the check to wait for a response\nfrom the target after sending data, though no assertions will be made about\nthe payload of the response.\n",
				Optional:            true,
			},
			"response_must_contain_mode": schema.StringAttribute{
				MarkdownDescription: "The response containment mode; either `exact` or `regex`. The modes are discussed\nunder `response_must_contain`.\n",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("exact"),
				Validators: []validator.String{
					validators.NativeServiceStringContainmentMode(),
				},
			},
			"send_data": schema.StringAttribute{
				MarkdownDescription: "Data to send to the target after connecting. Unlike in `tcp`, this\nfield is required. This field supports [escape codes](https://hund.io/help/documentation/text-field-escape-codes).\n",
				Required:            true,
			},
		}),
	}
}

func nativeServiceSchema(extension map[string]schema.Attribute) map[string]schema.Attribute {
	schema := map[string]schema.Attribute{
		"target": schema.StringAttribute{
			MarkdownDescription: "The host the check will make calls to.",
			Required:            true,
		},
		"consecutive_check_degraded_threshold": schema.Int64Attribute{
			MarkdownDescription: "The number of consecutive failed checks required before posting a \"degraded\"\nstatus.\n\n  Note that regardless of threshold settings, a component will post \"operational\"\nwhenever a check succeeds, thus resetting the consecutive check failure count.\n\n  When `null`, denotes that this check will not use a \"degraded\" stage\nwhen encountering check failures.\n\n  When 0, denotes that this check will post \"degraded\" upon the first check failure.\n",
			Optional:            true,
		},
		"consecutive_check_outage_threshold": schema.Int64Attribute{
			MarkdownDescription: "The number of consecutive failed checks required before posting an \"outage\"\nstatus. If `consecutive_check_degraded_threshold` is non-null, then the outage\nwill only be posted after degraded has posted according to its own threshold.\n\n  Note that regardless of threshold settings, a component will post \"operational\"\nwhenever a check succeeds, thus resetting the consecutive check failure count.\n\n  When 0, denotes that this check will post \"outage\" upon the first check failure\n(or the first check failure after \"degraded\" has been posted in case\n`consecutive_check_degraded_threshold` is set).\n",
			Optional:            true,
			Computed:            true,
			Default:             int64default.StaticInt64(1),
		},
		"frequency": schema.Int64Attribute{
			MarkdownDescription: "The frequency of the check in milliseconds. The maximum frequency is every 30\nseconds.\n\n-> Any frequency greater than every 60 seconds will force the component\nto become High-Frequency, at an additional cost. For specific pricing\ninformation, please visit the [pricing](https://hund.io/pricing) page.\n",
			Optional:            true,
			Computed:            true,
			Default:             int64default.StaticInt64(60_000),
		},
		"percentage_regions_failed_threshold": schema.Float64Attribute{
			MarkdownDescription: "The percentage of regions that must report a failed check before the entire\ncheck can be considered failed. Requiring at least two regions for this\nthreshold is recommended in order to confirm failures across regions.\n",
			Optional:            true,
			Computed:            true,
			Default:             float64default.StaticFloat64(0.5),
		},
		"regions": schema.SetAttribute{
			MarkdownDescription: "The regions you would like the target to be checked from. All regions are\nweighted equally when calculating the outcome of a check. Currently, a single\ncheck can use up to 8 regions simultaneously. Using at least two regions for a\nsingle check is recommended in order to confirm failures across regions.\n\n-> Each check may use up to **three** regions at no extra cost. Each region added to this check beyond the base three will incur an additional cost. For specific pricing information, please visit the [pricing](https://hund.io/pricing) page.\n",
			Required:            true,
			ElementType:         types.StringType,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(
					stringvalidator.OneOf(
						string(hundApiV1.AmsNl1),
						string(hundApiV1.FraDe1),
						string(hundApiV1.HelFi1),
						string(hundApiV1.LonGb1),
						string(hundApiV1.NjUs1),
						string(hundApiV1.ParFr1),
						string(hundApiV1.SinSg1),
						string(hundApiV1.SydAu1),
						string(hundApiV1.TxUs1),
						string(hundApiV1.WaUs1),
					),
				),
			},
		},
		"timeout": schema.Int64Attribute{
			MarkdownDescription: "The maximum number of milliseconds the check should wait on the host before\nfailing.\n",
			Optional:            true,
			Computed:            true,
			Default:             int64default.StaticInt64(15_000),
		},
	}

	for k, a := range extension {
		schema[k] = a
	}

	return schema
}

func watchdogServiceTypeChanged(plan types.Object, state types.Object) bool {
	stateAttrs := state.Attributes()

	for k, v := range plan.Attributes() {
		nullState := stateAttrs[k].IsNull()
		nullPlan := v.IsNull()

		if nullPlan != nullState {
			return true
		}
	}

	return false
}

func translationOriginalFieldMarkdownDescription(baseDesc string) string {
	return strings.TrimSuffix(baseDesc, ".") + ", in the default translation."
}

func translationFieldMarkdownDescription(baseDesc string) string {
	return strings.TrimSuffix(baseDesc, ".") + ", translated into multiple languages. Map keys express the language each string value is to be interpreted in. The `original` field of this map denotes the language used for the non-`_translations` version of this attribute."
}

func idFieldMarkdownDescription(objName string) string {
	return fmt.Sprintf("The ObjectId of this %s.", objName)
}

func createdAtFieldMarkdownDescription(objName string) string {
	return fmt.Sprintf("The timestamp at which this %s was created.", objName)
}

func updatedAtFieldMarkdownDescription(objName string) string {
	return fmt.Sprintf("The timestamp at which this %s was last updated.", objName)
}
