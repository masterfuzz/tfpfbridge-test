package pkg

import (
	framework "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/masterfuzz/tfpfbridge-test/terraform/check/internal/provider"
)

func NewProvider() framework.Provider {
	return provider.New("0.0.1")()
}
