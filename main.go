package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net"
	"os"
	"strings"

	"github.com/CovenantSQL/beacon/ipv6"
)

var (
	mode   string
	domain string
	trim   bool
	hex    bool
)

func main() {
	flag.StringVar(&mode, "mode", "ipv6", "storage type for data")
	flag.StringVar(&domain, "domain", "example.org", "domain used to store data")
	flag.BoolVar(&trim, "trim", false, "trim whitespace for input")
	flag.BoolVar(&hex, "hex", false, "output in hex")
	flag.Parse()

	if mode == "ipv6" {
		fi, err := os.Stdin.Stat()
		if err != nil {
			slog.Error(fmt.Sprintf("open stdin failed: %v", err))
			os.Exit(1)
		}
		if fi.Mode()&os.ModeCharDevice == 0 {
			reader := bufio.NewReader(os.Stdin)
			allInput, err := ioutil.ReadAll(reader)
			if err != nil {
				slog.Error(fmt.Sprintf("failed to read inputs: %v", err))
				os.Exit(1)
			}
			if trim {
				allInput = []byte(strings.TrimSpace(string(allInput)))
			}
			if len(allInput) == 0 {
				slog.Error("blank input")
				os.Exit(1)
			}
			fmt.Print("Generated IPv6 addr:\n;; AAAA Records:\n")
			ips, err := ipv6.ToIPv6(allInput)
			if err != nil {
				slog.Error(fmt.Sprintf("failed to convert IPv6: %v", err))
				os.Exit(1)
			}
			if len(ips) > 100 {
				slog.Error(fmt.Sprintf("generated IPv6 addr count above 100: %d", len(ips)))
				os.Exit(1)
			}
			if len(ips) == 0 {
				slog.Error(fmt.Sprintf("failed to generate IPv6 for %s", allInput))
				os.Exit(1)
			}
			for i, ip := range ips {
				fmt.Printf("%02d.%s	1	IN	AAAA	%s\n", i, domain, ip)
			}
			return
		} else {
			if domain == "example.org" {
				slog.Info("please specify the source domain")
				os.Exit(2)
			}
			f := func(host string) ([]net.IP, error) {
				return net.LookupIP(host)
			}
			out, err := ipv6.FromDomain(domain, f)
			if err != nil {
				slog.Error(fmt.Sprintf("failed to get data from %s: %v", domain, err))
				os.Exit(2)
			}
			slog.Info(fmt.Sprintf("#### %s ####", domain))
			if hex {
				fmt.Printf("%x\n", out)
			} else {
				fmt.Printf("%s\n", string(out))
			}
			slog.Info(fmt.Sprintf("#### %s ####", domain))
		}
	}
}
