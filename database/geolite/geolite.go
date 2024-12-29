package geolite

import (
	"net"

	"github.com/oschwald/geoip2-golang"

	"pkg/errors"
)

func NewClientGeoLite(filePath string) (customGeoLite *Reader, err error) {

	// Открываем файл с базой геолокации
	geoLite, err := geoip2.Open(filePath)
	if err != nil {
		return nil, errors.InternalServer.Wrap(err)
	}

	return &Reader{geoLite: geoLite}, nil
}

// Reader - обертка над geoip2.Reader, чтобы обезопасить работу с преинициализированным geoip2.Reader
type Reader struct {
	geoLite *geoip2.Reader
}

func (r *Reader) City(ip net.IP) (record *geoip2.City, err error) {
	if !r.IsInitialized() {
		return nil, errors.InternalServer.New("geoLite is nil")
	}
	return r.geoLite.City(ip)
}

func (r *Reader) Close() error {
	if !r.IsInitialized() {
		return errors.InternalServer.New("geoLite is nil")
	}
	return r.geoLite.Close()
}

func (r *Reader) IsInitialized() bool {
	return r.geoLite != nil
}
