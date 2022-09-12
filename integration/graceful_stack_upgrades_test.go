package integration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testStackUpgrades(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		pack   occam.Pack
		docker occam.Docker
	)

	it.Before(func() {
		pack = occam.NewPack()
		docker = occam.NewDocker()

		Expect(docker.Pull.Execute("paketobuildpacks/builder-jammy-buildpackless-base")).To(Succeed())
		Expect(docker.Pull.Execute("paketobuildpacks/run-jammy-base")).To(Succeed())

	})

	context("when building a node app that does not use a package manager and stacks change between builds", func() {
		var (
			image occam.Image

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
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		it("successfully builds an OCI image", func() {
			var err error
			var logs fmt.Stringer
			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(nodeBuildpack).
				WithPullPolicy("never").
				WithBuilder("paketobuildpacks/builder:buildpackless-base").
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(nodeBuildpack).
				WithPullPolicy("never").
				WithBuilder("paketobuildpacks/builder-jammy-buildpackless-base").
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

		})
	})

	context("when building a node app that uses npm and stacks change between builds", func() {
		var (
			image occam.Image

			name   string
			source string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())

			// Using vendored app to force layer reuse
			source, err = occam.Source(filepath.Join("testdata", "vendored"))
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		it("successfully builds an OCI image", func() {
			var err error
			var logs fmt.Stringer
			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(nodeBuildpack).
				WithPullPolicy("never").
				WithBuilder("paketobuildpacks/builder:buildpackless-base").
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(nodeBuildpack).
				WithPullPolicy("never").
				WithBuilder("paketobuildpacks/builder-jammy-buildpackless-base").
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

		})
	})

	context("when building a node app that uses yarn and stacks change between builds", func() {
		var (
			image occam.Image

			name   string
			source string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())

			source, err = occam.Source(filepath.Join("testdata", "yarn"))
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		it("successfully builds an OCI image", func() {
			var err error
			var logs fmt.Stringer
			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(nodeBuildpack).
				WithPullPolicy("never").
				WithBuilder("paketobuildpacks/builder:buildpackless-base").
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(nodeBuildpack).
				WithPullPolicy("never").
				WithBuilder("paketobuildpacks/builder-jammy-buildpackless-base").
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

		})
	})
}
