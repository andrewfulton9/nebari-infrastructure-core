package main

import (
	"fmt"
	"strings"
)

// RenderPage produces a complete markdown page for one provider.
// structs is the merged map of all StructDef values for this provider.
// rootName is the entry-point struct (e.g. "Config" or "NebariConfig").
func RenderPage(title, description string, structs map[string]StructDef, rootName string) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "# %s\n\n", title)
	if description != "" {
		fmt.Fprintf(&sb, "%s\n\n", description)
	}

	visited := make(map[string]bool)
	renderStruct(&sb, rootName, structs, visited, 2)

	return sb.String()
}

func renderStruct(sb *strings.Builder, name string, structs map[string]StructDef, visited map[string]bool, headingLevel int) {
	sd, ok := structs[name]
	if !ok || visited[name] {
		return
	}
	visited[name] = true

	if sd.Doc != "" {
		fmt.Fprintf(sb, "%s\n\n", sd.Doc)
	}

	sb.WriteString("| YAML Key | Type | Required | Description |\n")
	sb.WriteString("|----------|------|----------|-------------|\n")

	var nested []string
	seen := make(map[string]bool)

	for _, f := range sd.Fields {
		if f.Inline || f.YAMLKey == "" {
			continue
		}
		req := ""
		if f.Required {
			req = "**Yes**"
		}
		typeStr := linkifyType(f.GoType, structs)
		fmt.Fprintf(sb, "| `%s` | %s | %s | %s |\n", f.YAMLKey, typeStr, req, f.Doc)

		// Collect locally defined struct types referenced by this field.
		for typeName := range structs {
			if !visited[typeName] && !seen[typeName] && containsTypeName(f.GoType, typeName) {
				nested = append(nested, typeName)
				seen[typeName] = true
			}
		}
	}
	sb.WriteString("\n")

	for _, typeName := range nested {
		heading := strings.Repeat("#", headingLevel)
		fmt.Fprintf(sb, "%s %s\n\n", heading, typeName)
		renderStruct(sb, typeName, structs, visited, headingLevel+1)
	}
}

// linkifyType wraps any local struct name in the type string with a markdown anchor link.
func linkifyType(typeStr string, structs map[string]StructDef) string {
	for name := range structs {
		anchor := strings.ToLower(name)
		linked := fmt.Sprintf("[%s](#%s)", name, anchor)
		typeStr = strings.ReplaceAll(typeStr, name, linked)
	}
	return typeStr
}

// containsTypeName reports whether a formatted type string references the given struct name.
func containsTypeName(typeStr, name string) bool {
	return strings.Contains(typeStr, name)
}
