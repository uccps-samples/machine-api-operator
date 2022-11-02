package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	configv1 "github.com/uccps-samples/api/config/v1"
	"github.com/uccps-samples/machine-api-operator/pkg/apis/machine/v1beta1"
	vsphereapis "github.com/uccps-samples/machine-api-operator/pkg/apis/vsphereprovider"
	capimachine "github.com/uccps-samples/machine-api-operator/pkg/controller/machine"
	machine "github.com/uccps-samples/machine-api-operator/pkg/controller/vsphere"
	"github.com/uccps-samples/machine-api-operator/pkg/metrics"
	"github.com/uccps-samples/machine-api-operator/pkg/version"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

// The default durations for the leader electrion operations.
var (
	leaseDuration = 120 * time.Second
	renewDealine  = 110 * time.Second
	retryPeriod   = 90 * time.Second
)

func main() {
	var printVersion bool
	flag.BoolVar(&printVersion, "version", false, "print version and exit")

	klog.InitFlags(nil)
	watchNamespace := flag.String(
		"namespace",
		"",
		"Namespace that the controller watches to reconcile machine-api objects. If unspecified, the controller watches for machine-api objects across all namespaces.",
	)

	leaderElectResourceNamespace := flag.String(
		"leader-elect-resource-namespace",
		"",
		"The namespace of resource object that is used for locking during leader election. If unspecified and running in cluster, defaults to the service account namespace for the controller. Required for leader-election outside of a cluster.",
	)

	leaderElect := flag.Bool(
		"leader-elect",
		false,
		"Start a leader election client and gain leadership before executing the main loop. Enable this when running replicated components for high availability.",
	)

	leaderElectLeaseDuration := flag.Duration(
		"leader-elect-lease-duration",
		leaseDuration,
		"The duration that non-leader candidates will wait after observing a leadership renewal until attempting to acquire leadership of a led but unrenewed leader slot. This is effectively the maximum duration that a leader can be stopped before it is replaced by another candidate. This is only applicable if leader election is enabled.",
	)

	metricsAddress := flag.String(
		"metrics-bind-address",
		metrics.DefaultMachineMetricsAddress,
		"Address for hosting metrics",
	)

	flag.Set("logtostderr", "true")
	healthAddr := flag.String(
		"health-addr",
		":9440",
		"The address for health checking.",
	)
	flag.Parse()

	if printVersion {
		fmt.Println(version.String)
		os.Exit(0)
	}

	cfg := config.GetConfigOrDie()
	syncPeriod := 10 * time.Minute

	opts := manager.Options{
		MetricsBindAddress:      *metricsAddress,
		HealthProbeBindAddress:  *healthAddr,
		SyncPeriod:              &syncPeriod,
		LeaderElection:          *leaderElect,
		LeaderElectionNamespace: *leaderElectResourceNamespace,
		LeaderElectionID:        "cluster-api-provider-vsphere-leader",
		LeaseDuration:           leaderElectLeaseDuration,
		// Slow the default retry and renew election rate to reduce etcd writes at idle: BZ 1858400
		RetryPeriod:   &retryPeriod,
		RenewDeadline: &renewDealine,
	}

	if *watchNamespace != "" {
		opts.Namespace = *watchNamespace
		klog.Infof("Watching machine-api objects only in namespace %q for reconciliation.", opts.Namespace)
	}

	// Setup a Manager
	mgr, err := manager.New(cfg, opts)
	if err != nil {
		klog.Fatalf("Failed to set up overall controller manager: %v", err)
	}

	// Create a taskIDCache for create task IDs in case they are lost due to
	// network error or stale cache.
	taskIDCache := make(map[string]string)

	// Initialize machine actuator.
	machineActuator := machine.NewActuator(machine.ActuatorParams{
		Client:        mgr.GetClient(),
		APIReader:     mgr.GetAPIReader(),
		EventRecorder: mgr.GetEventRecorderFor("vspherecontroller"),
		TaskIDCache:   taskIDCache,
	})

	if err := configv1.AddToScheme(mgr.GetScheme()); err != nil {
		klog.Fatal(err)
	}

	if err := vsphereapis.AddToScheme(mgr.GetScheme()); err != nil {
		klog.Fatal(err)
	}

	if err := v1beta1.AddToScheme(mgr.GetScheme()); err != nil {
		klog.Fatal(err)
	}

	capimachine.AddWithActuator(mgr, machineActuator)

	if err := mgr.AddReadyzCheck("ping", healthz.Ping); err != nil {
		klog.Fatal(err)
	}

	if err := mgr.AddHealthzCheck("ping", healthz.Ping); err != nil {
		klog.Fatal(err)
	}

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		klog.Fatalf("Failed to run manager: %v", err)
	}
}
