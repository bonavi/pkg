package vast

import (
	"fmt"
)

type VASTRequest struct {
	IFA       string `schema:"ifa" json:"ifa"`
	UserAgent string `schema:"ua" json:"ua" validate:"required"`
	API       int    `schema:"api" json:"api"`
	Width     int    `schema:"w" json:"w"`
	Height    int    `schema:"h" json:"h"`

	IP string `schema:"-" json:"ip"`

	// App
	Bundle  *string `schema:"bundle" json:"bundle"`
	AppName *string `schema:"appname" json:"appname"`

	// Site
	Domain *string `schema:"domain" json:"domain"`
	UserID *string `schema:"uid" json:"uid"`
}

func (v VASTRequest) Validate() error {

	switch {
	case v.Domain != nil && v.UserID != nil:
		break
	case v.Bundle != nil && v.AppName != nil:
		break
	default:
		return fmt.Errorf("invalid vast request, needed (domain and uid) or (bundle and appname)")
	}

	return nil
}
