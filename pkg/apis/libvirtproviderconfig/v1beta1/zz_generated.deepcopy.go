package v1beta1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func (in *CloudInit) DeepCopyInto(out *CloudInit) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	return
}
func (in *CloudInit) DeepCopy() *CloudInit {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(CloudInit)
	in.DeepCopyInto(out)
	return out
}
func (in *Ignition) DeepCopyInto(out *Ignition) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	return
}
func (in *Ignition) DeepCopy() *Ignition {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(Ignition)
	in.DeepCopyInto(out)
	return out
}
func (in *LibvirtClusterProviderConfig) DeepCopyInto(out *LibvirtClusterProviderConfig) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	out.TypeMeta = in.TypeMeta
	return
}
func (in *LibvirtClusterProviderConfig) DeepCopy() *LibvirtClusterProviderConfig {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(LibvirtClusterProviderConfig)
	in.DeepCopyInto(out)
	return out
}
func (in *LibvirtClusterProviderConfig) DeepCopyObject() runtime.Object {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
func (in *LibvirtClusterProviderStatus) DeepCopyInto(out *LibvirtClusterProviderStatus) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	out.TypeMeta = in.TypeMeta
	return
}
func (in *LibvirtClusterProviderStatus) DeepCopy() *LibvirtClusterProviderStatus {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(LibvirtClusterProviderStatus)
	in.DeepCopyInto(out)
	return out
}
func (in *LibvirtClusterProviderStatus) DeepCopyObject() runtime.Object {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
func (in *LibvirtMachineProviderCondition) DeepCopyInto(out *LibvirtMachineProviderCondition) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	in.LastProbeTime.DeepCopyInto(&out.LastProbeTime)
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
	return
}
func (in *LibvirtMachineProviderCondition) DeepCopy() *LibvirtMachineProviderCondition {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(LibvirtMachineProviderCondition)
	in.DeepCopyInto(out)
	return out
}
func (in *LibvirtMachineProviderConfig) DeepCopyInto(out *LibvirtMachineProviderConfig) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	out.TypeMeta = in.TypeMeta
	if in.Ignition != nil {
		in, out := &in.Ignition, &out.Ignition
		*out = new(Ignition)
		**out = **in
	}
	if in.CloudInit != nil {
		in, out := &in.CloudInit, &out.CloudInit
		*out = new(CloudInit)
		**out = **in
	}
	if in.Volume != nil {
		in, out := &in.Volume, &out.Volume
		*out = new(Volume)
		**out = **in
	}
	return
}
func (in *LibvirtMachineProviderConfig) DeepCopy() *LibvirtMachineProviderConfig {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(LibvirtMachineProviderConfig)
	in.DeepCopyInto(out)
	return out
}
func (in *LibvirtMachineProviderConfig) DeepCopyObject() runtime.Object {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
func (in *LibvirtMachineProviderConfigList) DeepCopyInto(out *LibvirtMachineProviderConfigList) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]LibvirtMachineProviderConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}
func (in *LibvirtMachineProviderConfigList) DeepCopy() *LibvirtMachineProviderConfigList {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(LibvirtMachineProviderConfigList)
	in.DeepCopyInto(out)
	return out
}
func (in *LibvirtMachineProviderConfigList) DeepCopyObject() runtime.Object {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
func (in *LibvirtMachineProviderStatus) DeepCopyInto(out *LibvirtMachineProviderStatus) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	out.TypeMeta = in.TypeMeta
	if in.InstanceID != nil {
		in, out := &in.InstanceID, &out.InstanceID
		*out = new(string)
		**out = **in
	}
	if in.InstanceState != nil {
		in, out := &in.InstanceState, &out.InstanceState
		*out = new(string)
		**out = **in
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]LibvirtMachineProviderCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}
func (in *LibvirtMachineProviderStatus) DeepCopy() *LibvirtMachineProviderStatus {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(LibvirtMachineProviderStatus)
	in.DeepCopyInto(out)
	return out
}
func (in *LibvirtMachineProviderStatus) DeepCopyObject() runtime.Object {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
func (in *Volume) DeepCopyInto(out *Volume) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	return
}
func (in *Volume) DeepCopy() *Volume {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(Volume)
	in.DeepCopyInto(out)
	return out
}
