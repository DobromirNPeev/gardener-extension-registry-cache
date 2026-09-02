// SPDX-FileCopyrightText: Contributors to the Gardener project
//
// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"time"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/gardener/gardener-extension-registry-cache/pkg/apis/registry/v1alpha3"
	"github.com/gardener/gardener-extension-registry-cache/test/common"
	"github.com/gardener/gardener-extension-registry-cache/test/e2e"
)

var _ = Describe("Registry Cache Extension Tests", Label("cache"), Ordered, func() {
	f := e2e.DefaultShootCreationFramework()
	shoot := e2e.DefaultShoot("e2e-cache-tls")
	size := resource.MustParse("2Gi")
	common.AddOrUpdateRegistryCacheExtension(shoot, []v1alpha3.RegistryCache{
		{Upstream: "ghcr.io", Volume: &v1alpha3.Volume{Size: &size}, HTTP: &v1alpha3.HTTP{TLS: true}},
	})
	f.Shoot = shoot

	It("should create Shoot", func(ctx SpecContext) {
		Expect(f.CreateShootAndWaitForCreation(ctx, false)).To(Succeed())
		f.Verify()
	}, SpecTimeout(15*time.Minute))

	It("should disable TLS", func(ctx SpecContext) {
		Expect(f.UpdateShoot(ctx, f.Shoot, func(shoot *gardencorev1beta1.Shoot) error {
			common.AddOrUpdateRegistryCacheExtension(shoot, []v1alpha3.RegistryCache{
				{Upstream: "ghcr.io", Volume: &v1alpha3.Volume{Size: &size}, HTTP: &v1alpha3.HTTP{TLS: false}},
			})

			return nil
		})).To(Succeed())
	}, SpecTimeout(10*time.Minute))

	It("should verify registry-cache works with TLS disabled", func(ctx SpecContext) {
		common.VerifyRegistryCache(ctx, f.Logger, f.ShootFramework.ShootClient, common.GithubRegistryJitesoftAlpine3188Image, common.AlpinePodMutateFn)
	}, SpecTimeout(12*time.Minute))

	It("should enable TLS", func(ctx SpecContext) {
		Expect(f.UpdateShoot(ctx, f.Shoot, func(shoot *gardencorev1beta1.Shoot) error {
			common.AddOrUpdateRegistryCacheExtension(shoot, []v1alpha3.RegistryCache{
				{Upstream: "ghcr.io", Volume: &v1alpha3.Volume{Size: &size}, HTTP: &v1alpha3.HTTP{TLS: true}},
			})

			return nil
		})).To(Succeed())
	}, SpecTimeout(10*time.Minute))

	It("should verify registry-cache works with TLS enabled", func(ctx SpecContext) {
		// We are using ghcr.io/jitesoft/alpine:3.19.4 as ghcr.io/jitesoft/alpine:3.18.8 is already used in the test.
		// Hence, ghcr.io/jitesoft/alpine:3.18.8 will be present in the Node.
		common.VerifyRegistryCache(ctx, f.Logger, f.ShootFramework.ShootClient, common.GithubRegistryJitesoftAlpine3194Image, common.AlpinePodMutateFn)
	}, SpecTimeout(12*time.Minute))

	It("should delete Shoot", func(ctx SpecContext) {
		Expect(f.DeleteShootAndWaitForDeletion(ctx, f.Shoot)).To(Succeed())
	}, SpecTimeout(15*time.Minute))
})
