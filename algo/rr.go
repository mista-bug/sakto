package lb

import (
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

func RoundRobinInit(addresses []string, requestCount *uint64) *httputil.ReverseProxy {
	addressLen := uint32(len(addresses))
	counter := atomic.AddUint64(requestCount, 1)
	index := counter % uint64(addressLen)
	remote, err := url.Parse(addresses[index])
	if err != nil {
		panic(err)
	}

	return httputil.NewSingleHostReverseProxy(remote)
}
