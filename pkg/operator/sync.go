package operator

import (
	"context"
	"fmt"
	"time"

	"github.com/uccps-samples/library-go/pkg/operator/events"
	"github.com/uccps-samples/library-go/pkg/operator/resource/resourceapply"
	"github.com/uccps-samples/library-go/pkg/operator/resource/resourcehash"
	"github.com/uccps-samples/library-go/pkg/operator/resource/resourcemerge"
	machinecontroller "github.com/uccps-samples/machine-api-operator/pkg/controller/machine"
	"github.com/uccps-samples/machine-api-operator/pkg/metrics"
	"github.com/uccps-samples/machine-api-operator/pkg/util/conditions"
	mapiwebhooks "github.com/uccps-samples/machine-api-operator/pkg/webhooks"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/klog/v2"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	checkStatusRequeuePeriod            = 5 * time.Second
	deploymentMinimumAvailabilityTime   = 3 * time.Minute
	machineAPITerminationHandler        = "machine-api-termination-handler"
	machineExposeMetricsPort            = 8441
	machineSetExposeMetricsPort         = 8442
	machineHealthCheckExposeMetricsPort = 8444
	defaultMachineHealthPort            = 9440
	defaultMachineSetHealthPort         = 9441
	defaultMachineHealthCheckHealthPort = 9442
	kubeRBACConfigName                  = "config"
	certStoreName                       = "machine-api-controllers-tls"
	externalTrustBundleConfigMapName    = "mao-trusted-ca"
	hostKubeConfigPath                  = "/var/lib/kubelet/kubeconfig"
	hostKubePKIPath                     = "/var/lib/kubelet/pki"
	operatorStatusNoOpMessage           = "Cluster Machine API Operator is in NoOp mode"
)

var (
	// daemonsetMaxUnavailable must be set to "10%" to conform with other
	// daemonsets.
	daemonsetMaxUnavailable = intstr.FromString("10%")

	commonPodTemplateAnnotations = map[string]string{
		"target.workload.uccp.io/management": `{"effect": "PreferredDuringScheduling"}`,
	}
)

func (optr *Operator) syncAll(config *OperatorConfig) (reconcile.Result, error) {
	if err := optr.statusProgressing(); err != nil {
		klog.Errorf("Error syncing ClusterOperatorStatus: %v", err)
		return reconcile.Result{}, fmt.Errorf("error syncing ClusterOperatorStatus: %v", err)
	}

	if config.Controllers.Provider == clusterAPIControllerNoOp {
		klog.V(3).Info("Provider is NoOp, skipping synchronisation")
		if err := optr.statusAvailable(operatorStatusNoOpMessage); err != nil {
			klog.Errorf("Error syncing ClusterOperatorStatus: %v", err)
			return reconcile.Result{}, fmt.Errorf("error syncing ClusterOperatorStatus: %v", err)
		}
		return reconcile.Result{}, nil
	}

	errors := []error{}
	// Sync webhook configuration
	if err := optr.syncWebhookConfiguration(); err != nil {
		errors = append(errors, fmt.Errorf("Error syncing machine API webhook configurations: %w", err))
	}

	if err := optr.syncClusterAPIController(config); err != nil {
		errors = append(errors, fmt.Errorf("Error syncing machine-api-controller: %w", err))
	}

	// Sync Termination Handler DaemonSet if supported
	if config.Controllers.TerminationHandler != clusterAPIControllerNoOp {
		if err := optr.syncTerminationHandler(config); err != nil {
			errors = append(errors, fmt.Errorf("Error syncing termination handler: %w", err))
		}
	}

	if len(errors) > 0 {
		err := utilerrors.NewAggregate(errors)
		if err := optr.statusDegraded(err.Error()); err != nil {
			// Just log the error here.  We still want to
			// return the outer error.
			klog.Errorf("Error syncing ClusterOperatorStatus: %v", err)
		}
		klog.Errorf("Error syncing machine controller components: %v", err)
		return reconcile.Result{}, err
	}

	result, err := optr.checkRolloutStatus(config)
	if err != nil {
		if err := optr.statusDegraded(err.Error()); err != nil {
			// Just log the error here.  We still want to
			// return the outer error.
			klog.Errorf("Error syncing ClusterOperatorStatus: %v", err)
		}
		klog.Errorf("Error waiting for resource to sync: %v", err)
		return reconcile.Result{}, err
	}
	if result.Requeue || result.RequeueAfter > 0 {
		// The deployment is not yet rolled out, do not set the status to available yet
		return result, nil
	}

	klog.V(3).Info("Synced up all machine API configurations")

	message := fmt.Sprintf("Cluster Machine API Operator is available at %s", optr.printOperandVersions())
	if err := optr.statusAvailable(message); err != nil {
		klog.Errorf("Error syncing ClusterOperatorStatus: %v", err)
		return reconcile.Result{}, fmt.Errorf("error syncing ClusterOperatorStatus: %v", err)
	}
	return reconcile.Result{}, nil
}

func (optr *Operator) checkRolloutStatus(config *OperatorConfig) (reconcile.Result, error) {
	// Check for machine-controllers deployment
	result, err := optr.checkDeploymentRolloutStatus(newDeployment(config, nil))
	if err != nil {
		return reconcile.Result{}, err
	}
	if result.Requeue || result.RequeueAfter > 0 {
		return result, nil
	}

	if config.Controllers.TerminationHandler != clusterAPIControllerNoOp {
		// Check for termination handler
		result, err := optr.checkDaemonSetRolloutStatus(newTerminationDaemonSet(config))
		if err != nil {
			return reconcile.Result{}, err
		}
		if result.Requeue || result.RequeueAfter > 0 {
			return result, nil
		}
	}

	return reconcile.Result{}, nil
}

func (optr *Operator) syncClusterAPIController(config *OperatorConfig) error {
	controllersDeployment := newDeployment(config, nil)

	// we watch some resources so that our deployment will redeploy without explicitly and carefully ordered resource creation
	inputHashes, err := resourcehash.MultipleObjectHashStringMapForObjectReferences(
		optr.kubeClient,
		resourcehash.NewObjectRef().ForConfigMap().InNamespace(config.TargetNamespace).Named(externalTrustBundleConfigMapName),
	)
	if err != nil {
		return fmt.Errorf("invalid dependency reference: %q", err)
	}
	ensureDependecyAnnotations(inputHashes, controllersDeployment)

	expectedGeneration := resourcemerge.ExpectedDeploymentGeneration(controllersDeployment, optr.generations)
	d, updated, err := resourceapply.ApplyDeployment(context.TODO(), optr.kubeClient.AppsV1(),
		events.NewLoggingEventRecorder(optr.name), controllersDeployment, expectedGeneration)
	if err != nil {
		return err
	}
	if updated {
		resourcemerge.SetDeploymentGeneration(&optr.generations, d)
	}

	return nil
}

func (optr *Operator) syncTerminationHandler(config *OperatorConfig) error {
	terminationDaemonSet := newTerminationDaemonSet(config)
	expectedGeneration := resourcemerge.ExpectedDaemonSetGeneration(terminationDaemonSet, optr.generations)
	ds, updated, err := resourceapply.ApplyDaemonSet(context.TODO(), optr.kubeClient.AppsV1(),
		events.NewLoggingEventRecorder(optr.name), terminationDaemonSet, expectedGeneration)
	if err != nil {
		return err
	}
	if updated {
		resourcemerge.SetDaemonSetGeneration(&optr.generations, ds)
	}
	return nil
}

func (optr *Operator) syncWebhookConfiguration() error {
	if err := optr.syncValidatingWebhook(); err != nil {
		return err
	}

	return optr.syncMutatingWebhook()
}

func (optr *Operator) syncValidatingWebhook() error {
	expectedGeneration := resourcemerge.ExpectedValidatingWebhooksConfiguration(mapiwebhooks.NewValidatingWebhookConfiguration().Name, optr.generations)
	validatingWebhook, updated, err := resourceapply.ApplyValidatingWebhookConfiguration(context.TODO(), optr.kubeClient.AdmissionregistrationV1(),
		events.NewLoggingEventRecorder(optr.name),
		mapiwebhooks.NewValidatingWebhookConfiguration(), expectedGeneration)
	if err != nil {
		return err
	}
	if updated {
		resourcemerge.SetValidatingWebhooksConfigurationGeneration(&optr.generations, validatingWebhook)
	}

	return nil
}

func (optr *Operator) syncMutatingWebhook() error {
	expectedGeneration := resourcemerge.ExpectedMutatingWebhooksConfiguration(mapiwebhooks.NewMutatingWebhookConfiguration().Name, optr.generations)
	validatingWebhook, updated, err := resourceapply.ApplyMutatingWebhookConfiguration(context.TODO(), optr.kubeClient.AdmissionregistrationV1(),
		events.NewLoggingEventRecorder(optr.name),
		mapiwebhooks.NewMutatingWebhookConfiguration(), expectedGeneration)
	if err != nil {
		return err
	}
	if updated {
		resourcemerge.SetMutatingWebhooksConfigurationGeneration(&optr.generations, validatingWebhook)
	}

	return nil
}

func (optr *Operator) checkDeploymentRolloutStatus(resource *appsv1.Deployment) (reconcile.Result, error) {
	d, err := optr.kubeClient.AppsV1().Deployments(resource.Namespace).Get(context.Background(), resource.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return reconcile.Result{}, fmt.Errorf("deployment %s is not found", resource.Name)
	}
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("getting Deployment %s during rollout: %v", resource.Name, err)
	}

	if d.DeletionTimestamp != nil {
		return reconcile.Result{}, fmt.Errorf("deployment %s is being deleted", resource.Name)
	}

	if d.Generation > d.Status.ObservedGeneration || d.Status.UpdatedReplicas != d.Status.Replicas || d.Status.UnavailableReplicas > 0 {
		klog.V(3).Infof("deployment %s is not ready. status: (replicas: %d, updated: %d, ready: %d, unavailable: %d)", d.Name, d.Status.Replicas, d.Status.UpdatedReplicas, d.Status.ReadyReplicas, d.Status.UnavailableReplicas)
		return reconcile.Result{Requeue: true, RequeueAfter: checkStatusRequeuePeriod}, nil
	}

	c := conditions.GetDeploymentCondition(d, appsv1.DeploymentAvailable)
	if c == nil {
		klog.V(3).Infof("deployment %s is not reporting available yet", resource.Name)
		return reconcile.Result{Requeue: true, RequeueAfter: checkStatusRequeuePeriod}, nil
	}

	if c.Status == corev1.ConditionFalse {
		klog.V(3).Infof("deployment %s is reporting available=false", resource.Name)
		return reconcile.Result{Requeue: true, RequeueAfter: checkStatusRequeuePeriod}, nil
	}

	if c.LastTransitionTime.Time.Add(deploymentMinimumAvailabilityTime).After(time.Now()) {
		klog.V(3).Infof("deployment %s has been available for less than %s", resource.Name, deploymentMinimumAvailabilityTime)
		// Requeue at the deploymentMinimumAvailabilityTime mark so we don't spam retries
		nextCheck := c.LastTransitionTime.Time.Add(deploymentMinimumAvailabilityTime).Sub(time.Now())
		return reconcile.Result{Requeue: true, RequeueAfter: nextCheck}, nil
	}

	return reconcile.Result{}, nil
}

func (optr *Operator) checkDaemonSetRolloutStatus(resource *appsv1.DaemonSet) (reconcile.Result, error) {
	d, err := optr.kubeClient.AppsV1().DaemonSets(resource.Namespace).Get(context.Background(), resource.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return reconcile.Result{}, fmt.Errorf("daemonset %s is not found", resource.Name)
	}
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("getting DaemonSet %s during rollout: %v", resource.Name, err)
	}

	if d.DeletionTimestamp != nil {
		return reconcile.Result{}, fmt.Errorf("daemonset %s is being deleted", resource.Name)
	}

	if d.Generation > d.Status.ObservedGeneration || d.Status.UpdatedNumberScheduled != d.Status.DesiredNumberScheduled || d.Status.NumberUnavailable > 0 {
		klog.V(3).Infof("daemonset %s is not ready. status: (desired: %d, updated: %d, available: %d, unavailable: %d)", d.Name, d.Status.DesiredNumberScheduled, d.Status.UpdatedNumberScheduled, d.Status.NumberAvailable, d.Status.NumberUnavailable)
		return reconcile.Result{Requeue: true, RequeueAfter: checkStatusRequeuePeriod}, nil
	}

	return reconcile.Result{}, nil
}

func newDeployment(config *OperatorConfig, features map[string]bool) *appsv1.Deployment {
	replicas := int32(1)
	template := newPodTemplateSpec(config, features)

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "machine-api-controllers",
			Namespace: config.TargetNamespace,
			Annotations: map[string]string{
				maoOwnedAnnotation: "",
			},
			Labels: map[string]string{
				"api":     "clusterapi",
				"k8s-app": "controller",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"api":     "clusterapi",
					"k8s-app": "controller",
				},
			},
			Template: *template,
		},
	}
}

// List of the volumes needed by newKubeProxyContainer
func newRBACConfigVolumes() []corev1.Volume {
	var readOnly int32 = 420
	return []corev1.Volume{
		{
			Name: kubeRBACConfigName,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "kube-rbac-proxy",
					},
					DefaultMode: pointer.Int32Ptr(readOnly),
				},
			},
		},
		{
			Name: certStoreName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  certStoreName,
					DefaultMode: pointer.Int32Ptr(readOnly),
				},
			},
		},
		{
			Name: "trusted-ca",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					Items: []corev1.KeyToPath{{Key: "ca-bundle.crt", Path: "tls-ca-bundle.pem"}},
					LocalObjectReference: corev1.LocalObjectReference{
						Name: externalTrustBundleConfigMapName,
					},
					Optional: pointer.BoolPtr(true),
				},
			},
		},
	}
}

func newPodTemplateSpec(config *OperatorConfig, features map[string]bool) *corev1.PodTemplateSpec {
	containers := newContainers(config, features)
	proxyContainers := newKubeProxyContainers(config.Controllers.KubeRBACProxy)
	tolerations := []corev1.Toleration{
		{
			Key:    "node-role.kubernetes.io/master",
			Effect: corev1.TaintEffectNoSchedule,
		},
		{
			Key:      "CriticalAddonsOnly",
			Operator: corev1.TolerationOpExists,
		},
		{
			Key:               "node.kubernetes.io/not-ready",
			Effect:            corev1.TaintEffectNoExecute,
			Operator:          corev1.TolerationOpExists,
			TolerationSeconds: pointer.Int64Ptr(120),
		},
		{
			Key:               "node.kubernetes.io/unreachable",
			Effect:            corev1.TaintEffectNoExecute,
			Operator:          corev1.TolerationOpExists,
			TolerationSeconds: pointer.Int64Ptr(120),
		},
	}

	var readOnly int32 = 420
	volumes := []corev1.Volume{
		{
			Name: "cert",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  "machine-api-operator-webhook-cert",
					DefaultMode: pointer.Int32Ptr(readOnly),
					Items: []corev1.KeyToPath{
						{
							Key:  "tls.crt",
							Path: "tls.crt",
						},
						{
							Key:  "tls.key",
							Path: "tls.key",
						},
					},
				},
			},
		},
		{
			Name: "bound-sa-token",
			VolumeSource: corev1.VolumeSource{
				Projected: &corev1.ProjectedVolumeSource{
					Sources: []corev1.VolumeProjection{
						{
							ServiceAccountToken: &corev1.ServiceAccountTokenProjection{
								Audience: "uccp",
								Path:     "token",
							},
						},
					},
				},
			},
		},
	}
	volumes = append(volumes, newRBACConfigVolumes()...)

	return &corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: commonPodTemplateAnnotations,
			Labels: map[string]string{
				"api":     "clusterapi",
				"k8s-app": "controller",
			},
		},
		Spec: corev1.PodSpec{
			Containers:         append(containers, proxyContainers...),
			PriorityClassName:  "system-node-critical",
			NodeSelector:       map[string]string{"node-role.kubernetes.io/master": ""},
			ServiceAccountName: "machine-api-controllers",
			Tolerations:        tolerations,
			Volumes:            volumes,
		},
	}
}

func getProxyArgs(config *OperatorConfig) []corev1.EnvVar {
	var envVars []corev1.EnvVar

	if config.Proxy == nil {
		return envVars
	}
	if config.Proxy.Status.HTTPProxy != "" {
		envVars = append(envVars, corev1.EnvVar{
			Name:  "HTTP_PROXY",
			Value: config.Proxy.Spec.HTTPProxy,
		})
	}
	if config.Proxy.Status.HTTPSProxy != "" {
		envVars = append(envVars, corev1.EnvVar{
			Name:  "HTTPS_PROXY",
			Value: config.Proxy.Spec.HTTPSProxy,
		})
	}
	if config.Proxy.Status.NoProxy != "" {
		envVars = append(envVars, corev1.EnvVar{
			Name:  "NO_PROXY",
			Value: config.Proxy.Status.NoProxy,
		})
	}
	return envVars
}

func newContainers(config *OperatorConfig, features map[string]bool) []corev1.Container {
	resources := corev1.ResourceRequirements{
		Requests: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceMemory: resource.MustParse("20Mi"),
			corev1.ResourceCPU:    resource.MustParse("10m"),
		},
	}
	args := []string{
		"--logtostderr=true",
		"--v=3",
		"--leader-elect=true",
		"--leader-elect-lease-duration=120s",
		fmt.Sprintf("--namespace=%s", config.TargetNamespace),
	}

	proxyEnvArgs := getProxyArgs(config)

	containers := []corev1.Container{
		{
			Name:      "machineset-controller",
			Image:     config.Controllers.MachineSet,
			Command:   []string{"/machineset-controller"},
			Args:      args,
			Resources: resources,
			Env:       proxyEnvArgs,
			Ports: []corev1.ContainerPort{
				{
					Name:          "webhook-server",
					ContainerPort: 8443,
				},
				{
					Name:          "healthz",
					ContainerPort: defaultMachineSetHealthPort,
				},
			},
			ReadinessProbe: &corev1.Probe{
				Handler: corev1.Handler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: "/healthz",
						Port: intstr.Parse("healthz"),
					},
				},
			},
			LivenessProbe: &corev1.Probe{
				Handler: corev1.Handler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: "/readyz",
						Port: intstr.Parse("healthz"),
					},
				},
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					MountPath: "/etc/machine-api-operator/tls",
					Name:      "cert",
					ReadOnly:  true,
				},
			},
		},
		{
			Name:      "machine-controller",
			Image:     config.Controllers.Provider,
			Command:   []string{"/machine-controller-manager"},
			Args:      args,
			Resources: resources,
			Env: append(proxyEnvArgs, corev1.EnvVar{
				Name: "NODE_NAME",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "spec.nodeName",
					},
				},
			}),
			Ports: []corev1.ContainerPort{{
				Name:          "healthz",
				ContainerPort: defaultMachineHealthPort,
			}},
			ReadinessProbe: &corev1.Probe{
				Handler: corev1.Handler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: "/healthz",
						Port: intstr.Parse("healthz"),
					},
				},
			},
			LivenessProbe: &corev1.Probe{
				Handler: corev1.Handler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: "/readyz",
						Port: intstr.Parse("healthz"),
					},
				},
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					MountPath: "/etc/pki/ca-trust/extracted/pem",
					Name:      "trusted-ca",
					ReadOnly:  true,
				},
				{
					MountPath: "/var/run/secrets/uccp/serviceaccount",
					Name:      "bound-sa-token",
					ReadOnly:  true,
				},
			},
		},
		{
			Name:      "nodelink-controller",
			Image:     config.Controllers.NodeLink,
			Command:   []string{"/nodelink-controller"},
			Args:      args,
			Env:       proxyEnvArgs,
			Resources: resources,
		},
		{
			Name:      "machine-healthcheck-controller",
			Image:     config.Controllers.MachineHealthCheck,
			Command:   []string{"/machine-healthcheck"},
			Args:      args,
			Env:       proxyEnvArgs,
			Resources: resources,
			Ports: []corev1.ContainerPort{
				{
					Name:          "healthz",
					ContainerPort: defaultMachineHealthCheckHealthPort,
				},
			},
			ReadinessProbe: &corev1.Probe{
				Handler: corev1.Handler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: "/healthz",
						Port: intstr.Parse("healthz"),
					},
				},
			},
			LivenessProbe: &corev1.Probe{
				Handler: corev1.Handler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: "/readyz",
						Port: intstr.Parse("healthz"),
					},
				},
			},
		},
	}
	return containers
}

func newKubeProxyContainers(image string) []corev1.Container {
	return []corev1.Container{
		newKubeProxyContainer(image, "machineset-mtrc", metrics.DefaultMachineSetMetricsAddress, machineSetExposeMetricsPort),
		newKubeProxyContainer(image, "machine-mtrc", metrics.DefaultMachineMetricsAddress, machineExposeMetricsPort),
		newKubeProxyContainer(image, "mhc-mtrc", metrics.DefaultHealthCheckMetricsAddress, machineHealthCheckExposeMetricsPort),
	}
}

func newKubeProxyContainer(image, portName, upstreamPort string, exposePort int32) corev1.Container {
	configMountPath := "/etc/kube-rbac-proxy"
	tlsCertMountPath := "/etc/tls/private"
	resources := corev1.ResourceRequirements{
		Requests: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceMemory: resource.MustParse("20Mi"),
			corev1.ResourceCPU:    resource.MustParse("10m"),
		},
	}
	args := []string{
		fmt.Sprintf("--secure-listen-address=0.0.0.0:%d", exposePort),
		fmt.Sprintf("--upstream=http://localhost%s", upstreamPort),
		fmt.Sprintf("--config-file=%s/config-file.yaml", configMountPath),
		fmt.Sprintf("--tls-cert-file=%s/tls.crt", tlsCertMountPath),
		fmt.Sprintf("--tls-private-key-file=%s/tls.key", tlsCertMountPath),
		"--tls-cipher-suites=TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305",
		"--logtostderr=true",
		"--v=3",
	}
	ports := []corev1.ContainerPort{{
		Name:          portName,
		ContainerPort: exposePort,
	}}

	return corev1.Container{
		Name:      fmt.Sprintf("kube-rbac-proxy-%s", portName),
		Image:     image,
		Args:      args,
		Resources: resources,
		Ports:     ports,
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      kubeRBACConfigName,
				MountPath: configMountPath,
			},
			{
				Name:      certStoreName,
				MountPath: tlsCertMountPath,
			}},
	}
}

func newTerminationDaemonSet(config *OperatorConfig) *appsv1.DaemonSet {
	template := newTerminationPodTemplateSpec(config)

	return &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      machineAPITerminationHandler,
			Namespace: config.TargetNamespace,
			Annotations: map[string]string{
				maoOwnedAnnotation: "",
			},
			Labels: map[string]string{
				"api":     "clusterapi",
				"k8s-app": "termination-handler",
			},
		},
		Spec: appsv1.DaemonSetSpec{
			UpdateStrategy: appsv1.DaemonSetUpdateStrategy{
				RollingUpdate: &appsv1.RollingUpdateDaemonSet{
					MaxUnavailable: &daemonsetMaxUnavailable,
				},
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"api":     "clusterapi",
					"k8s-app": "termination-handler",
				},
			},
			Template: *template,
		},
	}
}

func newTerminationPodTemplateSpec(config *OperatorConfig) *corev1.PodTemplateSpec {
	containers := newTerminationContainers(config)

	return &corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: commonPodTemplateAnnotations,
			Labels: map[string]string{
				"api":     "clusterapi",
				"k8s-app": "termination-handler",
			},
		},
		Spec: corev1.PodSpec{
			Containers:                   containers,
			PriorityClassName:            "system-node-critical",
			NodeSelector:                 map[string]string{machinecontroller.MachineInterruptibleInstanceLabelName: ""},
			ServiceAccountName:           machineAPITerminationHandler,
			AutomountServiceAccountToken: pointer.BoolPtr(false),
			HostNetwork:                  true,
			Volumes: []corev1.Volume{
				{
					Name: "kubeconfig",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: hostKubeConfigPath,
						},
					},
				},
				{
					Name: "pki",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: hostKubePKIPath,
						},
					},
				},
			},
		},
	}
}

func newTerminationContainers(config *OperatorConfig) []corev1.Container {
	resources := corev1.ResourceRequirements{
		Requests: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceMemory: resource.MustParse("20Mi"),
			corev1.ResourceCPU:    resource.MustParse("10m"),
		},
	}
	terminationArgs := []string{
		"--logtostderr=true",
		"--v=3",
		"--node-name=$(NODE_NAME)",
		fmt.Sprintf("--namespace=%s", config.TargetNamespace),
		"--poll-interval-seconds=5",
	}

	proxyEnvArgs := getProxyArgs(config)

	return []corev1.Container{
		{
			Name:      "termination-handler",
			Image:     config.Controllers.TerminationHandler,
			Command:   []string{"/termination-handler"},
			Args:      terminationArgs,
			Resources: resources,
			Env: append(proxyEnvArgs, corev1.EnvVar{
				Name:  "KUBECONFIG",
				Value: hostKubeConfigPath,
			}, corev1.EnvVar{
				Name: "NODE_NAME",
				ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "spec.nodeName",
				},
				},
			}),

			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "kubeconfig",
					MountPath: hostKubeConfigPath,
					ReadOnly:  true,
				},
				{
					Name:      "pki",
					MountPath: hostKubePKIPath,
					ReadOnly:  true,
				},
			},
		},
	}
}

// ensureDependecyAnnotations uses inputHash map of external dependencies to force new generation of the deployment
// triggering the Kubernetes rollout as defined when the inputHash changes by adding it annotation to the deployment object.
func ensureDependecyAnnotations(inputHashes map[string]string, deployment *appsv1.Deployment) {
	for k, v := range inputHashes {
		annotationKey := fmt.Sprintf("operator.uccp.io/dep-%s", k)
		if deployment.Annotations == nil {
			deployment.Annotations = map[string]string{}
		}
		if deployment.Spec.Template.Annotations == nil {
			deployment.Spec.Template.Annotations = map[string]string{}
		}
		deployment.Annotations[annotationKey] = v
		deployment.Spec.Template.Annotations[annotationKey] = v
	}
}
