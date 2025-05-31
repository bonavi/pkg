package openrtb

import (
	"encoding/json"

	"pkg/pointer"
)

// Details of the device on which the content and impressions are displayed.
type Device struct {
	// Browser user agent string.
	//
	// Recommended.
	UserAgent string `json:"ua,omitempty" bson:"ua"`

	// Location of the device assumed to be the user’s current
	// location defined by a Geo object.
	//
	// Recommended.
	Geo *Geo `json:"geo,omitempty" bson:"geo"`

	// Standard “Do Not Track” flag as set in the header by the browser, where:
	//    0 = tracking is unrestricted;
	//    1 = do not track.
	//
	// Recommended.
	DNT int `json:"dnt,omitempty" bson:"dnt"`

	// “Limit Ad Tracking” signal commercially endorsed (e.g., iOS, Android), where:
	//    0 = tracking is unrestricted;
	//    1 = tracking must be limited per commercial guidelines.
	//
	// Recommended.
	LMT int `json:"lmt,omitempty" bson:"lmt"`

	// IPv4 address closest to device.
	//
	// Recommended.
	IP string `json:"ip,omitempty" bson:"ip"`

	// IP address closest to device as IPv6.
	IPv6 string `json:"ipv6,omitempty" bson:"ipv6"`

	// The general type of device.
	DeviceType DeviceType `json:"devicetype,omitempty" bson:"devicetype"`

	// Device make (e.g., “Apple”).
	Make string `json:"make,omitempty" bson:"make"`

	// Device model (e.g., “iPhone”).
	Model string `json:"model,omitempty" bson:"model"`

	// Device operating system (e.g., “iOS”).
	OS string `json:"os,omitempty" bson:"os"`

	// Device operating system version (e.g., “3.1.2”).
	OSVersion string `json:"osv,omitempty" bson:"osv"`

	// Hardware version of the device (e.g., “5S” for iPhone 5S).
	HWVersion string `json:"hwv,omitempty" bson:"hwv"`

	// Physical height of the screen in pixels.
	Height int `json:"h,omitempty" bson:"h"`

	// Physical width of the screen in pixels.
	Width int `json:"w,omitempty" bson:"w"`

	// Screen size as pixels per linear inch.
	PPI int `json:"ppi,omitempty" bson:"ppi"`

	// The ratio of physical pixels to device independent pixels.
	PixelRatio float64 `json:"pxratio,omitempty" bson:"pxratio"`

	// Support for JavaScript, where:
	//    0 = no;
	//    1 = yes.
	JS int `json:"js,omitempty" bson:"js"`

	// Indicates if the geolocation API will be available to JavaScript code running
	// in the banner, where:
	//    0 = no;
	//    1 = yes.
	GeoFetch int `json:"geofetch,omitempty" bson:"geofetch"`

	// Version of Flash supported by the browser.
	FlashVersion string `json:"flashver,omitempty" bson:"flashver"`

	// Browser language using ISO-639-1-alpha-2.
	Language string `json:"language,omitempty" bson:"language"`

	// Carrier or ISP (e.g., “VERIZON”) using exchange curated string
	// names which should be published to bidders a priori.
	Carrier string `json:"carrier,omitempty" bson:"carrier"`

	// Mobile carrier as the concatenated MCC-MNC code (e.g., “310-005” identifies Verizon
	// Wireless CDMA in the USA). Refer to https://en.wikipedia.org/wiki/Mobile_country_code
	// for further examples. Note that the dash between the MCC and MNC parts is required
	// to remove parsing ambiguity.
	MCCMNC string `json:"mccmnc,omitempty" bson:"mccmnc"`

	// Network connection type.
	ConnectionType ConnectionType `json:"connectiontype,omitempty" bson:"connectiontype"`

	// ID sanctioned for advertiser use in the clear (i.e., not hashed).
	IFA string `json:"ifa,omitempty" bson:"ifa"`

	// Hardware device ID (e.g., IMEI); hashed via SHA1.
	IDSHA1 string `json:"didsha1,omitempty" bson:"didsha1"`

	// Hardware device ID (e.g., IMEI); hashed via MD5.
	IDMD5 string `json:"didmd5,omitempty" bson:"didmd5"`

	// Platform device ID (e.g., Android ID); hashed via SHA1.
	PIDSHA1 string `json:"dpidsha1,omitempty" bson:"dpidsha1"`

	// Platform device ID (e.g., Android ID); hashed via MD5.
	PIDMD5 string `json:"dpidmd5,omitempty" bson:"dpidmd5"`

	// MAC address of the device; hashed via SHA1.
	MacSHA1 string `json:"macsha1,omitempty" bson:"macsha1"`

	// MAC address of the device; hashed via MD5.
	MacMD5 string `json:"macmd5,omitempty" bson:"macmd5"`

	// Placeholder for exchange-specific extensions to OpenRTB.
	Ext json.RawMessage `json:"ext,omitempty" bson:"ext"`
}

func (d *Device) Copy() Device {

	var geo *Geo
	if d.Geo != nil {
		geo = pointer.Pointer(d.Geo.Copy())
	}

	var ext []byte
	if len(d.Ext) != 0 {
		ext = make([]byte, len(d.Ext))
		copy(ext, d.Ext)
	}

	return Device{
		UserAgent:      d.UserAgent,
		Geo:            geo,
		DNT:            d.DNT,
		LMT:            d.LMT,
		IP:             d.IP,
		IPv6:           d.IPv6,
		DeviceType:     d.DeviceType,
		Make:           d.Make,
		Model:          d.Model,
		OS:             d.OS,
		OSVersion:      d.OSVersion,
		HWVersion:      d.HWVersion,
		Height:         d.Height,
		Width:          d.Width,
		PPI:            d.PPI,
		PixelRatio:     d.PixelRatio,
		JS:             d.JS,
		GeoFetch:       d.GeoFetch,
		FlashVersion:   d.FlashVersion,
		Language:       d.Language,
		Carrier:        d.Carrier,
		MCCMNC:         d.MCCMNC,
		ConnectionType: d.ConnectionType,
		IFA:            d.IFA,
		IDSHA1:         d.IDSHA1,
		IDMD5:          d.IDMD5,
		PIDSHA1:        d.PIDSHA1,
		PIDMD5:         d.PIDMD5,
		MacSHA1:        d.MacSHA1,
		MacMD5:         d.MacMD5,
		Ext:            ext,
	}
}
