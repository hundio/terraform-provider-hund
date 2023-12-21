package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hundio/terraform-provider-hund/internal/models"
)

func PlanModifyIssueTemplateIssueApplication(ctx context.Context, rootPath path.Path, state *models.IssueTemplateApplicationIssueModel, plan *models.IssueTemplateApplicationIssueModel, resp *resource.ModifyPlanResponse) {
	if state == nil || plan == nil {
		return
	}

	// if !state.IssueTemplateId.Equal(plan.IssueTemplateId) || !state.Label.Equal(plan.Label) {
	// }

	if !state.Title.Equal(plan.Title) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, rootPath.AtName("title_translations"), types.MapUnknown(types.StringType))...)
	}

	if !state.TitleTranslations.Equal(plan.TitleTranslations) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, rootPath.AtName("title"), types.StringUnknown())...)
	}

	if !state.Body.Equal(plan.Body) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, rootPath.AtName("body_translations"), types.MapUnknown(types.StringType))...)
	}

	if !state.BodyTranslations.Equal(plan.BodyTranslations) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, rootPath.AtName("body"), types.StringUnknown())...)
	}
}

func PlanModifyIssueTemplateUpdateApplication(ctx context.Context, rootPath path.Path, state *models.IssueTemplateApplicationUpdateModel, plan *models.IssueTemplateApplicationUpdateModel, resp *resource.ModifyPlanResponse) {
	if state == nil || plan == nil {
		return
	}

	if !state.Body.Equal(plan.Body) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, rootPath.AtName("body_translations"), types.MapUnknown(types.StringType))...)
	}

	if !state.BodyTranslations.Equal(plan.BodyTranslations) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, rootPath.AtName("body"), types.StringUnknown())...)
	}
}
