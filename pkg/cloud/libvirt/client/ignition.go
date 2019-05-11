package client

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"github.com/golang/glog"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
	providerconfigv1 "github.com/openshift/cluster-api-provider-libvirt/pkg/apis/libvirtproviderconfig/v1beta1"
)

func setIgnition(domainDef *libvirtxml.Domain, client *libvirtClient, ignition *providerconfigv1.Ignition, kubeClient kubernetes.Interface, machineNamespace, volumeName, poolName string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	glog.Info("Creating ignition file")
	ignitionDef := newIgnitionDef()
	if ignition.UserDataSecret == "" {
		return fmt.Errorf("ignition.userDataSecret not set")
	}
	secret, err := kubeClient.CoreV1().Secrets(machineNamespace).Get(ignition.UserDataSecret, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("can not retrieve user data secret '%v/%v' when constructing cloud init volume: %v", machineNamespace, ignition.UserDataSecret, err)
	}
	userDataSecret, ok := secret.Data["userData"]
	if !ok {
		return fmt.Errorf("can not retrieve user data secret '%v/%v' when constructing cloud init volume: key 'userData' not found in the secret", machineNamespace, ignition.UserDataSecret)
	}
	ignitionDef.Name = volumeName
	ignitionDef.PoolName = poolName
	ignitionDef.Content = string(userDataSecret)
	glog.Infof("Ignition: %+v", ignitionDef)
	ignitionVolumeName, err := ignitionDef.createAndUpload(client)
	if err != nil {
		return err
	}
	domainDef.QEMUCommandline = &libvirtxml.DomainQEMUCommandline{Args: []libvirtxml.DomainQEMUCommandlineArg{{Value: "-fw_cfg"}, {Value: fmt.Sprintf("name=opt/com.coreos/config,file=%s", ignitionVolumeName)}}}
	return nil
}

type defIgnition struct {
	Name		string
	PoolName	string
	Content		string
}

func newIgnitionDef() defIgnition {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return defIgnition{}
}
func (ign *defIgnition) createAndUpload(client *libvirtClient) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	volumeDef := newDefVolume(ign.Name)
	ignFile, err := ign.createFile()
	if err != nil {
		return "", err
	}
	defer func() {
		if err = os.Remove(ignFile); err != nil {
			glog.Infof("Error while removing tmp Ignition file: %s", err)
		}
	}()
	img, err := newImage(ignFile)
	if err != nil {
		return "", err
	}
	size, err := img.size()
	if err != nil {
		return "", err
	}
	volumeDef.Capacity.Unit = "B"
	volumeDef.Capacity.Value = size
	volumeDef.Target.Format.Type = "raw"
	return uploadVolume(ign.PoolName, client, volumeDef, img)
}
func (ign *defIgnition) createFile() (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	glog.Info("Creating Ignition temporary file")
	tempFile, err := ioutil.TempFile("", ign.Name)
	if err != nil {
		return "", fmt.Errorf("Cannot create tmp file for Ignition: %s", err)
	}
	defer tempFile.Close()
	var file bool
	file = true
	if _, err := os.Stat(ign.Content); err != nil {
		var js map[string]interface{}
		if errConf := json.Unmarshal([]byte(ign.Content), &js); errConf != nil {
			return "", fmt.Errorf("coreos_ignition 'content' is neither a file "+"nor a valid json object %s", ign.Content)
		}
		file = false
	}
	if !file {
		if _, err := tempFile.WriteString(ign.Content); err != nil {
			return "", fmt.Errorf("Cannot write Ignition object to temporary " + "ignition file")
		}
	} else if file {
		ignFile, err := os.Open(ign.Content)
		if err != nil {
			return "", fmt.Errorf("Error opening supplied Ignition file %s", ign.Content)
		}
		defer ignFile.Close()
		_, err = io.Copy(tempFile, ignFile)
		if err != nil {
			return "", fmt.Errorf("Error copying supplied Igition file to temporary file: %s", ign.Content)
		}
	}
	return tempFile.Name(), nil
}
