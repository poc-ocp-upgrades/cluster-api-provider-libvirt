package client

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"io"
	"strconv"
	"strings"
	"time"
	libvirt "github.com/libvirt/libvirt-go"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
)

const (
	defaultSize = 17706254336
)

var ErrVolumeNotFound = errors.New("Domain not found")
var waitSleepInterval = 1 * time.Second
var waitTimeout = 5 * time.Minute

func waitForSuccess(errorMessage string, f func() error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	start := time.Now()
	for {
		err := f()
		if err == nil {
			return nil
		}
		glog.Infof("%s. Re-trying.\n", err)
		time.Sleep(waitSleepInterval)
		if time.Since(start) > waitTimeout {
			return fmt.Errorf("%s: %s", errorMessage, err)
		}
	}
}
func newDefVolume(name string) libvirtxml.StorageVolume {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return libvirtxml.StorageVolume{Name: name, Target: &libvirtxml.StorageVolumeTarget{Format: &libvirtxml.StorageVolumeTargetFormat{Type: "qcow2"}, Permissions: &libvirtxml.StorageVolumeTargetPermissions{Mode: "644"}}, Capacity: &libvirtxml.StorageVolumeSize{Unit: "bytes", Value: 1}}
}
func newDefBackingStoreFromLibvirt(baseVolume *libvirt.StorageVol) (libvirtxml.StorageVolumeBackingStore, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	baseVolumeDef, err := newDefVolumeFromLibvirt(baseVolume)
	if err != nil {
		return libvirtxml.StorageVolumeBackingStore{}, fmt.Errorf("could not get volume: %s", err)
	}
	baseVolPath, err := baseVolume.GetPath()
	if err != nil {
		return libvirtxml.StorageVolumeBackingStore{}, fmt.Errorf("could not get base image path: %s", err)
	}
	backingStoreDef := libvirtxml.StorageVolumeBackingStore{Path: baseVolPath, Format: &libvirtxml.StorageVolumeTargetFormat{Type: baseVolumeDef.Target.Format.Type}}
	return backingStoreDef, nil
}
func newDefVolumeFromLibvirt(volume *libvirt.StorageVol) (libvirtxml.StorageVolume, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	name, err := volume.GetName()
	if err != nil {
		return libvirtxml.StorageVolume{}, fmt.Errorf("could not get name for volume: %s", err)
	}
	volumeDefXML, err := volume.GetXMLDesc(0)
	if err != nil {
		return libvirtxml.StorageVolume{}, fmt.Errorf("could not get XML description for volume %s: %s", name, err)
	}
	volumeDef, err := newDefVolumeFromXML(volumeDefXML)
	if err != nil {
		return libvirtxml.StorageVolume{}, fmt.Errorf("could not get a volume definition from XML for %s: %s", volumeDef.Name, err)
	}
	return volumeDef, nil
}
func newDefVolumeFromXML(s string) (libvirtxml.StorageVolume, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var volumeDef libvirtxml.StorageVolume
	err := xml.Unmarshal([]byte(s), &volumeDef)
	if err != nil {
		return libvirtxml.StorageVolume{}, err
	}
	return volumeDef, nil
}
func timeFromEpoch(str string) time.Time {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var s, ns int
	ts := strings.Split(str, ".")
	if len(ts) == 2 {
		ns, _ = strconv.Atoi(ts[1])
	}
	s, _ = strconv.Atoi(ts[0])
	return time.Unix(int64(s), int64(ns))
}
func uploadVolume(poolName string, client *libvirtClient, volumeDef libvirtxml.StorageVolume, img image) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pool, err := client.connection.LookupStoragePoolByName(poolName)
	if err != nil {
		return "", fmt.Errorf("can't find storage pool %q", poolName)
	}
	defer pool.Free()
	err = waitForSuccess("Error refreshing pool for volume", func() error {
		return pool.Refresh(0)
	})
	if err != nil {
		return "", fmt.Errorf("timeout when calling waitForSuccess: %v", err)
	}
	volumeDefXML, err := xml.Marshal(volumeDef)
	if err != nil {
		return "", fmt.Errorf("Error serializing libvirt volume: %s", err)
	}
	volume, err := pool.StorageVolCreateXML(string(volumeDefXML), 0)
	if err != nil {
		return "", fmt.Errorf("Error creating libvirt volume for device %s: %s", volumeDef.Name, err)
	}
	defer volume.Free()
	err = img.importImage(newCopier(client.connection, volume, volumeDef.Capacity.Value), volumeDef)
	if err != nil {
		return "", fmt.Errorf("Error while uploading volume %s: %s", img.string(), err)
	}
	volumeKey, err := volume.GetKey()
	if err != nil {
		return "", fmt.Errorf("Error retrieving volume key: %s", err)
	}
	glog.Infof("Volume ID: %s", volumeKey)
	return volumeKey, nil
}
func newCopier(virConn *libvirt.Connect, volume *libvirt.StorageVol, size uint64) func(src io.Reader) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	copier := func(src io.Reader) error {
		var bytesCopied int64
		stream, err := virConn.NewStream(0)
		if err != nil {
			return err
		}
		defer func() {
			if uint64(bytesCopied) != size {
				stream.Abort()
			} else {
				stream.Finish()
			}
			stream.Free()
		}()
		volume.Upload(stream, 0, size, 0)
		sio := newStreamIO(*stream)
		bytesCopied, err = io.Copy(sio, src)
		if err != nil {
			return err
		}
		glog.Infof("%d bytes uploaded\n", bytesCopied)
		return nil
	}
	return copier
}
