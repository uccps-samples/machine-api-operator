package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	osconfigv1 "github.com/uccps-samples/api/config/v1"
	"github.com/uccps-samples/machine-api-operator/pkg/metrics"
	"github.com/uccps-samples/machine-api-operator/pkg/operator"
	"github.com/uccps-samples/machine-api-operator/pkg/util"
	"github.com/uccps-samples/machine-api-operator/pkg/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	coreclientsetv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
)

const (
	// defaultMetricsPort is the default port to expose metrics.
	defaultMetricsPort = 8080
)

var (
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Starts Machine API Operator",
		Long:  "",
		Run:   runStartCmd,
	}

	startOpts struct {
		kubeconfig string
		imagesFile string
	}
)

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().StringVar(&startOpts.kubeconfig, "kubeconfig", "", "Kubeconfig file to access a remote cluster (testing only)")
	startCmd.PersistentFlags().StringVar(&startOpts.imagesFile, "images-json", "", "images.json file for MAO.")

	klog.InitFlags(nil)
	flag.Parse()
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}

func runStartCmd(cmd *cobra.Command, args []string) {
	flag.Set("logtostderr", "true")

	// To help debugging, immediately log version
	klog.Infof("Version: %+v", version.Version)

	if startOpts.imagesFile == "" {
		klog.Fatalf("--images-json should not be empty")
	}

	cb, err := NewClientBuilder(startOpts.kubeconfig)
	if err != nil {
		klog.Fatalf("error creating clients: %v", err)
	}
	stopCh := make(chan struct{})

	le := util.GetLeaderElectionConfig(cb.config, osconfigv1.LeaderElection{})

	leaderelection.RunOrDie(context.TODO(), leaderelection.LeaderElectionConfig{
		Lock:          CreateResourceLock(cb, componentNamespace, componentName),
		RenewDeadline: le.RenewDeadline.Duration,
		RetryPeriod:   le.RetryPeriod.Duration,
		LeaseDuration: le.LeaseDuration.Duration,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				ctrlCtx := CreateControllerContext(cb, stopCh, componentNamespace)
				startControllers(ctrlCtx)
				ctrlCtx.KubeNamespacedInformerFactory.Start(ctrlCtx.Stop)
				ctrlCtx.ConfigInformerFactory.Start(ctrlCtx.Stop)
				initMachineAPIInformers(ctrlCtx)
				startMetricsCollectionAndServer(ctrlCtx)
				close(ctrlCtx.InformersStarted)

				select {}
			},
			OnStoppedLeading: func() {
				klog.Fatalf("Leader election lost")
			},
		},
	})
	panic("unreachable")
}

func initMachineAPIInformers(ctx *ControllerContext) {
	mInformer := ctx.MachineInformerFactory.Machine().V1beta1().Machines().Informer()
	msInformer := ctx.MachineInformerFactory.Machine().V1beta1().MachineSets().Informer()
	ctx.MachineInformerFactory.Start(ctx.Stop)
	if !cache.WaitForCacheSync(ctx.Stop,
		mInformer.HasSynced,
		msInformer.HasSynced) {
		klog.Fatal("Failed to sync caches for Machine api informers")
	}
	klog.Info("Synced up machine api informer caches")
}

func initRecorder(kubeClient kubernetes.Interface) record.EventRecorder {
	eventRecorderScheme := runtime.NewScheme()
	osconfigv1.Install(eventRecorderScheme)
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&coreclientsetv1.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})
	return eventBroadcaster.NewRecorder(eventRecorderScheme, v1.EventSource{Component: "machineapioperator"})
}

func startControllers(ctx *ControllerContext) {
	kubeClient := ctx.ClientBuilder.KubeClientOrDie(componentName)
	recorder := initRecorder(kubeClient)
	go operator.New(
		componentNamespace, componentName,
		startOpts.imagesFile,
		config,
		ctx.KubeNamespacedInformerFactory.Apps().V1().Deployments(),
		ctx.KubeNamespacedInformerFactory.Apps().V1().DaemonSets(),
		ctx.ConfigInformerFactory.Config().V1().FeatureGates(),
		ctx.KubeNamespacedInformerFactory.Admissionregistration().V1().ValidatingWebhookConfigurations(),
		ctx.KubeNamespacedInformerFactory.Admissionregistration().V1().MutatingWebhookConfigurations(),
		ctx.ConfigInformerFactory.Config().V1().Proxies(),
		ctx.ClientBuilder.KubeClientOrDie(componentName),
		ctx.ClientBuilder.OpenshiftClientOrDie(componentName),
		ctx.ClientBuilder.DynamicClientOrDie(componentName),
		recorder,
	).Run(1, ctx.Stop)
}

func startMetricsCollectionAndServer(ctx *ControllerContext) {
	machineInformer := ctx.MachineInformerFactory.Machine().V1beta1().Machines()
	machinesetInformer := ctx.MachineInformerFactory.Machine().V1beta1().MachineSets()
	machineMetricsCollector := metrics.NewMachineCollector(
		machineInformer,
		machinesetInformer,
		componentNamespace)
	prometheus.MustRegister(machineMetricsCollector)
	metricsPort := defaultMetricsPort
	if port, ok := os.LookupEnv("METRICS_PORT"); ok {
		v, err := strconv.Atoi(port)
		if err != nil {
			klog.Fatalf("Error parsing METRICS_PORT (%q) environment variable: %v", port, err)
		}
		metricsPort = v
	}
	klog.V(4).Info("Starting server to serve prometheus metrics")
	go startHTTPMetricServer(fmt.Sprintf("localhost:%d", metricsPort))
}

func startHTTPMetricServer(metricsPort string) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    metricsPort,
		Handler: mux,
	}
	klog.Fatal(server.ListenAndServe())
}
