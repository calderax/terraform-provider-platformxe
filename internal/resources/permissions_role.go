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

var _ resource.Resource = &RoleResource{}
var _ resource.ResourceWithConfigure = &RoleResource{}

type RoleResource struct {
	client *platformxe.Client
}

type RoleResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Model       types.String `tfsdk:"model"`
}

func NewPermissionsRoleResource() resource.Resource {
	return &RoleResource{}
}

func (r *RoleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_permissions_role"
}

func (r *RoleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a PlatformXe permission role.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Role identifier.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Role name (1-100 chars).",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Role description.",
			},
			"model": schema.StringAttribute{
				Optional:    true,
				Description: "Permission model: SIMPLE or FULL. Defaults to SIMPLE.",
			},
		},
	}
}

func (r *RoleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RoleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	model := "SIMPLE"
	if !plan.Model.IsNull() {
		model = plan.Model.ValueString()
	}

	desc := ""
	if !plan.Description.IsNull() {
		desc = plan.Description.ValueString()
	}

	result, err := r.client.Permissions.CreateRole(plan.Name.ValueString(), desc, model)
	if err != nil {
		resp.Diagnostics.AddError("Create Role Failed", err.Error())
		return
	}

	if id, ok := result["id"].(string); ok {
		plan.ID = types.StringValue(id)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *RoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RoleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.Permissions.GetRole(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Role Failed", err.Error())
		return
	}

	if name, ok := result["name"].(string); ok {
		state.Name = types.StringValue(name)
	}
	if desc, ok := result["description"].(string); ok {
		state.Description = types.StringValue(desc)
	}
	if model, ok := result["model"].(string); ok {
		state.Model = types.StringValue(model)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *RoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan RoleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updates := map[string]interface{}{}
	if !plan.Name.IsNull() {
		updates["name"] = plan.Name.ValueString()
	}
	if !plan.Description.IsNull() {
		updates["description"] = plan.Description.ValueString()
	}

	_, err := r.client.Permissions.UpdateRole(plan.ID.ValueString(), updates)
	if err != nil {
		resp.Diagnostics.AddError("Update Role Failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *RoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RoleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Permissions.DeleteRole(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Role Failed", err.Error())
		return
	}
}
