package ports

import (
	"io/ioutil"
	"path"
	"strings"
)

const (
	sysfsTTYPath   = "/sys/class/tty/"
	sysfsUSBPrefix = "ttyUSB"
	devPath        = "/dev"
	baud           = 115200
)

// ScanUSB fetch all ports
func ScanUSB() ([]string, error) {
	files, err := ioutil.ReadDir(sysfsTTYPath)
	if err != nil {
		return nil, err
	}
	ports := make([]string, 0)
	for _, file := range files {
		if strings.HasPrefix(file.Name(), sysfsUSBPrefix) {
			ports = append(ports, path.Join(devPath, file.Name()))
		}
	}
	return ports, nil
}
