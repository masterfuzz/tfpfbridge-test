package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &HttpHealthResource{}
var _ resource.ResourceWithImportState = &HttpHealthResource{}

func NewHttpHealthResource() resource.Resource {
	return &HttpHealthResource{}
}

type HttpHealthResource struct{}

// Schema implements resource.Resource
func (*HttpHealthResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "HTTPS Healthcheck",

		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				MarkdownDescription: "URL",
				Required:            true,
			},
			"retries": schema.Int64Attribute{
				MarkdownDescription: "Retries",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.Int64{modifiers.DefaultInt64(5)},
			},
			"method": schema.StringAttribute{
				MarkdownDescription: "Method",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{modifiers.DefaultString("GET")},
			},
			"timeout": schema.Int64Attribute{
				MarkdownDescription: "Timeout",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.Int64{modifiers.DefaultInt64(5000)},
			},
			"interval": schema.Int64Attribute{
				MarkdownDescription: "Interval",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.Int64{modifiers.DefaultInt64(200)},
			},
			"status_code": schema.StringAttribute{
				MarkdownDescription: "Status Code",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{modifiers.DefaultString("200")},
			},
			"consecutive_successes": schema.Int64Attribute{
				MarkdownDescription: "Consecutive successes required",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.Int64{modifiers.DefaultInt64(200)},
			},
			"headers": schema.MapAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "HTTP Headers",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

type HttpHealthResourceModel struct {
	URL                  types.String `tfsdk:"url"`
	Id                   types.String `tfsdk:"id"`
	Retries              types.Int64  `tfsdk:"retries"`
	Method               types.String `tfsdk:"method"`
	Timeout              types.Int64  `tfsdk:"timeout"`
	Interval             types.Int64  `tfsdk:"interval"`
	StatusCode           types.String `tfsdk:"status_code"`
	ConsecutiveSuccesses types.Int64  `tfsdk:"consecutive_successes"`
	Headers              types.Map    `tfsdk:"headers"`
}

func (r *HttpHealthResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_http_health"
}

// func (r *HttpHealthResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
// 	// Prevent panic if the provider has not been configured.
// 	if req.ProviderData == nil {
// 		return
// 	}

// 	client, ok := req.ProviderData.(*http.Client)

// 	if !ok {
// 		resp.Diagnostics.AddError(
// 			"Unexpected Resource Configure Type",
// 			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
// 		)

// 		return
// 	}

// 	r.client = client
// }

// func inStringRange(v int64, r string) (bool) {
// 	for _, interval := range strings.Split(r, ",") {
// 		if strings.Contains(interval, "-") {
// 			lr := strings.Split(interval, "-")
// 			if len(lr) != 2 {
// 				// error
// 			}
// 			if v > lr[0]
// 		}

// 	}
// }

func (r *HttpHealthResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *HttpHealthResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint, err := url.Parse(data.URL.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to parse url %q, got error %s", data.URL.ValueString(), err))
		return
	}

	var checkCode func(int) bool
	if data.StatusCode.IsNull() {
		checkCode = func(c int) bool { return c < 400 }
	} else {
		v, err := strconv.Atoi(data.StatusCode.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error", fmt.Sprintf("Unable to parse status code pattern %s", err))
		}
		checkCode = func(c int) bool { return c == v }
	}

	// normalize headers
	headers := make(map[string][]string)
	if !data.Headers.IsNull() {
		tmp := make(map[string]string)
		resp.Diagnostics.Append(data.Headers.ElementsAs(ctx, &tmp, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		for k, v := range data.Headers.Elements() {
			headers[k] = []string{v.String()}
		}
	}

	window := helpers.RetryWindow{
		MaxTries:             int(data.Retries.ValueInt64()),
		Timeout:              time.Duration(data.Timeout.ValueInt64()) * time.Millisecond,
		Interval:             time.Duration(data.Interval.ValueInt64()) * time.Millisecond,
		ConsecutiveSuccesses: int(data.ConsecutiveSuccesses.ValueInt64()),
	}

	client := http.DefaultClient
	result := window.Do(func() bool {
		httpResponse, err := client.Do(&http.Request{
			URL:    endpoint,
			Method: data.Method.ValueString(),
			Header: headers,
		})
		if err != nil {
			resp.Diagnostics.AddWarning("Error connecting to healthcheck endpoint", fmt.Sprintf("%s", err))
			return false
		}

		return checkCode(httpResponse.StatusCode)
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

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HttpHealthResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *HttpHealthResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HttpHealthResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *HttpHealthResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HttpHealthResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *HttpHealthResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r *HttpHealthResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
