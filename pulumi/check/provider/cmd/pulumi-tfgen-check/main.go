package main

import (
	check "github.com/masterfuzz/tfpfbridge-test/pulumi/check/provider"
	"github.com/pulumi/pulumi-terraform-bridge/pkg/tfpfbridge/tfgen"
)

func main() {
	tfgen.Main("check", "0.0.1", check.MyProvider())
}
