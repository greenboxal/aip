package config

import (
	"errors"
	"fmt"
	"net"
	"net/netip"
	"net/url"
	"strconv"
	"strings"

	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
)

type NetworkConfig struct {
	GlobalPortOffset uint16 `json:"global_port_offset"`
}

type ListenEndpointID interface {
	Name() string
	Defaults() []string
	String() string
}

type listenEndpointID struct {
	name         string
	defaultValue []string
}

func (l listenEndpointID) String() string {
	return l.name
}

func (l listenEndpointID) Name() string {
	return l.name
}

func (l listenEndpointID) Defaults() []string {
	return l.defaultValue
}

func DefineListenEndpoint(name string, defaultValue ...string) ListenEndpointID {
	return listenEndpointID{name: name, defaultValue: defaultValue}
}

type NetworkManager struct {
	config *NetworkConfig

	listenAddresses map[ListenEndpointID][]multiaddr.Multiaddr
}

func NewNetworkManager(config *NetworkConfig) *NetworkManager {
	return &NetworkManager{
		config: config,
	}
}

func (l *NetworkManager) SetListenEndpoint(name ListenEndpointID, endpoint string) {
	var addrs []string

	if endpoint != "" {
		addrs = []string{endpoint}
	}

	l.SetListenEndpoints(name, addrs)
}

func (l *NetworkManager) SetListenEndpoints(name ListenEndpointID, endpoints []string) {
	addrs := make([]multiaddr.Multiaddr, 0, len(endpoints))

	for _, str := range endpoints {
		for _, addrStr := range strings.Split(str, ";") {
			ma, err := multiaddr.NewMultiaddr(addrStr)

			if err != nil {
				panic(err)
			}

			addrs = append(addrs, ma)
		}
	}

	l.SetListenMultiaddrs(name, addrs)
}

func (l *NetworkManager) SetListenMultiaddrs(name ListenEndpointID, addrs []multiaddr.Multiaddr) {
	l.listenAddresses[name] = addrs
}

func (l *NetworkManager) ListenMultiaddr(name ListenEndpointID) multiaddr.Multiaddr {
	addrs := l.ListenMultiaddrs(name)

	if len(addrs) == 0 {
		return nil
	}

	if len(addrs) > 1 {
		panic("cannot use multiple addresses")
	}

	return addrs[0]
}

func (l *NetworkManager) ListenMultiaddrs(name ListenEndpointID) []multiaddr.Multiaddr {
	result, ok := l.listenAddresses[name]

	if !ok {
		l.SetListenEndpoints(name, name.Defaults())

		return l.ListenMultiaddrs(name)
	}

	return result
}

func (l *NetworkManager) ListenEndpoint(name ListenEndpointID) net.Addr {
	addrs := l.ListenEndpoints(name)

	if len(addrs) == 0 {
		return nil
	}

	if len(addrs) > 1 {
		panic("cannot use multiple addresses")
	}

	return addrs[0]
}
func (l *NetworkManager) ListenEndpoints(name ListenEndpointID) []net.Addr {
	addrs, ok := l.listenAddresses[name]

	if !ok {
		if !ok {
			l.SetListenEndpoints(name, name.Defaults())

			return l.ListenEndpoints(name)
		}
	}

	results := make([]net.Addr, len(addrs))

	for i, v := range addrs {
		ma := l.patchMultiAddr(v)

		multiaddr.ForEach(ma, func(c multiaddr.Component) bool {
			switch c.Protocol().Code {
			case multiaddr.P_IP4:
			case multiaddr.P_IP6:
			case multiaddr.P_TCP:
			case multiaddr.P_UDP:
			case multiaddr.P_DNS:
			case multiaddr.P_DNS4:
			case multiaddr.P_DNS6:
				// Nothing

			default:
				ma = ma.Decapsulate(&c)
			}

			return true
		})

		a, err := manet.ToNetAddr(ma)

		if err != nil {
			endpointUrl, err := MultiaddrToUrl(v)

			if err != nil {
				panic(err)
			}

			udpPort, isUdp := ma.ValueForProtocol(multiaddr.P_UDP)

			if isUdp == nil && udpPort != "" {
				ip, err := net.ResolveUDPAddr("udp", endpointUrl.Host)

				if err != nil {
					panic(err)
				}

				a = ip
			} else {
				ip, err := net.ResolveTCPAddr("tcp", endpointUrl.Host)

				if err != nil {
					panic(err)
				}

				a = ip
			}
		}

		results[i] = a
	}

	return results
}

func (l *NetworkManager) patchAddresses(network string, endpoints []string) []net.Addr {
	result := make([]net.Addr, len(endpoints))

	for i, v := range endpoints {
		result[i] = l.patchAddress(network, v)
	}

	return result
}

func (l *NetworkManager) patchAddress(network, endpoint string) net.Addr {
	addr, err := netip.ParseAddrPort(endpoint)

	if err != nil {
		panic(err)
	}

	addr = netip.AddrPortFrom(addr.Addr(), addr.Port()+l.config.GlobalPortOffset)

	return localAddress{network: network, endpoint: addr.String()}
}

func (l *NetworkManager) patchMultiAddr(ma multiaddr.Multiaddr) multiaddr.Multiaddr {
	var components []multiaddr.Multiaddr

	multiaddr.ForEach(ma, func(c multiaddr.Component) bool {
		switch c.Protocol().Code {
		case multiaddr.P_TCP:
			c = l.patchMultiAddrComponent(c)
		case multiaddr.P_UDP:
			c = l.patchMultiAddrComponent(c)
		}

		components = append(components, &c)

		return true
	})

	return multiaddr.Join(components...)
}

func (l *NetworkManager) patchMultiAddrComponent(c multiaddr.Component) multiaddr.Component {
	port, err := strconv.ParseUint(c.Value(), 10, 16)

	if err != nil {
		panic(err)
	}

	port = port + uint64(l.config.GlobalPortOffset)

	patched, err := multiaddr.NewComponent(c.Protocol().Name, strconv.FormatUint(port, 10))

	if err != nil {
		panic(err)
	}

	return *patched
}

type localAddress struct {
	network  string
	endpoint string
}

func (l localAddress) Network() string {
	return l.network
}

func (l localAddress) String() string {
	return l.endpoint
}

func MultiaddrToUrl(addr multiaddr.Multiaddr) (result *url.URL, err error) {
	var hostname string
	var port string

	result = &url.URL{}

	multiaddr.ForEach(addr, func(c multiaddr.Component) bool {
		switch c.Protocol().Code {
		case multiaddr.P_IP4:
			hostname = c.Value()

		case multiaddr.P_IP6:
			hostname = fmt.Sprintf("[%s]", c.Value())

		case multiaddr.P_DNS:
			fallthrough
		case multiaddr.P_DNS4:
			fallthrough
		case multiaddr.P_DNS6:
			hostname = c.Value()

		case multiaddr.P_TCP:
			fallthrough
		case multiaddr.P_UDP:
			port = c.Value()

		case multiaddr.P_HTTP:
			if result.Scheme == "tls" {
				result.Scheme = "https"
			} else {
				result.Scheme = "http"
			}

			result.Path = c.Value()

		case multiaddr.P_TLS:
			if result.Scheme == "http" {
				result.Scheme = "https"
			} else {
				result.Scheme = "tls"
			}

		case multiaddr.P_HTTPS:
			result.Scheme = "https"

			result.Path = c.Value()

		default:
			err = errors.New("invalid protocol")
			return false
		}

		return true
	})

	result.Host = fmt.Sprintf("%s:%s", hostname, port)

	return
}
