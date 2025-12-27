package internal

import (
	"net"
	"net/http"
)

func StartFileServer(port, folder string) {
	go func() {
		http.Handle("/", http.FileServer(http.Dir(folder)))
		_ = http.ListenAndServe(port, nil)
	}()
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
