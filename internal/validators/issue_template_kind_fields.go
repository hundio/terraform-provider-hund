package validators

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func IssueTemplateKindFields() issueTemplateKindFields {
	return issueTemplateKindFields{}
}

type issueTemplateKindFields struct{}

var _ resource.ConfigValidator = &issueTemplateKindFields{}

func (v issueTemplateKindFields) Description(ctx context.Context) string {
	return "Validate an Issue Template does not set incompatible fields depending on `kind`."
}

func (v issueTemplateKindFields) MarkdownDescription(ctx context.Context) string {
	return "Validate an Issue Template does not set incompatible fields depending on `kind`."
}

func (v issueTemplateKindFields) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var kind types.String

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("kind"), &kind)...)

	if kind.ValueString() == "update" {
		var title, titleTranslations attr.Value

		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("title"), &title)...)
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("title_translations"), &titleTranslations)...)

		if !title.IsNull() || !titleTranslations.IsNull() {
			var errorPath path.Path

			if !title.IsNull() {
				errorPath = path.Root("title")
			} else {
				errorPath = path.Root("title_translations")
			}

			resp.Diagnostics.AddAttributeError(
				errorPath,
				"Incompatible Issue Template Field Given",
				"Issue Templates of kind `update` cannot set title fields.",
			)
		}
	}
}
