/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package machineset

import (
	"reflect"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	machinev1 "github.com/openshift/api/machine/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ reconcile.Reconciler = &ReconcileMachineSet{}

func TestMachineSetToMachines(t *testing.T) {
	machineSetList := &machinev1.MachineSetList{
		TypeMeta: metav1.TypeMeta{
			Kind: "MachineSetList",
		},
		Items: []machinev1.MachineSet{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "withMatchingLabels",
					Namespace: "test",
				},
				Spec: machinev1.MachineSetSpec{
					Selector: metav1.LabelSelector{
						MatchLabels: map[string]string{
							"foo":                             "bar",
							machinev1.MachineClusterLabelName: "test-cluster",
						},
					},
				},
			},
		},
	}
	controller := true
	m := machinev1.Machine{
		TypeMeta: metav1.TypeMeta{
			Kind: "Machine",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "withOwnerRef",
			Namespace: "test",
			Labels: map[string]string{
				machinev1.MachineClusterLabelName: "test-cluster",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					Name:       "Owner",
					Kind:       "MachineSet",
					Controller: &controller,
				},
			},
		},
	}
	m2 := machinev1.Machine{
		TypeMeta: metav1.TypeMeta{
			Kind: "Machine",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "noOwnerRefNoLabels",
			Namespace: "test",
			Labels: map[string]string{
				machinev1.MachineClusterLabelName: "test-cluster",
			},
		},
	}
	m3 := machinev1.Machine{
		TypeMeta: metav1.TypeMeta{
			Kind: "Machine",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "withMatchingLabels",
			Namespace: "test",
			Labels: map[string]string{
				"foo":                             "bar",
				machinev1.MachineClusterLabelName: "test-cluster",
			},
		},
	}
	testsCases := []struct {
		machine  machinev1.Machine
		object   client.Object
		expected []reconcile.Request
	}{
		{
			machine:  m,
			object:   &m,
			expected: []reconcile.Request{},
		},
		{
			machine:  m2,
			object:   &m2,
			expected: nil,
		},
		{
			machine: m3,
			object:  &m3,
			expected: []reconcile.Request{
				{NamespacedName: client.ObjectKey{Namespace: "test", Name: "withMatchingLabels"}},
			},
		},
	}

	machinev1.AddToScheme(scheme.Scheme)
	r := &ReconcileMachineSet{
		Client: fake.NewFakeClient(&m, &m2, &m3, machineSetList),
		scheme: scheme.Scheme,
	}

	for _, tc := range testsCases {
		got := r.MachineToMachineSets(tc.object)
		if !reflect.DeepEqual(got, tc.expected) {
			t.Errorf("Case %s. Got: %v, expected: %v", tc.machine.Name, got, tc.expected)
		}
	}
}

func TestShouldExcludeMachine(t *testing.T) {
	controller := true
	testCases := []struct {
		machineSet machinev1.MachineSet
		machine    machinev1.Machine
		expected   bool
	}{
		{
			machineSet: machinev1.MachineSet{},
			machine: machinev1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "withNoMatchingOwnerRef",
					Namespace: "test",
					OwnerReferences: []metav1.OwnerReference{
						{
							Name:       "Owner",
							Kind:       "MachineSet",
							Controller: &controller,
						},
					},
				},
			},
			expected: true,
		},
		{
			machineSet: machinev1.MachineSet{
				Spec: machinev1.MachineSetSpec{
					Selector: metav1.LabelSelector{
						MatchLabels: map[string]string{
							"foo": "bar",
						},
					},
				},
			},
			machine: machinev1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "withMatchingLabels",
					Namespace: "test",
					Labels: map[string]string{
						"foo": "bar",
					},
				},
			},
			expected: false,
		},
		{
			machineSet: machinev1.MachineSet{},
			machine: machinev1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "withDeletionTimestamp",
					Namespace:         "test",
					DeletionTimestamp: &metav1.Time{Time: time.Now()},
					Labels: map[string]string{
						"foo": "bar",
					},
				},
			},
			expected: true,
		},
	}

	for _, tc := range testCases {
		got := shouldExcludeMachine(&tc.machineSet, &tc.machine)
		if got != tc.expected {
			t.Errorf("Case %s. Got: %v, expected: %v", tc.machine.Name, got, tc.expected)
		}
	}
}

func TestAdoptOrphan(t *testing.T) {
	m := machinev1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name: "orphanMachine",
		},
	}
	ms := machinev1.MachineSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: "adoptOrphanMachine",
		},
	}
	controller := true
	blockOwnerDeletion := true
	testCases := []struct {
		machineSet machinev1.MachineSet
		machine    machinev1.Machine
		expected   []metav1.OwnerReference
	}{
		{
			machine:    m,
			machineSet: ms,
			expected: []metav1.OwnerReference{
				{
					APIVersion:         machinev1.SchemeGroupVersion.String(),
					Kind:               "MachineSet",
					Name:               "adoptOrphanMachine",
					UID:                "",
					Controller:         &controller,
					BlockOwnerDeletion: &blockOwnerDeletion,
				},
			},
		},
	}

	machinev1.AddToScheme(scheme.Scheme)
	r := &ReconcileMachineSet{
		Client: fake.NewFakeClient(&m),
		scheme: scheme.Scheme,
	}
	for _, tc := range testCases {
		r.adoptOrphan(&tc.machineSet, &tc.machine)
		got := tc.machine.GetOwnerReferences()
		if !reflect.DeepEqual(got, tc.expected) {
			t.Errorf("Case %s. Got: %+v, expected: %+v", tc.machine.Name, got, tc.expected)
		}
	}
}

var _ = Describe("MachineSet Reconcile", func() {
	var r *ReconcileMachineSet
	var result reconcile.Result
	var reconcileErr error
	var rec *record.FakeRecorder

	BeforeEach(func() {
		Expect(machinev1.AddToScheme(scheme.Scheme)).To(Succeed())
		rec = record.NewFakeRecorder(32)

		r = &ReconcileMachineSet{
			scheme:   scheme.Scheme,
			recorder: rec,
		}
	})

	JustBeforeEach(func() {
		request := reconcile.Request{NamespacedName: types.NamespacedName{Name: "machineset1", Namespace: "default"}}
		result, reconcileErr = r.Reconcile(ctx, request)
	})

	Context("ignore machine sets marked for deletion", func() {
		BeforeEach(func() {
			dt := metav1.Now()

			ms := &machinev1.MachineSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "machineset1",
					Namespace:         "default",
					DeletionTimestamp: &dt,
				},
				Spec: machinev1.MachineSetSpec{
					Template: machinev1.MachineTemplateSpec{},
				}}

			r.Client = fake.NewFakeClientWithScheme(scheme.Scheme, ms)
		})

		It("returns an empty result", func() {
			Expect(result).To(Equal(reconcile.Result{}))
		})

		It("does not return an error", func() {
			Expect(reconcileErr).NotTo(HaveOccurred())
		})
	})

	Context("record event if reconcile fails", func() {
		BeforeEach(func() {
			var replicas int32
			ms := &machinev1.MachineSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "machineset1",
					Namespace: "default",
				},
				Spec: machinev1.MachineSetSpec{
					Replicas: &replicas,
				},
			}

			ms.Spec.Selector.MatchLabels = map[string]string{
				"--$-invalid": "true",
			}

			r.Client = fake.NewFakeClientWithScheme(scheme.Scheme, ms)
		})

		It("did something with events", func() {
			Eventually(rec.Events).Should(Receive())
		})
	})
})
