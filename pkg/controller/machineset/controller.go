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
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	machinev1beta1 "github.com/uccps-samples/machine-api-operator/pkg/apis/machine/v1beta1"
	"github.com/uccps-samples/machine-api-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var (
	controllerKind = machinev1beta1.SchemeGroupVersion.WithKind("MachineSet")

	// stateConfirmationTimeout is the amount of time allowed to wait for desired state.
	stateConfirmationTimeout = 10 * time.Second

	// stateConfirmationInterval is the amount of time between polling for the desired state.
	// The polling is against a local memory cache.
	stateConfirmationInterval = 100 * time.Millisecond

	// controllerName is the name of this controller
	controllerName = "machineset_controller"
)

// Add creates a new MachineSet Controller and adds it to the Manager with default RBAC.
// The Manager will set fields on the Controller and Start it when the Manager is Started.
func Add(mgr manager.Manager, opts manager.Options) error {
	r := newReconciler(mgr)
	return add(mgr, r, r.MachineToMachineSets)
}

// newReconciler returns a new reconcile.Reconciler.
func newReconciler(mgr manager.Manager) *ReconcileMachineSet {
	return &ReconcileMachineSet{Client: mgr.GetClient(), scheme: mgr.GetScheme(), recorder: mgr.GetEventRecorderFor(controllerName)}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler.
func add(mgr manager.Manager, r reconcile.Reconciler, mapFn handler.ToRequestsFunc) error {
	// Create a new controller.
	c, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to MachineSet.
	err = c.Watch(
		&source.Kind{Type: &machinev1beta1.MachineSet{}},
		&handler.EnqueueRequestForObject{},
	)
	if err != nil {
		return err
	}

	// Map Machine changes to MachineSets using ControllerRef.
	err = c.Watch(
		&source.Kind{Type: &machinev1beta1.Machine{}},
		&handler.EnqueueRequestForOwner{IsController: true, OwnerType: &machinev1beta1.MachineSet{}},
	)
	if err != nil {
		return err
	}

	// Map Machine changes to MachineSets by machining labels.
	return c.Watch(
		&source.Kind{Type: &machinev1beta1.Machine{}},
		&handler.EnqueueRequestsFromMapFunc{ToRequests: mapFn},
	)
}

// ReconcileMachineSet reconciles a MachineSet object
type ReconcileMachineSet struct {
	client.Client
	scheme   *runtime.Scheme
	recorder record.EventRecorder
}

func (r *ReconcileMachineSet) MachineToMachineSets(o handler.MapObject) []reconcile.Request {
	result := []reconcile.Request{}
	m := &machinev1beta1.Machine{}
	key := client.ObjectKey{Namespace: o.Meta.GetNamespace(), Name: o.Meta.GetName()}
	err := r.Client.Get(context.Background(), key, m)
	if err != nil {
		klog.Errorf("Unable to retrieve Machine %v from store: %v", key, err)
		return nil
	}

	for _, ref := range m.ObjectMeta.OwnerReferences {
		if ref.Controller != nil && *ref.Controller {
			return result
		}
	}

	mss := r.getMachineSetsForMachine(m)
	if len(mss) == 0 {
		klog.V(4).Infof("Found no machine set for machine: %v", m.Name)
		return nil
	}

	for _, ms := range mss {
		name := client.ObjectKey{Namespace: ms.Namespace, Name: ms.Name}
		result = append(result, reconcile.Request{NamespacedName: name})
	}

	return result
}

// Reconcile reads that state of the cluster for a MachineSet object and makes changes based on the state read
// and what is in the MachineSet.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=machine.uccp.io,resources=machinesets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=machine.uccp.io,resources=machines,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileMachineSet) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the MachineSet instance
	ctx := context.TODO()
	machineSet := &machinev1beta1.MachineSet{}
	if err := r.Get(ctx, request.NamespacedName, machineSet); err != nil {
		if apierrors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Ignore deleted MachineSets, this can happen when foregroundDeletion
	// is enabled
	if machineSet.DeletionTimestamp != nil {
		return reconcile.Result{}, nil
	}

	result, err := r.reconcile(ctx, machineSet)
	if err != nil {
		klog.Errorf("Failed to reconcile MachineSet %q: %v", request.NamespacedName, err)
		r.recorder.Eventf(machineSet, corev1.EventTypeWarning, "ReconcileError", "%v", err)
	}
	return result, err
}

func (r *ReconcileMachineSet) reconcile(ctx context.Context, machineSet *machinev1beta1.MachineSet) (reconcile.Result, error) {
	klog.V(4).Infof("Reconcile machineset %v", machineSet.Name)
	if errList := machineSet.Validate(); len(errList) > 0 {
		err := fmt.Errorf("%q machineset validation failed: %v", machineSet.Name, errList.ToAggregate().Error())
		klog.Error(err)
		return reconcile.Result{}, err
	}

	allMachines := &machinev1beta1.MachineList{}

	if err := r.Client.List(context.Background(), allMachines, client.InNamespace(machineSet.Namespace)); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to list machines: %w", err)
	}

	// Make sure that label selector can match template's labels.
	// TODO(vincepri): Move to a validation (admission) webhook when supported.
	selector, err := metav1.LabelSelectorAsSelector(&machineSet.Spec.Selector)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to parse MachineSet %q label selector: %w", machineSet.Name, err)
	}

	if !selector.Matches(labels.Set(machineSet.Spec.Template.Labels)) {
		return reconcile.Result{}, fmt.Errorf("failed validation on MachineSet %q label selector, cannot match any machines ", machineSet.Name)
	}

	// Filter out irrelevant machines (deleting/mismatch labels) and claim orphaned machines.
	var machineNames []string
	machineSetMachines := make(map[string]*machinev1beta1.Machine)
	for idx := range allMachines.Items {
		machine := &allMachines.Items[idx]
		if shouldExcludeMachine(machineSet, machine) {
			continue
		}

		// Attempt to adopt machine if it meets previous conditions and it has no controller references.
		if metav1.GetControllerOf(machine) == nil {
			if err := r.adoptOrphan(machineSet, machine); err != nil {
				klog.Warningf("Failed to adopt MachineSet %q into MachineSet %q: %v", machine.Name, machineSet.Name, err)
				continue
			}
		}
		machineNames = append(machineNames, machine.Name)
		machineSetMachines[machine.Name] = machine
	}
	// sort the filteredMachines from the oldest to the youngest
	sort.Strings(machineNames)

	var filteredMachines []*machinev1beta1.Machine
	for _, machineName := range machineNames {
		filteredMachines = append(filteredMachines, machineSetMachines[machineName])
	}

	syncErr := r.syncReplicas(machineSet, filteredMachines)

	ms := machineSet.DeepCopy()
	newStatus := r.calculateStatus(ms, filteredMachines)

	// Always updates status as machines come up or die.
	updatedMS, err := updateMachineSetStatus(r.Client, machineSet, newStatus)
	if err != nil {
		if syncErr != nil {
			return reconcile.Result{}, fmt.Errorf("failed to sync machines: %v. failed to update machine set status: %w", syncErr, err)
		}
		return reconcile.Result{}, fmt.Errorf("failed to update machine set status: %w", err)
	}

	if syncErr != nil {
		return reconcile.Result{}, fmt.Errorf("failed to sync machines: %w", syncErr)
	}

	var replicas int32
	if updatedMS.Spec.Replicas != nil {
		replicas = *updatedMS.Spec.Replicas
	}

	// Resync the MachineSet after MinReadySeconds as a last line of defense to guard against clock-skew.
	// Clock-skew is an issue as it may impact whether an available replica is counted as a ready replica.
	// A replica is available if the amount of time since last transition exceeds MinReadySeconds.
	// If there was a clock skew, checking whether the amount of time since last transition to ready state
	// exceeds MinReadySeconds could be incorrect.
	// To avoid an available replica stuck in the ready state, we force a reconcile after MinReadySeconds,
	// at which point it should confirm any available replica to be available.
	if updatedMS.Spec.MinReadySeconds > 0 &&
		updatedMS.Status.ReadyReplicas == replicas &&
		updatedMS.Status.AvailableReplicas != replicas {

		return reconcile.Result{Requeue: true}, nil
	}

	return reconcile.Result{}, nil
}

// syncReplicas essentially scales machine resources up and down.
func (r *ReconcileMachineSet) syncReplicas(ms *machinev1beta1.MachineSet, machines []*machinev1beta1.Machine) error {
	if ms.Spec.Replicas == nil {
		return fmt.Errorf("the Replicas field in Spec for machineset %v is nil, this should not be allowed", ms.Name)
	}

	diff := len(machines) - int(*(ms.Spec.Replicas))

	if diff < 0 {
		diff *= -1
		klog.Infof("Too few replicas for %v %s/%s, need %d, creating %d",
			controllerKind, ms.Namespace, ms.Name, *(ms.Spec.Replicas), diff)

		var machineList []*machinev1beta1.Machine
		var errstrings []string
		for i := 0; i < diff; i++ {
			klog.Infof("Creating machine %d of %d, ( spec.replicas(%d) > currentMachineCount(%d) )",
				i+1, diff, *(ms.Spec.Replicas), len(machines))

			machine := r.createMachine(ms)
			if err := r.Client.Create(context.Background(), machine); err != nil {
				klog.Errorf("Unable to create Machine %q: %v", machine.Name, err)
				errstrings = append(errstrings, err.Error())
				continue
			}

			machineList = append(machineList, machine)
		}

		if len(errstrings) > 0 {
			return errors.New(strings.Join(errstrings, "; "))
		}

		return r.waitForMachineCreation(machineList)
	} else if diff > 0 {
		klog.Infof("Too many replicas for %v %s/%s, need %d, deleting %d",
			controllerKind, ms.Namespace, ms.Name, *(ms.Spec.Replicas), diff)

		deletePriorityFunc, err := getDeletePriorityFunc(ms)
		if err != nil {
			return err
		}
		klog.Infof("Found %s delete policy", ms.Spec.DeletePolicy)
		// Choose which Machines to delete.
		machinesToDelete := getMachinesToDeletePrioritized(machines, diff, deletePriorityFunc)

		// TODO: Add cap to limit concurrent delete calls.
		errCh := make(chan error, diff)
		var wg sync.WaitGroup
		wg.Add(diff)
		for _, machine := range machinesToDelete {
			go func(targetMachine *machinev1beta1.Machine) {
				defer wg.Done()
				err := r.Client.Delete(context.Background(), targetMachine)
				if err != nil {
					klog.Errorf("Unable to delete Machine %s: %v", targetMachine.Name, err)
					errCh <- err
				}
			}(machine)
		}
		wg.Wait()

		select {
		case err := <-errCh:
			// all errors have been reported before and they're likely to be the same, so we'll only return the first one we hit.
			if err != nil {
				return err
			}
		default:
		}

		return r.waitForMachineDeletion(machinesToDelete)
	}

	return nil
}

// createMachine creates a machine resource.
// the name of the newly created resource is going to be created by the API server, we set the generateName field
func (r *ReconcileMachineSet) createMachine(machineSet *machinev1beta1.MachineSet) *machinev1beta1.Machine {
	gv := machinev1beta1.SchemeGroupVersion
	machine := &machinev1beta1.Machine{
		TypeMeta: metav1.TypeMeta{
			Kind:       gv.WithKind("Machine").Kind,
			APIVersion: gv.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:      machineSet.Spec.Template.ObjectMeta.Labels,
			Annotations: machineSet.Spec.Template.ObjectMeta.Annotations,
		},
		Spec: machineSet.Spec.Template.Spec,
	}
	machine.ObjectMeta.GenerateName = fmt.Sprintf("%s-", machineSet.Name)
	machine.ObjectMeta.OwnerReferences = []metav1.OwnerReference{*metav1.NewControllerRef(machineSet, controllerKind)}
	machine.Namespace = machineSet.Namespace

	return machine
}

// shouldExcludeMachine returns true if the machine should be filtered out, false otherwise.
func shouldExcludeMachine(machineSet *machinev1beta1.MachineSet, machine *machinev1beta1.Machine) bool {
	// Ignore inactive machines.
	if metav1.GetControllerOf(machine) != nil && !metav1.IsControlledBy(machine, machineSet) {
		klog.V(4).Infof("%s not controlled by %v", machine.Name, machineSet.Name)
		return true
	}

	if machine.ObjectMeta.DeletionTimestamp != nil {
		return true
	}

	if !hasMatchingLabels(machineSet, machine) {
		return true
	}

	return false
}

func (r *ReconcileMachineSet) adoptOrphan(machineSet *machinev1beta1.MachineSet, machine *machinev1beta1.Machine) error {
	newRef := *metav1.NewControllerRef(machineSet, controllerKind)
	machine.OwnerReferences = append(machine.OwnerReferences, newRef)
	return r.Client.Update(context.Background(), machine)
}

func (r *ReconcileMachineSet) waitForMachineCreation(machineList []*machinev1beta1.Machine) error {
	for _, machine := range machineList {
		pollErr := util.PollImmediate(stateConfirmationInterval, stateConfirmationTimeout, func() (bool, error) {
			key := client.ObjectKey{Namespace: machine.Namespace, Name: machine.Name}

			if err := r.Client.Get(context.Background(), key, &machinev1beta1.Machine{}); err != nil {
				if apierrors.IsNotFound(err) {
					return false, nil
				}
				klog.Error(err)
				return false, err
			}

			return true, nil
		})

		if pollErr != nil {
			klog.Error(pollErr)
			return fmt.Errorf("failed waiting for machine object to be created: %w", pollErr)
		}
	}

	return nil
}

func (r *ReconcileMachineSet) waitForMachineDeletion(machineList []*machinev1beta1.Machine) error {
	for _, machine := range machineList {
		pollErr := util.PollImmediate(stateConfirmationInterval, stateConfirmationTimeout, func() (bool, error) {
			m := &machinev1beta1.Machine{}
			key := client.ObjectKey{Namespace: machine.Namespace, Name: machine.Name}

			err := r.Client.Get(context.Background(), key, m)
			if apierrors.IsNotFound(err) || !m.DeletionTimestamp.IsZero() {
				return true, nil
			}

			return false, err
		})

		if pollErr != nil {
			klog.Error(pollErr)
			return fmt.Errorf("failed waiting for machine object to be deleted: %w", pollErr)
		}
	}
	return nil
}
