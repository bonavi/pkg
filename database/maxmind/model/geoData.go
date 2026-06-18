package model

type Geodata struct {
	City         GeoDataCity           `maxminddb:"city"`
	Country      GeoDataCountry        `maxminddb:"country"`
	Location     GeodataLocation       `maxminddb:"location"`
	Subdivisions []GeoDataSubdivisions `maxminddb:"subdivisions"`
	Traits       GeoDataTraits         `maxminddb:"traits"`
}

type GeoDataCity struct {
	Names map[string]string `maxminddb:"names"`
}

type GeoDataCountry struct {
	IsoCode string            `maxminddb:"iso_code"`
	Names   map[string]string `maxminddb:"names"`
}

type GeodataLocation struct {
	AccuracyRadius uint16  `maxminddb:"accuracy_radius"`
	Latitude       float64 `maxminddb:"latitude"`
	Longitude      float64 `maxminddb:"longitude"`
	TimeZone       string  `maxminddb:"time_zone"`
}

type GeoDataSubdivisions struct {
	Names map[string]string `maxminddb:"names"`
}

type GeoDataTraits struct {
	IsAnonymousProxy bool `maxminddb:"is_anonymous_proxy"`
}

func (gd *Geodata) GetLocation(Extra bool) (*Location, string) {
	loc := Location{
		CityRu:    gd.City.Names["ru"],
		CityEn:    gd.City.Names["en"],
		CountryRu: gd.Country.Names["ru"],
		CountryEn: gd.Country.Names["en"],
	}

	if len(gd.Subdivisions) != 0 {
		loc.RegionRu = gd.Subdivisions[0].Names["ru"]
		loc.RegionEn = gd.Subdivisions[0].Names["en"]
	}

	if loc.CityEn == "" && loc.CityRu == "" && loc.RegionRu == "" && loc.RegionEn == "" && loc.CountryRu == "" && loc.CountryEn == "" {
		return nil, gd.Country.IsoCode
	}

	if Extra {
		loc.Extra = &LocationExtra{
			TimeZone:  gd.Location.TimeZone,
			Latitude:  gd.Location.Latitude,
			Longitude: gd.Location.Longitude,
			Radius:    gd.Location.AccuracyRadius,
			Proxy:     gd.Traits.IsAnonymousProxy,
		}
	}

	return &loc, gd.Country.IsoCode
}
