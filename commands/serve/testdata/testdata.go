package testdata

import "embed"

var files embed.FS

//go:embed route_guide_db.json
var Data []byte
