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
	f.Shoot = e2e.DefaultShoot("e2e-cache-def")

	It("should create Shoot", func(ctx SpecContext) {
		Expect(f.CreateShootAndWaitForCreation(ctx, false)).To(Succeed())
		f.Verify()
	}, SpecTimeout(15*time.Minute))

	It("should enable the registry-cache extension", func(ctx SpecContext) {
		size := resource.MustParse("2Gi")
		Expect(f.UpdateShoot(ctx, f.Shoot, func(shoot *gardencorev1beta1.Shoot) error {
			common.AddOrUpdateRegistryCacheExtension(shoot, []v1alpha3.RegistryCache{
				{Upstream: "ghcr.io", Volume: &v1alpha3.Volume{Size: &size}},
			})

			return nil
		})).To(Succeed())
	}, SpecTimeout(10*time.Minute))

	It("[ghcr.io] should verify registry-cache works", func(ctx SpecContext) {
		common.VerifyRegistryCache(ctx, f.Logger, f.ShootFramework.ShootClient, common.GithubRegistryJitesoftAlpine3188Image, common.AlpinePodMutateFn)
	}, SpecTimeout(12*time.Minute))

	It("should add the registry.gitlab.com upstream to the registry-cache extension", func(ctx SpecContext) {
		size := resource.MustParse("2Gi")
		Expect(f.UpdateShoot(ctx, f.Shoot, func(shoot *gardencorev1beta1.Shoot) error {
			common.AddOrUpdateRegistryCacheExtension(shoot, []v1alpha3.RegistryCache{
				{Upstream: "ghcr.io", Volume: &v1alpha3.Volume{Size: &size}},
				{Upstream: "registry.gitlab.com", Volume: &v1alpha3.Volume{Size: &size}},
			})

			return nil
		})).To(Succeed())
	}, SpecTimeout(10*time.Minute))

	It("[registry.gitlab.com] should verify registry-cache works", func(ctx SpecContext) {
		common.VerifyRegistryCache(ctx, f.Logger, f.ShootFramework.ShootClient, common.GitlabRegistryJitesoftAlpine31710Image, common.AlpinePodMutateFn)
	}, SpecTimeout(12*time.Minute))

	It("should remove the registry.gitlab.com upstream from the registry-cache extension", func(ctx SpecContext) {
		size := resource.MustParse("2Gi")
		Expect(f.UpdateShoot(ctx, f.Shoot, func(shoot *gardencorev1beta1.Shoot) error {
			common.AddOrUpdateRegistryCacheExtension(shoot, []v1alpha3.RegistryCache{
				{Upstream: "ghcr.io", Volume: &v1alpha3.Volume{Size: &size}},
			})

			return nil
		})).To(Succeed())
	}, SpecTimeout(10*time.Minute))

	It("[registry.gitlab.com] should verify registry configuration is removed", func(ctx SpecContext) {
		common.VerifyHostsTOMLFilesDeletedForAllNodes(ctx, f.Logger, f.ShootFramework.ShootClient, []string{"registry.gitlab.com"})
	}, SpecTimeout(2*time.Minute))

	It("should disable the registry-cache extension", func(ctx SpecContext) {
		Expect(f.UpdateShoot(ctx, f.Shoot, func(shoot *gardencorev1beta1.Shoot) error {
			common.RemoveExtension(shoot, "registry-cache")

			return nil
		})).To(Succeed())
	}, SpecTimeout(10*time.Minute))

	It("[ghcr.io] should verify registry configuration is removed", func(ctx SpecContext) {
		common.VerifyHostsTOMLFilesDeletedForAllNodes(ctx, f.Logger, f.ShootFramework.ShootClient, []string{"ghcr.io"})
	}, SpecTimeout(2*time.Minute))

	It("should delete Shoot", func(ctx SpecContext) {
		Expect(f.DeleteShootAndWaitForDeletion(ctx, f.Shoot)).To(Succeed())
	}, SpecTimeout(15*time.Minute))
})
