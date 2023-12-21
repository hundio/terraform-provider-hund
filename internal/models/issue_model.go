package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
)

// IssueModel describes the data source data model.
type IssueModel struct {
	Id                   types.String   `tfsdk:"id"`
	CreatedAt            types.String   `tfsdk:"created_at"`
	UpdatedAt            types.String   `tfsdk:"updated_at"`
	BeganAt              types.String   `tfsdk:"began_at"`
	EndedAt              types.String   `tfsdk:"ended_at"`
	CancelledAt          types.String   `tfsdk:"cancelled_at"`
	Title                types.String   `tfsdk:"title"`
	Body                 types.String   `tfsdk:"body"`
	BodyHtml             types.String   `tfsdk:"body_html"`
	TitleTranslations    types.Map      `tfsdk:"title_translations"`
	BodyTranslations     types.Map      `tfsdk:"body_translations"`
	BodyHtmlTranslations types.Map      `tfsdk:"body_html_translations"`
	Label                types.String   `tfsdk:"label"`
	ComponentIds         []types.String `tfsdk:"component_ids"`
	Updates              []UpdateModel  `tfsdk:"updates"`
	Schedule             *ScheduleModel `tfsdk:"schedule"`
	Duration             types.Int64    `tfsdk:"duration"`
	OpenGraphImageUrl    types.String   `tfsdk:"open_graph_image_url"`
	Priority             types.Int64    `tfsdk:"priority"`
	Resolved             types.Bool     `tfsdk:"resolved"`
	Retrospective        types.Bool     `tfsdk:"retrospective"`
	Scheduled            types.Bool     `tfsdk:"scheduled"`
	Specialization       types.String   `tfsdk:"specialization"`
	Standing             types.Bool     `tfsdk:"standing"`
	StateOverride        types.Int64    `tfsdk:"state_override"`

	Template *IssueTemplateApplicationIssueModel `tfsdk:"template"`

	ArchiveOnDestroy types.Bool `tfsdk:"archive_on_destroy"`
}

func ToIssueModel(ctx context.Context, issue hundApiV1.Issue) (IssueModel, diag.Diagnostics) {
	diag := diag.Diagnostics{}

	model := IssueModel{
		Id:                types.StringValue(issue.Id),
		CreatedAt:         types.StringValue(hundApiV1.ToStringTimestamp(issue.CreatedAt)),
		UpdatedAt:         types.StringValue(hundApiV1.ToStringTimestamp(issue.UpdatedAt)),
		BeganAt:           types.StringValue(hundApiV1.ToStringTimestamp(issue.BeganAt)),
		EndedAt:           types.StringPointerValue(hundApiV1.ToStringTimestampPtr(issue.EndedAt)),
		CancelledAt:       types.StringPointerValue(hundApiV1.ToStringTimestampPtr(issue.CancelledAt)),
		Duration:          types.Int64Value(int64(issue.Duration)),
		OpenGraphImageUrl: types.StringPointerValue(issue.OpenGraphImageUrl),
		Priority:          types.Int64Value(int64(issue.Priority)),
		Resolved:          types.BoolValue(issue.Resolved),
		Retrospective:     types.BoolValue(issue.Retrospective),
		Scheduled:         types.BoolValue(issue.Scheduled),
		Specialization:    types.StringValue((string)(issue.Specialization)),
		Standing:          types.BoolValue(issue.Standing),
		StateOverride:     types.Int64PointerValue(hundApiV1.ToInt64Ptr((*int)(issue.StateOverride))),
		Label:             types.StringPointerValue((*string)(issue.Label)),
	}

	title, titleMap, diag0 := hundApiV1.FromI18nString(issue.Title)
	diag.Append(diag0...)

	body, bodyMap, diag0 := hundApiV1.FromI18nString(issue.Body)
	diag.Append(diag0...)

	bodyHtml, bodyHtmlMap, diag0 := hundApiV1.FromI18nString(issue.BodyHtml)
	diag.Append(diag0...)

	if diag.HasError() {
		return model, diag
	}

	model.Title = *title
	model.TitleTranslations = *titleMap

	model.Body = *body
	model.BodyTranslations = *bodyMap

	model.BodyHtml = *bodyHtml
	model.BodyHtmlTranslations = *bodyHtmlMap

	if issue.Schedule != nil {
		model.Schedule = &ScheduleModel{
			Id:                  types.StringValue(issue.Schedule.Id),
			Started:             types.BoolValue(issue.Schedule.Started),
			Ended:               types.BoolValue(issue.Schedule.Ended),
			Notified:            types.BoolValue(issue.Schedule.Notified),
			StartsAt:            types.StringValue(hundApiV1.ToStringTimestamp(issue.Schedule.StartsAt)),
			EndsAt:              types.StringValue(hundApiV1.ToStringTimestamp(issue.Schedule.EndsAt)),
			NotifySubscribersAt: types.StringPointerValue(hundApiV1.ToStringTimestampPtr(issue.Schedule.NotifySubscribersAt)),
		}
	}

	for _, component := range issue.Components.Data {
		model.ComponentIds = append(model.ComponentIds, types.StringValue(component.Id))
		// _, nameMap, diag0 := hundApiV1.FromI18nString(component.Name)
		// diag.Append(diag0...)

		// var components []ComponentModel
		// components = append(components, ComponentModel{
		// 	Id:   types.StringValue(component.Id),
		// 	Name: *nameMap,
		// })

		// model.Components, diag0 = types.ListValue(ComponentModel, components)
	}

	model.Updates = []UpdateModel{}

	for _, update := range issue.Updates.Data {
		updateModel, diag0 := ToUpdateModel(ctx, update)
		diag.Append(diag0...)

		model.Updates = append(model.Updates, updateModel)
	}

	if issue.Template != nil {
		model.Template, diag0 = ToIssueTemplateApplicationIssueModel(ctx, *issue.Template)
		diag.Append(diag0...)
	}

	return model, diag
}

func ToUpdateModel(ctx context.Context, update hundApiV1.Update) (UpdateModel, diag.Diagnostics) {
	diag := diag.Diagnostics{}

	model := UpdateModel{
		Id:             types.StringValue(update.Id),
		IssueId:        types.StringValue(update.Issue),
		CreatedAt:      types.StringValue(hundApiV1.ToStringTimestamp(update.CreatedAt)),
		UpdatedAt:      types.StringValue(hundApiV1.ToStringTimestamp(update.UpdatedAt)),
		EffectiveAfter: types.StringValue(hundApiV1.ToStringTimestamp(update.EffectiveAfter)),
		Effective:      types.BoolValue(update.Effective),
		Reopening:      types.BoolValue(update.Reopening),
		Label:          types.StringPointerValue((*string)(update.Label)),
		StateOverride:  types.Int64PointerValue(hundApiV1.ToInt64Ptr((*int)(update.StateOverride))),
	}

	if update.Body != nil {
		body, bodyMap, diag0 := hundApiV1.FromI18nString(*update.Body)
		diag.Append(diag0...)

		model.Body = *body
		model.BodyTranslations = *bodyMap
	}

	bodyHtml, bodyHtmlMap, diag0 := hundApiV1.FromI18nString(update.BodyHtml)
	diag.Append(diag0...)

	model.BodyHtml = *bodyHtml
	model.BodyHtmlTranslations = *bodyHtmlMap

	if update.Template != nil {
		model.Template, diag0 = ToIssueTemplateApplicationUpdateModel(ctx, *update.Template)
		diag.Append(diag0...)
	}

	return model, diag
}

type ScheduleModel struct {
	Id                  types.String `tfsdk:"id"`
	Started             types.Bool   `tfsdk:"started"`
	Ended               types.Bool   `tfsdk:"ended"`
	Notified            types.Bool   `tfsdk:"notified"`
	StartsAt            types.String `tfsdk:"starts_at"`
	EndsAt              types.String `tfsdk:"ends_at"`
	NotifySubscribersAt types.String `tfsdk:"notify_subscribers_at"`
}

type UpdateModel struct {
	Id                   types.String `tfsdk:"id"`
	IssueId              types.String `tfsdk:"issue_id"`
	CreatedAt            types.String `tfsdk:"created_at"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
	Effective            types.Bool   `tfsdk:"effective"`
	Reopening            types.Bool   `tfsdk:"reopening"`
	EffectiveAfter       types.String `tfsdk:"effective_after"`
	Body                 types.String `tfsdk:"body"`
	BodyHtml             types.String `tfsdk:"body_html"`
	BodyTranslations     types.Map    `tfsdk:"body_translations"`
	BodyHtmlTranslations types.Map    `tfsdk:"body_html_translations"`
	Label                types.String `tfsdk:"label"`
	StateOverride        types.Int64  `tfsdk:"state_override"`

	Template *IssueTemplateApplicationUpdateModel `tfsdk:"template"`

	ArchiveOnDestroy types.Bool `tfsdk:"archive_on_destroy"`
}

type UpdateList struct{}

var _ planmodifier.List = UpdateList{}

func (u UpdateList) Description(ctx context.Context) string {
	return "Plan modifier for Update lists."
}

func (u UpdateList) MarkdownDescription(ctx context.Context) string {
	return "Plan modifier for Update lists."
}

func (u UpdateList) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if !req.PlanValue.IsUnknown() {
		return
	}

	var endedAt types.String
	diag0 := req.Plan.GetAttribute(ctx, path.Root("ended_at"), &endedAt)
	resp.Diagnostics.Append(diag0...)

	attr := req.Config.Schema.GetAttributes()["updates"]

	schema := attr.GetType().(types.ListType)

	var diag diag.Diagnostics
	if endedAt.IsUnknown() || endedAt.IsNull() {
		resp.PlanValue, diag = types.ListValueFrom(ctx, schema.ElemType, []UpdateModel{})
	} else {
		resp.PlanValue, diag = types.ListValueFrom(ctx, schema.ElemType, []UpdateModel{
			{
				Id:                   types.StringUnknown(),
				IssueId:              types.StringUnknown(),
				CreatedAt:            types.StringUnknown(),
				UpdatedAt:            types.StringUnknown(),
				Effective:            types.BoolUnknown(),
				Reopening:            types.BoolUnknown(),
				Body:                 types.StringUnknown(),
				BodyHtml:             types.StringUnknown(),
				BodyTranslations:     types.MapUnknown(types.StringType),
				BodyHtmlTranslations: types.MapUnknown(types.StringType),
				Label:                types.StringUnknown(),
				StateOverride:        types.Int64Unknown(),
				EffectiveAfter:       endedAt,
				// ArchiveOnDestroy: basetypes.BoolValue{},
			},
		})
	}

	resp.Diagnostics.Append(diag...)
}
