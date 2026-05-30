package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/api"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/store"
)

var version = "dev"

func main() {
	var storeFile, port, host, uiHost, uiPort, tokensFlag string
	var showVersion bool

	flag.StringVarP(&host, "host", "h", "", "network interface to bind to (default: all interfaces)")
	flag.StringVarP(&port, "port", "p", "8080", "port to listen on")
	flag.StringVarP(&storeFile, "store-file", "s", "", "path to JSON persistence file (default: in-memory)")
	flag.StringVarP(&tokensFlag, "tokens", "t", "", "comma-separated list of allowed API tokens (default: all tokens allowed)")
	flag.StringVarP(&uiHost, "ui-host", "", "", "network interface to bind the UI to (default: same as --host)")
	flag.StringVarP(&uiPort, "uiport", "", "", "port for the UI (default: same as --port)")
	flag.BoolVarP(&showVersion, "version", "v", false, "print version and exit")
	flag.Parse()

	if showVersion {
		fmt.Println(version)
		return
	}

	var allowedTokens []string
	if tokensFlag != "" {
		for t := range strings.SplitSeq(tokensFlag, ",") {
			if t = strings.TrimSpace(t); t != "" {
				allowedTokens = append(allowedTokens, t)
			}
		}
	}

	var s store.Store
	if storeFile != "" {
		fs, err := store.NewJSONFileStore(storeFile)
		if err != nil {
			log.Fatalf("failed to open store file: %v", err)
		}
		log.Printf("using JSON file store: %s", storeFile)
		s = fs
	} else {
		s = store.NewMemoryStore()
		log.Println("using in-memory store")
	}

	resolvedUIHost := uiHost
	if resolvedUIHost == "" {
		resolvedUIHost = host
	}
	resolvedUIPort := uiPort
	if resolvedUIPort == "" {
		resolvedUIPort = port
	}

	apiAddr := net.JoinHostPort(host, port)
	uiAddr  := net.JoinHostPort(resolvedUIHost, resolvedUIPort)

	printStartupURLs(host, port, resolvedUIHost, resolvedUIPort)

	if apiAddr == uiAddr {
		log.Fatal(http.ListenAndServe(apiAddr, api.NewServer(s, allowedTokens)))
	} else {
		errCh := make(chan error, 2)
		go func() { errCh <- http.ListenAndServe(apiAddr, api.NewAPIHandler(s, allowedTokens)) }()
		go func() { errCh <- http.ListenAndServe(uiAddr, api.NewUIHandler()) }()
		log.Fatal(<-errCh)
	}
}

func printStartupURLs(apiHost, apiPort, uiHost, uiPort string) {
	fmt.Printf("\nNullCloud backend %s\n\n", version)

	fmt.Println("  API:")
	for _, ip := range listAddresses(apiHost) {
		fmt.Printf("    http://%s/v1\n", net.JoinHostPort(ip, apiPort))
	}

	fmt.Println("\n  UI:")
	for _, ip := range listAddresses(uiHost) {
		fmt.Printf("    http://%s/ui\n", net.JoinHostPort(ip, uiPort))
	}

	fmt.Println()
}

// listAddresses returns the IP addresses that correspond to a bind host.
// An empty host or 0.0.0.0 expands to localhost plus all non-loopback IPv4 addresses.
func listAddresses(host string) []string {
	if host != "" && host != "0.0.0.0" && host != "::" {
		return []string{host}
	}
	ips := []string{"localhost"}
	ifaces, err := net.Interfaces()
	if err != nil {
		return ips
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.To4() == nil || ip.IsLoopback() {
				continue
			}
			ips = append(ips, ip.String())
		}
	}
	return ips
}
