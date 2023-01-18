package modifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ tfsdk.AttributePlanModifier = defaultInt64Modifier{}

func DefaultInt64(def int64) defaultInt64Modifier {
	return defaultInt64Modifier{Default: types.Int64Value(def)}
}

type defaultInt64Modifier struct {
	Default types.Int64
}

func (m defaultInt64Modifier) String() string {
	return m.Default.String()
}

func (m defaultInt64Modifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%s`", m)
}

func (m defaultInt64Modifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%s`", m)
}

func (m defaultInt64Modifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	if !req.AttributeConfig.IsNull() {
		return
	}

	resp.AttributePlan = m.Default
}
