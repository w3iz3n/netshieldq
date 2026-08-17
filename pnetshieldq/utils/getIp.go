package utils

import (
	"fmt"
	"net"
	"net/http"
)

func getIPHandler(w http.ResponseWriter, r *http.Request) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Fprintf(w, "Error getting IP: %s", err)
		return
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		fmt.Fprintf(w, "Unable to parse user IP")
		return
	}
	fmt.Println(userIP)
	fmt.Fprintf(w, " %s", userIP)
}
