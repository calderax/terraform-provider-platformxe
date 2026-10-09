// =============================================================================
// Copyright 2026 Caldera Technologies Ltd.
// Proprietary and confidential.
// Unauthorized copying or distribution is prohibited.
// =============================================================================

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	platformxe "github.com/calderax/platformxe-go"
)

var _ resource.Resource = &TemplateResource{}
var _ resource.ResourceWithConfigure = &TemplateResource{}

type TemplateResource struct {
	client *platformxe.Client
}

type TemplateResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Subject types.String `tfsdk:"subject"`
	HTML    types.String `tfsdk:"html"`
}

func NewTemplatesTemplateResource() resource.Resource {
	return &TemplateResource{}
}

func (r *TemplateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_templates_template"
}

func (r *TemplateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a PlatformXe email/message template.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Template identifier.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Template name.",
			},
			"subject": schema.StringAttribute{
				Required:    true,
				Description: "Email subject line (supports template variables).",
			},
			"html": schema.StringAttribute{
				Required:    true,
				Description: "HTML body of the template.",
			},
		},
	}
}

func (r *TemplateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*platformxe.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", "Expected *platformxe.Client")
		return
	}
	r.client = client
}

func (r *TemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := map[string]interface{}{
		"name":    plan.Name.ValueString(),
		"subject": plan.Subject.ValueString(),
		"html":    plan.HTML.ValueString(),
	}

	result, err := r.client.Templates.Create(input)
	if err != nil {
		resp.Diagnostics.AddError("Create Template Failed", err.Error())
		return
	}

	if id, ok := result["id"].(string); ok {
		plan.ID = types.StringValue(id)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *TemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.Templates.Get(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Template Failed", err.Error())
		return
	}

	if v, ok := result["name"].(string); ok {
		state.Name = types.StringValue(v)
	}
	if v, ok := result["subject"].(string); ok {
		state.Subject = types.StringValue(v)
	}
	if v, ok := result["html"].(string); ok {
		state.HTML = types.StringValue(v)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *TemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := map[string]interface{}{
		"name":    plan.Name.ValueString(),
		"subject": plan.Subject.ValueString(),
		"html":    plan.HTML.ValueString(),
	}

	_, err := r.client.Templates.Update(plan.ID.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("Update Template Failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *TemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Templates.Delete(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Template Failed", err.Error())
		return
	}
}
