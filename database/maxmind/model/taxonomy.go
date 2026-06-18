package model

type TaxonomyType uint8

const (
	TaxonomyTypeUnspecified TaxonomyType = 0
	TaxonomyTypeCountry     TaxonomyType = 1
	TaxonomyTypeRegion      TaxonomyType = 2
	TaxonomyTypeCity        TaxonomyType = 3
)

type Taxonomy struct {
	ParentId uint64
	NameRu   string
	NameEn   string
	Type     TaxonomyType
}
