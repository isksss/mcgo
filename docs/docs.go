package docs

import (
	_ "embed"
)

//go:embed config/paper.toml
var PaperToml []byte

//go:embed config/velocity.toml
var VelocityToml []byte
