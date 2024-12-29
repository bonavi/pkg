package openrtb

import (
	"encoding/json"

	"pkg/pointer"
)

// Inventory contains inventory specific attributes.
type Inventory struct {
	// Exchange-specific app ID.
	//
	// Recommended.
	ID string `json:"id" bson:"id"`

	// App name (may be aliased at the publisher’s request).
	Name string `json:"name" bson:"name"`

	// Domain of the app (e.g., “mygame.foo.com”).
	Domain string `json:"domain" bson:"domain"`

	// Array of IAB content categories of the app.
	Categories []ContentCategory `json:"cat" bson:"cat"`

	// Array of IAB content categories that describe the current section of the app.
	SectionCategories []ContentCategory `json:"sectioncat" bson:"sectioncat"`

	// Array of IAB content categories that describe the current page or view of the app.
	PageCategory []ContentCategory `json:"pagecat" bson:"pagecat"`

	// Indicates if the app has a privacy policy, where:
	//    0 = no;
	//    1 = yes.
	PrivacyPolicy int `json:"privacypolicy" bson:"privacypolicy"`

	// Details about the Publisher of the app.
	Publisher *Publisher `json:"publisher" bson:"publisher"`

	// Details about the Content within the app.
	Content *Content `json:"content" bson:"content"`

	// Comma separated list of keywords about the app.
	//
	// FIXME: keywords can be a string or an array strings.
	Keywords string `json:"keywords" bson:"keywords"`

	// Placeholder for exchange-specific extensions to OpenRTB.
	Ext json.RawMessage `json:"ext,omitempty" bson:"ext"`
}

func (i *Inventory) Copy() Inventory {

	var categories []ContentCategory
	if len(i.Categories) != 0 {
		categories = make([]ContentCategory, len(i.Categories))
		copy(categories, i.Categories)
	}

	var sectionCategories []ContentCategory
	if len(i.SectionCategories) != 0 {
		sectionCategories = make([]ContentCategory, len(i.SectionCategories))
		copy(sectionCategories, i.SectionCategories)
	}

	var pageCategory []ContentCategory
	if len(i.PageCategory) != 0 {
		pageCategory = make([]ContentCategory, len(i.PageCategory))
		copy(pageCategory, i.PageCategory)
	}

	var publisher *Publisher
	if i.Publisher != nil {
		publisher = pointer.Pointer(i.Publisher.Copy())
	}

	var content *Content
	if i.Content != nil {
		content = pointer.Pointer(i.Content.Copy())
	}

	var ext []byte
	if len(i.Ext) != 0 {
		ext = make([]byte, len(i.Ext))
		copy(ext, i.Ext)
	}

	return Inventory{
		ID:                i.ID,
		Name:              i.Name,
		Domain:            i.Domain,
		Categories:        i.Categories,
		SectionCategories: i.SectionCategories,
		PageCategory:      i.PageCategory,
		PrivacyPolicy:     i.PrivacyPolicy,
		Publisher:         publisher,
		Content:           content,
		Keywords:          i.Keywords,
		Ext:               ext,
	}
}

// Details of the application calling for the impression.
type App struct {
	Inventory

	// A platform-specific application identifier intended to be unique to the app
	// and independent of the exchange. On Android, this should be a bundle or package
	// name (e.g., com.foo.mygame). On iOS, it is typically a numeric ID.
	Bundle string `json:"bundle" bson:"bundle"`

	// App store URL for an installed app; for IQG 2.1 compliance.
	StoreURL string `json:"storeurl" bson:"storeurl"`

	// Application version.
	Version string `json:"ver" bson:"ver"`

	//    0 = app is free;
	//    1 = the app is a paid version.
	Paid int `json:"paid" bson:"paid"`
}

func (a *App) Copy() App {

	inventory := a.Inventory.Copy()

	return App{
		Inventory: inventory,
		Bundle:    a.Bundle,
		StoreURL:  a.StoreURL,
		Version:   a.Version,
		Paid:      a.Paid,
	}
}

// Details of the website calling for the impression.
type Site struct {
	Inventory

	// URL of the page where the impression will be shown.
	Page string `json:"page" bson:"page"`

	// Referrer URL that caused navigation to the current page.
	Refferer string `json:"ref" bson:"ref"`

	// Search string that caused navigation to the current page.
	Search string `json:"search" bson:"search"`

	// Indicates if the site has been programmed to optimize layout when viewed on mobile
	// devices, where:
	//   0 = no;
	//   1 = yes.
	Mobile int `json:"mobile" bson:"mobile"`
}

func (s *Site) Copy() Site {

	inventory := s.Inventory.Copy()

	return Site{
		Inventory: inventory,
		Page:      s.Page,
		Refferer:  s.Refferer,
		Search:    s.Search,
		Mobile:    s.Mobile,
	}
}
