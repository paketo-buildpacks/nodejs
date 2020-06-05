package integration_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
	. "github.com/paketo-buildpacks/occam/matchers"
)

func testNPM(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect     = NewWithT(t).Expect
		Eventually = NewWithT(t).Eventually

		pack   occam.Pack
		docker occam.Docker
	)

	it.Before(func() {
		pack = occam.NewPack()
		docker = occam.NewDocker()
	})

	context("when the node_modules are not vendored", func() {
		var (
			image     occam.Image
			container occam.Container

			name string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
		})

		it("builds a working OCI image for a simple app", func() {
			var err error
			image, _, err = pack.Build.
				WithBuildpacks(nodeBuildpack).
				WithNoPull().
				Execute(name, filepath.Join("testdata", "npm"))
			Expect(err).NotTo(HaveOccurred())

			container, err = docker.Container.Run.Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			Eventually(container, "5s").Should(BeAvailable(), ContainerLogs(container.ID))

			response, err := http.Get(fmt.Sprintf("http://localhost:%s/env", container.HostPort()))
			Expect(err).NotTo(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK))

			var env struct {
				NpmConfigLoglevel   string `json:"NPM_CONFIG_LOGLEVEL"`
				NpmConfigProduction string `json:"NPM_CONFIG_PRODUCTION"`
			}
			Expect(json.NewDecoder(response.Body).Decode(&env)).To(Succeed())
			Expect(env.NpmConfigLoglevel).To(Equal("error"))
			Expect(env.NpmConfigProduction).To(Equal("true"))
		})
	})
}
