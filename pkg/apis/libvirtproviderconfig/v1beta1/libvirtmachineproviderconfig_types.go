package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ClusterIDLabel		= "machine.openshift.io/cluster-api-cluster"
	MachineRoleLabel	= "machine.openshift.io/cluster-api-machine-role"
	MachineTypeLabel	= "machine.openshift.io/cluster-api-machine-type"
)

type LibvirtMachineProviderConfig struct {
	metav1.TypeMeta			`json:",inline"`
	DomainMemory			int		`json:"domainMemory"`
	DomainVcpu			int		`json:"domainVcpu"`
	IgnKey				string		`json:"ignKey"`
	Ignition			*Ignition	`json:"ignition"`
	CloudInit			*CloudInit	`json:"cloudInit"`
	Volume				*Volume		`json:"volume"`
	NetworkInterfaceName		string		`json:"networkInterfaceName"`
	NetworkInterfaceHostname	string		`json:"networkInterfaceHostname"`
	NetworkInterfaceAddress		string		`json:"networkInterfaceAddress"`
	NetworkUUID			string		`json:"networkUUID"`
	Autostart			bool		`json:"autostart"`
	URI				string		`json:"uri"`
}
type Ignition struct {
	UserDataSecret string `json:"userDataSecret"`
}
type CloudInit struct {
	UserDataSecret	string	`json:"userDataSecret"`
	SSHAccess	bool	`json:"sshAccess"`
}
type Volume struct {
	PoolName	string	`json:"poolName"`
	BaseVolumeID	string	`json:"baseVolumeID"`
	VolumeName	string	`json:"volumeName"`
}
type LibvirtClusterProviderConfig struct {
	metav1.TypeMeta `json:",inline"`
}
type LibvirtMachineProviderStatus struct {
	metav1.TypeMeta	`json:",inline"`
	InstanceID	*string					`json:"instanceID"`
	InstanceState	*string					`json:"instanceState"`
	Conditions	[]LibvirtMachineProviderCondition	`json:"conditions"`
}
type LibvirtMachineProviderConditionType string

const (
	MachineCreated LibvirtMachineProviderConditionType = "MachineCreated"
)

type LibvirtMachineProviderCondition struct {
	Type			LibvirtMachineProviderConditionType	`json:"type"`
	Status			corev1.ConditionStatus			`json:"status"`
	LastProbeTime		metav1.Time				`json:"lastProbeTime"`
	LastTransitionTime	metav1.Time				`json:"lastTransitionTime"`
	Reason			string					`json:"reason"`
	Message			string					`json:"message"`
}
type LibvirtClusterProviderStatus struct {
	metav1.TypeMeta `json:",inline"`
}
type LibvirtMachineProviderConfigList struct {
	metav1.TypeMeta	`json:",inline"`
	metav1.ListMeta	`json:"metadata,omitempty"`
	Items		[]LibvirtMachineProviderConfig	`json:"items"`
}

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	SchemeBuilder.Register(&LibvirtMachineProviderConfig{}, &LibvirtMachineProviderConfigList{}, &LibvirtMachineProviderStatus{})
}
