package main_test

import (
	"testing"
	"time"

	"github.com/onsi/gomega/gexec"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
)

var entrypoint string

func TestEntrypoint(t *testing.T) {
	var Expect = NewWithT(t).Expect

	SetDefaultEventuallyTimeout(5 * time.Second)

	var err error
	entrypoint, err = gexec.Build("github.com/thitch97/nodejs/actions/dispatch/entrypoint")
	Expect(err).NotTo(HaveOccurred())

	suite := spec.New("entrypoint", spec.Report(report.Terminal{}), spec.Parallel())
	suite("SendDispatch", testSendDispatch)
	suite.Run(t)
}
