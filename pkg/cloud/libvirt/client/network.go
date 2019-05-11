package client

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"github.com/golang/glog"
	libvirt "github.com/libvirt/libvirt-go"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
)

const (
	netModeIsolated	= "none"
	netModeNat		= "nat"
	netModeRoute	= "route"
	netModeBridge	= "bridge"
)

type Network interface {
	GetXMLDesc(flags libvirt.NetworkXMLFlags) (string, error)
}

func newDefNetworkfromLibvirt(network Network) (libvirtxml.Network, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	networkXMLDesc, err := network.GetXMLDesc(0)
	if err != nil {
		return libvirtxml.Network{}, fmt.Errorf("Error retrieving libvirt domain XML description: %s", err)
	}
	networkDef := libvirtxml.Network{}
	err = xml.Unmarshal([]byte(networkXMLDesc), &networkDef)
	if err != nil {
		return libvirtxml.Network{}, fmt.Errorf("Error reading libvirt network XML description: %s", err)
	}
	return networkDef, nil
}
func HasDHCP(net libvirtxml.Network) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if net.Forward != nil {
		if net.Forward.Mode == "nat" || net.Forward.Mode == "route" || net.Forward.Mode == "" {
			return true
		}
	}
	return false
}
func updateOrAddHost(n *libvirt.Network, ip, mac, name string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err := updateHost(n, ip, mac, name)
	if virErr, ok := err.(libvirt.Error); ok && virErr.Code == libvirt.ERR_OPERATION_INVALID && virErr.Domain == libvirt.FROM_NETWORK {
		return addHost(n, ip, mac, name)
	}
	return err
}
func addHost(n *libvirt.Network, ip, mac, name string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	xmlDesc, err := getHostXMLDesc(ip, mac, name)
	if err != nil {
		return fmt.Errorf("error getting host xml desc: %v", err)
	}
	glog.Infof("Adding host with XML:\n%s", xmlDesc)
	return n.Update(libvirt.NETWORK_UPDATE_COMMAND_ADD_LAST, libvirt.NETWORK_SECTION_IP_DHCP_HOST, -1, xmlDesc, libvirt.NETWORK_UPDATE_AFFECT_CURRENT)
}
func getHostXMLDesc(ip, mac, name string) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	networkDHCPHost := libvirtxml.NetworkDHCPHost{IP: ip, MAC: mac, Name: name}
	tmp := struct {
		XMLName	xml.Name	`xml:"host"`
		libvirtxml.NetworkDHCPHost
	}{xml.Name{}, networkDHCPHost}
	xml, err := xmlMarshallIndented(tmp)
	if err != nil {
		return "", fmt.Errorf("could not marshall: %v", err)
	}
	return xml, nil
}
func updateHost(n *libvirt.Network, ip, mac, name string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	xmlDesc, err := getHostXMLDesc(ip, mac, name)
	if err != nil {
		return fmt.Errorf("error getting host xml desc: %v", err)
	}
	glog.Infof("Updating host with XML:\n%s", xmlDesc)
	return n.Update(libvirt.NETWORK_UPDATE_COMMAND_MODIFY, libvirt.NETWORK_SECTION_IP_DHCP_HOST, -1, xmlDesc, libvirt.NETWORK_UPDATE_AFFECT_CURRENT)
}
func randomMACAddress() (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	buf[0] = (buf[0] | 2) & 0xfe
	buf[0] |= 2
	if buf[0] == 0xfe {
		buf[0] = 0xee
	}
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5]), nil
}
