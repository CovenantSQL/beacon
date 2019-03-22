package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/CovenantSQL/beacon/ipv6"
	log "github.com/Sirupsen/logrus"
	"io"
	"net"
	"os"
	"strings"
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
			log.Errorf("open stdin failed: %v", err)
			os.Exit(1)
		}
		if fi.Mode()&os.ModeCharDevice == 0 && fi.Size() > 0 {
			reader := bufio.NewReader(os.Stdin)
			var allInput []byte

			for {
				input, err := reader.ReadByte()
				if err != nil && err == io.EOF {
					break
				}
				allInput = append(allInput, input)
			}
			if trim {
				allInput = []byte(strings.TrimSpace(string(allInput)))
			}
			fmt.Print("Generated IPv6 addr:\n;; AAAA Records:\n")
			ips, err := ipv6.ToIPv6(allInput)
			if err != nil {
				log.Errorf("failed to convert IPv6: %v", err)
				os.Exit(1)
			}
			if len(ips) > 100 {
				log.Errorf("generated IPv6 addr count above 100: %d", len(ips))
				os.Exit(1)
			}
			if len(ips) == 0 {
				log.Errorf("failed to generate IPv6 for %s", allInput)
				os.Exit(1)
			}
			for i, ip := range ips {
				fmt.Printf("%02d.%s	1	IN	AAAA	%s\n", i, domain, ip)
			}
			return
		} else {
			if domain == "example.org" {
				log.Println("please specify the source domain")
				os.Exit(2)
			}
			allIPv6 := make([]net.IP, 0, 4)
			for i := 0; ; i++ {
				ips, err := net.LookupIP(fmt.Sprintf("%02d.%s", i, domain))
				if err != nil {
					if _, ok := err.(*net.DNSError); ok && strings.Contains(err.Error(), "no such host") {
						break
					} else {
						log.Errorf("DNS error: %v", err)
						os.Exit(2)
					}
				} else {
					if len(ips) == 0 {
						log.Error("empty IP list")
						os.Exit(2)
					}
					if len(ips[0]) != net.IPv6len {
						log.Errorf("unexpected IP: %s", ips[0])
						os.Exit(2)
					}
					allIPv6 = append(allIPv6, ips[0])
				}

			}
			out, err := ipv6.FromIPv6(allIPv6)
			if err != nil {
				log.Errorf("convert from IPv6 failed: %v", err)
				os.Exit(2)
			}
			log.Infof("#### %s ####\n", domain)
			if hex {
				fmt.Printf("%x\n", out)
			} else {
				fmt.Printf("%s\n", string(out))
			}
			log.Infof("#### %s ####\n", domain)
		}
	}
}
