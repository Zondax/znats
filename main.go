package main

import (
	"github.com/zondax/znats/internal/commands"
	"github.com/zondax/znats/internal/conf"
	"github.com/zondax/znats/internal/version"
	"strings"
)

import (
	"github.com/zondax/golem/pkg/cli"
)

func main() {
	appName := "znats"
	envPrefix := strings.ReplaceAll(appName, "-", "_")

	appSettings := cli.AppSettings{
		Name:        appName,
		Description: "Nats wrapper for Golem",
		ConfigPath:  "$HOME/.znats/",
		EnvPrefix:   envPrefix,
		GitVersion:  version.GitVersion,
		GitRevision: version.GitRevision,
	}

	// Define application level features
	cli := cli.New[conf.Config](appSettings)
	defer cli.Close()

	cli.GetRoot().AddCommand(commands.GetStartCommand(cli))

	cli.Run()
}
