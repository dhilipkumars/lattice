package command_factory_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"

	"github.com/cloudfoundry-incubator/lattice/ltc/exit_handler"
	"github.com/cloudfoundry-incubator/lattice/ltc/exit_handler/fake_exit_handler"
	"github.com/cloudfoundry-incubator/lattice/ltc/logs/command_factory"
	"github.com/cloudfoundry-incubator/lattice/ltc/logs/console_tailed_logs_outputter/fake_tailed_logs_outputter"
	"github.com/cloudfoundry-incubator/lattice/ltc/logs/reserved_app_ids"
	"github.com/cloudfoundry-incubator/lattice/ltc/output"
	"github.com/cloudfoundry-incubator/lattice/ltc/test_helpers"
        "github.com/cloudfoundry-incubator/lattice/ltc/app_examiner/fake_app_examiner"
	"github.com/codegangsta/cli"
)

var _ = Describe("CommandFactory", func() {
	var (
		outputBuffer            *gbytes.Buffer
		fakeTailedLogsOutputter *fake_tailed_logs_outputter.FakeTailedLogsOutputter
		signalChan              chan os.Signal
		exitHandler             exit_handler.ExitHandler
	)

	BeforeEach(func() {
		outputBuffer = gbytes.NewBuffer()
		fakeTailedLogsOutputter = fake_tailed_logs_outputter.NewFakeTailedLogsOutputter()
		signalChan = make(chan os.Signal)
		exitHandler = &fake_exit_handler.FakeExitHandler{}
	})

	Describe("LogsCommand", func() {

		var logsCommand cli.Command
		var appExaminer  *fake_app_examiner.FakeAppExaminer

		BeforeEach(func() {
			commandFactory := command_factory.NewLogsCommandFactory(output.New(outputBuffer), fakeTailedLogsOutputter, exitHandler)
			appExaminer = &fake_app_examiner.FakeAppExaminer{}
			logsCommand = commandFactory.MakeLogsCommand(appExaminer)
		})

		It("tails logs", func() {
			args := []string{
				"my-app-guid",
			}

			test_helpers.AsyncExecuteCommandWithArgs(logsCommand, args)

			Eventually(fakeTailedLogsOutputter.OutputTailedLogsCallCount).Should(Equal(1))
			Expect(fakeTailedLogsOutputter.OutputTailedLogsArgsForCall(0)).To(Equal("my-app-guid"))
		})

		It("handles invalid appguids", func() {
			test_helpers.ExecuteCommandWithArgs(logsCommand, []string{})

			Expect(outputBuffer).To(test_helpers.Say("Incorrect Usage"))
			Expect(fakeTailedLogsOutputter.OutputTailedLogsCallCount()).To(Equal(0))
		})
	})

	Describe("DebugLogsCommand", func() {

		var debugLogsCommand cli.Command

		BeforeEach(func() {
			commandFactory := command_factory.NewLogsCommandFactory(output.New(outputBuffer), fakeTailedLogsOutputter, exitHandler)
			debugLogsCommand = commandFactory.MakeDebugLogsCommand()
		})

		It("tails logs from the lattice-debug stream", func() {
			test_helpers.AsyncExecuteCommandWithArgs(debugLogsCommand, []string{})

			Eventually(fakeTailedLogsOutputter.OutputTailedLogsCallCount).Should(Equal(1))
			Expect(fakeTailedLogsOutputter.OutputTailedLogsArgsForCall(0)).To(Equal(reserved_app_ids.LatticeDebugLogStreamAppId))
		})
	})

})
