package client

import (
	"testing"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
)

func TestSetCoreOSIgnition(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	testCases := []struct {
		ignKey		string
		expected	string
		errorMessage	string
	}{{ignKey: "myIgnitionConfig", expected: "name=opt/com.coreos/config,file=myIgnitionConfig", errorMessage: ""}, {ignKey: "", expected: "", errorMessage: "error setting coreos ignition, ignKey is empty"}}
TestCases:
	for i, tc := range testCases {
		domainDef := libvirtxml.Domain{}
		err := setCoreOSIgnition(&domainDef, tc.ignKey)
		if err != nil {
			if err.Error() != tc.errorMessage {
				t.Errorf("test case %d: failed to return error when key is empty. Got: %s, Expected: %s", i, err.Error(), tc.errorMessage)
			}
		} else {
			for i, v := range domainDef.QEMUCommandline.Args {
				if v.Value == "-fw_cfg" {
					if domainDef.QEMUCommandline.Args[i+1].Value == "name=opt/com.coreos/config,file=myIgnitionConfig" {
						continue TestCases
					}
				}
			}
			t.Errorf("test case %d: failed to setCoreOSIgnition for key %s. Expected: %s", i, tc.ignKey, tc.expected)
		}
	}
}
