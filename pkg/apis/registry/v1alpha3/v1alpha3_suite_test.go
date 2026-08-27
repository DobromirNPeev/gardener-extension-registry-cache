// SPDX-FileCopyrightText: Contributors to the Gardener project
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha3_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestV1alpha3(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "v1alpha3 Suite")
}
