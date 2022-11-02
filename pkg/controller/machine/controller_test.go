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

package machine

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	machinev1 "github.com/uccps-samples/machine-api-operator/pkg/apis/machine/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	_ reconcile.Reconciler = &ReconcileMachine{}
)

func TestReconcileRequest(t *testing.T) {
	machineProvisioning := machinev1.Machine{
		TypeMeta: metav1.TypeMeta{
			Kind: "Machine",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:       "create",
			Namespace:  "default",
			Finalizers: []string{machinev1.MachineFinalizer, metav1.FinalizerDeleteDependents},
			Labels: map[string]string{
				machinev1.MachineClusterIDLabel: "testcluster",
			},
		},
		Spec: machinev1.MachineSpec{
			ProviderSpec: machinev1.ProviderSpec{
				Value: &runtime.RawExtension{
					Raw: []byte("{}"),
				},
			},
		},
	}
	machineProvisioned := machinev1.Machine{
		TypeMeta: metav1.TypeMeta{
			Kind: "Machine",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:       "update",
			Namespace:  "default",
			Finalizers: []string{machinev1.MachineFinalizer, metav1.FinalizerDeleteDependents},
			Labels: map[string]string{
				machinev1.MachineClusterIDLabel: "testcluster",
			},
		},
		Spec: machinev1.MachineSpec{
			ProviderSpec: machinev1.ProviderSpec{
				Value: &runtime.RawExtension{
					Raw: []byte("{}"),
				},
			},
		},
		Status: machinev1.MachineStatus{
			Addresses: []corev1.NodeAddress{
				{
					Type:    corev1.NodeInternalIP,
					Address: "0.0.0.0",
				},
			},
		},
	}
	time := metav1.Now()
	machineDeleting := machinev1.Machine{
		TypeMeta: metav1.TypeMeta{
			Kind: "Machine",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              "delete",
			Namespace:         "default",
			Finalizers:        []string{machinev1.MachineFinalizer, metav1.FinalizerDeleteDependents},
			DeletionTimestamp: &time,
			Labels: map[string]string{
				machinev1.MachineClusterIDLabel: "testcluster",
			},
		},
		Spec: machinev1.MachineSpec{
			ProviderSpec: machinev1.ProviderSpec{
				Value: &runtime.RawExtension{
					Raw: []byte("{}"),
				},
			},
		},
	}
	machineFailed := machinev1.Machine{
		TypeMeta: metav1.TypeMeta{
			Kind: "Machine",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:       "failed",
			Namespace:  "default",
			Finalizers: []string{machinev1.MachineFinalizer, metav1.FinalizerDeleteDependents},
			Labels: map[string]string{
				machinev1.MachineClusterIDLabel: "testcluster",
			},
		},
		Spec: machinev1.MachineSpec{
			ProviderID: pointer.StringPtr("providerID"),
			ProviderSpec: machinev1.ProviderSpec{
				Value: &runtime.RawExtension{
					Raw: []byte("{}"),
				},
			},
		},
		Status: machinev1.MachineStatus{
			Addresses: []corev1.NodeAddress{
				{
					Type:    corev1.NodeInternalIP,
					Address: "0.0.0.0",
				},
			},
		},
	}
	machineRunning := machinev1.Machine{
		TypeMeta: metav1.TypeMeta{
			Kind: "Machine",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:       "running",
			Namespace:  "default",
			Finalizers: []string{machinev1.MachineFinalizer, metav1.FinalizerDeleteDependents},
			Labels: map[string]string{
				machinev1.MachineClusterIDLabel: "testcluster",
			},
		},
		Spec: machinev1.MachineSpec{
			ProviderID: pointer.StringPtr("providerID"),
			ProviderSpec: machinev1.ProviderSpec{
				Value: &runtime.RawExtension{
					Raw: []byte("{}"),
				},
			},
		},
		Status: machinev1.MachineStatus{
			Addresses: []corev1.NodeAddress{
				{
					Type:    corev1.NodeInternalIP,
					Address: "0.0.0.0",
				},
			},
			NodeRef: &corev1.ObjectReference{
				Name: "a node",
			},
		},
	}

	type expected struct {
		createCallCount int64
		existCallCount  int64
		updateCallCount int64
		deleteCallCount int64
		result          reconcile.Result
		error           bool
		phase           string
	}
	testCases := []struct {
		request     reconcile.Request
		existsValue bool
		expected    expected
	}{
		{
			request:     reconcile.Request{NamespacedName: types.NamespacedName{Name: machineProvisioning.Name, Namespace: machineProvisioning.Namespace}},
			existsValue: false,
			expected: expected{
				createCallCount: 1,
				existCallCount:  1,
				updateCallCount: 0,
				deleteCallCount: 0,
				result:          reconcile.Result{RequeueAfter: requeueAfter},
				error:           false,
				phase:           phaseProvisioning,
			},
		},
		{
			request:     reconcile.Request{NamespacedName: types.NamespacedName{Name: machineProvisioned.Name, Namespace: machineProvisioned.Namespace}},
			existsValue: true,
			expected: expected{
				createCallCount: 0,
				existCallCount:  1,
				updateCallCount: 1,
				deleteCallCount: 0,
				result:          reconcile.Result{RequeueAfter: requeueAfter},
				error:           false,
				phase:           phaseProvisioned,
			},
		},
		{
			request:     reconcile.Request{NamespacedName: types.NamespacedName{Name: machineDeleting.Name, Namespace: machineDeleting.Namespace}},
			existsValue: false,
			expected: expected{
				createCallCount: 0,
				existCallCount:  1,
				updateCallCount: 0,
				deleteCallCount: 1,
				result:          reconcile.Result{},
				error:           false,
				phase:           phaseDeleting,
			},
		},
		{
			request:     reconcile.Request{NamespacedName: types.NamespacedName{Name: machineDeleting.Name, Namespace: machineDeleting.Namespace}},
			existsValue: true,
			expected: expected{
				createCallCount: 0,
				existCallCount:  1,
				updateCallCount: 0,
				deleteCallCount: 1,
				result:          reconcile.Result{RequeueAfter: requeueAfter},
				error:           false,
				phase:           phaseDeleting,
			},
		},
		{
			request:     reconcile.Request{NamespacedName: types.NamespacedName{Name: machineFailed.Name, Namespace: machineFailed.Namespace}},
			existsValue: false,
			expected: expected{
				createCallCount: 0,
				existCallCount:  1,
				updateCallCount: 0,
				deleteCallCount: 0,
				result:          reconcile.Result{},
				error:           false,
				phase:           phaseFailed, // A machine which does not exist but has providerID or addresses
			},
		},
		{
			request:     reconcile.Request{NamespacedName: types.NamespacedName{Name: machineRunning.Name, Namespace: machineRunning.Namespace}},
			existsValue: true,
			expected: expected{
				createCallCount: 0,
				existCallCount:  1,
				updateCallCount: 1,
				deleteCallCount: 0,
				result:          reconcile.Result{},
				error:           false,
				phase:           phaseRunning,
			},
		},
	}

	for _, tc := range testCases {
		act := newTestActuator()
		act.ExistsValue = tc.existsValue
		machinev1.AddToScheme(scheme.Scheme)
		r := &ReconcileMachine{
			Client: fake.NewFakeClientWithScheme(scheme.Scheme,
				&machineProvisioning,
				&machineProvisioned,
				&machineDeleting,
				&machineFailed,
				&machineRunning,
			),
			scheme:   scheme.Scheme,
			actuator: act,
		}

		result, err := r.Reconcile(tc.request)
		gotError := (err != nil)
		if tc.expected.error != gotError {
			var errorExpectation string
			if !tc.expected.error {
				errorExpectation = "no"
			}
			t.Errorf("Case: %s. Expected %s error, got: %v", tc.request.Name, errorExpectation, err)
		}

		if !reflect.DeepEqual(result, tc.expected.result) {
			t.Errorf("Case %s. Got: %v, expected %v", tc.request.Name, result, tc.expected.result)
		}

		if act.CreateCallCount != tc.expected.createCallCount {
			t.Errorf("Case %s. Got: %d createCallCount, expected %d", tc.request.Name, act.CreateCallCount, tc.expected.createCallCount)
		}

		if act.UpdateCallCount != tc.expected.updateCallCount {
			t.Errorf("Case %s. Got: %d updateCallCount, expected %d", tc.request.Name, act.UpdateCallCount, tc.expected.updateCallCount)
		}

		if act.ExistsCallCount != tc.expected.existCallCount {
			t.Errorf("Case %s. Got: %d existCallCount, expected %d", tc.request.Name, act.ExistsCallCount, tc.expected.existCallCount)
		}

		if act.DeleteCallCount != tc.expected.deleteCallCount {
			t.Errorf("Case %s. Got: %d deleteCallCount, expected %d", tc.request.Name, act.DeleteCallCount, tc.expected.deleteCallCount)
		}

		machine := &machinev1.Machine{}
		if err := r.Client.Get(context.TODO(), tc.request.NamespacedName, machine); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if tc.expected.phase != stringPointerDeref(machine.Status.Phase) {
			t.Errorf("Case %s. Got: %v, expected: %v", tc.request.Name, stringPointerDeref(machine.Status.Phase), tc.expected.phase)
		}
	}
}

func TestSetPhase(t *testing.T) {
	testCases := []struct {
		name                   string
		phase                  string
		errorMessage           string
		annotations            map[string]string
		existingProviderStatus string
		expectedProviderStatus string
	}{
		{
			name:         "when updating the phase to Running",
			phase:        phaseRunning,
			errorMessage: "",
			annotations:  nil,
		},
		{
			name:         "when updating the phase to Failed",
			phase:        phaseFailed,
			errorMessage: "test",
			annotations: map[string]string{
				MachineInstanceStateAnnotationName: unknownInstanceState,
			},
		},
		{
			name:         "when updating the phase to Failed with instanceState Set",
			phase:        phaseFailed,
			errorMessage: "test",
			annotations: map[string]string{
				MachineInstanceStateAnnotationName: unknownInstanceState,
			},
			existingProviderStatus: `{"instanceState":"Running"}`,
			expectedProviderStatus: `{"instanceState":"Unknown"}`,
		},
		{
			name:         "when updating the phase to Failed with vmState Set",
			phase:        phaseFailed,
			errorMessage: "test",
			annotations: map[string]string{
				MachineInstanceStateAnnotationName: unknownInstanceState,
			},
			existingProviderStatus: `{"vmState":"Running"}`,
			expectedProviderStatus: `{"vmState":"Unknown"}`,
		},
		{
			name:         "when updating the phase to Failed with state Set",
			phase:        phaseFailed,
			errorMessage: "test",
			annotations: map[string]string{
				MachineInstanceStateAnnotationName: unknownInstanceState,
			},
			existingProviderStatus: `{"state":"Running"}`,
			expectedProviderStatus: `{"state":"Running"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			k8sClient, err := client.New(cfg, client.Options{})
			g.Expect(err).ToNot(HaveOccurred())
			machinev1.AddToScheme(scheme.Scheme)
			reconciler := &ReconcileMachine{
				Client: k8sClient,
				scheme: scheme.Scheme,
			}

			// Set up the test namespace
			name := "test"
			namespace := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: name,
				},
			}
			g.Expect(k8sClient.Create(ctx, namespace)).To(Succeed())

			// Set up the test machine
			machine := &machinev1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: name,
					Namespace:    namespace.Name,
				},
			}

			g.Expect(k8sClient.Create(ctx, machine)).To(Succeed())
			defer k8sClient.Delete(ctx, machine)

			if tc.existingProviderStatus != "" {
				machine.Status.ProviderStatus = &runtime.RawExtension{
					Raw: []byte(tc.existingProviderStatus),
				}
			}

			g.Expect(k8sClient.Status().Update(ctx, machine)).To(Succeed())

			namespacedName := types.NamespacedName{
				Namespace: machine.Namespace,
				Name:      machine.Name,
			}

			// Set the phase to Running initially
			g.Expect(reconciler.setPhase(machine, phaseRunning, "")).To(Succeed())
			// validate persisted object
			got := machinev1.Machine{}
			g.Expect(reconciler.Client.Get(context.TODO(), namespacedName, &got)).To(Succeed())
			g.Expect(got.Status.Phase).ToNot(BeNil())
			g.Expect(*got.Status.Phase).To(Equal(phaseRunning))
			lastUpdated := got.Status.LastUpdated
			g.Expect(lastUpdated).ToNot(BeNil())
			// validate passed object
			g.Expect(machine.Status.Phase).ToNot(BeNil())
			g.Expect(*machine.Status.Phase).To(Equal(phaseRunning))
			objectLastUpdated := machine.Status.LastUpdated
			g.Expect(objectLastUpdated).ToNot(BeNil())

			// Modify the phase and verify the result
			g.Expect(reconciler.setPhase(machine, tc.phase, tc.errorMessage)).To(Succeed())
			// validate the persisted object
			got = machinev1.Machine{}
			g.Expect(reconciler.Client.Get(context.TODO(), namespacedName, &got)).To(Succeed())

			if tc.phase == phaseRunning {
				// The phase hasn't changed so the lastUpdated shouldn't either
				g.Expect(*got.Status.LastUpdated).To(Equal(*lastUpdated))
				g.Expect(*machine.Status.LastUpdated).To(Equal(*objectLastUpdated))
			}

			if tc.errorMessage != "" {
				g.Expect(got.Status.ErrorMessage).ToNot(BeNil())
				g.Expect(*got.Status.ErrorMessage).To(Equal(tc.errorMessage))
				g.Expect(machine.Status.ErrorMessage).ToNot(BeNil())
				g.Expect(*machine.Status.ErrorMessage).To(Equal(tc.errorMessage))
			}

			g.Expect(*got.Status.Phase).To(Equal(tc.phase))
			g.Expect(*machine.Status.Phase).To(Equal(tc.phase))

			g.Expect(got.GetAnnotations()).To(Equal(tc.annotations))
			g.Expect(machine.GetAnnotations()).To(Equal(tc.annotations))

			if tc.existingProviderStatus != "" {
				g.Expect(got.Status.ProviderStatus).ToNot(BeNil())
				g.Expect(got.Status.ProviderStatus.Raw).To(BeEquivalentTo(tc.expectedProviderStatus))
			}
		})
	}
}

func TestMachineIsProvisioned(t *testing.T) {
	name := "test"
	namespace := "test"
	providerID := "providerID"

	testCases := []struct {
		machine  *machinev1.Machine
		expected bool
	}{
		{
			machine: &machinev1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Status: machinev1.MachineStatus{},
			},
			expected: false,
		},
		{
			machine: &machinev1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Status: machinev1.MachineStatus{
					Addresses: []corev1.NodeAddress{
						{
							Type:    corev1.NodeInternalIP,
							Address: "0.0.0.0",
						},
					},
				},
			},
			expected: true,
		},
		{
			machine: &machinev1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: machinev1.MachineSpec{
					ProviderID: &providerID,
				},
				Status: machinev1.MachineStatus{},
			},
			expected: true,
		},
	}

	for _, tc := range testCases {
		if got := machineIsProvisioned(tc.machine); got != tc.expected {
			t.Errorf("Got: %v, expected: %v", got, tc.expected)
		}
	}
}

func TestStringPointerDeref(t *testing.T) {
	value := "test"
	testCases := []struct {
		stringPointer *string
		expected      string
	}{
		{
			stringPointer: nil,
			expected:      "",
		},
		{
			stringPointer: &value,
			expected:      value,
		},
	}
	for _, tc := range testCases {
		if got := stringPointerDeref(tc.stringPointer); got != tc.expected {
			t.Errorf("Got: %v, expected: %v", got, tc.expected)
		}
	}
}

func TestMachineIsFailed(t *testing.T) {
	testCases := []struct {
		machine  *machinev1.Machine
		expected bool
	}{
		{
			machine: &machinev1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "fromNilPhase",
					Namespace: "test",
				},
				Status: machinev1.MachineStatus{},
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		if got := machineIsFailed(tc.machine); got {
			t.Errorf("Expected: %v, got: %v", got, tc.expected)
		}
	}
}

func TestNodeIsUnreachable(t *testing.T) {
	testCases := []struct {
		name     string
		node     *corev1.Node
		expected bool
	}{
		{
			name: "Node should be unreachable",
			node: &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "Node",
					Namespace: "test",
				},
				Status: corev1.NodeStatus{
					Conditions: []corev1.NodeCondition{
						{
							Type:   corev1.NodeReady,
							Status: corev1.ConditionUnknown,
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "Node should not be unreachable",
			node: &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "Node",
					Namespace: "test",
				},
				Status: corev1.NodeStatus{
					Conditions: []corev1.NodeCondition{
						{
							Type:   corev1.NodeReady,
							Status: corev1.ConditionTrue,
						},
					},
				},
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if actual := nodeIsUnreachable(tc.node); actual != tc.expected {
				t.Errorf("Expected: %v, got: %v", actual, tc.expected)
			}
		})
	}
}

func TestIsInvalidMachineConfigurationError(t *testing.T) {
	invalidMachineConfigurationError := InvalidMachineConfiguration("invalidConfiguration")
	createError := CreateMachine("createFailed")

	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "With an InvalidMachineConfigurationError",
			err:      invalidMachineConfigurationError,
			expected: true,
		},
		{
			name:     "With a CreateError",
			err:      createError,
			expected: false,
		},
		{
			name:     "With a wrapped InvalidMachineConfigurationError",
			err:      fmt.Errorf("Wrap: %w", invalidMachineConfigurationError),
			expected: true,
		},
		{
			name:     "With a wrapped CreateError",
			err:      fmt.Errorf("Wrap: %w", createError),
			expected: false,
		},
		{
			name:     "With a double wrapped InvalidMachineConfigurationError",
			err:      fmt.Errorf("Wrap: %w", fmt.Errorf("Wrap: %w", invalidMachineConfigurationError)),
			expected: true,
		},
		{
			name:     "With a double wrapped CreateError",
			err:      fmt.Errorf("Wrap: %w", fmt.Errorf("Wrap: %w", createError)),
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if actual := isInvalidMachineConfigurationError(tc.err); actual != tc.expected {
				t.Errorf("Case: %s, got: %v, expected: %v", tc.name, actual, tc.expected)
			}
		})
	}
}

func TestDelayIfRequeueAfterError(t *testing.T) {
	requeueAfter30s := &RequeueAfterError{RequeueAfter: 30 * time.Second}
	requeueAfter1m := &RequeueAfterError{RequeueAfter: time.Minute}
	createError := CreateMachine("createFailed")
	wrappedCreateError := fmt.Errorf("Wrap: %w", createError)
	doubleWrappedCreateError := fmt.Errorf("Wrap: %w", fmt.Errorf("Wrap: %w", createError))

	testCases := []struct {
		name           string
		err            error
		expectedErr    error
		expectedResult reconcile.Result
	}{
		{
			name:           "with a RequeAfterError (30s)",
			err:            requeueAfter30s,
			expectedErr:    nil,
			expectedResult: reconcile.Result{Requeue: true, RequeueAfter: 30 * time.Second},
		},
		{
			name:           "with a RequeAfterError (1m)",
			err:            requeueAfter1m,
			expectedErr:    nil,
			expectedResult: reconcile.Result{Requeue: true, RequeueAfter: time.Minute},
		},
		{
			name:           "with a CreateError",
			err:            createError,
			expectedErr:    createError,
			expectedResult: reconcile.Result{},
		},
		{
			name:           "with a wrapped RequeAfterError (30s)",
			err:            fmt.Errorf("Wrap: %w", requeueAfter30s),
			expectedErr:    nil,
			expectedResult: reconcile.Result{Requeue: true, RequeueAfter: 30 * time.Second},
		},
		{
			name:           "with a wrapped RequeAfterError (1m)",
			err:            fmt.Errorf("Wrap: %w", requeueAfter1m),
			expectedErr:    nil,
			expectedResult: reconcile.Result{Requeue: true, RequeueAfter: time.Minute},
		},
		{
			name:           "with a wrapped CreateError",
			err:            wrappedCreateError,
			expectedErr:    wrappedCreateError,
			expectedResult: reconcile.Result{},
		},
		{
			name:           "with a double wrapped RequeAfterError (30s)",
			err:            fmt.Errorf("Wrap: %w", fmt.Errorf("Wrap: %w", requeueAfter30s)),
			expectedErr:    nil,
			expectedResult: reconcile.Result{Requeue: true, RequeueAfter: 30 * time.Second},
		},
		{
			name:           "with a double wrapped RequeAfterError (1m)",
			err:            fmt.Errorf("Wrap: %w", fmt.Errorf("Wrap: %w", requeueAfter1m)),
			expectedErr:    nil,
			expectedResult: reconcile.Result{Requeue: true, RequeueAfter: time.Minute},
		},
		{
			name:           "with a double wrapped CreateError",
			err:            doubleWrappedCreateError,
			expectedErr:    doubleWrappedCreateError,
			expectedResult: reconcile.Result{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := delayIfRequeueAfterError(tc.err)
			if err != tc.expectedErr {
				t.Errorf("Case: %s, got: %v, expected: %v", tc.name, err, tc.expectedErr)
			}
			if result != tc.expectedResult {
				t.Errorf("Case: %s, got: %v, expected: %v", tc.name, result, tc.expectedResult)
			}
		})
	}
}
