---
page_title: "platformxe_templates_template Resource"
description: "Manages a PlatformXe content template."
---

# platformxe_templates_template

Manages a reusable content template for email and other messaging.

## Example Usage

```hcl
resource "platformxe_templates_template" "welcome" {
  name      = "welcome-email"
  subject   = "Welcome to {{company}}"
  html      = file("${path.module}/templates/welcome.html")
  variables = ["company", "name"]
}
```

## Argument Reference

- `name` (String, Required) -- Template identifier.
- `subject` (String, Optional) -- Email subject line. Supports `{{variable}}` interpolation.
- `html` (String, Required) -- HTML body. Supports `{{variable}}` interpolation.
- `variables` (List of String, Optional) -- Declared template variables.

## Attribute Reference

- `id` (String) -- Template ID, assigned by PlatformXe.
- `version` (Number) -- Current template version.

## Import

Templates can be imported using their ID:

```bash
terraform import platformxe_templates_template.welcome tmpl_abc123
```
