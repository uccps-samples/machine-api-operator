package operator

import (
	"testing"

	. "github.com/onsi/gomega"
	osconfigv1 "github.com/uccps-samples/api/config/v1"
	fakeos "github.com/uccps-samples/client-go/config/clientset/versioned/fake"
	"golang.org/x/net/context"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	fakedynamic "k8s.io/client-go/dynamic/fake"
	fakekube "k8s.io/client-go/kubernetes/fake"
	"sigs.k8s.io/yaml"
)

var yamlContent = `
apiVersion: metal3.io/v1alpha1
kind: Provisioning
metadata:
  name: test
spec:
  provisioningInterface: "ensp0"
  provisioningIP: "172.30.20.3"
  provisioningNetworkCIDR: "172.30.20.0/24"
  provisioningDHCPExternal: false
  provisioningDHCPRange: "172.30.20.11, 172.30.20.101"
  provisioningOSDownloadURL: "http://172.22.0.1/images/rhcos-44.81.202001171431.0-openstack.x86_64.qcow2.gz?sha256=e98f83a2b9d4043719664a2be75fe8134dc6ca1fdbde807996622f8cc7ecd234"
`
var (
	expectedProvisioningInterface   = "ensp0"
	expectedProvisioningIP          = "172.30.20.3"
	expectedProvisioningNetworkCIDR = "172.30.20.0/24"
	expectedProvisioningDHCPRange   = "172.30.20.11, 172.30.20.101"
	expectedOSImageURL              = "http://172.22.0.1/images/rhcos-44.81.202001171431.0-openstack.x86_64.qcow2.gz?sha256=e98f83a2b9d4043719664a2be75fe8134dc6ca1fdbde807996622f8cc7ecd234"
	expectedProvisioningIPCIDR      = "172.30.20.3/24"
	expectedDeployKernelURL         = "http://172.30.20.3:6180/images/ironic-python-agent.kernel"
	expectedDeployRamdiskURL        = "http://172.30.20.3:6180/images/ironic-python-agent.initramfs"
	expectedIronicEndpoint          = "http://172.30.20.3:6385/v1/"
	expectedIronicInspectorEndpoint = "http://172.30.20.3:5050/v1/"
	expectedHttpPort                = "6180"
	expectedProvisioningNetwork     = "Managed"
)

var (
	Managed = `
apiVersion: metal3.io/v1alpha1
kind: Provisioning
metadata:
  name: test
spec:
  provisioningInterface: "ensp0"
  provisioningIP: "172.30.20.3"
  provisioningNetworkCIDR: "172.30.20.0/24"
  provisioningDHCPExternal: false
  provisioningDHCPRange: "172.30.20.11, 172.30.20.101"
  provisioningOSDownloadURL: "http://172.22.0.1/images/rhcos-44.81.202001171431.0-openstack.x86_64.qcow2.gz?sha256=e98f83a2b9d4043719664a2be75fe8134dc6ca1fdbde807996622f8cc7ecd234"
  provisioningNetwork: "Managed"
`
	Unmanaged = `
apiVersion: metal3.io/v1alpha1
kind: Provisioning
metadata:
  name: test
spec:
  provisioningInterface: "ensp0"
  provisioningIP: "172.30.20.3"
  provisioningNetworkCIDR: "172.30.20.0/24"
  provisioningDHCPRange: ""
  provisioningOSDownloadURL: "http://172.22.0.1/images/rhcos-44.81.202001171431.0-openstack.x86_64.qcow2.gz?sha256=e98f83a2b9d4043719664a2be75fe8134dc6ca1fdbde807996622f8cc7ecd234"
  provisioningNetwork: "Unmanaged"
`
	Disabled = `
apiVersion: metal3.io/v1alpha1
kind: Provisioning
metadata:
  name: test
spec:
  provisioningInterface: ""
  provisioningIP: "172.30.20.3"
  provisioningNetworkCIDR: "172.30.20.0/24"
  provisioningDHCPExternal: false
  provisioningDHCPRange: ""
  provisioningOSDownloadURL: "http://172.22.0.1/images/rhcos-44.81.202001171431.0-openstack.x86_64.qcow2.gz?sha256=e98f83a2b9d4043719664a2be75fe8134dc6ca1fdbde807996622f8cc7ecd234"
  provisioningNetwork: "Disabled"
`
)

func TestGenerateRandomPassword(t *testing.T) {
	pwd, err := generateRandomPassword()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if pwd == "" {
		t.Errorf("Expected a valid string but got null")
	}
}

func newOperatorWithBaremetalConfig() *OperatorConfig {
	return &OperatorConfig{
		targetNamespace,
		Controllers{
			"docker.io/openshift/origin-aws-machine-controllers:v4.0.0",
			"docker.io/openshift/origin-machine-api-operator:v4.0.0",
			"docker.io/openshift/origin-machine-api-operator:v4.0.0",
			"docker.io/openshift/origin-machine-api-operator:v4.0.0",
			"docker.io/openshift/origin-kube-rbac-proxy:v4.0.0",
			"docker.io/openshift/origin-aws-machine-controllers:v4.0.0",
		},
		BaremetalControllers{
			"quay.io/openshift/origin-baremetal-operator:v4.2.0",
			"quay.io/openshift/origin-ironic:v4.2.0",
			"quay.io/openshift/origin-ironic-inspector:v4.2.0",
			"quay.io/openshift/origin-ironic-ipa-downloader:v4.2.0",
			"quay.io/openshift/origin-ironic-machine-os-downloader:v4.2.0",
			"quay.io/openshift/origin-ironic-static-ip-manager:v4.2.0",
		},
		&osconfigv1.Proxy{},
	}
}

//Testing the case where the password does already exist
func TestCreateMariadbPasswordSecret(t *testing.T) {
	kubeClient := fakekube.NewSimpleClientset(nil...)
	operatorConfig := newOperatorWithBaremetalConfig()
	client := kubeClient.CoreV1()

	// First create a mariadb password secret
	if err := createMariadbPasswordSecret(kubeClient.CoreV1(), operatorConfig); err != nil {
		t.Fatalf("Failed to create first Mariadb password. %s ", err)
	}
	// Read and get Mariadb password from Secret just created.
	oldMaridbPassword, err := client.Secrets(operatorConfig.TargetNamespace).Get(context.Background(), baremetalSecretName, metav1.GetOptions{})
	if err != nil {
		t.Fatal("Failure getting the first Mariadb password that just got created.")
	}
	oldPassword, ok := oldMaridbPassword.StringData[baremetalSecretKey]
	if !ok || oldPassword == "" {
		t.Fatal("Failure reading first Mariadb password from Secret.")
	}

	// The pasword definitely exists. Try creating again.
	if err := createMariadbPasswordSecret(kubeClient.CoreV1(), operatorConfig); err != nil {
		t.Fatal("Failure creating second Mariadb password.")
	}
	newMaridbPassword, err := client.Secrets(operatorConfig.TargetNamespace).Get(context.Background(), baremetalSecretName, metav1.GetOptions{})
	if err != nil {
		t.Fatal("Failure getting the second Mariadb password.")
	}
	newPassword, ok := newMaridbPassword.StringData[baremetalSecretKey]
	if !ok || newPassword == "" {
		t.Fatal("Failure reading second Mariadb password from Secret.")
	}
	if oldPassword != newPassword {
		t.Fatalf("Both passwords do not match.")
	} else {
		t.Logf("First Mariadb password is being preserved over re-creation as expected.")
	}
}

func TestGetBaremetalProvisioningConfig(t *testing.T) {
	testConfigResource := "test"
	u := &unstructured.Unstructured{Object: map[string]interface{}{}}
	if err := yaml.Unmarshal([]byte(yamlContent), &u); err != nil {
		t.Errorf("failed to unmarshall input yaml content:%v", err)
	}
	dynamicClient := fakedynamic.NewSimpleDynamicClient(runtime.NewScheme(), u)
	baremetalConfig, err := getBaremetalProvisioningConfig(dynamicClient, testConfigResource)
	if err != nil {
		t.Logf("Unstructed Config:  %+v", u)
		t.Fatalf("Failed to get Baremetal Provisioning Interface from CR %s", testConfigResource)
	}
	if baremetalConfig.ProvisioningInterface != expectedProvisioningInterface ||
		baremetalConfig.ProvisioningIp != expectedProvisioningIP ||
		baremetalConfig.ProvisioningNetworkCIDR != expectedProvisioningNetworkCIDR ||
		baremetalConfig.ProvisioningDHCPRange != expectedProvisioningDHCPRange ||
		baremetalConfig.ProvisioningNetwork != provisioningNetworkManaged {
		t.Logf("Expected: ProvisioningInterface: %s, ProvisioningIP: %s, ProvisioningNetworkCIDR: %s, ProvisioningNetwork: %s, expectedProvisioningDHCPRange: %s, Got: %+v", expectedProvisioningInterface, expectedProvisioningIP, expectedProvisioningNetworkCIDR, expectedProvisioningNetwork, expectedProvisioningDHCPRange, baremetalConfig)
		t.Fatalf("failed getBaremetalProvisioningConfig. One or more BaremetalProvisioningConfig items do not match the expected config.")
	}
}

func TestGetIncorrectBaremetalProvisioningCR(t *testing.T) {
	incorrectConfigResource := "test1"
	u := &unstructured.Unstructured{Object: map[string]interface{}{}}
	if err := yaml.Unmarshal([]byte(yamlContent), &u); err != nil {
		t.Errorf("failed to unmarshall input yaml content:%v", err)
	}
	dynamicClient := fakedynamic.NewSimpleDynamicClient(runtime.NewScheme(), u)
	baremetalConfig, err := getBaremetalProvisioningConfig(dynamicClient, incorrectConfigResource)
	if err == nil && baremetalConfig == nil {
		t.Logf("Unable to get Baremetal Provisioning Config from CR %s as expected", incorrectConfigResource)
	} else {
		t.Errorf("BaremetalProvisioningConfig is not expected to be set.")
	}
}

func TestGetMetal3DeploymentConfig(t *testing.T) {
	testConfigResource := "test"
	u := &unstructured.Unstructured{Object: map[string]interface{}{}}
	if err := yaml.Unmarshal([]byte(yamlContent), &u); err != nil {
		t.Errorf("failed to unmarshall input yaml content:%v", err)
	}
	dynamicClient := fakedynamic.NewSimpleDynamicClient(runtime.NewScheme(), u)
	baremetalConfig, err := getBaremetalProvisioningConfig(dynamicClient, testConfigResource)
	if err != nil {
		t.Logf("Unstructed Config:  %+v", u)
		t.Errorf("Failed to get Baremetal Provisioning Config from CR %s", testConfigResource)
	}
	actualCacheURL := getMetal3DeploymentConfig("CACHEURL", *baremetalConfig)
	if actualCacheURL != nil {
		t.Errorf("CacheURL is found to be %s. CACHEURL is not expected.", *actualCacheURL)
	} else {
		t.Logf("CacheURL is not available as expected.")
	}
	actualOSImageURL := getMetal3DeploymentConfig("RHCOS_IMAGE_URL", *baremetalConfig)
	if actualOSImageURL != nil {
		t.Logf("Actual OS Image Download URL is %s, Expected is %s", *actualOSImageURL, expectedOSImageURL)
		if *actualOSImageURL != expectedOSImageURL {
			t.Errorf("Actual %s and Expected %s OS Image Download URLs do not match", *actualOSImageURL, expectedOSImageURL)
		}
	} else {
		t.Errorf("OS Image Download URL is not available.")
	}
	actualProvisioningIPCIDR := getMetal3DeploymentConfig("PROVISIONING_IP", *baremetalConfig)
	if actualProvisioningIPCIDR != nil {
		t.Logf("Actual ProvisioningIP with CIDR is %s, Expected is %s", *actualProvisioningIPCIDR, expectedProvisioningIPCIDR)
		if *actualProvisioningIPCIDR != expectedProvisioningIPCIDR {
			t.Errorf("Actual %s and Expected %s Provisioning IPs with CIDR do not match", *actualProvisioningIPCIDR, expectedProvisioningIPCIDR)
		}
	} else {
		t.Errorf("Provisioning IP with CIDR is not available.")
	}
	actualProvisioningInterface := getMetal3DeploymentConfig("PROVISIONING_INTERFACE", *baremetalConfig)
	if actualProvisioningInterface != nil {
		t.Logf("Actual Provisioning Interface is %s, Expected is %s", *actualProvisioningInterface, expectedProvisioningInterface)
		if *actualProvisioningInterface != expectedProvisioningInterface {
			t.Errorf("Actual %s and Expected %s Provisioning Interfaces do not match", *actualProvisioningIPCIDR, expectedProvisioningIPCIDR)
		}
	} else {
		t.Errorf("Provisioning Interface is not available.")
	}
	actualDeployKernelURL := getMetal3DeploymentConfig("DEPLOY_KERNEL_URL", *baremetalConfig)
	if actualDeployKernelURL != nil {
		t.Logf("Actual Deploy Kernel URL is %s, Expected is %s", *actualDeployKernelURL, expectedDeployKernelURL)
		if *actualDeployKernelURL != expectedDeployKernelURL {
			t.Errorf("Actual %s and Expected %s Deploy Kernel URLs do not match", *actualDeployKernelURL, expectedDeployKernelURL)
		}
	} else {
		t.Errorf("Deploy Kernel URL is not available.")
	}
	actualDeployRamdiskURL := getMetal3DeploymentConfig("DEPLOY_RAMDISK_URL", *baremetalConfig)
	if actualDeployRamdiskURL != nil {
		t.Logf("Actual Deploy Ramdisk URL is %s, Expected is %s", *actualDeployRamdiskURL, expectedDeployRamdiskURL)
		if *actualDeployRamdiskURL != expectedDeployRamdiskURL {
			t.Errorf("Actual %s and Expected %s Deploy Ramdisk URLs do not match", *actualDeployRamdiskURL, expectedDeployRamdiskURL)
		}
	} else {
		t.Errorf("Deploy Ramdisk URL is not available.")
	}
	actualIronicEndpoint := getMetal3DeploymentConfig("IRONIC_ENDPOINT", *baremetalConfig)
	if actualIronicEndpoint != nil {
		t.Logf("Actual Ironic Endpoint is %s, Expected is %s", *actualIronicEndpoint, expectedIronicEndpoint)
		if *actualIronicEndpoint != expectedIronicEndpoint {
			t.Errorf("Actual %s and Expected %s Ironic Endpoints do not match", *actualIronicEndpoint, expectedIronicEndpoint)
		}
	} else {
		t.Errorf("Ironic Endpoint is not available.")
	}
	actualIronicInspectorEndpoint := getMetal3DeploymentConfig("IRONIC_INSPECTOR_ENDPOINT", *baremetalConfig)
	if actualIronicInspectorEndpoint != nil {
		t.Logf("Actual Ironic Inspector Endpoint is %s, Expected is %s", *actualIronicInspectorEndpoint, expectedIronicInspectorEndpoint)
		if *actualIronicInspectorEndpoint != expectedIronicInspectorEndpoint {
			t.Errorf("Actual %s and Expected %s Ironic Inspector Endpoints do not match", *actualIronicInspectorEndpoint, expectedIronicInspectorEndpoint)
		}
	} else {
		t.Errorf("Ironic Inspector Endpoint is not available.")
	}
	actualHttpPort := getMetal3DeploymentConfig("HTTP_PORT", *baremetalConfig)
	t.Logf("Actual Http Port is %s, Expected is %s", *actualHttpPort, expectedHttpPort)
	if *actualHttpPort != expectedHttpPort {
		t.Errorf("Actual %s and Expected %s Http Ports do not match", *actualHttpPort, expectedHttpPort)
	}
	actualDHCPRange := getMetal3DeploymentConfig("DHCP_RANGE", *baremetalConfig)
	if actualDHCPRange != nil {
		t.Logf("Actual DHCP Range is %s, Expected is %s", *actualDHCPRange, expectedProvisioningDHCPRange)
		if *actualDHCPRange != expectedProvisioningDHCPRange {
			t.Errorf("Actual %s and Expected %s DHCP Range do not match", *actualDHCPRange, expectedProvisioningDHCPRange)
		}
	} else {
		t.Errorf("Provisioning DHCP Range is not available.")
	}
}

// Test interpretation of different provisioning config
func TestBaremetalProvisionigConfig(t *testing.T) {
	testConfigResource := "test"
	configNames := []string{"PROVISIONING_IP", "PROVISIONING_INTERFACE", "DEPLOY_KERNEL_URL", "DEPLOY_RAMDISK_URL", "IRONIC_ENDPOINT", "IRONIC_INSPECTOR_ENDPOINT", "HTTP_PORT", "DHCP_RANGE", "RHCOS_IMAGE_URL"}

	testCases := []struct {
		name                 string
		config               string
		expectedConfigValues []string
	}{
		{
			name:                 "Managed",
			config:               Managed,
			expectedConfigValues: []string{"172.30.20.3/24", "ensp0", "http://172.30.20.3:6180/images/ironic-python-agent.kernel", "http://172.30.20.3:6180/images/ironic-python-agent.initramfs", "http://172.30.20.3:6385/v1/", "http://172.30.20.3:5050/v1/", expectedHttpPort, "172.30.20.11, 172.30.20.101", "http://172.22.0.1/images/rhcos-44.81.202001171431.0-openstack.x86_64.qcow2.gz?sha256=e98f83a2b9d4043719664a2be75fe8134dc6ca1fdbde807996622f8cc7ecd234"},
		},
		{
			name:                 "Unmanaged",
			config:               Unmanaged,
			expectedConfigValues: []string{"172.30.20.3/24", "ensp0", "http://172.30.20.3:6180/images/ironic-python-agent.kernel", "http://172.30.20.3:6180/images/ironic-python-agent.initramfs", "http://172.30.20.3:6385/v1/", "http://172.30.20.3:5050/v1/", expectedHttpPort, "", "http://172.22.0.1/images/rhcos-44.81.202001171431.0-openstack.x86_64.qcow2.gz?sha256=e98f83a2b9d4043719664a2be75fe8134dc6ca1fdbde807996622f8cc7ecd234"},
		},
		{
			name:                 "Disabled",
			config:               Disabled,
			expectedConfigValues: []string{"172.30.20.3/24", "", "http://172.30.20.3:6180/images/ironic-python-agent.kernel", "http://172.30.20.3:6180/images/ironic-python-agent.initramfs", "http://172.30.20.3:6385/v1/", "http://172.30.20.3:5050/v1/", expectedHttpPort, "", "http://172.22.0.1/images/rhcos-44.81.202001171431.0-openstack.x86_64.qcow2.gz?sha256=e98f83a2b9d4043719664a2be75fe8134dc6ca1fdbde807996622f8cc7ecd234"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			t.Logf("Testing tc : %s", tc.name)
			u := &unstructured.Unstructured{Object: map[string]interface{}{}}
			if err := yaml.Unmarshal([]byte(tc.config), &u); err != nil {
				t.Errorf("failed to unmarshall input yaml content:%v", err)
			}
			dynamicClient := fakedynamic.NewSimpleDynamicClient(runtime.NewScheme(), u)
			baremetalConfig, err := getBaremetalProvisioningConfig(dynamicClient, testConfigResource)
			if err != nil {
				t.Errorf("getBaremetalProvisioningConfig returned err: %v", err)
			}
			g.Expect(err).To(BeNil())

			for i, envVar := range configNames {
				actualConfigValue := getMetal3DeploymentConfig(envVar, *baremetalConfig)

				g.Expect(*actualConfigValue).To(Equal(tc.expectedConfigValues[i]))
			}
		})
	}
}

func TestSyncBaremetalControllers(t *testing.T) {
	operatorConfig := newOperatorWithBaremetalConfig()

	stop := make(chan struct{})
	defer close(stop)
	optr := newFakeOperator(nil, nil, stop)

	u := &unstructured.Unstructured{Object: map[string]interface{}{}}
	if err := yaml.Unmarshal([]byte(yamlContent), &u); err != nil {
		t.Errorf("failed to unmarshall input yaml content:%v", err)
	}
	dynamicClient := fakedynamic.NewSimpleDynamicClient(runtime.NewScheme(), u)

	optr.dynamicClient = dynamicClient

	testCases := []struct {
		name          string
		configCRName  string
		expectedError error
	}{
		{
			name:          "InvalidCRName",
			configCRName:  "baremetal-disabled",
			expectedError: apierrors.NewNotFound(provisioningGR, "baremetal-disabled"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			err := optr.syncBaremetalControllers(operatorConfig, tc.configCRName)
			g.Expect(err).To(Equal(tc.expectedError))
		})
	}
}

func TestCheckMetal3DeploymentOwned(t *testing.T) {
	kubeClient := fakekube.NewSimpleClientset(nil...)
	operatorConfig := newOperatorWithBaremetalConfig()
	client := kubeClient.AppsV1()

	testCases := []struct {
		testCase      string
		deployment    *appsv1.Deployment
		expected      bool
		expectedError bool
	}{
		{
			testCase: "Only maoOwnedAnnotation",
			deployment: &appsv1.Deployment{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Deployment",
					APIVersion: "apps/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: baremetalDeploymentName,
					Annotations: map[string]string{
						maoOwnedAnnotation: "",
					},
				},
			},
			expected: true,
		},
		{
			testCase: "Only cboOwnedAnnotation",
			deployment: &appsv1.Deployment{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Deployment",
					APIVersion: "apps/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: baremetalDeploymentName,
					Annotations: map[string]string{
						cboOwnedAnnotation: "",
					},
				},
			},
			expected: false,
		},
		{
			testCase: "Both cboOwnedAnnotation and maoOwnedAnnotation",
			deployment: &appsv1.Deployment{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Deployment",
					APIVersion: "apps/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: baremetalDeploymentName,
					Annotations: map[string]string{
						cboOwnedAnnotation: "",
						maoOwnedAnnotation: "",
					},
				},
			},
			expected: false,
		},
		{
			testCase: "No cboOwnedAnnotation or maoOwnedAnnotation",
			deployment: &appsv1.Deployment{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Deployment",
					APIVersion: "apps/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:        baremetalDeploymentName,
					Annotations: map[string]string{},
				},
			},
			expected: true,
		},
	}
	for _, tc := range testCases {
		t.Run(string(tc.testCase), func(t *testing.T) {

			_, err := client.Deployments("test-namespace").Create(context.Background(), tc.deployment, metav1.CreateOptions{})
			if err != nil {
				t.Fatalf("Could not create metal3 test deployment.\n")
			}
			maoOwned, err := checkMetal3DeploymentMAOOwned(client, operatorConfig)
			if maoOwned != tc.expected {
				t.Errorf("Expected: %v, got: %v", tc.expected, maoOwned)
			}
			if tc.expectedError != (err != nil) {
				t.Errorf("ExpectedError: %v, got: %v", tc.expectedError, err)
			}
			err = client.Deployments("test-namespace").Delete(context.Background(), baremetalDeploymentName, metav1.DeleteOptions{})
			if err != nil {
				t.Errorf("Could not delete metal3 test deployment.\n")
			}
		})
	}

}

func TestCheckForBaremetalClusterOperator(t *testing.T) {
	testCases := []struct {
		testCase        string
		clusterOperator *osconfigv1.ClusterOperator
		expected        bool
		expectedError   bool
	}{
		{
			testCase: cboClusterOperatorName,
			clusterOperator: &osconfigv1.ClusterOperator{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ClusterOperator",
					APIVersion: "config.uccp.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: cboClusterOperatorName,
				},
				Status: osconfigv1.ClusterOperatorStatus{
					RelatedObjects: []osconfigv1.ObjectReference{
						{
							Group:    "",
							Resource: "namespaces",
							Name:     "uccp-machine-api",
						},
					},
				},
			},
			expected: true,
		},
		{
			testCase: "invalidCO",
			clusterOperator: &osconfigv1.ClusterOperator{
				ObjectMeta: metav1.ObjectMeta{
					Name: "invalidCO",
				},
			},
			expected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(string(tc.testCase), func(t *testing.T) {
			var osClient *fakeos.Clientset
			osClient = fakeos.NewSimpleClientset(tc.clusterOperator)
			_, err := osClient.ConfigV1().ClusterOperators().Create(context.Background(), tc.clusterOperator, metav1.CreateOptions{})
			if err != nil && !apierrors.IsAlreadyExists(err) {
				t.Fatalf("Unable to create ClusterOperator for test: %v", err)
			}
			exists, err := checkForBaremetalClusterOperator(osClient)
			if exists != tc.expected {
				t.Errorf("Expected: %v, got: %v", tc.expected, exists)
			}
			if tc.expectedError != (err != nil) {
				t.Errorf("ExpectedError: %v, got: %v", tc.expectedError, err)
			}
			err = osClient.ConfigV1().ClusterOperators().Delete(context.Background(), tc.testCase, metav1.DeleteOptions{})
		})
	}
}
