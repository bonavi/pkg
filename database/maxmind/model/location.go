package model

type Location struct {
	Id        uint64 // Hash from Country + Region + City
	CityEn    string // Город
	CityRu    string // Город
	RegionEn  string // Область/Край
	RegionRu  string // Область/Край
	CountryEn string // Страна
	CountryRu string // Страна
	Extra     *LocationExtra
}

type LocationExtra struct {
	TimeZone  string  // Часовой пояс
	Latitude  float64 // Широта
	Longitude float64 // Долгота
	Radius    uint16  // Радиус
	Proxy     bool    // Прокси
}

func (loc *Location) Names() [][3]string {
	var strs [][3]string

	if loc.CityRu != "" {
		strs = append(strs, [3]string{
			loc.CityRu, loc.CityEn, loc.CountryRu + loc.RegionRu + loc.CityRu})
	}

	if loc.RegionRu != "" {
		strs = append(strs, [3]string{loc.RegionRu, loc.RegionEn, loc.CountryRu + loc.RegionRu})
	}

	return append(strs, [3]string{loc.CountryRu, loc.CountryEn, loc.CountryRu})
}
