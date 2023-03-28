package util

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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

// MD5
func MD5(bt []byte) string {
	return fmt.Sprintf("%x", md5.Sum(bt))
}

// MD5File
func MD5File(r io.Reader) string {
	bf := make([]byte, 4096)
	hashMd5 := md5.New()
	for {
		n, err := r.Read(bf)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil && err != io.EOF {
			break
		}
		hashMd5.Write(bf[:n])
	}
	return fmt.Sprintf("%x", hashMd5.Sum(nil))
}

// Base64Encode
func Base64Encode(data string) string {
	return base64.URLEncoding.EncodeToString([]byte(data))
}

// Base64Decode
func Base64Decode(data string) string {
	b, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return ""
	}
	return string(b)
}

// Now
// zh - yyyy-MM-dd HH:mm:ss
func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
