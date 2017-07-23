package main

import (
	"code.cloudfoundry.org/buildpackapplifecycle/buildpackrunner"
	"github.com/cloudfoundry/libbuildpack"
	"fmt"
)

func WriteStartCommand(stagingInfoFile string, outputFile string) error {
	var stagingInfo buildpackrunner.DeaStagingInfo

	err := libbuildpack.NewYAML().Load(stagingInfoFile, &stagingInfo)
	if err != nil {
		return err
	}
	var webStartCommand = map[string]string{
		"web": fmt.Sprintf(
			"$HOME/gobis-server -c $HOME/gobis-config.yml & PORT=%s %s",
			getPort(),
			stagingInfo.StartCommand,
		),
	}

	release := buildpackrunner.Release{
		DefaultProcessTypes: webStartCommand,
	}

	return libbuildpack.NewYAML().Write(outputFile, &release)
}
