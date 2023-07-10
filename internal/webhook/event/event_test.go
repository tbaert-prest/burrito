package event_test

import (
	"context"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	configv1alpha1 "github.com/padok-team/burrito/api/v1alpha1"
	"github.com/padok-team/burrito/internal/annotations"
	utils "github.com/padok-team/burrito/internal/testing"
	"github.com/padok-team/burrito/internal/webhook/event"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment

const testTime = "Sun May  8 11:21:53 UTC 2023"

func TestLayer(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Webhook Handler Suite")
}

// type MockClock struct{}

// func (m *MockClock) Now() time.Time {
// 	t, _ := time.Parse(time.UnixDate, testTime)
// 	return t
// }

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("../../..", "manifests", "crds")},
		ErrorIfCRDPathMissing: true,
	}
	var err error
	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = configv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	utils.LoadResources(k8sClient, "testdata")
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())
})

var PushEventNoChanges = event.PushEvent{
	URL:      "https://github.com/padok-team/burrito-examples",
	Revision: "main",
	ChangeInfo: event.ChangeInfo{
		ShaBefore: "b3231e8771591b3864b3c582e85955c1f76aaded",
		ShaAfter:  "6c193d9cad1ddafdb31ff9f733630da9705bfd64",
	},
	Changes: []string{
		"README.md",
	},
}

var PushEventLayerPathChanges = event.PushEvent{
	URL:      "https://github.com/padok-team/burrito-examples",
	Revision: "main",
	ChangeInfo: event.ChangeInfo{
		ShaBefore: "b3231e8771591b3864b3c582e85955c1f76aaded",
		ShaAfter:  "6c193d9cad1ddafdb31ff9f733630da9705bfd64",
	},
	Changes: []string{
		"layer-path-changed/main.tf",
	},
}

var PushEventAdditionalPathChanges = event.PushEvent{
	URL:      "https://github.com/padok-team/burrito-examples",
	Revision: "main",
	ChangeInfo: event.ChangeInfo{
		ShaBefore: "b3231e8771591b3864b3c582e85955c1f76aaded",
		ShaAfter:  "6c193d9cad1ddafdb31ff9f733630da9705bfd64",
	},
	Changes: []string{
		"modules/module-changed/variables.tf",
		"terragrunt/layer-path-changed/module.hcl",
	},
}

var PushEventMultiplePathChanges = event.PushEvent{
	URL:      "https://github.com/padok-team/burrito-examples",
	Revision: "main",
	ChangeInfo: event.ChangeInfo{
		ShaBefore: "b3231e8771591b3864b3c582e85955c1f76aaded",
		ShaAfter:  "6c193d9cad1ddafdb31ff9f733630da9705bfd64",
	},
	Changes: []string{
		"modules/random-pets/variables.tf",
		"nominal-case-one/main.tf",
		"terragrunt/nominal-case-two/prod/inputs.hcl",
		"terragrunt/nominal-case-two/module.hcl",
	},
}

var PullRequestEventNotAffected = event.PullRequestEvent{
	Provider: "github",
	URL:      "https://github.com/example/repo",
	Revision: "feature/branch",
	Base:     "main",
	Action:   "opened",
	ID:       "42",
	Commit:   "5b2c5e5c6699bf2bf93138205565b85193996572",
}

var PullRequestEventSingleAffected = event.PullRequestEvent{
	Provider: "github",
	URL:      "https://github.com/padok-team/burrito-examples",
	Revision: "feature/branch",
	Base:     "main",
	Action:   "opened",
	ID:       "42",
	Commit:   "5b2c5e5c6699bf2bf93138205565b85193996572",
}

var PullRequestEventMultipleAffected = event.PullRequestEvent{
	Provider: "github",
	URL:      "https://github.com/example/other-repo",
	Revision: "feature/branch",
	Base:     "main",
	Action:   "opened",
	ID:       "42",
	Commit:   "5b2c5e5c6699bf2bf93138205565b85193996572",
}

var _ = Describe("Webhook", func() {
	var handleErr error
	Describe("Push Event", func() {
		Describe("Layer", func() {
			Describe("No paths are relevant to layer", Ordered, func() {
				BeforeAll(func() {
					handleErr = PushEventNoChanges.Handle(k8sClient)
				})
				It("should have only set the LastBranchCommit annotation", func() {
					layer := &configv1alpha1.TerraformLayer{}
					err := k8sClient.Get(context.TODO(), types.NamespacedName{
						Namespace: "default",
						Name:      "no-path-changed-1",
					}, layer)
					Expect(err).NotTo(HaveOccurred())
					Expect(handleErr).NotTo(HaveOccurred())
					//TODO: Maybe implement a test on a layer that already has the annotation
					_, ok := layer.Annotations[annotations.LastRelevantCommit]
					Expect(ok).To(BeFalse())
					Expect(layer.Annotations[annotations.LastBranchCommit]).To(Equal(PushEventNoChanges.ChangeInfo.ShaAfter))
				})
			})
			Describe("Layer path has been modified", Ordered, func() {
				BeforeAll(func() {
					handleErr = PushEventLayerPathChanges.Handle(k8sClient)
				})
				It("should have updated the LastBranchCommit and LastRelevantCommit annotations", func() {
					layer := &configv1alpha1.TerraformLayer{}
					err := k8sClient.Get(context.TODO(), types.NamespacedName{
						Namespace: "default",
						Name:      "layer-path-changed-1",
					}, layer)
					Expect(err).NotTo(HaveOccurred())
					Expect(handleErr).NotTo(HaveOccurred())
					Expect(layer.Annotations[annotations.LastBranchCommit]).To(Equal(PushEventLayerPathChanges.ChangeInfo.ShaAfter))
					Expect(layer.Annotations[annotations.LastRelevantCommit]).To(Equal(PushEventLayerPathChanges.ChangeInfo.ShaAfter))
				})
			})
			Describe("Additional path has been modified", Ordered, func() {
				BeforeAll(func() {
					handleErr = PushEventAdditionalPathChanges.Handle(k8sClient)
				})
				It("should have updated commit annotations for a absolute change path", func() {
					layer := &configv1alpha1.TerraformLayer{}
					err := k8sClient.Get(context.TODO(), types.NamespacedName{
						Namespace: "default",
						Name:      "layer-additional-paths-1",
					}, layer)
					Expect(err).NotTo(HaveOccurred())
					Expect(handleErr).NotTo(HaveOccurred())
					Expect(layer.Annotations[annotations.LastBranchCommit]).To(Equal(PushEventLayerPathChanges.ChangeInfo.ShaAfter))
					Expect(layer.Annotations[annotations.LastRelevantCommit]).To(Equal(PushEventLayerPathChanges.ChangeInfo.ShaAfter))
				})
				// TODO: make this test pass
				// It("should have updated commit annotations for a relative change path", func() {
				// 	layer := &configv1alpha1.TerraformLayer{}
				// 	err := k8sClient.Get(context.TODO(), types.NamespacedName{
				// 		Namespace: "default",
				// 		Name:      "layer-additional-paths-2",
				// 	}, layer)
				// 	Expect(err).NotTo(HaveOccurred())
				// 	Expect(handleErr).NotTo(HaveOccurred())
				// 	Expect(layer.Annotations[annotations.LastBranchCommit]).To(Equal(PushEventLayerPathChanges.ChangeInfo.ShaAfter))
				// 	Expect(layer.Annotations[annotations.LastRelevantCommit]).To(Equal(PushEventLayerPathChanges.ChangeInfo.ShaAfter))
				// })
			})
			Describe("Multiple paths have been modified", Ordered, func() {

			})
		})
		Describe("PullRequest", func() {
			Describe("A single pull request have been affected", Ordered, func() {

			})
			Describe("Multiple pull request have been affected", Ordered, func() {

			})
		})

		// BeforeAll(func() {

		// 	layer = &configv1alpha1.TerraformLayer{}
		// 	getErr = k8sClient.Get(context.TODO(), types.NamespacedName{
		// 		Namespace: "default",
		// 		Name:      "test",
		// 	}, layer)
		// })
		// It("should exists", func() {
		// 	Expect(getErr).NotTo(HaveOccurred())
		// })
		// It("should not return an error when adding first annotation", func() {
		// 	err := annotations.Add(context.TODO(), k8sClient, layer, map[string]string{annotations.LastPlanSum: "AuP6pMNxWsbSZKnxZvxD842wy0qaF9JCX8HW1nFeL1I"})
		// 	Expect(err).NotTo(HaveOccurred())
		// })
		// It("should not return an error when adding second annotation", func() {
		// 	err := annotations.Add(context.TODO(), k8sClient, layer, map[string]string{annotations.LastApplySum: "AuP6pMNxWsbSZKnxZvxD842wy0qaF9JCX8HW1nFeL1I"})
		// 	Expect(err).NotTo(HaveOccurred())
		// })
		// It("should not return an error when removing second annotation", func() {
		// 	err := annotations.Remove(context.TODO(), k8sClient, layer, annotations.LastApplySum)
		// 	Expect(err).NotTo(HaveOccurred())
		// })
	})
	// Describe("Pull Request Event", Ordered, func() {
	// 	Describe("No Pull request have been created")
	// 	Describe("A single pull request have been created")
	// 	Describe("Multiple pull requests have been created")
	// })
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
