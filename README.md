# Terraform Provider for PlatformXe

> **Available to paid tenants only.** The Terraform provider requires a Basic, Pro, or Enterprise plan. [Sign up at platformxe.com](https://platformxe.com/portal/register) and upgrade to get access.

Manage PlatformXe infrastructure as code — roles, policies, webhooks, workflows, sending domains, event subscriptions, and federation groups.

## Installation

Add the provider to your Terraform configuration:

```hcl
terraform {
  required_providers {
    platformxe = {
      source = "calderax/platformxe"
    }
  }
}
```

Then run:

```bash
terraform init
```

## Provider Configuration

```hcl
provider "platformxe" {
  api_key  = var.platformxe_api_key   # or set PLATFORMXE_API_KEY env var
  base_url = "https://platformxe.com" # optional, this is the default
}
```

| Attribute  | Type   | Required | Description |
|------------|--------|----------|-------------|
| `api_key`  | string | Yes*     | API key for authentication. Can also be set via `PLATFORMXE_API_KEY` environment variable. |
| `base_url` | string | No       | API base URL. Defaults to `https://platformxe.com`. |

## Resources (10)

| Resource | Description |
|----------|-------------|
| `platformxe_role` | Permission roles (SIMPLE or FULL model) |
| `platformxe_resource_policy` | ABAC resource policies with path, action, condition, priority |
| `platformxe_permission_override` | Per-user permission overrides with optional expiry |
| `platformxe_webhook` | Webhook endpoints with auto-generated signing secret |
| `platformxe_template` | Email/message templates with subject and HTML body |
| `platformxe_workflow` | Event-driven workflow automations with trigger config and actions |
| `platformxe_sending_domain` | Sending domains for email delivery (returns DNS records for verification) |
| `platformxe_event_subscription` | Event subscriptions with webhook delivery |
| `platformxe_federation_group` | Federation groups for multi-app permission orchestration (Enterprise) |
| `platformxe_federation_member` | Members within a federation group (organization_id + prefix) |

## Data Sources (3)

| Data Source | Description |
|-------------|-------------|
| `platformxe_tenant` | Read-only tenant details resolved from API key (id, name, slug, plan, region) |
| `platformxe_modules` | List registered permission modules |
| `platformxe_identity_providers` | Identity provider health and circuit breaker status |

## Examples

See [`examples/main.tf`](examples/main.tf) for a complete working example that demonstrates all resources.

## Development

Build the provider:

```bash
cd packages/terraform-provider
go mod tidy
go build -o terraform-provider-platformxe
```

Install locally for testing:

```bash
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/calderax/platformxe/0.1.0/$(go env GOOS)_$(go env GOARCH)
cp terraform-provider-platformxe ~/.terraform.d/plugins/registry.terraform.io/calderax/platformxe/0.1.0/$(go env GOOS)_$(go env GOARCH)/
```

## License

Proprietary. Built by Caldera Technologies Ltd.
