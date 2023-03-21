package util

import (
	"net"
	"sync"
)

// local
type local struct {
	ip           Result         // local host ip
	Cache        map[string]any // local cache
	PrivateFuncs PrivateFuncs   // private funcs
}

// PrivateFuncs
type PrivateFuncs []LocalFunc

// LocalFunc
type LocalFunc interface {
	Do(*local, func() error) error
}

var (
	localCtx  *local    // public local
	localOnce sync.Once // local init use lock once
)

// generateIp
func (l *local) generateIp() {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		l.ip.err = err
		return
	}

	for i := 0; i < len(netInterfaces); i++ {
		netInterface := netInterfaces[i]
		flags := netInterface.Flags
		if flags&net.FlagUp != 0 && flags&net.FlagLoopback == 0 {
			addrs, err := netInterface.Addrs()
			if err != nil {
				l.ip.err = err
				return
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
					l.ip.ctx = ipnet.IP.String()
					break
				}
			}
		}
	}
}

// Local
func Local() *local {
	localOnce.Do(func() {
		localCtx = &local{
			Cache: make(map[string]any, 0),
		}
		localCtx.generateIp()
	})
	return localCtx
}

func (l *local) IP() (string, error) {
	return l.ip.String(), l.ip.err
}
