package main

import (
	"strings"
	"testing"
)

func TestRenderPage_TitleAndDescription(t *testing.T) {
	structs := map[string]StructDef{
		"Config": {Name: "Config", Fields: []FieldDef{
			{YAMLKey: "region", GoType: "string", Required: true},
		}},
	}
	out := RenderPage("AWS Config", "Intro text.", structs, "Config")

	if !strings.Contains(out, "# AWS Config") {
		t.Error("missing page title")
	}
	if !strings.Contains(out, "Intro text.") {
		t.Error("missing description")
	}
}

func TestRenderPage_RequiredMarker(t *testing.T) {
	structs := map[string]StructDef{
		"Config": {Name: "Config", Fields: []FieldDef{
			{YAMLKey: "region", GoType: "string", Required: true},
			{YAMLKey: "tags", GoType: "map[string]string", Required: false},
		}},
	}
	out := RenderPage("Test", "", structs, "Config")

	if !strings.Contains(out, "**Yes**") {
		t.Error("required field should show **Yes**")
	}
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.Contains(line, "`tags`") && strings.Contains(line, "**Yes**") {
			t.Error("optional 'tags' field should not show **Yes**")
		}
	}
}

func TestRenderPage_NestedStruct(t *testing.T) {
	structs := map[string]StructDef{
		"Config": {Name: "Config", Fields: []FieldDef{
			{YAMLKey: "node_groups", GoType: "map[string]NodeGroup", Required: true},
		}},
		"NodeGroup": {Name: "NodeGroup", Fields: []FieldDef{
			{YAMLKey: "instance", GoType: "string", Required: true},
		}},
	}
	out := RenderPage("Test", "", structs, "Config")

	if !strings.Contains(out, "## NodeGroup") {
		t.Error("NodeGroup subsection missing")
	}
	if !strings.Contains(out, "| `instance`") {
		t.Error("NodeGroup.instance field missing from output")
	}
}

func TestRenderPage_InlineFieldSkipped(t *testing.T) {
	structs := map[string]StructDef{
		"Config": {Name: "Config", Fields: []FieldDef{
			{YAMLKey: "region", GoType: "string", Required: true},
			{Inline: true},
		}},
	}
	out := RenderPage("Test", "", structs, "Config")

	if !strings.Contains(out, "| `region`") {
		t.Error("region field missing")
	}
}

func TestRenderPage_NoCircularRecursion(t *testing.T) {
	structs := map[string]StructDef{
		"Config": {Name: "Config", Fields: []FieldDef{
			{YAMLKey: "sub", GoType: "Sub"},
		}},
		"Sub": {Name: "Sub", Fields: []FieldDef{
			{YAMLKey: "name", GoType: "string"},
		}},
	}
	// Must not infinite-loop or panic.
	out := RenderPage("Test", "", structs, "Config")
	if out == "" {
		t.Error("expected non-empty output")
	}
}

func TestRenderPage_TypeLinkified(t *testing.T) {
	structs := map[string]StructDef{
		"Config": {Name: "Config", Fields: []FieldDef{
			{YAMLKey: "node_groups", GoType: "map[string]NodeGroup", Required: true},
		}},
		"NodeGroup": {Name: "NodeGroup", Fields: []FieldDef{
			{YAMLKey: "instance", GoType: "string", Required: true},
		}},
	}
	out := RenderPage("Test", "", structs, "Config")

	if !strings.Contains(out, "[NodeGroup]") {
		t.Error("NodeGroup type should be linkified in the table")
	}
}
