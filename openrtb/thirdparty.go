package openrtb

import "encoding/json"

// ThirdParty abstract attributes.
type ThirdParty struct {
	// Content producer or originator ID. Useful if content is syndicated and may be
	// posted on a site using embed tags.
	ID string `json:"id,omitempty" bson:"id"`

	// Content producer or originator name (e.g., “Warner Bros”)
	Name string `json:"name,omitempty" bson:"name"`

	// Array of IAB content categories that describe the content producer.
	Categories []ContentCategory `json:"cat,omitempty" bson:"cat"`

	// Highest level domain of the content producer (e.g., “producer.com”).
	Domain string `json:"domain,omitempty" bson:"domain"`

	// Placeholder for exchange-specific extensions to OpenRTB.
	Ext json.RawMessage `json:"ext,omitempty" bson:"ext"`
}

func (t *ThirdParty) Copy() ThirdParty {

	var categories []ContentCategory
	if len(t.Categories) != 0 {
		categories = make([]ContentCategory, len(t.Categories))
		copy(categories, t.Categories)
	}

	var ext []byte
	if len(t.Ext) != 0 {
		ext = make([]byte, len(t.Ext))
		copy(ext, t.Ext)
	}

	return ThirdParty{
		ID:         t.ID,
		Name:       t.Name,
		Categories: categories,
		Domain:     t.Domain,
		Ext:        ext,
	}
}

// Entity that controls the content of and distributes the site or app.
type Publisher ThirdParty

func (p *Publisher) Copy() Publisher {

	// Конвертируем указатель на Publisher в указатель на ThirdParty
	thirdParty := (*ThirdParty)(p)

	return Publisher(thirdParty.Copy())
}

// Producer of the content; not necessarily the publisher (e.g., syndication).
type Producer ThirdParty

func (p *Producer) Copy() Producer {

	// Конвертируем указатель на Producer в указатель на ThirdParty
	thirdParty := (*ThirdParty)(p)

	return Producer(thirdParty.Copy())
}
