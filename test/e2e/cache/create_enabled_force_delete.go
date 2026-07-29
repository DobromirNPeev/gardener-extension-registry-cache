// SPDX-FileCopyrightText: Contributors to the Gardener project
//
// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/gardener/gardener-extension-registry-cache/pkg/apis/registry/v1alpha3"
	"github.com/gardener/gardener-extension-registry-cache/test/common"
	"github.com/gardener/gardener-extension-registry-cache/test/e2e"
)

var _ = Describe("Registry Cache Extension Tests", Label("cache"), Ordered, func() {
	f := e2e.DefaultShootCreationFramework()
	shoot := e2e.DefaultShoot("e2e-cache-fd")
	size := resource.MustParse("2Gi")
	common.AddOrUpdateRegistryCacheExtension(shoot, []v1alpha3.RegistryCache{
		{
			Upstream:          "gardener-extension-registry-cache-tests.common.repositories.cloud.sap",
			Volume:            &v1alpha3.Volume{Size: &size},
			ServiceNameSuffix: new("gardener-extension-registry-cache-tests"),
		},
	})
	f.Shoot = shoot

	It("should create Shoot", func(ctx SpecContext) {
		Expect(f.CreateShootAndWaitForCreation(ctx, false)).To(Succeed())
		f.Verify()
	}, SpecTimeout(15*time.Minute))

	It("should force delete Shoot", func(ctx SpecContext) {
		Expect(f.ForceDeleteShootAndWaitForDeletion(ctx, f.Shoot)).To(Succeed())
	}, SpecTimeout(10*time.Minute))
})
