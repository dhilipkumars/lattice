package receptor_json_runner_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestReceptorJsonRunner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ReceptorJsonRunner Suite")
}
