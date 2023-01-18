package check

import (
	"unicode"

	framework "github.com/hashicorp/terraform-plugin-framework/provider"
	// "github.com/pulumi/pulumi-terraform-bridge/pkg/tfbridge"
	provider "github.com/masterfuzz/tfpfbridge-test/terraform/check/pkg"
	"github.com/pulumi/pulumi-terraform-bridge/pkg/tfpfbridge/info"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

const checkPkg = "check"
const checkMod = "index"

func getProvider() framework.Provider {
	return provider.NewProvider()
}

func checkMember(mod string, mem string) tokens.ModuleMember {
	return tokens.ModuleMember(checkPkg + ":" + mod + ":" + mem)
}

func checkType(mod string, typ string) tokens.Type {
	return tokens.Type(checkMember(mod, typ))
}

func checkResourceTok(mod string, res string) tokens.Type {
	fn := string(unicode.ToLower(rune(res[0]))) + res[1:]
	return checkType(mod+"/"+fn, res)
}

func MyProvider() info.ProviderInfo {
	return info.ProviderInfo{
		P:                 getProvider,
		Name:              "check",
		GitHubOrg:         "tetrateio",
		TFProviderVersion: "dev",
		Version:           "dev",
		Resources: map[string]*info.ResourceInfo{
			"check_http_health":   {Tok: checkResourceTok(checkMod, "HttpHealth")},
			"check_local_command": {Tok: checkResourceTok(checkMod, "LocalCommand")},
		},
		JavaScript: &info.JavaScriptInfo{
			Dependencies: map[string]string{
				"@pulumi/pulumi": "^3.0.0",
			},
			DevDependencies: map[string]string{
				"@types/node": "^10.0.0", // so we can access strongly typed node definitions.
			},
		},
	}
}
