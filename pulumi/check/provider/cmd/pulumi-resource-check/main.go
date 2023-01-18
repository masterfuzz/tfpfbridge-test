package main

import (
	_ "embed"

	bridge "github.com/pulumi/pulumi-terraform-bridge/pkg/tfpfbridge"
	check "github.com/tetrateio/tetrate/cloud/providers/check/provider"
)

//go:embed schema.json
var pulumiSchema []byte

//go:embed renames.json
var pulumiRenames []byte

func main() {
	bridge.Main("check", "0.0.1", check.MyProvider(), pulumiSchema, pulumiRenames)
}
