package command_factory

import (
	"github.com/cloudfoundry-incubator/lattice/ltc/exit_handler"
	"github.com/cloudfoundry-incubator/lattice/ltc/logs/console_tailed_logs_outputter"
	"github.com/cloudfoundry-incubator/lattice/ltc/logs/reserved_app_ids"
	"github.com/cloudfoundry-incubator/lattice/ltc/output"
	"github.com/cloudfoundry-incubator/lattice/ltc/app_examiner"
	"github.com/codegangsta/cli"
)

type logsCommandFactory struct {
	output              *output.Output
	tailedLogsOutputter console_tailed_logs_outputter.TailedLogsOutputter
	exitHandler         exit_handler.ExitHandler
	app		    app_examiner.AppExaminer
}

func NewLogsCommandFactory(output *output.Output, tailedLogsOutputter console_tailed_logs_outputter.TailedLogsOutputter, exitHandler exit_handler.ExitHandler) *logsCommandFactory {
	return &logsCommandFactory{
		output:              output,
		tailedLogsOutputter: tailedLogsOutputter,
		exitHandler:         exitHandler,
	}
}

func (factory *logsCommandFactory) MakeLogsCommand( app app_examiner.AppExaminer) cli.Command {
	var logsCommand = cli.Command{
		Name:        "logs",
		ShortName:   "lo",
		Usage:       "Streams logs from the specified application",
		Description: "ltc logs APP_NAME",
		Action:      factory.tailLogs,
		Flags:       []cli.Flag{},
	}

	factory.app = app

	return logsCommand
}

func (factory *logsCommandFactory) MakeDebugLogsCommand() cli.Command {
	return cli.Command{
		Name:        "debug-logs",
		Usage:       "Streams logs from the lattice cluster components",
		Description: "ltc debug-logs",
		Action:      factory.tailDebugLogs,
	}
}

func (factory *logsCommandFactory) tailLogs(context *cli.Context) {
	appGuid := context.Args().First()

	if appGuid == "" {
		factory.output.IncorrectUsage("")
		return
	}

// Check if there is really such app before we start waiting for its logs. 
	_, err := factory.app.AppStatus (appGuid)
	
	if err != nil {
		factory.output.Say(err.Error())	
		return 
	}	

	factory.tailedLogsOutputter.OutputTailedLogs(appGuid)
}

func (factory *logsCommandFactory) tailDebugLogs(context *cli.Context) {
	factory.tailedLogsOutputter.OutputTailedLogs(reserved_app_ids.LatticeDebugLogStreamAppId)
}
