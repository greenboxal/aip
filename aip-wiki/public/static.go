package public

import "embed"

//go:embed css/*
var content embed.FS

func Content() embed.FS {
	return content
}
