package rdy_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRdy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rdy Suite")
}
