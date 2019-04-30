package client

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
	"math/rand"
	"github.com/davecgh/go-spew/spew"
	"github.com/golang/glog"
	libvirt "github.com/libvirt/libvirt-go"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
	"github.com/openshift/cluster-api-provider-libvirt/lib/cidr"
)

var ErrLibVirtConIsNil = errors.New("the libvirt connection was nil")
var ErrDomainNotFound = errors.New("Domain not found")

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rand.Seed(time.Now().UTC().UnixNano())
}

type pendingMapping struct {
	mac		string
	hostname	string
	network		*libvirt.Network
}

func newDomainDef() libvirtxml.Domain {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var serialPort uint
	domainDef := libvirtxml.Domain{OS: &libvirtxml.DomainOS{Type: &libvirtxml.DomainOSType{Type: "hvm"}}, Memory: &libvirtxml.DomainMemory{Unit: "MiB", Value: 512}, VCPU: &libvirtxml.DomainVCPU{Placement: "static", Value: 1}, CPU: &libvirtxml.DomainCPU{}, Devices: &libvirtxml.DomainDeviceList{Graphics: []libvirtxml.DomainGraphic{{Spice: &libvirtxml.DomainGraphicSpice{AutoPort: "yes"}}}, Channels: []libvirtxml.DomainChannel{{Target: &libvirtxml.DomainChannelTarget{VirtIO: &libvirtxml.DomainChannelTargetVirtIO{Name: "org.qemu.guest_agent.0"}}}}, RNGs: []libvirtxml.DomainRNG{{Model: "virtio", Backend: &libvirtxml.DomainRNGBackend{Random: &libvirtxml.DomainRNGBackendRandom{}}}}, Consoles: []libvirtxml.DomainConsole{{Source: &libvirtxml.DomainChardevSource{Pty: &libvirtxml.DomainChardevSourcePty{}}, Target: &libvirtxml.DomainConsoleTarget{Type: "virtio", Port: &serialPort}}}}, Features: &libvirtxml.DomainFeatureList{PAE: &libvirtxml.DomainFeature{}, ACPI: &libvirtxml.DomainFeature{}, APIC: &libvirtxml.DomainFeatureAPIC{}}}
	if v := os.Getenv("TERRAFORM_LIBVIRT_TEST_DOMAIN_TYPE"); v != "" {
		domainDef.Type = v
	} else {
		domainDef.Type = "kvm"
	}
	return domainDef
}
func getHostArchitecture(virConn *libvirt.Connect) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type HostCapabilities struct {
		XMLName	xml.Name	`xml:"capabilities"`
		Host	struct {
			XMLName	xml.Name	`xml:"host"`
			CPU	struct {
				XMLName	xml.Name	`xml:"cpu"`
				Arch	string		`xml:"arch"`
			}
		}
	}
	info, err := virConn.GetCapabilities()
	if err != nil {
		return "", err
	}
	capabilities := HostCapabilities{}
	xml.Unmarshal([]byte(info), &capabilities)
	return capabilities.Host.CPU.Arch, nil
}
func getHostCapabilities(virConn *libvirt.Connect) (libvirtxml.Caps, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	caps := libvirtxml.Caps{}
	capsXML, err := virConn.GetCapabilities()
	if err != nil {
		return caps, err
	}
	xml.Unmarshal([]byte(capsXML), &caps)
	glog.Infof("Capabilities of host \n %+v", caps)
	return caps, nil
}
func getGuestForArchType(caps libvirtxml.Caps, arch string, virttype string) (libvirtxml.CapsGuest, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, guest := range caps.Guests {
		glog.Infof("Checking for %s/%s against %s/%s", arch, virttype, guest.Arch.Name, guest.OSType)
		if guest.Arch.Name == arch && guest.OSType == virttype {
			glog.Infof("Found %d machines in guest for %s/%s", len(guest.Arch.Machines), arch, virttype)
			return guest, nil
		}
	}
	return libvirtxml.CapsGuest{}, fmt.Errorf("Could not find any guests for architecure type %s/%s", virttype, arch)
}
func getCanonicalMachineName(caps libvirtxml.Caps, arch string, virttype string, targetmachine string) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	glog.Info("Get machine name")
	guest, err := getGuestForArchType(caps, arch, virttype)
	if err != nil {
		return "", err
	}
	for _, machine := range guest.Arch.Machines {
		if machine.Name == targetmachine {
			if machine.Canonical != "" {
				return machine.Canonical, nil
			}
			return machine.Name, nil
		}
	}
	return "", fmt.Errorf("Cannot find machine type %s for %s/%s in %v", targetmachine, virttype, arch, caps)
}
func newDomainDefForConnection(virConn *libvirt.Connect) (libvirtxml.Domain, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	d := newDomainDef()
	arch, err := getHostArchitecture(virConn)
	if err != nil {
		return d, err
	}
	d.OS.Type.Arch = arch
	caps, err := getHostCapabilities(virConn)
	if err != nil {
		return d, err
	}
	guest, err := getGuestForArchType(caps, d.OS.Type.Arch, d.OS.Type.Type)
	if err != nil {
		return d, err
	}
	d.Devices.Emulator = guest.Arch.Emulator
	if len(guest.Arch.Machines) > 0 {
		d.OS.Type.Machine = guest.Arch.Machines[0].Name
	}
	canonicalmachine, err := getCanonicalMachineName(caps, d.OS.Type.Arch, d.OS.Type.Type, d.OS.Type.Machine)
	if err != nil {
		return d, err
	}
	d.OS.Type.Machine = canonicalmachine
	return d, nil
}
func setCoreOSIgnition(domainDef *libvirtxml.Domain, ignKey string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if ignKey == "" {
		return fmt.Errorf("error setting coreos ignition, ignKey is empty")
	}
	domainDef.QEMUCommandline = &libvirtxml.DomainQEMUCommandline{Args: []libvirtxml.DomainQEMUCommandlineArg{{Value: "-fw_cfg"}, {Value: fmt.Sprintf("name=opt/com.coreos/config,file=%s", ignKey)}}}
	return nil
}
func newDefDisk(i int) libvirtxml.DomainDisk {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return libvirtxml.DomainDisk{Device: "disk", Target: &libvirtxml.DomainDiskTarget{Bus: "virtio", Dev: fmt.Sprintf("vd%s", diskLetterForIndex(i))}, Driver: &libvirtxml.DomainDiskDriver{Name: "qemu", Type: "qcow2"}}
}

var diskLetters = []rune("abcdefghijklmnopqrstuvwxyz")

const oui = "05abcd"

func diskLetterForIndex(i int) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	q := i / len(diskLetters)
	r := i % len(diskLetters)
	letter := diskLetters[r]
	if q == 0 {
		return fmt.Sprintf("%c", letter)
	}
	return fmt.Sprintf("%s%c", diskLetterForIndex(q-1), letter)
}
func randomWWN(strlen int) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	const chars = "abcdef0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return oui + string(result)
}
func setDisks(domainDef *libvirtxml.Domain, diskVolume *libvirt.StorageVol) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	disk := newDefDisk(0)
	glog.Info("Getting disk volume")
	diskVolumeFile, err := diskVolume.GetPath()
	if err != nil {
		return fmt.Errorf("Error retrieving volume file: %s", err)
	}
	glog.Info("Constructing domain disk source")
	disk.Source = &libvirtxml.DomainDiskSource{File: &libvirtxml.DomainDiskSourceFile{File: diskVolumeFile}}
	domainDef.Devices.Disks = append(domainDef.Devices.Disks, disk)
	return nil
}
func xmlMarshallIndented(b interface{}) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	buf := new(bytes.Buffer)
	enc := xml.NewEncoder(buf)
	enc.Indent("  ", "    ")
	if err := enc.Encode(b); err != nil {
		return "", fmt.Errorf("could not marshall this:\n%s", spew.Sdump(b))
	}
	return buf.String(), nil
}
func setNetworkInterfaces(domainDef *libvirtxml.Domain, virConn *libvirt.Connect, partialNetIfaces map[string]*pendingMapping, waitForLeases *[]*libvirtxml.DomainInterface, networkInterfaceHostname string, networkInterfaceName string, networkInterfaceAddress string, offset int) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for i := 0; i < 1; i++ {
		netIface := libvirtxml.DomainInterface{Model: &libvirtxml.DomainInterfaceModel{Type: "virtio"}}
		var err error
		mac, err := randomMACAddress()
		if err != nil {
			return fmt.Errorf("Error generating mac address: %s", err)
		}
		netIface.MAC = &libvirtxml.DomainInterfaceMAC{Address: mac}
		if networkInterfaceName != "" {
			network, err := virConn.LookupNetworkByName(networkInterfaceName)
			if err != nil {
				return fmt.Errorf("Can't retrieve network name %s", networkInterfaceName)
			}
			defer network.Free()
			networkName, err := network.GetName()
			if err != nil {
				return fmt.Errorf("Error retrieving network name: %s", err)
			}
			networkDef, err := newDefNetworkfromLibvirt(network)
			if err != nil {
				return fmt.Errorf("Error retrieving network definition: %v", err)
			}
			if HasDHCP(networkDef) {
				hostname := domainDef.Name
				if networkInterfaceHostname != "" {
					hostname = networkInterfaceHostname
				}
				glog.Infof("Networkaddress: %v", networkInterfaceAddress)
				if networkInterfaceAddress != "" {
					_, networkCIDR, err := net.ParseCIDR(networkInterfaceAddress)
					if err != nil {
						return fmt.Errorf("failed to parse libvirt network ipRange: %v", err)
					}
					var ip net.IP
					if ip, err = cidr.GenerateIP(networkCIDR, offset); err != nil {
						return fmt.Errorf("failed to generate ip: %v", err)
					}
					glog.Infof("Adding IP/MAC/host=%s/%s/%s to %s", ip.String(), mac, hostname, networkName)
					if err := updateOrAddHost(network, ip.String(), mac, hostname); err != nil {
						return err
					}
				} else {
					wait := false
					for _, iface := range *waitForLeases {
						if iface == &netIface {
							wait = true
							break
						}
					}
					if !wait {
						return fmt.Errorf("Cannot map '%s': we are not waiting for DHCP lease and no IP has been provided", hostname)
					}
					glog.Infof("Do not have an IP for '%s' yet: will wait until DHCP provides one...", hostname)
					partialNetIfaces[strings.ToUpper(mac)] = &pendingMapping{mac: strings.ToUpper(mac), hostname: hostname, network: network}
				}
			}
			netIface.Source = &libvirtxml.DomainInterfaceSource{Network: &libvirtxml.DomainInterfaceSourceNetwork{Network: networkName}}
		}
		domainDef.Devices.Interfaces = append(domainDef.Devices.Interfaces, netIface)
	}
	return nil
}

type Config struct{ URI string }

func domainDefInit(domainDef *libvirtxml.Domain, name string, memory, vcpu int) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if name != "" {
		domainDef.Name = name
	} else {
		return fmt.Errorf("machine does not have an name set")
	}
	if memory != 0 {
		domainDef.Memory = &libvirtxml.DomainMemory{Value: uint(memory), Unit: "MiB"}
	} else {
		return fmt.Errorf("machine does not have an DomainMemory set")
	}
	if vcpu != 0 {
		domainDef.VCPU = &libvirtxml.DomainVCPU{Value: vcpu}
	} else {
		return fmt.Errorf("machine does not have an DomainVcpu set")
	}
	domainDef.CPU.Mode = "host-passthrough"
	return nil
}
