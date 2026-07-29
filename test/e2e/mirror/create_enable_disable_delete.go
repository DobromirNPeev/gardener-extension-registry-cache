// SPDX-FileCopyrightText: Contributors to the Gardener project
//
// SPDX-License-Identifier: Apache-2.0

package mirror

import (
	"time"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gardener/gardener-extension-registry-cache/pkg/apis/mirror/v1alpha1"
	"github.com/gardener/gardener-extension-registry-cache/test/common"
	"github.com/gardener/gardener-extension-registry-cache/test/e2e"
)

var _ = Describe("Registry Mirror Extension Tests", Label("mirror"), Ordered, func() {
	f := e2e.DefaultShootCreationFramework()
	f.Shoot = e2e.DefaultShoot("e2e-mirror-def")

	BeforeAll(func() {
		DeferCleanup(func(ctx SpecContext) {
			Expect(e2e.DeleteShootIfExists(ctx, f)).To(Succeed())
		}, NodeTimeout(15*time.Minute))
	})

	It("should create Shoot", func(ctx SpecContext) {
		Expect(f.CreateShootAndWaitForCreation(ctx, false)).To(Succeed())
		f.Verify()
	}, SpecTimeout(15*time.Minute))

	It("should enable the registry-mirror extension", func(ctx SpecContext) {
		Expect(f.UpdateShoot(ctx, f.Shoot, func(shoot *gardencorev1beta1.Shoot) error {
			common.AddOrUpdateRegistryMirrorExtension(shoot, []v1alpha1.MirrorConfiguration{
				{
					Upstream: "docker.io",
					Hosts: []v1alpha1.MirrorHost{
						{Host: "https://mirror.gcr.io"},
					},
				},
				{
					Upstream: "public.ecr.aws",
					Hosts: []v1alpha1.MirrorHost{
						{Host: "https://public-mirror.example.com"},
						{Host: "https://private-mirror.internal", Capabilities: []v1alpha1.MirrorHostCapability{v1alpha1.MirrorHostCapabilityPull, v1alpha1.MirrorHostCapabilityResolve}},
						{Host: "https://harbor.example.com/v2/k8s", OverridePath: new(true)},
					},
				},
			})

			return nil
		})).To(Succeed())
	}, SpecTimeout(10*time.Minute))

	It("should verify registry mirror configuration is applied", func(ctx SpecContext) {
		common.VerifyHostsTOMLFilesCreatedForAllNodes(ctx, f.Logger, f.ShootFramework.ShootClient, map[string]string{
			"docker.io":      dockerHostsTOML,
			"public.ecr.aws": ecrHostsTOML,
		})
	}, SpecTimeout(1*time.Minute))

	It("should disable the registry-mirror extension", func(ctx SpecContext) {
		Expect(f.UpdateShoot(ctx, f.Shoot, func(shoot *gardencorev1beta1.Shoot) error {
			common.RemoveExtension(shoot, "registry-mirror")

			return nil
		})).To(Succeed())
	}, SpecTimeout(10*time.Minute))

	It("should verify registry mirror configuration is removed", func(ctx SpecContext) {
		common.VerifyHostsTOMLFilesDeletedForAllNodes(ctx, f.Logger, f.ShootFramework.ShootClient, []string{"docker.io"})
	}, SpecTimeout(1*time.Minute))

	It("should delete Shoot", func(ctx SpecContext) {
		Expect(f.DeleteShootAndWaitForDeletion(ctx, f.Shoot)).To(Succeed())
	}, SpecTimeout(15*time.Minute))
})

const (
	dockerHostsTOML = `# managed by gardener-node-agent
server = "https://registry-1.docker.io"

[host."https://mirror.gcr.io"]
  capabilities = ["pull"]

`

	ecrHostsTOML = `# managed by gardener-node-agent
server = "https://public.ecr.aws"

[host."https://public-mirror.example.com"]
  capabilities = ["pull"]

[host."https://private-mirror.internal"]
  capabilities = ["pull","resolve"]

[host."https://harbor.example.com/v2/k8s"]
  capabilities = ["pull"]
  override_path = true

`
)
