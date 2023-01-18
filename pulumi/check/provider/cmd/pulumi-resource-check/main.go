package main

import (
	_ "embed"

	check "github.com/masterfuzz/tfpfbridge-test/pulumi/check/provider"
	bridge "github.com/pulumi/pulumi-terraform-bridge/pkg/tfpfbridge"
)

//go:embed schema.json
var pulumiSchema []byte

//go:embed renames.json
var pulumiRenames []byte

func main() {
	bridge.Main("check", "0.0.1", check.MyProvider(), pulumiSchema, pulumiRenames)
}
