package util

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// ToJSON
func ToJSON(obj any) string {
	bt, _ := json.Marshal(obj)
	return string(bt)
}

// IP
func IP() (string, error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for i := 0; i < len(netInterfaces); i++ {
		netInterface := netInterfaces[i]
		flags := netInterface.Flags
		if flags&net.FlagUp != 0 && flags&net.FlagLoopback == 0 {
			addrs, err := netInterface.Addrs()
			if err != nil {
				return "", err
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
					return ipnet.IP.String(), nil
				}
			}
		}
	}
	return "", nil
}

// Now
// zh - yyyy-MM-dd HH:mm:ss
func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// PanicTrue
func PanicTrue(guard bool, format string, params ...any) {
	if guard {
		panic(fmt.Sprintf(format, params...))
	}
}

// PanicFalse
func PanicFalse(guard bool, format string, params ...any) {
	if !guard {
		panic(fmt.Sprintf(format, params...))
	}
}
