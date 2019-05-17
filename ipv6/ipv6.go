package ipv6

import (
	"fmt"
	"net"
	"sort"
	"sync"

	"github.com/pkg/errors"
)

func ToIPv6(in []byte) (ips []net.IP, err error) {
	if len(in)%net.IPv6len != 0 {
		return nil, errors.New("must be n * 16 length")
	}
	ipCount := len(in) / net.IPv6len
	ips = make([]net.IP, ipCount)
	for i := 0; i < ipCount; i++ {
		ips[i] = make(net.IP, net.IPv6len)
		copy(ips[i], in[i*net.IPv6len:(i+1)*net.IPv6len])
	}
	return
}

func FromIPv6(ips []net.IP) (out []byte, err error) {
	ipCount := len(ips)
	out = make([]byte, ipCount*net.IPv6len)
	for i := 0; i < ipCount; i++ {
		copy(out[i*net.IPv6len:(i+1)*net.IPv6len], ips[i])
	}

	return
}

func FromDomain(domain string) (out []byte, err error) {
	concurrentNum := 5
	retryCount := 3
	allIPv6 := make([]net.IP, 0, 4)
	var ipsKeys []int
	ipsMap := make(map[int][]net.IP)
	var ipsErr error
	wg := new(sync.WaitGroup)
	var syncLock sync.Mutex

	for i := 0; ; i++ {
		// Concurrent by group
		var successCount int
		wg.Add(concurrentNum)
		for j := 0; j < concurrentNum; j++ {
			go func(i, j int) {
				defer wg.Done()
				index := i*concurrentNum + j
				for a := 0; a < retryCount; a++ {
					ips, err := net.LookupIP(fmt.Sprintf("%02d.%s", index, domain))
					if err == nil {
						syncLock.Lock()
						ipsMap[index] = ips
						successCount++
						syncLock.Unlock()
						break
					} else {
						if j == 0 {
							ipsErr = err
						}
					}
				}
			}(i, j)
		}
		wg.Wait()

		if len(ipsMap) != 0 {
			if successCount < concurrentNum {
				break
			}
		} else {
			return nil, ipsErr
		}
	}

	for k := range ipsMap {
		ipsKeys = append(ipsKeys, k)
	}
	sort.Ints(ipsKeys)

	for _, key := range ipsKeys {
		ips := ipsMap[key]
		if len(ips) == 0 {
			return nil, errors.New("empty IP list")
		}
		if len(ips[0]) != net.IPv6len {
			return nil, errors.Errorf("unexpected IP: %s", ips[0])
		}
		allIPv6 = append(allIPv6, ips[0])
	}

	out, err = FromIPv6(allIPv6)
	if err != nil {
		return nil, errors.Errorf("convert from IPv6 failed: %v", err)
	}
	return
}
