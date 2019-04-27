package machines

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"
	"github.com/golang/glog"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/apimachinery/pkg/util/wait"
	"github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/framework"
	"github.com/openshift/cluster-api-actuator-pkg/pkg/manifests"
	"github.com/openshift/cluster-api-provider-libvirt/test/utils"
	machinev1beta1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	poolTimeout				= 20 * time.Second
	pollInterval				= 1 * time.Second
	poolClusterAPIDeploymentTimeout		= 10 * time.Minute
	timeoutPoolMachineRunningInterval	= 10 * time.Minute
)

func TestCart(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Machine Suite")
}
func BuildPKSecret(secretName, namespace, pkLoc string) (*apiv1.Secret, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	pkBytes, err := ioutil.ReadFile(pkLoc)
	if err != nil {
		return nil, fmt.Errorf("unable to read %v: %v", pkLoc, err)
	}
	return &apiv1.Secret{ObjectMeta: metav1.ObjectMeta{Name: secretName, Namespace: namespace}, Data: map[string][]byte{"privatekey": pkBytes}}, nil
}
func createSecretAndWait(f *framework.Framework, secret *apiv1.Secret) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	_, err := f.KubeClient.CoreV1().Secrets(secret.Namespace).Create(secret)
	Expect(err).NotTo(HaveOccurred())
	err = wait.Poll(framework.PollInterval, framework.PoolTimeout, func() (bool, error) {
		if _, err := f.KubeClient.CoreV1().Secrets(secret.Namespace).Get(secret.Name, metav1.GetOptions{}); err != nil {
			return false, nil
		}
		return true, nil
	})
	Expect(err).NotTo(HaveOccurred())
}

var _ = framework.SigKubeDescribe("Machines", func() {
	f, err := framework.NewFramework()
	if err != nil {
		panic(fmt.Errorf("unable to create framework: %v", err))
	}
	var testNamespace *apiv1.Namespace
	machinesToDelete := framework.InitMachinesToDelete()
	BeforeEach(func() {
		f.BeforeEach()
		testNamespace = &apiv1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "namespace-" + string(uuid.NewUUID())}}
		By(fmt.Sprintf("Creating %q namespace", testNamespace.Name))
		_, err = f.KubeClient.CoreV1().Namespaces().Create(testNamespace)
		Expect(err).NotTo(HaveOccurred())
		if f.LibvirtPK != "" {
			libvirtPKSecret, err := BuildPKSecret("libvirt-private-key", testNamespace.Name, f.LibvirtPK)
			Expect(err).NotTo(HaveOccurred())
			glog.V(2).Infof("Creating %q secret", libvirtPKSecret.Name)
			_, err = f.KubeClient.CoreV1().Secrets(libvirtPKSecret.Namespace).Create(libvirtPKSecret)
			Expect(err).NotTo(HaveOccurred())
		}
		f.DeployClusterAPIStack(testNamespace.Name, "libvirt-private-key")
	})
	AfterEach(func() {
		if testNamespace != nil {
			f.DestroyClusterAPIStack(testNamespace.Name, "libvirt-private-key")
			glog.V(2).Infof(testNamespace.Name+": %#v", testNamespace)
			By(fmt.Sprintf("Destroying %q namespace", testNamespace.Name))
			f.KubeClient.CoreV1().Namespaces().Delete(testNamespace.Name, &metav1.DeleteOptions{})
		}
	})
	Context("libvirt actuator", func() {
		It("can create domain", func() {
			clusterID := framework.ClusterID
			if clusterID == "" {
				clusterID = "cluster-" + string(uuid.NewUUID())
			}
			cluster := &machinev1beta1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: clusterID, Namespace: testNamespace.Name}, Spec: machinev1beta1.ClusterSpec{ClusterNetwork: machinev1beta1.ClusterNetworkingConfig{Services: machinev1beta1.NetworkRanges{CIDRBlocks: []string{"10.0.0.1/24"}}, Pods: machinev1beta1.NetworkRanges{CIDRBlocks: []string{"10.0.0.1/24"}}, ServiceDomain: "example.com"}}}
			f.CreateClusterAndWait(cluster)
			testMachineProviderSpec, err := utils.TestingMachineProviderSpec(f.LibvirtURI, cluster.Name)
			Expect(err).NotTo(HaveOccurred())
			testMachine := manifests.TestingMachine(cluster.Name, cluster.Namespace, testMachineProviderSpec)
			lcw, err := NewLibvirtClient("qemu:///system")
			Expect(err).NotTo(HaveOccurred())
			f.CreateMachineAndWait(testMachine, lcw)
			machinesToDelete.AddMachine(testMachine, f, lcw)
			f.DeleteMachineAndWait(testMachine, lcw)
		})
	})
	It("Can deploy compute nodes through machineset", func() {
		clusterID := framework.ClusterID
		if clusterID == "" {
			clusterUUID := string(uuid.NewUUID())
			clusterID = "cluster-" + clusterUUID[:8]
		}
		cluster := &machinev1beta1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: clusterID, Namespace: testNamespace.Name}, Spec: machinev1beta1.ClusterSpec{ClusterNetwork: machinev1beta1.ClusterNetworkingConfig{Services: machinev1beta1.NetworkRanges{CIDRBlocks: []string{"10.0.0.1/24"}}, Pods: machinev1beta1.NetworkRanges{CIDRBlocks: []string{"10.0.0.1/24"}}, ServiceDomain: "example.com"}}}
		f.CreateClusterAndWait(cluster)
		masterUserDataSecret, err := manifests.MasterMachineUserDataSecret("masteruserdatasecret", testNamespace.Name, []string{"127.0.0.1"})
		Expect(err).NotTo(HaveOccurred())
		createSecretAndWait(f, masterUserDataSecret)
		masterMachineProviderSpec, err := utils.MasterMachineProviderSpec(masterUserDataSecret.Name, f.LibvirtURI)
		Expect(err).NotTo(HaveOccurred())
		masterMachine := manifests.MasterMachine(cluster.Name, cluster.Namespace, masterMachineProviderSpec)
		lcw, err := NewLibvirtClient("qemu:///system")
		Expect(err).NotTo(HaveOccurred())
		f.CreateMachineAndWait(masterMachine, lcw)
		machinesToDelete.AddMachine(masterMachine, f, lcw)
		var masterMachinePrivateIP string
		err = wait.Poll(pollInterval, poolTimeout, func() (bool, error) {
			privateIP, err := lcw.GetPrivateIP(masterMachine)
			if err != nil {
				return false, nil
			}
			masterMachinePrivateIP = privateIP
			return true, nil
		})
		if err != nil {
			glog.Errorf("Unable to get instance ip address: %v", err)
		}
		Expect(err).NotTo(HaveOccurred())
		glog.V(2).Infof("Master machine running at %v", masterMachinePrivateIP)
		By("Collecting master kubeconfig")
		restConfig, err := f.GetMasterMachineRestConfig(masterMachine, lcw)
		Expect(err).NotTo(HaveOccurred())
		By("Upload actuator image to the master guest")
		err = f.UploadDockerImageToInstance(f.MachineControllerImage, masterMachinePrivateIP)
		Expect(err).NotTo(HaveOccurred())
		if f.MachineManagerImage != f.MachineControllerImage {
			err = f.UploadDockerImageToInstance(f.MachineManagerImage, masterMachinePrivateIP)
			Expect(err).NotTo(HaveOccurred())
		}
		sshConfig, err := framework.DefaultSSHConfig()
		Expect(err).NotTo(HaveOccurred())
		clusterFramework, err := framework.NewFrameworkFromConfig(restConfig, sshConfig)
		Expect(err).NotTo(HaveOccurred())
		By(fmt.Sprintf("Creating %q namespace", testNamespace.Name))
		_, err = clusterFramework.KubeClient.CoreV1().Namespaces().Create(testNamespace)
		Expect(err).NotTo(HaveOccurred())
		if f.LibvirtPK != "" {
			libvirtPKSecret, err := BuildPKSecret("libvirt-private-key", testNamespace.Name, f.LibvirtPK)
			Expect(err).NotTo(HaveOccurred())
			glog.V(2).Infof("Creating %q secret", libvirtPKSecret.Name)
			_, err = clusterFramework.KubeClient.CoreV1().Secrets(libvirtPKSecret.Namespace).Create(libvirtPKSecret)
			Expect(err).NotTo(HaveOccurred())
		}
		clusterFramework.DeployClusterAPIStack(testNamespace.Name, "libvirt-private-key")
		By("Deploy worker nodes through machineset")
		masterPrivateIP := masterMachinePrivateIP
		clusterFramework.CreateClusterAndWait(cluster)
		workerUserDataSecret, err := manifests.WorkerMachineUserDataSecret("workeruserdatasecret", testNamespace.Name, masterPrivateIP)
		Expect(err).NotTo(HaveOccurred())
		createSecretAndWait(clusterFramework, workerUserDataSecret)
		workerMachineSetProviderSpec, err := utils.WorkerMachineProviderSpec(workerUserDataSecret.Name, f.LibvirtURI)
		Expect(err).NotTo(HaveOccurred())
		workerMachineSet := manifests.WorkerMachineSet(cluster.Name, cluster.Namespace, workerMachineSetProviderSpec)
		clusterFramework.CreateMachineSetAndWait(workerMachineSet, lcw)
		machinesToDelete.AddMachineSet(workerMachineSet, clusterFramework, lcw)
		By("Checking master and worker nodes are ready")
		err = clusterFramework.WaitForNodesToGetReady(2)
		Expect(err).NotTo(HaveOccurred())
		By("Both master and worker nodes are ready")
		By("Checking compute node role and node linking")
		err = wait.Poll(framework.PollInterval, 5*framework.PoolTimeout, func() (bool, error) {
			items, err := clusterFramework.KubeClient.CoreV1().Nodes().List(metav1.ListOptions{})
			if err != nil {
				return false, fmt.Errorf("unable to list nodes: %v", err)
			}
			var nonMasterNodes []apiv1.Node
			for _, node := range items.Items {
				if _, isMaster := node.Labels["node-role.kubernetes.io/master"]; isMaster {
					continue
				}
				nonMasterNodes = append(nonMasterNodes, node)
			}
			glog.V(2).Infof("Non-master nodes to check: %#v", nonMasterNodes)
			machines, err := clusterFramework.CAPIClient.MachineV1beta1().Machines(workerMachineSet.Namespace).List(metav1.ListOptions{LabelSelector: labels.SelectorFromSet(workerMachineSet.Spec.Selector.MatchLabels).String()})
			Expect(err).NotTo(HaveOccurred())
			matches := make(map[string]string)
			for _, machine := range machines.Items {
				if machine.Status.NodeRef != nil {
					matches[machine.Status.NodeRef.Name] = machine.Name
				}
			}
			glog.V(2).Infof("Machine-node matches: %#v\n", matches)
			for _, node := range nonMasterNodes {
				_, isCompute := node.Labels["node-role.kubernetes.io/compute"]
				if !isCompute {
					glog.V(2).Infof("node %q does not have the compute role assigned", node.Name)
					return false, nil
				}
				glog.V(2).Infof("node %q role set to 'node-role.kubernetes.io/compute'", node.Name)
				matchingMachine, found := matches[node.Name]
				if !found {
					glog.V(2).Infof("node %q is not linked with a machine", node.Name)
					return false, nil
				}
				glog.V(2).Infof("node %q is linked with %q machine", node.Name, matchingMachine)
			}
			return true, nil
		})
		Expect(err).NotTo(HaveOccurred())
		By("Destroying worker machines")
		clusterFramework.DeleteMachineSetAndWait(workerMachineSet, lcw)
		By("Destroying master machine")
		f.DeleteMachineAndWait(masterMachine, lcw)
	})
})
