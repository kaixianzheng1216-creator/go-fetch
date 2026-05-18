// Package migrations embeds database migration SQL files.
package migrations

import "embed"

// FS contains embedded SQL migration files.
//
//go:embed *.sql
var FS embed.FS
