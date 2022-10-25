module github.com/uccps-samples/machine-api-operator

go 1.16

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/go-logr/logr v1.2.0
	github.com/google/gofuzz v1.1.0
	github.com/google/uuid v1.1.2
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.14.0
	github.com/openshift/cluster-api-provider-gcp v0.0.1-0.20210615203611-a02074e8d5bb
	github.com/operator-framework/operator-sdk v0.5.1-0.20190301204940-c2efe6f74e7b
	github.com/prometheus/client_golang v1.11.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/uccps-samples/api v0.0.0-20221025085120-91bb4873437a
	github.com/uccps-samples/client-go v0.0.0-20221025085836-81245c36050d
	github.com/uccps-samples/library-go v0.0.0-20221025090733-45335f921569
	github.com/vmware/govmomi v0.22.2
	golang.org/x/net v0.0.0-20210825183410-e898025ed96a
	gopkg.in/gcfg.v1 v1.2.3
	k8s.io/api v0.22.1
	k8s.io/apimachinery v0.22.1
	k8s.io/apiserver v0.22.1
	k8s.io/client-go v0.22.1
	k8s.io/code-generator v0.22.1
	k8s.io/klog/v2 v2.30.0
	k8s.io/kubectl v0.22.0
	k8s.io/utils v0.0.0-20210930125809-cb0fa318a74b
	sigs.k8s.io/cluster-api-provider-aws v0.0.0-00010101000000-000000000000
	sigs.k8s.io/cluster-api-provider-azure v0.0.0-00010101000000-000000000000
	sigs.k8s.io/controller-runtime v0.9.3
	sigs.k8s.io/controller-tools v0.6.2
	sigs.k8s.io/yaml v1.2.0
)

replace sigs.k8s.io/cluster-api-provider-aws => github.com/openshift/cluster-api-provider-aws v0.2.1-0.20210622023641-c69a3acaee27

replace sigs.k8s.io/cluster-api-provider-azure => github.com/openshift/cluster-api-provider-azure v0.1.0-alpha.3.0.20210816141152-a7c40345b994
