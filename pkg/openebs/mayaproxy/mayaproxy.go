package mayaproxy

import (
	"encoding/json"
	"github.com/golang/glog"
	"net/http"
	"bytes"
	"io/ioutil"
	mayav1 "github.com/kubernetes-incubator/external-storage/openebs/types/v1"
	"gopkg.in/yaml.v2"
	"errors"
)

// CreateVolume requests mapi server to create an openebs volume. It returns an error if volume creation fails
func (mayaConfig MayaConfig) CreateVolume(spec mayav1.VolumeSpec) error {
	spec.Kind = "PersistentVolumeClaim"
	spec.APIVersion = "v1"

	// Marshal serializes the value provided into a YAML document
	yamlValue, _ := yaml.Marshal(spec)

	glog.V(4).Infof("[DEBUG] volume Spec Created:\n%v\n", string(yamlValue))

	url, err := mayaConfig.GetVolumeURL(versionLatest)
	glog.V(4).Infof("[DEBUG] create volume URL %v", url.String())
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(yamlValue))
	req.Header.Add("Content-Type", "application/yaml")
	c := &http.Client{
		Timeout: timeout,
	}
	resp, err := c.Do(req)

	if err != nil {
		glog.Errorf("Error when connecting maya-apiserver %v", err)
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("Unable to read response from maya-apiserver %v", err)
		return err
	}
	code := resp.StatusCode
	if code != http.StatusOK {
		glog.Errorf("Error response from maya-apiserver: %v", http.StatusText(code))
		return errors.New("Error response from maya-apiserver")
	}

	glog.Infof("volume Successfully Created:\n%v\n", string(data))
	return nil
}

// DeleteVolume requests mapi server to delete an openebs volume. It returns an error if volume deletion fails
func (mayaConfig MayaConfig) DeleteVolume(volumeName string) error {
	glog.V(2).Infof("[DEBUG] Delete Volume :%v", string(volumeName))

	url, err := mayaConfig.GetVolumeDeleteURL(versionLatest, volumeName)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	c := &http.Client{
		Timeout: timeout,
	}
	resp, err := c.Do(req)

	if err != nil {
		glog.Errorf("Error when connecting to maya-apiserver  %v", err)
		return err
	}
	defer resp.Body.Close()

	code := resp.StatusCode
	if code != http.StatusOK {
		glog.Errorf("HTTP Status error from maya-apiserver: %v\n", http.StatusText(code))
		return err
	}
	glog.Info("volume Deletion Successfully initiated")
	return nil
}

// ListVolume requests mapi server to GET the details
// of a volume and returns it by filling into *mayav1.Volume.
// If the volume does not exist or can't be retrieved then it returns an error
func (mayaConfig MayaConfig) ListVolume(volumeName string) (*mayav1.Volume, error) {
	var volume mayav1.Volume

	glog.V(2).Infof("[DEBUG] Get details for Volume :%v", string(volumeName))

	url, err := mayaConfig.GetVolumeInfoURL(versionLatest, volumeName)
	if err != nil {
		return nil, err
	}

	glog.V(2).Infof("[DEBUG] Requesting for volume details at %s", url.String())
	req, err := http.NewRequest("GET", url.String(), nil)
	c := &http.Client{
		Timeout: timeout,
	}
	resp, err := c.Do(req)
	if err != nil {
		glog.Errorf("Error when connecting to maya-apiserver %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	code := resp.StatusCode
	if code != http.StatusOK {
		glog.Errorf("HTTP Status error from maya-apiserver: %v\n", http.StatusText(code))
		return nil, errors.New("HTTP Status error from maya-apiserver: " + http.StatusText(code))
	}

	// Fill the obtained json into volume
	json.NewDecoder(resp.Body).Decode(&volume)
	glog.V(2).Infof("volume Details Successfully Retrieved %v", volume)

	return &volume, nil
}