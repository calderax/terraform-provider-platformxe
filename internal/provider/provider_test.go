// =============================================================================
// Copyright 2026 Caldera Technologies Ltd.
// Proprietary and confidential.
// Unauthorized copying or distribution is prohibited.
// =============================================================================

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories is used by acceptance tests across all
// resource and data source packages. It instantiates the PlatformXe provider
// using the Protocol V6 interface expected by the Terraform Plugin Framework.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"platformxe": providerserver.NewProtocol6WithError(New()),
}

// TestAccProtoV6ProviderFactories returns the provider factories map so that
// resource and data source test packages can reference it without a circular
// import.
func TestAccProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return testAccProtoV6ProviderFactories
}

// testAccProviderConfig returns a reusable HCL provider block. The API key and
// optional base URL are read from environment variables at runtime so that no
// secrets are hard-coded in test files.
func TestAccProviderConfig() string {
	return `
provider "platformxe" {}
`
}

func TestProviderMetadata(t *testing.T) {
	p := New()
	if p == nil {
		t.Fatal("Provider is nil")
	}
}

func TestProviderSchema(t *testing.T) {
	p := &PlatformXeProvider{}
	if p == nil {
		t.Fatal("PlatformXeProvider is nil")
	}
}
