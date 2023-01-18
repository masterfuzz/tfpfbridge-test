package provider

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/tetrateio/tetrate/cloud/providers/terraform/check/internal/helpers"
	"github.com/tetrateio/tetrate/cloud/providers/terraform/check/internal/modifiers"
)

var _ resource.Resource = &LocalCommandResource{}
var _ resource.ResourceWithImportState = &LocalCommandResource{}

type LocalCommandResource struct{}

type LocalCommandResourceModel struct {
	Command              types.String `tfsdk:"command"`
	Timeout              types.Int64  `tfsdk:"timeout"`
	Retries              types.Int64  `tfsdk:"retries"`
	Interval             types.Int64  `tfsdk:"interval"`
	ConsecutiveSuccesses types.Int64  `tfsdk:"consecutive_successes"`
	WorkDir              types.String `tfsdk:"working_directory"`
	Id                   types.String `tfsdk:"id"`
}

// ImportState implements resource.ResourceWithImportState
func (*LocalCommandResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create implements resource.Resource
func (*LocalCommandResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *LocalCommandResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	window := helpers.RetryWindow{
		MaxTries:             int(data.Retries.ValueInt64()),
		Timeout:              time.Duration(data.Timeout.ValueInt64()) * time.Millisecond,
		Interval:             time.Duration(data.Interval.ValueInt64()) * time.Millisecond,
		ConsecutiveSuccesses: int(data.ConsecutiveSuccesses.ValueInt64()),
	}

	result := window.Do(func() bool {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd := exec.Command("sh", "-c", data.Command.ValueString())
		cmd.Dir = data.WorkDir.ValueString()
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			resp.Diagnostics.AddWarning("Error running command", fmt.Sprintf("%s", err))
			return false
		}
		tflog.Debug(ctx, stdout.String())
		tflog.Debug(ctx, stderr.String())
		return true
	})

	switch result {
	case helpers.Success:
		break
	case helpers.TimeoutExceeded:
		resp.Diagnostics.AddError("Timeout exceeded", fmt.Sprintf("Timeout of %d milliseconds exceeded", data.Timeout.ValueInt64()))
		return
	case helpers.RetriesExceeded:
		resp.Diagnostics.AddError("Retries exceeded", fmt.Sprintf("All %d attempts failed", data.Retries.ValueInt64()))
		return
	}

	data.Id = types.StringValue("example-id")
	tflog.Trace(ctx, "created resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Delete implements resource.Resource
func (*LocalCommandResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *LocalCommandResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Metadata implements resource.Resource
func (*LocalCommandResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_command"
}

// Read implements resource.Resource
func (*LocalCommandResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *LocalCommandResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Schema implements resource.Resource
func (*LocalCommandResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Local Command",

		Attributes: map[string]tfsdk.Attribute{
			"command": {
				Type:                types.StringType,
				MarkdownDescription: "Command",
				Required:            true,
			},
			"timeout": {
				MarkdownDescription: "Timeout",
				Type:                types.Int64Type,
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []tfsdk.AttributePlanModifier{modifiers.DefaultInt64(5000)},
			},
			"retries": {
				MarkdownDescription: "Retries",
				Type:                types.Int64Type,
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []tfsdk.AttributePlanModifier{modifiers.DefaultInt64(5)},
			},
			"interval": {
				MarkdownDescription: "Interval",
				Type:                types.Int64Type,
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []tfsdk.AttributePlanModifier{modifiers.DefaultInt64(200)},
			},
			"consecutive_successes": {
				MarkdownDescription: "Consecutive Successes",
				Type:                types.Int64Type,
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []tfsdk.AttributePlanModifier{modifiers.DefaultInt64(1)},
			},
			"working_directory": {
				MarkdownDescription: "Working Directory",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []tfsdk.AttributePlanModifier{modifiers.DefaultString("")},
			},
			"id": {
				Type:                types.StringType,
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers:       tfsdk.AttributePlanModifiers{resource.UseStateForUnknown()},
				// tfsdk.AttributePlanModifier{
				// 	stringplanmodifier.UseStateForUnknown(),
				// },
			},
		}}, diag.Diagnostics{}
}

// Update implements resource.Resource
func (*LocalCommandResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *LocalCommandResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func NewLocalCommandResource() resource.Resource {
	return &LocalCommandResource{}
}