package pkg

import (
	framework "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/tetrateio/tetrate/cloud/providers/terraform/check/internal/provider"
)

func NewProvider() framework.Provider {
	return provider.New("0.0.1")()
}
