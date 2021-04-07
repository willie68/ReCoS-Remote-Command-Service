package web

import (
	"embed"
)

//go:embed webclient/*
var WebClientAssets embed.FS

//go:embed webadmin/*
var WebAdminAssets embed.FS

//go:embed zoneinfo/zoneinfo.zip
var ZoneinfoAssets []byte
