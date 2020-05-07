package hole

import "net"

// NAT ...
type NAT interface {
	Mapping() (port int, err error)
	StopMapping() (err error)
	Remapping() (port int, err error)
	GetExternalAddress() (addr net.IP, err error)
}
