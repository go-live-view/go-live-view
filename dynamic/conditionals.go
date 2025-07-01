package dynamic

import (
	"github.com/go-live-view/go-live-view/rend"
)

// If renders the then node only if condition is true
func If(condition bool, then rend.Node) rend.Node {
	if !condition || then == nil {
		return nil
	}
	return Wrap(then)
}

// IfFunc allows lazy evaluation
func IfFunc(condition bool, then func() rend.Node) rend.Node {
	if !condition || then == nil {
		return nil
	}
	node := then()
	if node == nil {
		return nil
	}
	return Wrap(node)
}

// IfElse renders then if condition is true, otherwise renders elseValue
func IfElse(condition bool, then, elseValue rend.Node) rend.Node {
	if condition {
		return If(true, then)
	}
	return If(true, elseValue)
}

// IfElseFunc allows lazy evaluation for both branches
func IfElseFunc(condition bool, then, elseValue func() rend.Node) rend.Node {
	if condition {
		return IfFunc(true, then)
	}
	return IfFunc(true, elseValue)
}

// Unless renders the then node only if condition is false (opposite of If)
func Unless(condition bool, then rend.Node) rend.Node {
	return If(!condition, then)
}

// UnlessFunc allows lazy evaluation
func UnlessFunc(condition bool, then func() rend.Node) rend.Node {
	return IfFunc(!condition, then)
}

// Switch provides switch-case conditional rendering
func Switch[T comparable](value T, cases ...SwitchCase[T]) rend.Node {
	for _, c := range cases {
		if c.Value == value {
			return If(true, c.Result)
		}
	}
	return nil
}

// SwitchCase represents a case in a switch statement
// Result must be a rend.Node
// (for lazy, user can pass IfFunc or similar)
type SwitchCase[T comparable] struct {
	Value  T
	Result rend.Node
}

// Case creates a switch case
func Case[T comparable](value T, result rend.Node) SwitchCase[T] {
	return SwitchCase[T]{Value: value, Result: result}
}

// SwitchDefault provides switch-case with default value
func SwitchDefault[T comparable](value T, defaultValue rend.Node, cases ...SwitchCase[T]) rend.Node {
	for _, c := range cases {
		if c.Value == value {
			return If(true, c.Result)
		}
	}
	return If(true, defaultValue)
}

// When provides multiple condition checking (like elsif chain)
type WhenCase struct {
	Condition bool
	Result    rend.Node
}

func When(conditions ...WhenCase) rend.Node {
	for _, c := range conditions {
		if c.Condition {
			return If(true, c.Result)
		}
	}
	return nil
}

// WhenIs creates a when case
func WhenIs(condition bool, result rend.Node) WhenCase {
	return WhenCase{Condition: condition, Result: result}
}

// WhenDefault provides When with a default fallback
func WhenDefault(defaultValue rend.Node, conditions ...WhenCase) rend.Node {
	for _, c := range conditions {
		if c.Condition {
			return If(true, c.Result)
		}
	}
	return If(true, defaultValue)
}

// Show/Hide helpers for common visibility patterns
func Show(condition bool, content rend.Node) rend.Node {
	return If(condition, content)
}

func Hide(condition bool, content rend.Node) rend.Node {
	return Unless(condition, content)
}
