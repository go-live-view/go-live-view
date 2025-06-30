package dynamic

import (
	"github.com/go-live-view/go-live-view/rend"
)

// NotNil renders content only if the condition is not nil.
// This is useful for pointer safety - it will render content for any non-nil value,
// including falsy values like false, 0, or empty strings.
// For truthiness checking, use Presence() instead.
func NotNil(cond any, cb func() rend.Node) rend.Node {
	if cond == nil {
		return nil
	}

	return Wrap(cb())
}
