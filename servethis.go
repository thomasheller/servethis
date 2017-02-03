package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	dir := flag.String("dir", "", "directory (default current working directory)")

	flag.Parse()

	if *dir == "" {
		cwd, err := os.Getwd()

		if err != nil {
			log.Fatalf("Error fetching current working directory: %s\n", err)
		}

		dir = &cwd
	}

	listner, err := net.Listen("tcp4", ":0")
	if err != nil {
		log.Fatal(err)
	}

	_, port, err := net.SplitHostPort(listner.Addr().String())
	if err != nil {
		log.Fatal("Unable to parse listening port")
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal("Error retrieving interfaces")
	}

	ifaces := []string{}
	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil || ip.To4() == nil || ip.IsLoopback() || ip.IsMulticast() {
			continue
		}
		ifaces = append(ifaces, ip.String())
	}
	if len(ifaces) == 0 {
		log.Fatal("No valid ipv4 interfaces found!")
	}

	log.Printf("Serving folder: %s  Ctrl+C to exit", *dir)
	for _, iface := range ifaces {
		log.Printf("Listening on: http://%s:%s/", iface, port)
	}

	log.Fatal(http.Serve(listner, http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Served %s to %s", r.URL, r.RemoteAddr)
			http.FileServer(http.Dir(*dir)).ServeHTTP(w, r)
		})))
}
