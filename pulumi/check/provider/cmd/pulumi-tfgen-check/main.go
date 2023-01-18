package main

import (
	"github.com/pulumi/pulumi-terraform-bridge/pkg/tfpfbridge/tfgen"
	check "github.com/tetrateio/tetrate/cloud/providers/check/provider"
)

func main() {
	tfgen.Main("check", "0.0.1", check.MyProvider())
}
