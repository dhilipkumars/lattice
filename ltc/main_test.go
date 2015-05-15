package main_test

import (
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var (
	cli string
)

var _ = BeforeSuite(func() {
	var err error
	cli, err = gexec.Build("github.com/cloudfoundry-incubator/lattice/ltc")
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

var _ = Describe("lattice-cli", func() {
	It("compiles and displays help text", func() {
		command := exec.Command(cli)

		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)

		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))
		Eventually(session.Out).Should(gbytes.Say("ltc - Command line interface for Lattice."))
	})
	Describe("exit codes", func() {
		It("exits non-zero when an unknown command is invoked", func() {
			command := exec.Command(cli, "unknownCommand")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session, 3*time.Second).Should(gbytes.Say("not a registered command"))
			Eventually(session).Should(gexec.Exit(0))
		})
		It("exits non-zero when known command is invoked with invalid option", func() {
			command := exec.Command(cli, "status", "--badFlag")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session, 3*time.Second).Should(gexec.Exit(1))
		})
	})
	Describe("Flag verification", func() {
		It("informs user for any incorrect provided flags", func() {
			command := exec.Command(cli, "create", "--instances", "--bad-flag")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session.Out).Should(gbytes.Say("\"--bad-flag\""))
			Consistently(session.Out).ShouldNot(gbytes.Say("\"--instances\""))
		})
		It("checks flags with prefix '--'", func() {
			command := exec.Command(cli, "create", "not-a-flag", "--invalid-flag")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session.Out).Should(gbytes.Say("Unknown flag \"--invalid-flag\""))
			Consistently(session.Out).ShouldNot(gbytes.Say("Unknown flag \"not-a-flag\""))
		})
		It("checks flags with prefix '-'", func() {
			command := exec.Command(cli, "create", "not-a-flag", "-invalid-flag")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session.Out).Should(gbytes.Say("\"-invalid-flag\""))
			Consistently(session.Out).ShouldNot(gbytes.Say("\"not-a-flag\""))
		})
		It("checks flags but ignores the value after '=' ", func() {
			command := exec.Command(cli, "create", "-i=1", "-invalid-flag=blarg")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session.Out).Should(gbytes.Say("\"-invalid-flag\""))
			Consistently(session.Out).ShouldNot(gbytes.Say("Unknown flag \"-p\""))
		})
		It("outputs all unknown flags in single sentence", func() {
			command := exec.Command(cli, "create", "--bad-flag1", "--bad-flag2", "--bad-flag3")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session.Out).Should(gbytes.Say("\"--bad-flag1\", \"--bad-flag2\", \"--bad-flag3\""))
		})
		It("only checks input flags against flags from the provided command", func() {
			command := exec.Command(cli, "create", "--instances", "--skip-ssl-validation")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session.Out).Should(gbytes.Say("\"--skip-ssl-validation\""))
		})
		It("accepts -h and --h flags for all commands", func() {
			command := exec.Command(cli, "create", "-h")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Consistently(session.Out).ShouldNot(gbytes.Say("Unknown flag \"-h\""))
			command = exec.Command(cli, "target", "--h")
			session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Consistently(session.Out).ShouldNot(gbytes.Say("Unknown flag \"--h\""))
		})
		Context("When a negative integer is preceeded by a valid flag", func() {
			It("skips validation for negative integer flag values", func() {
				command := exec.Command(cli, "create", "-i", "-10")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session.Out).ShouldNot(gbytes.Say("\"-10\""))
			})
		})
		Context("When a negative integer is preceeded by a invalid flag", func() {
			It("validates the negative integer as a flag", func() {
				command := exec.Command(cli, "create", "-badflag", "-10")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session.Out).Should(gbytes.Say("\"-badflag\""))
				Eventually(session.Out).Should(gbytes.Say("\"-10\""))
			})
		})
	})
})
