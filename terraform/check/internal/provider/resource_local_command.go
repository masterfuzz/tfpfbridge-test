package provider

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/masterfuzz/tfpfbridge-test/terraform/check/internal/helpers"
	"github.com/masterfuzz/tfpfbridge-test/terraform/check/internal/modifiers"
)

var _ resource.Resource = &LocalCommandResource{}
var _ resource.ResourceWithImportState = &LocalCommandResource{}

type LocalCommandResource struct{}

// Schema implements resource.Resource
func (*LocalCommandResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Local Command",

		Attributes: map[string]schema.Attribute{
			"command": schema.StringAttribute{
				MarkdownDescription: "Command",
				Required:            true,
			},
			"timeout": schema.Int64Attribute{
				MarkdownDescription: "Timeout",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.Int64{modifiers.DefaultInt64(500)},
			},
			"retries": schema.Int64Attribute{
				MarkdownDescription: "Retries",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.Int64{modifiers.DefaultInt64(5)},
			},
			"interval": schema.Int64Attribute{
				MarkdownDescription: "Interval",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.Int64{modifiers.DefaultInt64(200)},
			},
			"consecutive_successes": schema.Int64Attribute{
				MarkdownDescription: "Consecutive Successes",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.Int64{modifiers.DefaultInt64(200)},
			},
			"working_directory": schema.StringAttribute{
				MarkdownDescription: "Working Directory",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{modifiers.DefaultString(".")},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		}}
}

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
			fmt.Printf("error running command %s", err.Error())
			resp.Diagnostics.AddWarning("Error running command", fmt.Sprintf("%s", err))
			tflog.Error(ctx, fmt.Sprintf("Error running command %s", err))
			return false
		}
		tflog.Debug(ctx, fmt.Sprintf("Command stdout: %s", stdout.String()))
		tflog.Debug(ctx, fmt.Sprintf("Command stdout: %s", stderr.String()))
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
