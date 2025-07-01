package js

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJS(t *testing.T) {
	tt := []struct {
		name     string
		ops      []Operation
		expected string
	}{
		{
			name:     "empty operations",
			ops:      []Operation{},
			expected: "[]",
		},
		{
			name:     "single operation",
			ops:      []Operation{{"push", map[string]string{"event": "test"}}},
			expected: `[[&#34;push&#34;,{&#34;event&#34;:&#34;test&#34;}]]`,
		},
		{
			name: "multiple operations",
			ops: []Operation{
				{"show", nil},
				{"hide", nil},
			},
			expected: `[[&#34;show&#34;,null],[&#34;hide&#34;,null]]`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := JS(tc.ops...)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestPush(t *testing.T) {
	tt := []struct {
		name     string
		event    string
		args     *PushArgs
		expected Operation
	}{
		{
			name:  "basic push",
			event: "test_event",
			args:  &PushArgs{},
			expected: Operation{"push", &pushArgs{
				Event:    "test_event",
				PushArgs: &PushArgs{},
			}},
		},
		{
			name:  "push with args",
			event: "submit",
			args: &PushArgs{
				Target:  "#form",
				Loading: ".spinner",
				Value:   map[string]any{"field": "value"},
			},
			expected: Operation{"push", &pushArgs{
				Event: "submit",
				PushArgs: &PushArgs{
					Target:  "#form",
					Loading: ".spinner",
					Value:   map[string]any{"field": "value"},
				},
			}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := Push(tc.event, tc.args)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestDispatch(t *testing.T) {
	tt := []struct {
		name     string
		event    string
		args     *DispatchArgs
		expected Operation
	}{
		{
			name:  "basic dispatch",
			event: "click",
			args:  &DispatchArgs{},
			expected: Operation{"dispatch", &dispatchArgs{
				Event: "click",
				DispatchArgs: &DispatchArgs{
					Bubbles: boolPtr(true),
				},
			}},
		},
		{
			name:  "dispatch with args",
			event: "custom",
			args: &DispatchArgs{
				To:     "#element",
				Detail: map[string]any{"key": "value"},
			},
			expected: Operation{"dispatch", &dispatchArgs{
				Event: "custom",
				DispatchArgs: &DispatchArgs{
					To:      "#element",
					Detail:  map[string]any{"key": "value"},
					Bubbles: boolPtr(true),
				},
			}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := Dispatch(tc.event, tc.args)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestShow(t *testing.T) {
	tt := []struct {
		name     string
		args     *ShowArgs
		expected Operation
	}{
		{
			name: "basic show",
			args: &ShowArgs{},
			expected: Operation{"show", &showArgs{
				ShowArgs: &ShowArgs{
					Time:     intPtr(defaultTransitionTime),
					Blocking: boolPtr(true),
				},
				Transition: [3][]string{{}, {}, {}},
			}},
		},
		{
			name: "show with args",
			args: &ShowArgs{
				To:      "#modal",
				Time:    intPtr(500),
				Display: "flex",
			},
			expected: Operation{"show", &showArgs{
				ShowArgs: &ShowArgs{
					To:       "#modal",
					Time:     intPtr(500),
					Display:  "flex",
					Blocking: boolPtr(true),
				},
				Transition: [3][]string{{}, {}, {}},
			}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := Show(tc.args)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestHide(t *testing.T) {
	tt := []struct {
		name     string
		args     *HideArgs
		expected Operation
	}{
		{
			name: "basic hide",
			args: &HideArgs{},
			expected: Operation{"hide", &hideArgs{
				HideArgs: &HideArgs{
					Time:     intPtr(defaultTransitionTime),
					Blocking: boolPtr(true),
				},
				Transition: [3][]string{{}, {}, {}},
			}},
		},
		{
			name: "hide with args",
			args: &HideArgs{
				To:   "#modal",
				Time: intPtr(200),
			},
			expected: Operation{"hide", &hideArgs{
				HideArgs: &HideArgs{
					To:       "#modal",
					Time:     intPtr(200),
					Blocking: boolPtr(true),
				},
				Transition: [3][]string{{}, {}, {}},
			}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := Hide(tc.args)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFocus(t *testing.T) {
	tt := []struct {
		name     string
		args     *FocusArgs
		expected Operation
	}{
		{
			name:     "basic focus",
			args:     &FocusArgs{},
			expected: Operation{"focus", &FocusArgs{}},
		},
		{
			name:     "focus with target",
			args:     &FocusArgs{To: "#input"},
			expected: Operation{"focus", &FocusArgs{To: "#input"}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := Focus(tc.args)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFocusFirst(t *testing.T) {
	tt := []struct {
		name     string
		args     *FocusFirstArgs
		expected Operation
	}{
		{
			name:     "basic focus first",
			args:     &FocusFirstArgs{},
			expected: Operation{"focus_first", &FocusFirstArgs{}},
		},
		{
			name:     "focus first with target",
			args:     &FocusFirstArgs{To: ".form"},
			expected: Operation{"focus_first", &FocusFirstArgs{To: ".form"}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := FocusFirst(tc.args)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestPushFocus(t *testing.T) {
	tt := []struct {
		name     string
		args     *PushFocusArgs
		expected Operation
	}{
		{
			name:     "basic push focus",
			args:     &PushFocusArgs{},
			expected: Operation{"push_focus", &PushFocusArgs{}},
		},
		{
			name:     "push focus with target",
			args:     &PushFocusArgs{To: "#element"},
			expected: Operation{"push_focus", &PushFocusArgs{To: "#element"}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := PushFocus(tc.args)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestPopFocus(t *testing.T) {
	expected := Operation{"pop_focus", []any{}}
	result := PopFocus()
	assert.Equal(t, expected, result)
}

func TestNavigate(t *testing.T) {
	tt := []struct {
		name     string
		href     string
		args     *NavigateArgs
		expected Operation
	}{
		{
			name: "basic navigate",
			href: "/home",
			args: &NavigateArgs{},
			expected: Operation{"navigate", &navigateArgs{
				Href:         "/home",
				NavigateArgs: &NavigateArgs{},
			}},
		},
		{
			name: "navigate with replace",
			href: "/profile",
			args: &NavigateArgs{Replace: true},
			expected: Operation{"navigate", &navigateArgs{
				Href:         "/profile",
				NavigateArgs: &NavigateArgs{Replace: true},
			}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := Navigate(tc.href, tc.args)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestPatch(t *testing.T) {
	tt := []struct {
		name     string
		href     string
		args     *PatchArgs
		expected Operation
	}{
		{
			name: "basic patch",
			href: "/users",
			args: &PatchArgs{},
			expected: Operation{"patch", &patchArgs{
				Href:      "/users",
				PatchArgs: &PatchArgs{},
			}},
		},
		{
			name: "patch with replace",
			href: "/search",
			args: &PatchArgs{Replace: true},
			expected: Operation{"patch", &patchArgs{
				Href:      "/search",
				PatchArgs: &PatchArgs{Replace: true},
			}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := Patch(tc.href, tc.args)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestExec(t *testing.T) {
	tt := []struct {
		name     string
		attr     string
		args     *ExecArgs
		expected Operation
	}{
		{
			name: "basic exec",
			attr: "data-action",
			args: &ExecArgs{},
			expected: Operation{"exec", &execArgs{
				Attr:     "data-action",
				ExecArgs: &ExecArgs{},
			}},
		},
		{
			name: "exec with target",
			attr: "onclick",
			args: &ExecArgs{To: "#button"},
			expected: Operation{"exec", &execArgs{
				Attr:     "onclick",
				ExecArgs: &ExecArgs{To: "#button"},
			}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := Exec(tc.attr, tc.args)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestToggle(t *testing.T) {
	tt := []struct {
		name     string
		event    string
		args     *ToggleArgs
		expected Operation
	}{
		{
			name:  "basic toggle",
			event: "click",
			args:  &ToggleArgs{},
			expected: Operation{"toggle", &toggleArgs{
				Event: "click",
				ToggleArgs: &ToggleArgs{
					Time:     intPtr(defaultTransitionTime),
					Display:  "block",
					Blocking: boolPtr(true),
				},
				In:  [3][]string{{}, {}, {}},
				Out: [3][]string{{}, {}, {}},
			}},
		},
		{
			name:  "toggle with args",
			event: "hover",
			args: &ToggleArgs{
				To:      "#element",
				Display: "flex",
			},
			expected: Operation{"toggle", &toggleArgs{
				Event: "hover",
				ToggleArgs: &ToggleArgs{
					To:       "#element",
					Time:     intPtr(defaultTransitionTime),
					Display:  "flex",
					Blocking: boolPtr(true),
				},
				In:  [3][]string{{}, {}, {}},
				Out: [3][]string{{}, {}, {}},
			}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := Toggle(tc.event, tc.args)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestAddClass(t *testing.T) {
	tt := []struct {
		name     string
		names    string
		args     *AddClassArgs
		expected Operation
	}{
		{
			name:  "basic add class",
			names: "active",
			args:  &AddClassArgs{},
			expected: Operation{"add_class", &addClassArgs{
				Names: []string{"active"},
				AddClassArgs: &AddClassArgs{
					Time:     intPtr(defaultTransitionTime),
					Blocking: boolPtr(true),
				},
				Transition: [3][]string{{}, {}, {}},
			}},
		},
		{
			name:  "add multiple classes",
			names: "active visible",
			args: &AddClassArgs{
				To: "#element",
			},
			expected: Operation{"add_class", &addClassArgs{
				Names: []string{"active", "visible"},
				AddClassArgs: &AddClassArgs{
					To:       "#element",
					Time:     intPtr(defaultTransitionTime),
					Blocking: boolPtr(true),
				},
				Transition: [3][]string{{}, {}, {}},
			}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := AddClass(tc.names, tc.args)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestToggleClass(t *testing.T) {
	tt := []struct {
		name     string
		names    string
		args     *ToggleClassArgs
		expected Operation
	}{
		{
			name:  "basic toggle class",
			names: "hidden",
			args:  &ToggleClassArgs{},
			expected: Operation{"toggle_class", &toggleClassArgs{
				Names: []string{"hidden"},
				ToggleClassArgs: &ToggleClassArgs{
					Time:     intPtr(defaultTransitionTime),
					Blocking: boolPtr(true),
				},
				Transition: [3][]string{{}, {}, {}},
			}},
		},
		{
			name:  "toggle multiple classes",
			names: "active visible",
			args: &ToggleClassArgs{
				To: ".target",
			},
			expected: Operation{"toggle_class", &toggleClassArgs{
				Names: []string{"active", "visible"},
				ToggleClassArgs: &ToggleClassArgs{
					To:       ".target",
					Time:     intPtr(defaultTransitionTime),
					Blocking: boolPtr(true),
				},
				Transition: [3][]string{{}, {}, {}},
			}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := ToggleClass(tc.names, tc.args)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRemoveClass(t *testing.T) {
	tt := []struct {
		name     string
		names    string
		args     *RemoveClassArgs
		expected Operation
	}{
		{
			name:  "basic remove class",
			names: "error",
			args:  &RemoveClassArgs{},
			expected: Operation{"remove_class", &removeClassArgs{
				Names: []string{"error"},
				RemoveClassArgs: &RemoveClassArgs{
					Time:     intPtr(defaultTransitionTime),
					Blocking: boolPtr(true),
				},
				Transition: [3][]string{{}, {}, {}},
			}},
		},
		{
			name:  "remove multiple classes",
			names: "error warning",
			args: &RemoveClassArgs{
				To: "#form",
			},
			expected: Operation{"remove_class", &removeClassArgs{
				Names: []string{"error", "warning"},
				RemoveClassArgs: &RemoveClassArgs{
					To:       "#form",
					Time:     intPtr(defaultTransitionTime),
					Blocking: boolPtr(true),
				},
				Transition: [3][]string{{}, {}, {}},
			}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := RemoveClass(tc.names, tc.args)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestTransition(t *testing.T) {
	tt := []struct {
		name       string
		transition [3]string
		args       *TransitionArgs
		expected   Operation
	}{
		{
			name:       "basic transition",
			transition: [3]string{"fade-in", "fade-out", ""},
			args:       &TransitionArgs{},
			expected: Operation{"transition", &transitionArgs{
				Transition: [3][]string{
					{"fade-in"},
					{"fade-out"},
					{},
				},
				TransitionArgs: &TransitionArgs{
					Time:     intPtr(defaultTransitionTime),
					Blocking: boolPtr(true),
				},
			}},
		},
		{
			name:       "transition with args",
			transition: [3]string{"slide-left", "slide-right", "slide-end"},
			args: &TransitionArgs{
				To: "#panel",
			},
			expected: Operation{"transition", &transitionArgs{
				Transition: [3][]string{
					{"slide-left"},
					{"slide-right"},
					{"slide-end"},
				},
				TransitionArgs: &TransitionArgs{
					To:       "#panel",
					Time:     intPtr(defaultTransitionTime),
					Blocking: boolPtr(true),
				},
			}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := Transition(tc.transition, tc.args)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestSetAttr(t *testing.T) {
	tt := []struct {
		name     string
		key      string
		value    any
		args     *SetAttrArgs
		expected Operation
	}{
		{
			name:  "basic set attr",
			key:   "disabled",
			value: true,
			args:  &SetAttrArgs{},
			expected: Operation{"set_attr", &setAttrArgs{
				Attr:        [2]any{"disabled", true},
				SetAttrArgs: &SetAttrArgs{},
			}},
		},
		{
			name:  "set attr with target",
			key:   "class",
			value: "active",
			args:  &SetAttrArgs{To: "#button"},
			expected: Operation{"set_attr", &setAttrArgs{
				Attr:        [2]any{"class", "active"},
				SetAttrArgs: &SetAttrArgs{To: "#button"},
			}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := SetAttr(tc.key, tc.value, tc.args)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRemoveAttr(t *testing.T) {
	tt := []struct {
		name     string
		attr     string
		args     *RemoveAttrArgs
		expected Operation
	}{
		{
			name: "basic remove attr",
			attr: "disabled",
			args: &RemoveAttrArgs{},
			expected: Operation{"remove_attr", &removeAttrArgs{
				Attr:           "disabled",
				RemoveAttrArgs: &RemoveAttrArgs{},
			}},
		},
		{
			name: "remove attr with target",
			attr: "class",
			args: &RemoveAttrArgs{To: "#input"},
			expected: Operation{"remove_attr", &removeAttrArgs{
				Attr:           "class",
				RemoveAttrArgs: &RemoveAttrArgs{To: "#input"},
			}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := RemoveAttr(tc.attr, tc.args)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestToggleAttr(t *testing.T) {
	tt := []struct {
		name     string
		attr     string
		val      string
		args     *ToggleAttrArgs
		expected Operation
	}{
		{
			name: "basic toggle attr",
			attr: "aria-expanded",
			val:  "true",
			args: &ToggleAttrArgs{},
			expected: Operation{"toggle_attr", &toggleAttrArgs{
				Attr:           [2]any{"aria-expanded", "true"},
				ToggleAttrArgs: &ToggleAttrArgs{},
			}},
		},
		{
			name: "toggle attr with target",
			attr: "checked",
			val:  "checked",
			args: &ToggleAttrArgs{To: "#checkbox"},
			expected: Operation{"toggle_attr", &toggleAttrArgs{
				Attr:           [2]any{"checked", "checked"},
				ToggleAttrArgs: &ToggleAttrArgs{To: "#checkbox"},
			}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := ToggleAttr(tc.attr, tc.val, tc.args)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestToggleAttrs(t *testing.T) {
	tt := []struct {
		name     string
		attr     string
		val1     string
		val2     string
		args     *ToggleAttrsArgs
		expected Operation
	}{
		{
			name: "basic toggle attrs",
			attr: "data-state",
			val1: "open",
			val2: "closed",
			args: &ToggleAttrsArgs{},
			expected: Operation{"toggle_attr", &toggleAttrsArgs{
				Attrs:           [3]any{"data-state", "open", "closed"},
				ToggleAttrsArgs: &ToggleAttrsArgs{},
			}},
		},
		{
			name: "toggle attrs with target",
			attr: "aria-pressed",
			val1: "true",
			val2: "false",
			args: &ToggleAttrsArgs{To: "#toggle-btn"},
			expected: Operation{"toggle_attr", &toggleAttrsArgs{
				Attrs:           [3]any{"aria-pressed", "true", "false"},
				ToggleAttrsArgs: &ToggleAttrsArgs{To: "#toggle-btn"},
			}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := ToggleAttrs(tc.attr, tc.val1, tc.val2, tc.args)
			assert.Equal(t, tc.expected, result)
		})
	}
}
