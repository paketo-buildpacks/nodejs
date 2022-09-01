package integration_test

import (
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
)

var nodeBuildpack string

func TestIntegration(t *testing.T) {
	var pack occam.Pack
	pack = occam.NewPack()

	Expect := NewWithT(t).Expect

	output, err := exec.Command("bash", "-c", "../scripts/package.sh --version 1.2.3").CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), string(output))

	nodeBuildpack, err = filepath.Abs("../build/buildpackage.cnb")
	Expect(err).NotTo(HaveOccurred())

	SetDefaultEventuallyTimeout(10 * time.Second)

	suite := spec.New("Integration", spec.Parallel(), spec.Report(report.Terminal{}))
	suite("NPM", testNPM)
	suite("NodeStart", testNodeStart)
	suite("ReproducibleBuilds", testReproducibleBuilds)
	suite("Yarn", testYarn)
	suite.Run(t)

	builder, _ := pack.Builder.Inspect.Execute()
	if builder.BuilderName != "paketobuildpacks/builder-jammy-buildpackless-base" {

		spec.Run(t, "StackUpgrades", testStackUpgrades, spec.Report(report.Terminal{}))

	}
}
