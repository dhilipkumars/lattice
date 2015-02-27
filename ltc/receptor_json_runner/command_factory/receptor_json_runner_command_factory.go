package command_factory

import (
	"fmt"

	"github.com/cloudfoundry-incubator/lattice/ltc/logs/console_tailed_logs_outputter"
	"github.com/cloudfoundry-incubator/lattice/ltc/terminal"
	"github.com/codegangsta/cli"
)

//appRunner:             config.AppRunner,
//dockerMetadataFetcher: config.DockerMetadataFetcher,
//output:                config.Output,
//timeout:               config.Timeout,
//domain:                config.Domain,
//env:                   config.Env,
//clock:                 config.Clock,
//tailedLogsOutputter:   config.TailedLogsOutputter,

type ReceptorJsonCommandFactory struct {
	ui  terminal.UI
	tlo console_tailed_logs_outputter.TailedLogsOutputter
}

func NewReceptorJsonRunnerCommandFactory(ui terminal.UI, tlo console_tailed_logs_outputter.TailedLogsOutputter) *ReceptorJsonCommandFactory {
	return &ReceptorJsonCommandFactory{
		ui:  ui,
		tlo: tlo,
	}
}

func (commandFactory *ReceptorJsonCommandFactory) MakeCreateAppFromJsonCommand() cli.Command {
	var createCommand = cli.Command{
		Name:        "create-from-json",
		ShortName:   "cfj",
		Usage:       "ltc create-from-json /path/to/file.json",
		Description: `Create a LRP on lattice from JSON object`,
		Action:      commandFactory.createAppFromJson,
	}
	return createCommand
}

func (cmd *ReceptorJsonCommandFactory) createAppFromJson(c *cli.Context) {
	cmd.ui.Say(fmt.Sprintf("Attempting to Create LRP from %s", c.Args().First()))
}
