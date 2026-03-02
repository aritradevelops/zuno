package docker

import "embed"

//go:embed templates/*
var templates embed.FS

const pathToDocker = "./docker"
