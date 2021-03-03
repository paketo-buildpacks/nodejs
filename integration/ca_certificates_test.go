package integration_test

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
	. "github.com/paketo-buildpacks/occam/matchers"
)

func testCaCerts(t *testing.T, context spec.G, it spec.S) {
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

	context("the app uses ca certificates", func() {
		var (
			image     occam.Image
			container occam.Container

			name   string
			source string
			client *http.Client
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())
			source, err = occam.Source(filepath.Join("testdata", "ca-certs"))
			Expect(err).NotTo(HaveOccurred())

			caCert, err := ioutil.ReadFile(fmt.Sprintf("%s/client/ca.pem", source))
			Expect(err).ToNot(HaveOccurred())

			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)

			cert, err := tls.LoadX509KeyPair(fmt.Sprintf("%s/client/cert.pem", source), fmt.Sprintf("%s/client/key.pem", source))
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

		it.After(func() {
			Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		it("builds a working OCI image and requests made with a client-side cert succeed", func() {
			var err error
			var logs fmt.Stringer
			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(nodeBuildpack).
				WithPullPolicy("never").
				Execute(name, filepath.Join(source, "server"))
			Expect(err).NotTo(HaveOccurred())

			Expect(logs).To(ContainLines(ContainSubstring("CA Certificates Buildpack")))
			Expect(logs).To(ContainLines(ContainSubstring("Node Engine Buildpack")))
			Expect(logs).To(ContainLines(ContainSubstring("Node Start Buildpack")))

			container, err = docker.Container.Run.
				WithPublish("8080").
				WithEnv(map[string]string{
					"PORT":                 "8080",
					"SERVICE_BINDING_ROOT": "/bindings",
					"NODE_OPTIONS":         "--use-openssl-ca",
				}).
				WithVolume(fmt.Sprintf("%s/server/binding:/bindings/ca-certificates", source)).
				Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				cLogs, err := docker.Container.Logs.Execute(container.ID)
				Expect(err).NotTo(HaveOccurred())
				return cLogs.String()
			}).Should(
				ContainSubstring("Added 1 additional CA certificate(s) to system truststore"),
			)

			// give the app 1 second to get set up before we curl it
			time.Sleep(1 * time.Second)

			request, err := http.NewRequest("GET", fmt.Sprintf("https://localhost:%s", container.HostPort("8080")), nil)
			Expect(err).NotTo(HaveOccurred())

			response, err := client.Do(request)
			Expect(err).NotTo(HaveOccurred())
			defer response.Body.Close()

			Expect(response.StatusCode).To(Equal(http.StatusOK))

			content, err := ioutil.ReadAll(response.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(ContainSubstring("Hello, world!"))
		})
	})
}
