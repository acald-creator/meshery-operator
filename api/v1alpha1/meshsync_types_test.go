package v1alpha1

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

/*
Here we do not want test CURD of kubernetes service, we want to confirm that:
* our CRDs resouces group version could be support by specific kubernetes server
* the specific value of the CRDs field can be workwell without pruning by the api-server to etcd
*/
var _ = Describe("The test case for the meshsync CRDs", func() {

	context := context.Background()

	// Define the constants used in the tests
	const (
		URL          string = "https://layer5.io"
		str          string = "healthy"
		Reason       string = "Testcase"
		Message      string = "Message for testcase"
		PublishingTo string = "Publish for testcase"
		FileManager  string = "testcase-meshsync"
		Kind         string = "MeshSync"
		APIVersion   string = "meshery.layer5.io/v1alpha1"
	)

	// Define the namespaced name for the test object
	typeNamespace := types.NamespacedName{
		Namespace: "default",
		Name:      "default",
	}

	// Define the test object
	meshSync := &MeshSync{
		TypeMeta: metav1.TypeMeta{
			APIVersion: APIVersion,
			Kind:       Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "default",
		},
		Spec: MeshSyncSpec{
			Size: 2,
			Broker: MeshsyncBroker{
				Custom: CustomMeshsyncBroker{
					URL: URL,
				},
				Native: NativeMeshsyncBroker{
					Namespace: "default",
					Name:      "default",
				},
			},
		},
	}

	// Create the meshSync object before each test
	BeforeEach(func() {
		err := fakeClient.Create(context, meshSync)
		Expect(err).NotTo(HaveOccurred())
	})

	// Delete the meshSync object after each test
	AfterEach(func() {
		err := fakeClient.Delete(context, meshSync)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("The CURD case for the meshsync CRDs", func() {

		It("The meshsync CRDs create acticity should be succeed", func() {
			By("Get the meshsync CRDs")
			mesheSyncGet := &MeshSync{}
			err := fakeClient.Get(context, typeNamespace, mesheSyncGet)
			Expect(err).NotTo(HaveOccurred())

			By("Confirm the URL equal to https://layer5.io")
			url := mesheSyncGet.Spec.Broker.Custom.URL
			Expect(url == URL).Should(BeTrue())
		})

		It("The meshsync CRDs update the spec of the resources", func() {
			By("Update the size of the meshsync CRDs")
			meshSync.Spec.Size = 5
			err := fakeClient.Update(context, meshSync, &client.UpdateOptions{FieldManager: "testcase-meshsync"})
			Expect(err).NotTo(HaveOccurred())

			By("Get the latest version of meshsync CRDs")
			meshSyncGet := &MeshSync{}
			err = fakeClient.Get(context, typeNamespace, meshSyncGet)
			Expect(err).NotTo(HaveOccurred())

			By("Confirm the size equal to 5")
			Expect(meshSyncGet.Spec.Size == 5).Should(BeTrue())
		})

		It("The meshsync CRDs update the status of the resources", func() {
			meshSync.Status = MeshSyncStatus{
				PublishingTo: PublishingTo,
				Conditions: []Condition{
					{
						Type:               ConditionType(str),
						Status:             ConditionStatus(str),
						ObservedGeneration: 1,
						Reason:             Reason,
						Message:            Message,
						LastTransitionTime: metav1.Now(),
						LastProbeTime:      metav1.Now(),
					},
				},
			}
			By("Update the status of the meshsync CRDs")
			err := fakeClient.Status().Update(context, meshSync, &client.SubResourceUpdateOptions{
				UpdateOptions: client.UpdateOptions{FieldManager: FileManager},
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(meshSync.Status.PublishingTo == PublishingTo).Should(BeTrue())
		})

		It("The meshsyncList CRDs should be support by the kubernetes server", func() {
			By("Confirm the meshsync CRDs support by the kubernetes server")
			meshSyncList := &MeshSyncList{}
			err := fakeClient.List(context, meshSyncList, &client.ListOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(meshSyncList.Items) > 0).Should(BeTrue())
		})

	})

	// Test coverage for delete the meshsync CRDs
	Context("The test coverage for delete the meshsync CRDs", func() {
		It("The meshsync CRDs remove the resources ", func() {
			By("Delete the meshsync CRDs")
			err := fakeClient.Delete(context, meshSync)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
