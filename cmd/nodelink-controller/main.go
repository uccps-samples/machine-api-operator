package main

import (
	"flag"
	"runtime"
	"time"

	mapiv1 "github.com/uccps-samples/machine-api-operator/pkg/apis/machine/v1beta1"
	"github.com/uccps-samples/machine-api-operator/pkg/controller"
	"github.com/uccps-samples/machine-api-operator/pkg/controller/nodelink"
	sdkVersion "github.com/operator-framework/operator-sdk/version"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

// The default durations for the leader electrion operations.
var (
	leaseDuration = 120 * time.Second
	renewDealine  = 110 * time.Second
	retryPeriod   = 90 * time.Second
)

func printVersion() {
	klog.Infof("Go Version: %s", runtime.Version())
	klog.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	klog.Infof("operator-sdk Version: %v", sdkVersion.Version)
}

func main() {
	printVersion()

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

	klog.InitFlags(nil)
	flag.Set("logtostderr", "true")
	flag.Parse()

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		klog.Fatal(err)
	}

	opts := manager.Options{
		// Disable metrics serving
		MetricsBindAddress:      "0",
		LeaderElection:          *leaderElect,
		LeaderElectionNamespace: *leaderElectResourceNamespace,
		LeaderElectionID:        "cluster-api-provider-nodelink-leader",
		LeaseDuration:           leaderElectLeaseDuration,
		// Slow the default retry and renew election rate to reduce etcd writes at idle: BZ 1858400
		RetryPeriod:   &retryPeriod,
		RenewDeadline: &renewDealine,
	}
	if *watchNamespace != "" {
		opts.Namespace = *watchNamespace
		klog.Infof("Watching machine-api objects only in namespace %q for reconciliation.", opts.Namespace)
	}
	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(cfg, opts)
	if err != nil {
		klog.Fatal(err)
	}

	klog.Infof("Registering Components.")

	// Setup Scheme for all resources
	if err := mapiv1.AddToScheme(mgr.GetScheme()); err != nil {
		klog.Fatal(err)
	}

	// Setup all Controllers
	if err := controller.AddToManager(mgr, opts, nodelink.Add); err != nil {
		klog.Fatal(err)
	}

	klog.Info("Starting the Cmd.")

	// Start the Cmd
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		klog.Fatal(err)
	}
}
