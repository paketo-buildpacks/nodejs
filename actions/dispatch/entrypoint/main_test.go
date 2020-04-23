package main_test

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"os/exec"
	"testing"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testSendDispatch(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect     = NewWithT(t).Expect
		Eventually = NewWithT(t).Eventually
	)

	context("when given a release event payload", func() {
		var (
			eventPath string
			api       *httptest.Server
			requests  []*http.Request
		)

		it.Before(func() {
			requests = []*http.Request{}

			file, err := ioutil.TempFile("", "event.json")
			Expect(err).NotTo(HaveOccurred())

			_, err = file.WriteString(`{
				"repository": {
					"full_name": "thitch97/nodejs"
				},
				"release": {
					"name": "Release v1.2.3"
				}
			}`)
			Expect(err).NotTo(HaveOccurred())

			Expect(file.Close()).To(Succeed())

			eventPath = file.Name()

			api = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				dump, _ := httputil.DumpRequest(req, true)
				receivedRequest, _ := http.ReadRequest(bufio.NewReader(bytes.NewBuffer(dump)))

				requests = append(requests, receivedRequest)

				if req.Header.Get("Authorization") != "token some-github-token" {
					w.WriteHeader(http.StatusForbidden)
					return
				}

				switch req.URL.Path {
				case "/repos/some-org/some-repo/dispatches":
					w.WriteHeader(http.StatusNoContent)

				case "/repos/loop-org/loop-repo/dispatches":
					w.Header().Set("Location", "/repos/loop-org/loop-repo/dispatches")
					w.WriteHeader(http.StatusFound)

				case "/repos/fail-org/fail-repo/dispatches":
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error": "server-error"}`))

				default:
					t.Fatal(fmt.Sprintf("unknown request: %s", dump))
				}
			}))
		})

		it.After(func() {
			api.Close()

			Expect(os.RemoveAll(eventPath)).To(Succeed())
		})

		it("sends a repository_dispatch webhook to a repo", func() {
			command := exec.Command(
				entrypoint,
				"--endpoint", api.URL,
				"--repo", "some-org/some-repo",
				"--token", "some-github-token",
			)
			command.Env = append(command.Env, fmt.Sprintf("GITHUB_EVENT_PATH=%s", eventPath))
			buffer := gbytes.NewBuffer()

			session, err := gexec.Start(command, buffer, buffer)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(0), fmt.Sprintf("output:\n%s\n", buffer.Contents()))

			Expect(buffer).To(gbytes.Say(`Dispatching`))
			Expect(buffer).To(gbytes.Say(`Repository: thitch97/nodejs`))
			Expect(buffer).To(gbytes.Say(`Release: Release v1.2.3`))
			Expect(buffer).To(gbytes.Say(`Success!`))

			Expect(requests).To(HaveLen(1))
			request := requests[0]

			Expect(request.Method).To(Equal("POST"))
			Expect(request.URL.Path).To(Equal("/repos/some-org/some-repo/dispatches"))

			body, err := ioutil.ReadAll(request.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(MatchJSON(`{
				"event_type": "update-buildpack-toml",
				"client_payload": {
					"repo": "thitch97/nodejs",
					"release": "Release v1.2.3"
				}
			}`))
		})

		context("failure cases", func() {
			context("when the event path does not exist", func() {
				it.Before(func() {
					Expect(os.RemoveAll(eventPath)).To(Succeed())
				})

				it("prints an error message and exits non-zero", func() {
					command := exec.Command(
						entrypoint,
						"--endpoint", api.URL,
						"--repo", "some-org/some-repo",
						"--token", "some-github-token",
					)
					command.Env = append(command.Env, fmt.Sprintf("GITHUB_EVENT_PATH=%s", eventPath))
					buffer := gbytes.NewBuffer()

					session, err := gexec.Start(command, buffer, buffer)
					Expect(err).NotTo(HaveOccurred())

					Eventually(session).Should(gexec.Exit(1), fmt.Sprintf("output:\n%s\n", buffer.Contents()))

					Expect(buffer).To(gbytes.Say(`Dispatching`))
					Expect(buffer).To(gbytes.Say(`Error: failed to read \$GITHUB_EVENT_PATH:`))
					Expect(buffer).To(gbytes.Say(`no such file or directory`))
				})
			})

			context("when the event file contains malformed json", func() {
				it.Before(func() {
					Expect(ioutil.WriteFile(eventPath, []byte("%%%"), 0644)).To(Succeed())
				})

				it("prints an error message and exits non-zero", func() {
					command := exec.Command(
						entrypoint,
						"--endpoint", api.URL,
						"--repo", "some-org/some-repo",
						"--token", "some-github-token",
					)
					command.Env = append(command.Env, fmt.Sprintf("GITHUB_EVENT_PATH=%s", eventPath))
					buffer := gbytes.NewBuffer()

					session, err := gexec.Start(command, buffer, buffer)
					Expect(err).NotTo(HaveOccurred())

					Eventually(session).Should(gexec.Exit(1), fmt.Sprintf("output:\n%s\n", buffer.Contents()))

					Expect(buffer).To(gbytes.Say(`Dispatching`))
					Expect(buffer).To(gbytes.Say(`Error: failed to decode \$GITHUB_EVENT_PATH:`))
					Expect(buffer).To(gbytes.Say(`invalid character`))
				})
			})

			context("when the dispatch request cannot be created", func() {
				it("prints an error message and exits non-zero", func() {
					command := exec.Command(
						entrypoint,
						"--endpoint", "%%%",
						"--repo", "some-org/some-repo",
						"--token", "some-github-token",
					)
					command.Env = append(command.Env, fmt.Sprintf("GITHUB_EVENT_PATH=%s", eventPath))
					buffer := gbytes.NewBuffer()

					session, err := gexec.Start(command, buffer, buffer)
					Expect(err).NotTo(HaveOccurred())

					Eventually(session).Should(gexec.Exit(1), fmt.Sprintf("output:\n%s\n", buffer.Contents()))

					Expect(buffer).To(gbytes.Say(`Dispatching`))
					Expect(buffer).To(gbytes.Say(`Error: failed to create dispatch request`))
					Expect(buffer).To(gbytes.Say(`invalid URL escape`))
				})
			})

			context("when the dispatch request cannot be completed", func() {
				it("prints an error message and exits non-zero", func() {
					command := exec.Command(
						entrypoint,
						"--endpoint", api.URL,
						"--repo", "loop-org/loop-repo",
						"--token", "some-github-token",
					)
					command.Env = append(command.Env, fmt.Sprintf("GITHUB_EVENT_PATH=%s", eventPath))
					buffer := gbytes.NewBuffer()

					session, err := gexec.Start(command, buffer, buffer)
					Expect(err).NotTo(HaveOccurred())

					Eventually(session).Should(gexec.Exit(1), fmt.Sprintf("output:\n%s\n", buffer.Contents()))

					Expect(buffer).To(gbytes.Say(`Dispatching`))
					Expect(buffer).To(gbytes.Say(`Error: failed to complete dispatch request`))
					Expect(buffer).To(gbytes.Say(`stopped after 10 redirects`))
				})
			})

			context("when the dispatch request response is not success", func() {
				it("prints an error message and exits non-zero", func() {
					command := exec.Command(
						entrypoint,
						"--endpoint", api.URL,
						"--repo", "fail-org/fail-repo",
						"--token", "some-github-token",
					)
					command.Env = append(command.Env, fmt.Sprintf("GITHUB_EVENT_PATH=%s", eventPath))
					buffer := gbytes.NewBuffer()

					session, err := gexec.Start(command, buffer, buffer)
					Expect(err).NotTo(HaveOccurred())

					Eventually(session).Should(gexec.Exit(1), fmt.Sprintf("output:\n%s\n", buffer.Contents()))

					Expect(buffer).To(gbytes.Say(`Dispatching`))
					Expect(buffer).To(gbytes.Say(`Error: unexpected response from dispatch request`))
					Expect(buffer).To(gbytes.Say(`500 Internal Server Error`))
					Expect(buffer).To(gbytes.Say(`{"error": "server-error"}`))
				})
			})
		})
	})
}
