package web

import "embed"

// Assets holds our static frontend files.
// We use '*' to grab everything in this folder (index.html, style.css, etc).
//
//go:embed *
var Assets embed.FS
