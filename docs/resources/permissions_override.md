---
page_title: "platformxe_permissions_override Resource"
description: "Manages a PlatformXe per-user permission override."
---

# platformxe_permissions_override

Manages a per-user permission override that takes precedence over role-based permissions.

## Example Usage

```hcl
resource "platformxe_permissions_override" "block_admin" {
  admin_id = "usr_123"
  path     = "admin/settings"
  action   = "WRITE"
  effect   = "DENY"
  reason   = "Temporary restriction during audit"
}
```

## Argument Reference

- `admin_id` (String, Required) -- Target user ID.
- `path` (String, Required) -- Permission path.
- `action` (String, Required) -- Action (e.g., `READ`, `WRITE`, `DELETE`).
- `effect` (String, Required) -- `ALLOW` or `DENY`.
- `reason` (String, Optional) -- Human-readable reason for the override.

## Attribute Reference

- `id` (String) -- Override ID, assigned by PlatformXe.

## Import

Permission overrides can be imported using their ID:

```bash
terraform import platformxe_permissions_override.block_admin ovr_abc123
```
