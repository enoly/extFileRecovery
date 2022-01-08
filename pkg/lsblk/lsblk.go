package lsblk

import (
	"encoding/json"
	"os/exec"
)

type BlockDevice struct {
	Name   string
	Fstype string
}

type DeviceList struct {
	Blockdevices []BlockDevice
}

func GetDeviceList() (*DeviceList, error) {
	out, err := exec.Command("lsblk", "-J", "-p", "-l", "-o", "NAME,FSTYPE").Output()
	if err != nil {
		return nil, err
	}

	var devices DeviceList
	if err := json.Unmarshal(out, &devices); err != nil {
		return nil, err
	}

	return &devices, nil
}
