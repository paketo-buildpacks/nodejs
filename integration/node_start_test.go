package integration_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
	. "github.com/paketo-buildpacks/occam/matchers"
)

func testNodeStart(t *testing.T, context spec.G, it spec.S) {
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

	context("when building a node app that does not use a package manager", func() {
		var (
			image     occam.Image
			container occam.Container

			name   string
			source string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())
			source, err = occam.Source(filepath.Join("testdata", "no_package_manager"))
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		it("should build a working OCI image and run the app", func() {
			var err error
			var logs fmt.Stringer
			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(nodeBuildpack).
				WithPullPolicy("never").
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred())

			Expect(logs).To(ContainLines(ContainSubstring("Node Engine Buildpack")))
			Expect(logs).To(ContainLines(ContainSubstring("Node Start Buildpack")))
			Expect(logs).NotTo(ContainLines(ContainSubstring("Procfile Buildpack")))
			Expect(logs).NotTo(ContainLines(ContainSubstring("Environment Variables Buildpack")))

			container, err = docker.Container.Run.
				WithEnv(map[string]string{"PORT": "8080"}).
				WithPublish("8080").
				WithPublishAll().
				Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			Eventually(container).Should(BeAvailable())

			response, err := http.Get(fmt.Sprintf("http://localhost:%s", container.HostPort("8080")))
			Expect(err).NotTo(HaveOccurred())
			defer response.Body.Close()
			Expect(response.StatusCode).To(Equal(http.StatusOK))

			content, err := ioutil.ReadAll(response.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(ContainSubstring("hello world"))
		})

		context("when using optional utility buildpacks", func() {
			it.Before(func() {
				Expect(ioutil.WriteFile(filepath.Join(source, "Procfile"), []byte("web: node server.js"), 0644)).To(Succeed())
			})
			it("should build a working OCI image and run the app with the start command from the Procfile and other utility buildpacks", func() {
				var err error
				var logs fmt.Stringer
				image, logs, err = pack.WithNoColor().Build.
					WithBuildpacks(nodeBuildpack).
					WithPullPolicy("never").
					WithEnv(map[string]string{"BPE_SOME_VARIABLE": "some-value"}).
					Execute(name, source)
				Expect(err).NotTo(HaveOccurred())

				Expect(logs).To(ContainLines(ContainSubstring("Node Engine Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("Node Start Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("Procfile Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("web: node server.js")))
				Expect(logs).To(ContainLines(ContainSubstring("Environment Variables Buildpack")))

				Expect(image.Buildpacks[3].Layers["environment-variables"].Metadata["variables"]).To(Equal(map[string]interface{}{"SOME_VARIABLE": "some-value"}))

				container, err = docker.Container.Run.
					WithEnv(map[string]string{"PORT": "8080"}).
					WithPublish("8080").
					WithPublishAll().
					Execute(image.ID)
				Expect(err).NotTo(HaveOccurred())

				Eventually(container).Should(BeAvailable())

				response, err := http.Get(fmt.Sprintf("http://localhost:%s", container.HostPort("8080")))
				Expect(err).NotTo(HaveOccurred())
				defer response.Body.Close()
				Expect(response.StatusCode).To(Equal(http.StatusOK))

				content, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(content)).To(ContainSubstring("hello world"))
			})
		})
	})
}
