package config

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type Service struct {
	Name string `yaml:"name"`

	// A prefix matcher to select service based on the path part of the url
	Matcher  string   `yaml:"matcher"`
	Replicas []string `yaml:"replicas"`
}

// Config is a representation of the configuration
// given to the LB from a config source
type Config struct {
	Services []Service `yaml:"services"`

	// Strategy is the name of strategy to be used in load balancing between instances
	Strategy string `yaml:"strategy"`
}

// Server is an instance of a running server
type Server struct {
	URL   *url.URL
	Proxy *httputil.ReverseProxy
}

type ServerList struct {
	// Servers are the replicas
	Servers []*Server

	// Name is the name of the service
	Name string
	// the current server to forward the request to.
	// the next server should be (current + 1) % len(servers)
	current uint32
}

func (s *Server) Forward(rw http.ResponseWriter, r *http.Request) {
	s.Proxy.ServeHTTP(rw, r)
}

func (sl *ServerList) Next() uint32 {
	nxt := atomic.AddUint32(&sl.current, uint32(1))
	lenS := uint32(len(sl.Servers))
	return nxt % lenS
}
