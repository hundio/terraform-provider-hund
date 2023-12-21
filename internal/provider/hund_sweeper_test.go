package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

// sharedClientForDomain returns a common provider client configured for the specified domain
func sharedClientForDomain(domain string) (*hundApiV1.Client, error) {
	domainEnv := os.Getenv("HUND_DOMAIN")

	if domainEnv != "" {
		domain = domainEnv
	}

	key := os.Getenv("HUND_KEY")

	return newProviderClient("test", domain, key)
}
