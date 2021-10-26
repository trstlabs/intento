package docs

import "embed"

// Docs are Trustless Hub docs.
//go:embed *.md */*.md
var Docs embed.FS
