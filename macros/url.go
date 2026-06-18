package macros

import (
	"fmt"
	"net/url"
)

// ValidateHTTPURL проверяет, что строка является валидным HTTP(S) URL
func ValidateHTTPURL(s string) error {
	u, err := url.Parse(s)
	if err != nil {
		return err
	}
	if (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return fmt.Errorf("invalid http(s) URL")
	}
	return nil
}