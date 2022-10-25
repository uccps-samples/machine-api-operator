package testing

import (
	"fmt"
	"time"

	mapiv1 "github.com/uccps-samples/machine-api-operator/pkg/apis/machine/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/utils/pointer"
)

const (
	// Namespace contains the default namespace for machine-api components
	Namespace = "uccp-machine-api"
	// MachineAnnotationKey contains default machine node annotation
	MachineAnnotationKey = "machine.uccp.io/machine"
)

var (
	// KnownDate contains date that can be used under tests
	KnownDate = metav1.Time{Time: time.Date(1985, 06, 03, 0, 0, 0, 0, time.Local)}
)

// FooBar returns foo:bar map that can be used as default label
func FooBar() map[string]string {
	return map[string]string{"foo": "bar"}
}

// NewSelector returns new LabelSelector
func NewSelector(labels map[string]string) *metav1.LabelSelector {
	return &metav1.LabelSelector{MatchLabels: labels}
}

// NewSelectorFooBar returns new foo:bar label selector
func NewSelectorFooBar() *metav1.LabelSelector {
	return NewSelector(FooBar())
}

func NewExternalRemediationTemplate() *unstructured.Unstructured {
	// Create remediation template resource.
	infraRemediationTmpl := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "InfrastructureRemediationTemplate",
			"apiVersion": "infrastructure.machine.uccp.io/v1alpha3",
			"spec": map[string]interface{}{
				"template": map[string]interface{}{},
			},
			"metadata": map[string]interface{}{
				"namespace": Namespace,
			},
		},
	}

	return infraRemediationTmpl
}

func NewExternalRemediationMachine() *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "InfrastructureRemediation",
			"apiVersion": "infrastructure.machine.uccp.io/v1alpha3",
			"metadata": map[string]interface{}{
				"name":      "Machine",
				"namespace": Namespace,
			},
		},
	}
}

// NewNode returns new node object that can be used for testing
func NewNode(name string, ready bool) *corev1.Node {
	nodeReadyStatus := corev1.ConditionTrue
	if !ready {
		nodeReadyStatus = corev1.ConditionUnknown
	}

	return &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: metav1.NamespaceNone,
			Annotations: map[string]string{
				MachineAnnotationKey: fmt.Sprintf("%s/%s", Namespace, "fakeMachine"),
			},
			Labels: map[string]string{},
			UID:    uuid.NewUUID(),
		},
		TypeMeta: metav1.TypeMeta{
			Kind: "Node",
		},
		Status: corev1.NodeStatus{
			Conditions: []corev1.NodeCondition{
				{
					Type:               corev1.NodeReady,
					Status:             nodeReadyStatus,
					LastTransitionTime: KnownDate,
				},
			},
		},
	}
}

// NewMachine returns new machine object that can be used for testing
func NewMachine(name string, nodeName string) *mapiv1.Machine {
	m := &mapiv1.Machine{
		TypeMeta: metav1.TypeMeta{Kind: "Machine"},
		ObjectMeta: metav1.ObjectMeta{
			Annotations: make(map[string]string),
			Name:        name,
			Namespace:   Namespace,
			Labels:      FooBar(),
			UID:         uuid.NewUUID(),
			OwnerReferences: []metav1.OwnerReference{
				{
					Kind:       "MachineSet",
					Controller: pointer.BoolPtr(true),
				},
			},
			// the following line is to account for a change in the fake client, see https://github.com/kubernetes-sigs/controller-runtime/pull/1306
			ResourceVersion: "999",
		},
		Spec: mapiv1.MachineSpec{},
	}
	if nodeName != "" {
		m.Status = mapiv1.MachineStatus{
			NodeRef: &corev1.ObjectReference{
				Name:      nodeName,
				Namespace: metav1.NamespaceNone,
			},
		}
	}
	return m
}

// NewMachineHealthCheck returns new MachineHealthCheck object that can be used for testing
func NewMachineHealthCheck(name string) *mapiv1.MachineHealthCheck {
	return &mapiv1.MachineHealthCheck{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: Namespace,
			// the following line is to account for a change in the fake client, see https://github.com/kubernetes-sigs/controller-runtime/pull/1306
			ResourceVersion: "999",
		},
		TypeMeta: metav1.TypeMeta{
			Kind: "MachineHealthCheck",
		},
		Spec: mapiv1.MachineHealthCheckSpec{
			Selector: *NewSelectorFooBar(),
			UnhealthyConditions: []mapiv1.UnhealthyCondition{
				{
					Type:    "Ready",
					Status:  "Unknown",
					Timeout: metav1.Duration{Duration: 300 * time.Second},
				},
				{
					Type:    "Ready",
					Status:  "False",
					Timeout: metav1.Duration{Duration: 300 * time.Second},
				},
			},
		},
		Status: mapiv1.MachineHealthCheckStatus{},
	}
}
