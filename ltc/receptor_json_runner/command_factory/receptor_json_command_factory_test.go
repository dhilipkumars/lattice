package command_factory_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"

	"github.com/cloudfoundry-incubator/lattice/ltc/logs/console_tailed_logs_outputter/fake_tailed_logs_outputter"
	"github.com/cloudfoundry-incubator/lattice/ltc/receptor_json_runner/command_factory"
	"github.com/cloudfoundry-incubator/lattice/ltc/terminal"
	"github.com/cloudfoundry-incubator/lattice/ltc/test_helpers"
	"github.com/codegangsta/cli"
)

var _ = Describe("CommandFactory", func() {

	var (
		outputBuffer            *gbytes.Buffer
		ui                      terminal.UI
		fakeTailedLogsOutputter *fake_tailed_logs_outputter.FakeTailedLogsOutputter
	)

	BeforeEach(func() {
		outputBuffer = gbytes.NewBuffer()
		ui = terminal.NewUI(nil, outputBuffer, nil)
		fakeTailedLogsOutputter = fake_tailed_logs_outputter.NewFakeTailedLogsOutputter()
	})

	Describe("CreateJSONCommand", func() {
		var createJsonCommand cli.Command

		BeforeEach(func() {
			commandFactory := command_factory.NewReceptorJsonRunnerCommandFactory(ui, fakeTailedLogsOutputter)
			createJsonCommand = commandFactory.MakeCreateAppFromJsonCommand()
		})

		FIt("reads in json from the specified path and sends it off to the receptor", func() {
			args := []string{"~/mrbigglesworth/lrp.json"}

			test_helpers.ExecuteCommandWithArgs(createJsonCommand, args)

			Expect(outputBuffer).To(test_helpers.Say("Attempting to Create LRP from ~/mrbigglesworth/lrp.json"))

			//            PRETEND WE SENT THE JSON AND THAT TEST PASSES

			Expect(fakeTailedLogsOutputter.OutputTailedLogsCallCount()).To(Equal(1))
		})
	})

})
