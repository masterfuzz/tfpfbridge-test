package modifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ planmodifier.String = defaultStringModifier{}

func DefaultString(def string) defaultStringModifier {
	return defaultStringModifier{Default: types.StringValue(def)}
}

func NullableString() defaultStringModifier {
	return defaultStringModifier{Default: types.StringNull()}
}

type defaultStringModifier struct {
	Default types.String
}

// PlanModifyString implements planmodifier.String
func (m defaultStringModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if !req.ConfigValue.IsNull() {
		return
	}

	resp.PlanValue = m.Default
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
