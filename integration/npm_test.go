package integration_test

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
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

	context("when building a node app that uses npm", func() {
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
			source, err = occam.Source(filepath.Join("testdata", "npm"))
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		it("builds a working OCI image for a simple app", func() {
			var err error
			var logs fmt.Stringer
			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(nodeBuildpack).
				WithPullPolicy("never").
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred())

			Expect(logs).To(ContainLines(ContainSubstring("Node Engine Buildpack")))
			Expect(logs).To(ContainLines(ContainSubstring("NPM Install Buildpack")))
			Expect(logs).To(ContainLines(ContainSubstring("NPM Start Buildpack")))
			Expect(logs).NotTo(ContainLines(ContainSubstring("Procfile Buildpack")))
			Expect(logs).NotTo(ContainLines(ContainSubstring("Environment Variables Buildpack")))
			Expect(logs).NotTo(ContainLines(ContainSubstring("Image Labels Buildpack")))

			container, err = docker.Container.Run.
				WithEnv(map[string]string{"PORT": "8080"}).
				WithPublish("8080").
				WithPublishAll().
				Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			Eventually(container, "5s").Should(BeAvailable())

			response, err := http.Get(fmt.Sprintf("http://localhost:%s/env", container.HostPort("8080")))
			Expect(err).NotTo(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK))

			var env struct {
				NpmConfigLoglevel string `json:"NPM_CONFIG_LOGLEVEL"`
			}
			Expect(json.NewDecoder(response.Body).Decode(&env)).To(Succeed())
			Expect(env.NpmConfigLoglevel).To(Equal("error"))

		})

		context("when using optional utility buildpacks", func() {
			it.Before(func() {
				Expect(ioutil.WriteFile(filepath.Join(source, "Procfile"), []byte("web: node server.js"), 0644)).To(Succeed())
			})

			it("builds a working OCI image for a simple app and uses the Procfile start command and other utility buildpacks", func() {
				var err error
				var logs fmt.Stringer
				image, logs, err = pack.WithNoColor().Build.
					WithBuildpacks(nodeBuildpack).
					WithPullPolicy("never").
					WithEnv(map[string]string{
						"BPE_SOME_VARIABLE": "some-value",
						"BP_IMAGE_LABELS":   "some-label=some-value",
					}).
					Execute(name, source)
				Expect(err).NotTo(HaveOccurred())

				Expect(logs).To(ContainLines(ContainSubstring("Node Engine Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("NPM Install Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("NPM Start Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("Procfile Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("web: node server.js")))
				Expect(logs).To(ContainLines(ContainSubstring("Environment Variables Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("Image Labels Buildpack")))

				Expect(image.Buildpacks[5].Layers["environment-variables"].Metadata["variables"]).To(Equal(map[string]interface{}{"SOME_VARIABLE": "some-value"}))
				Expect(image.Labels["some-label"]).To(Equal("some-value"))

				container, err = docker.Container.Run.
					WithEnv(map[string]string{"PORT": "8080"}).
					WithPublish("8080").
					WithPublishAll().
					Execute(image.ID)
				Expect(err).NotTo(HaveOccurred())

				Eventually(container, "5s").Should(BeAvailable())

				response, err := http.Get(fmt.Sprintf("http://localhost:%s/env", container.HostPort("8080")))
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(http.StatusOK))

				var env struct {
					NpmConfigLoglevel string `json:"NPM_CONFIG_LOGLEVEL"`
				}
				Expect(json.NewDecoder(response.Body).Decode(&env)).To(Succeed())
				Expect(env.NpmConfigLoglevel).To(Equal("error"))
			})
		})

		context("when using CA certificates", func() {
			var (
				client *http.Client
			)

			it.Before(func() {
				var err error
				source, err = occam.Source(filepath.Join("testdata", "ca_cert_apps"))
				Expect(err).NotTo(HaveOccurred())

				caCert, err := ioutil.ReadFile(fmt.Sprintf("%s/client-certs/ca.pem", source))
				Expect(err).ToNot(HaveOccurred())

				caCertPool := x509.NewCertPool()
				caCertPool.AppendCertsFromPEM(caCert)

				cert, err := tls.LoadX509KeyPair(fmt.Sprintf("%s/client-certs/cert.pem", source), fmt.Sprintf("%s/client-certs/key.pem", source))
				Expect(err).ToNot(HaveOccurred())

				client = &http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{
							RootCAs:      caCertPool,
							Certificates: []tls.Certificate{cert},
							MinVersion:   tls.VersionTLS12,
						},
					},
				}
			})

			it("builds a working OCI image and uses a client-side CA cert for requests", func() {
				var err error
				var logs fmt.Stringer
				image, logs, err = pack.WithNoColor().Build.
					WithBuildpacks(nodeBuildpack).
					WithPullPolicy("never").
					Execute(name, filepath.Join(source, "npm_server"))
				Expect(err).NotTo(HaveOccurred())

				Expect(logs).To(ContainLines(ContainSubstring("CA Certificates Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("Node Engine Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("NPM Install Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("NPM Start Buildpack")))

				container, err = docker.Container.Run.
					WithPublish("8080").
					WithEnv(map[string]string{
						"PORT":                 "8080",
						"SERVICE_BINDING_ROOT": "/bindings",
						"NODE_OPTIONS":         "--use-openssl-ca",
					}).
					WithVolume(fmt.Sprintf("%s/binding:/bindings/ca-certificates", source)).
					Execute(image.ID)
				Expect(err).NotTo(HaveOccurred())

				Eventually(func() string {
					cLogs, err := docker.Container.Logs.Execute(container.ID)
					Expect(err).NotTo(HaveOccurred())
					return cLogs.String()
				}).Should(
					ContainSubstring("Added 1 additional CA certificate(s) to system truststore"),
				)

				request, err := http.NewRequest("GET", fmt.Sprintf("https://localhost:%s/env", container.HostPort("8080")), nil)
				Expect(err).NotTo(HaveOccurred())

				var response *http.Response
				Eventually(func() error {
					var err error
					response, err = client.Do(request)
					return err
				}).Should(BeNil())
				Expect(err).NotTo(HaveOccurred())
				defer response.Body.Close()

				Expect(response.StatusCode).To(Equal(http.StatusOK))

				var env struct {
					NpmConfigLoglevel string `json:"NPM_CONFIG_LOGLEVEL"`
				}

				Expect(json.NewDecoder(response.Body).Decode(&env)).To(Succeed())
				Expect(env.NpmConfigLoglevel).To(Equal("error"))
			})
		})
	})
}
