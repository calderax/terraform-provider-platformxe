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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	platformxe "github.com/calderax/platformxe-go"
)

var _ resource.Resource = &SubscriptionResource{}
var _ resource.ResourceWithConfigure = &SubscriptionResource{}

type SubscriptionResource struct {
	client *platformxe.Client
}

type SubscriptionResourceModel struct {
	ID         types.String `tfsdk:"id"`
	EventTypes types.List   `tfsdk:"event_types"`
	WebhookURL types.String `tfsdk:"webhook_url"`
	IsActive   types.Bool   `tfsdk:"is_active"`
}

func NewEventsSubscriptionResource() resource.Resource {
	return &SubscriptionResource{}
}

func (r *SubscriptionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_events_subscription"
}

func (r *SubscriptionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a PlatformXe event subscription.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Subscription identifier.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"event_types": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "List of event types to subscribe to.",
			},
			"webhook_url": schema.StringAttribute{
				Required:    true,
				Description: "URL to receive event notifications.",
			},
			"is_active": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether the subscription is active. Defaults to true.",
			},
		},
	}
}

func (r *SubscriptionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SubscriptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SubscriptionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var eventTypes []string
	resp.Diagnostics.Append(plan.EventTypes.ElementsAs(ctx, &eventTypes, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := map[string]interface{}{
		"eventTypes": eventTypes,
		"webhookUrl": plan.WebhookURL.ValueString(),
		"isActive":   plan.IsActive.ValueBool(),
	}

	result, err := r.client.Subscriptions.Create(input)
	if err != nil {
		resp.Diagnostics.AddError("Create Subscription Failed", err.Error())
		return
	}

	if id, ok := result["id"].(string); ok {
		plan.ID = types.StringValue(id)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *SubscriptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SubscriptionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.Subscriptions.Get(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Subscription Failed", err.Error())
		return
	}

	if v, ok := result["webhookUrl"].(string); ok {
		state.WebhookURL = types.StringValue(v)
	}
	if v, ok := result["isActive"].(bool); ok {
		state.IsActive = types.BoolValue(v)
	}

	if evts, ok := result["eventTypes"].([]interface{}); ok {
		eventVals := make([]types.String, len(evts))
		for i, e := range evts {
			if s, ok := e.(string); ok {
				eventVals[i] = types.StringValue(s)
			}
		}
		listVal, diags := types.ListValueFrom(ctx, types.StringType, eventVals)
		resp.Diagnostics.Append(diags...)
		state.EventTypes = listVal
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *SubscriptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SubscriptionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var eventTypes []string
	resp.Diagnostics.Append(plan.EventTypes.ElementsAs(ctx, &eventTypes, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := map[string]interface{}{
		"eventTypes": eventTypes,
		"webhookUrl": plan.WebhookURL.ValueString(),
		"isActive":   plan.IsActive.ValueBool(),
	}

	_, err := r.client.Subscriptions.Update(plan.ID.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("Update Subscription Failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *SubscriptionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SubscriptionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Subscriptions.Delete(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Subscription Failed", err.Error())
		return
	}
}
