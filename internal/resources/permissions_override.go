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

var _ resource.Resource = &OverrideResource{}
var _ resource.ResourceWithConfigure = &OverrideResource{}

type OverrideResource struct {
	client *platformxe.Client
}

type OverrideResourceModel struct {
	ID        types.String `tfsdk:"id"`
	AdminID   types.String `tfsdk:"admin_id"`
	Path      types.String `tfsdk:"path"`
	Action    types.String `tfsdk:"action"`
	Effect    types.String `tfsdk:"effect"`
	Reason    types.String `tfsdk:"reason"`
	ExpiresAt types.String `tfsdk:"expires_at"`
}

func NewPermissionsOverrideResource() resource.Resource {
	return &OverrideResource{}
}

func (r *OverrideResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_permissions_override"
}

func (r *OverrideResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a PlatformXe admin permission override.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Override identifier.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"admin_id": schema.StringAttribute{
				Required:    true,
				Description: "The admin user ID this override applies to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"path": schema.StringAttribute{
				Required:    true,
				Description: "Resource path the override applies to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"action": schema.StringAttribute{
				Required:    true,
				Description: "Action the override governs.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"effect": schema.StringAttribute{
				Required:    true,
				Description: "Override effect: ALLOW or DENY.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"reason": schema.StringAttribute{
				Optional:    true,
				Description: "Reason for the override.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"expires_at": schema.StringAttribute{
				Optional:    true,
				Description: "ISO 8601 expiration timestamp for the override.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *OverrideResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OverrideResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OverrideResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := map[string]interface{}{
		"adminId": plan.AdminID.ValueString(),
		"path":    plan.Path.ValueString(),
		"action":  plan.Action.ValueString(),
		"effect":  plan.Effect.ValueString(),
	}
	if !plan.Reason.IsNull() {
		input["reason"] = plan.Reason.ValueString()
	}
	if !plan.ExpiresAt.IsNull() {
		input["expiresAt"] = plan.ExpiresAt.ValueString()
	}

	result, err := r.client.Permissions.CreateOverride(input)
	if err != nil {
		resp.Diagnostics.AddError("Create Override Failed", err.Error())
		return
	}

	if id, ok := result["id"].(string); ok {
		plan.ID = types.StringValue(id)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *OverrideResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OverrideResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.Permissions.ListOverrides(state.AdminID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Override Failed", err.Error())
		return
	}

	overrides, ok := result["overrides"].([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Read Override Failed", "Unexpected response format")
		return
	}

	var found map[string]interface{}
	for _, o := range overrides {
		override, ok := o.(map[string]interface{})
		if !ok {
			continue
		}
		if id, ok := override["id"].(string); ok && id == state.ID.ValueString() {
			found = override
			break
		}
	}

	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if v, ok := found["path"].(string); ok {
		state.Path = types.StringValue(v)
	}
	if v, ok := found["action"].(string); ok {
		state.Action = types.StringValue(v)
	}
	if v, ok := found["effect"].(string); ok {
		state.Effect = types.StringValue(v)
	}
	if v, ok := found["reason"].(string); ok {
		state.Reason = types.StringValue(v)
	}
	if v, ok := found["expiresAt"].(string); ok {
		state.ExpiresAt = types.StringValue(v)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *OverrideResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Overrides are immutable — all attributes have RequiresReplace.
	// This method should never be called, but is required by the interface.
	resp.Diagnostics.AddError("Update Not Supported", "Permission overrides are immutable. Terraform will delete and recreate as needed.")
}

func (r *OverrideResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OverrideResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Permissions.DeleteOverride(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Override Failed", err.Error())
		return
	}
}
