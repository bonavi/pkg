package geolite

import (
	"net"
	"sync/atomic"

	"github.com/oschwald/geoip2-golang"

	"pkg/errors"
)

// Reader - обертка над geoip2.Reader, чтобы обезопасить работу с преинициализированным geoip2.Reader
type Reader struct {
	geoLite atomic.Pointer[geoip2.Reader]
}

func (r *Reader) Init(filePath string) error {
	geoLite, err := geoip2.Open(filePath)
	if err != nil {
		return errors.Default.Wrap(err)
	}
	r.geoLite.Store(geoLite)
	return nil
}

func (r *Reader) City(ip net.IP) (record *geoip2.City, err error) {
	if !r.IsInitialized() {
		return nil, errors.Default.New("geoLite is nil")
	}
	return r.geoLite.Load().City(ip)
}

func (r *Reader) Close() error {
	if !r.IsInitialized() {
		return errors.Default.New("geoLite is nil")
	}
	return r.geoLite.Load().Close()
}

func (r *Reader) IsInitialized() bool {
	return r.geoLite.Load() != nil
}
