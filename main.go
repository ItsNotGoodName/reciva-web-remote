package main

import (
	"fmt"

	"github.com/ItsNotGoodName/reciva-web-remote/cmd"
	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/build"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/model"
)

var (
	builtBy    = "unknown"
	commit     = ""
	date       = ""
	releaseURL = ""
	summary    = "dev"
	version    = "dev"
)

//go:generate swag fmt -d http,internal/model,internal/state
//go:generate swag init -g api.go -d http,internal/model,internal/state,internal/pubsub --outputTypes go,json --output ./docs/swagger
//go:generate npm run swag --prefix web
func main() {
	// Create config
	cfg := config.NewConfig(config.WithFlag)

	// Show version and exit
	if cfg.ShowVersion {
		fmt.Println(version)
		return
	}

	// Show info and exit
	if cfg.ShowInfo {
		fmt.Printf("Version: %s\nCommit: %s\nDate: %s\nBuilt by: %s\nRelease url: %s\nSummary: %s\n", version, commit, date, builtBy, releaseURL, summary)
		return
	}

	build.CurrentBuild = model.Build{
		BuiltBy:    builtBy,
		Commit:     commit,
		Date:       date,
		ReleaseURL: releaseURL,
		Summary:    summary,
		Version:    version,
	}

	cmd.Server(cfg)
}
