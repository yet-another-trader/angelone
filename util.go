package angelone

import (
	"errors"
	"io"
	"net"
	"net/http"
)

func localIPAddr() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	localIP := ""

	for _, address := range addrs {
		if ip, ok := address.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				localIP = ip.IP.String()
				return "", nil
			}
		}
	}

	if localIP == "" {
		return "", errors.New("could not find local IP address")
	}

	return localIP, nil
}

func publicIPAddr() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func macAddr() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		mac := iface.HardwareAddr.String()
		if mac != "" {
			return mac, nil
		}
	}

	return "", errors.New("no active network interface found")
}

func must(f func() error) {
	if err := f(); err != nil {
		panic(err)
	}
}
