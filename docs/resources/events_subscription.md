---
page_title: "platformxe_events_subscription Resource"
description: "Manages a PlatformXe event subscription."
---

# platformxe_events_subscription

Manages an event subscription for receiving real-time event notifications.

## Example Usage

```hcl
resource "platformxe_events_subscription" "email_events" {
  event = "email.*"
  url   = "https://example.com/events"
}
```

## Argument Reference

- `event` (String, Required) -- Event pattern to subscribe to. Supports wildcards (`*`).
- `url` (String, Required) -- Delivery URL for event notifications.

## Attribute Reference

- `id` (String) -- Subscription ID, assigned by PlatformXe.

## Import

Event subscriptions can be imported using their ID:

```bash
terraform import platformxe_events_subscription.email_events sub_abc123
```
