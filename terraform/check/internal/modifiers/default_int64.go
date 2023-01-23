package modifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ planmodifier.Int64 = defaultInt64Modifier{}

func DefaultInt64(def int64) defaultInt64Modifier {
	return defaultInt64Modifier{Default: types.Int64Value(def)}
}

type defaultInt64Modifier struct {
	Default types.Int64
}

// PlanModifyInt64 implements planmodifier.Int64
func (m defaultInt64Modifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if !req.ConfigValue.IsNull() {
		return
	}

	resp.PlanValue = m.Default
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
