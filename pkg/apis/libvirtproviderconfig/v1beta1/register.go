package v1beta1

import (
	"bytes"
	"fmt"
	"github.com/ghodss/yaml"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/controller-runtime/pkg/runtime/scheme"
	machinev1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
)

var (
	SchemeGroupVersion	= schema.GroupVersion{Group: "libvirtproviderconfig.openshift.io", Version: "v1beta1"}
	SchemeBuilder		= &scheme.Builder{GroupVersion: SchemeGroupVersion}
)

type LibvirtProviderConfigCodec struct {
	encoder	runtime.Encoder
	decoder	runtime.Decoder
}

func NewScheme() (*runtime.Scheme, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return SchemeBuilder.Build()
}
func NewCodec() (*LibvirtProviderConfigCodec, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	scheme, err := NewScheme()
	if err != nil {
		return nil, err
	}
	codecFactory := serializer.NewCodecFactory(scheme)
	encoder, err := newEncoder(&codecFactory)
	if err != nil {
		return nil, err
	}
	codec := LibvirtProviderConfigCodec{encoder: encoder, decoder: codecFactory.UniversalDecoder(SchemeGroupVersion)}
	return &codec, nil
}
func (codec *LibvirtProviderConfigCodec) DecodeFromProviderSpec(providerConfig machinev1.ProviderSpec, out runtime.Object) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if providerConfig.Value != nil {
		if err := yaml.Unmarshal(providerConfig.Value.Raw, out); err != nil {
			return fmt.Errorf("decoding failure: %v", err)
		}
	}
	return nil
}
func (codec *LibvirtProviderConfigCodec) EncodeToProviderSpec(in runtime.Object) (*machinev1.ProviderSpec, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var buf bytes.Buffer
	if err := codec.encoder.Encode(in, &buf); err != nil {
		return nil, fmt.Errorf("encoding failed: %v", err)
	}
	return &machinev1.ProviderSpec{Value: &runtime.RawExtension{Raw: buf.Bytes()}}, nil
}
func (codec *LibvirtProviderConfigCodec) EncodeProviderStatus(in runtime.Object) (*runtime.RawExtension, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var buf bytes.Buffer
	if err := codec.encoder.Encode(in, &buf); err != nil {
		return nil, fmt.Errorf("encoding failed: %v", err)
	}
	return &runtime.RawExtension{Raw: buf.Bytes()}, nil
}
func (codec *LibvirtProviderConfigCodec) DecodeProviderStatus(providerStatus *runtime.RawExtension, out runtime.Object) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if providerStatus != nil {
		if err := yaml.Unmarshal(providerStatus.Raw, out); err != nil {
			return fmt.Errorf("decoding failure: %v", err)
		}
	}
	return nil
}
func newEncoder(codecFactory *serializer.CodecFactory) (runtime.Encoder, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	serializerInfos := codecFactory.SupportedMediaTypes()
	if len(serializerInfos) == 0 {
		return nil, fmt.Errorf("unable to find any serlializers")
	}
	encoder := codecFactory.EncoderForVersion(serializerInfos[0].Serializer, SchemeGroupVersion)
	return encoder, nil
}
