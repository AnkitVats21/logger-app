package main

import "embed"

// TemplateFS embeds the entire templates directory tree into the binary at compile time.
//
//go:embed templates
var TemplateFS embed.FS
