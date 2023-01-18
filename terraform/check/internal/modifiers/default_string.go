package modifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ tfsdk.AttributePlanModifier = defaultStringModifier{}

func DefaultString(def string) defaultStringModifier {
	return defaultStringModifier{Default: types.StringValue(def)}
}

func NullableString() defaultStringModifier {
	return defaultStringModifier{Default: types.StringNull()}
}

type defaultStringModifier struct {
	Default types.String
}

func (m defaultStringModifier) String() string {
	return m.Default.String()
}

func (m defaultStringModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%s`", m)
}

func (m defaultStringModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%s`", m)
}

func (m defaultStringModifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	if !req.AttributeConfig.IsNull() {
		return
	}

	resp.AttributePlan = m.Default
}
