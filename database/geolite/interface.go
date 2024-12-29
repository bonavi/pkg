package geolite

import (
	"net"

	"github.com/oschwald/geoip2-golang"
)

type GeoLite interface {
	City(ip net.IP) (record *geoip2.City, err error)
	Close() error
}
