package integration_test

import (
	"encoding/json"
	"os"
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

var settings struct {
	Extensions struct {
		UbiNodejsExtension struct {
			Online string
		}
	}

	Config struct {
		UbiNodejsExtension string `json:"ubi-nodejs-extension"`
	}
}

func TestIntegration(t *testing.T) {
	Expect := NewWithT(t).Expect

	output, err := exec.Command("bash", "-c", "../scripts/package.sh --version 1.2.3").CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), string(output))

	pack := occam.NewPack()
	builder, err := pack.Builder.Inspect.Execute()
	Expect(err).NotTo(HaveOccurred())

	file, err := os.Open("../integration.json")
	Expect(err).NotTo(HaveOccurred())

	Expect(json.NewDecoder(file).Decode(&settings.Config)).To(Succeed())
	Expect(file.Close()).To(Succeed())

	buildpackStore := occam.NewBuildpackStore()

	if builder.BuilderName == "paketobuildpacks/builder-ubi8-buildpackless-base" || builder.BuilderName == "paketobuildpacks/ubi-9-builder-buildpackless" {
		settings.Extensions.UbiNodejsExtension.Online, err = buildpackStore.Get.
			Execute(settings.Config.UbiNodejsExtension)
		Expect(err).ToNot(HaveOccurred())
	}

	nodeBuildpack, err = filepath.Abs("../build/buildpackage.cnb")
	Expect(err).NotTo(HaveOccurred())

	SetDefaultEventuallyTimeout(10 * time.Second)

	suite := spec.New("Integration", spec.Parallel(), spec.Report(report.Terminal{}))
	suite("NodeStart", testNodeStart)
	suite("NPM", testNPM)
	suite("ReproducibleBuilds", testReproducibleBuilds)
	suite("Yarn", testYarn)
	suite.Run(t)
}
