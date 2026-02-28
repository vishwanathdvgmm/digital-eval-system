// Package ui provides the embedded frontend static files.
package ui

import "embed"

//go:embed dist/*
var StaticFiles embed.FS
