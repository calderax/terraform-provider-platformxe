module github.com/calderax/terraform-provider-platformxe

go 1.22

require (
	github.com/calderax/platformxe-go v0.0.0
	github.com/hashicorp/terraform-plugin-framework v1.12.0
	github.com/hashicorp/terraform-plugin-go v0.25.0
	github.com/hashicorp/terraform-plugin-testing v1.10.0
)

replace github.com/calderax/platformxe-go => ../sdk-go
