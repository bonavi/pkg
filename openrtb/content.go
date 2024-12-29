package openrtb

import (
	"encoding/json"

	"pkg/pointer"
)

// Details about the published content itself, within which the ad will be shown.
type Content struct {
	// ID uniquely identifying the content.
	ID string `json:"id" bson:"id"`

	// Episode number.
	Episode int `json:"episode" bson:"episode"`

	// Content title.
	//
	// Video Examples: “Search Committee” (television), “A New Hope” (movie),
	// or “Endgame” (made for web).
	//
	// Non-Video Example: “Why an Antarctic Glacier Is Melting So Quickly”
	// (Time magazine article).
	Title string `json:"title" bson:"title"`

	// Content series.
	//
	// Video Examples: “The Office” (television), “Star Wars” (movie),
	// or “Arby ‘N’ The Chief” (made for web).
	//
	// Non-Video Example: “Ecocentric” (Time Magazine blog).
	Series string `json:"series" bson:"series"`

	// Content season (e.g., “Season 3”).
	Season string `json:"season" bson:"season"`

	// Artist credited with the content.
	Artist string `json:"artist" bson:"artist"`

	// Genre that best describes the content (e.g., rock, pop, etc).
	Genre string `json:"genre" bson:"genre"`

	// Album to which the content belongs; typically for audio.
	Album string `json:"album" bson:"album"`

	// International Standard Recording Code conforming to ISO-3901.
	ISRC string `json:"isrc" bson:"isrc"`

	// Details about the content Producer.
	Producer *Producer `json:"producer" bson:"producer"`

	// URL of the content, for buy-side contextualization or review.
	URL string `json:"url" bson:"url"`

	// Array of IAB content categories that describe the content producer.
	Categories []ContentCategory `json:"cat" bson:"cat"`

	// Production quality.
	ProductionQuality ProductionQuality `json:"prodq" bson:"prodq"`

	// Video quality.
	//
	// Deprecated: deprecated in favor of prodq.
	VideoQuality ProductionQuality `json:"videoquality" bson:"videoquality"`

	// Type of content (game, video, text, etc.).
	Context ContentContext `json:"context" bson:"context"`

	// Content rating (e.g., MPAA).
	ContentRating string `json:"contentrating" bson:"contentrating"`

	// User rating of the content (e.g., number of stars, likes, etc.).
	UserRating string `json:"userrating" bson:"userrating"`

	// Media rating per IQG guidelines.
	MediaRating IQGRating `json:"qagmediarating" bson:"qagmediarating"`

	// Comma separated list of keywords describing the content.
	//
	// FIXME: keywords can be a string or an array strings.
	Keywords string `json:"keywords" bson:"keywords"`

	// 0 = not live, 1 = content is live (e.g., stream, live blog).
	LiveStream int `json:"livestream" bson:"livestream"`

	// 0 = indirect, 1 = direct.
	SourceRelationship int `json:"sourcerelationship" bson:"sourcerelationship"`

	// Length of content in seconds; appropriate for video or audio.
	Length int `json:"len" bson:"len"`

	// Content language using ISO-639-1-alpha-2.
	Language string `json:"language" bson:"language"`

	// Indicator of whether or not the content is embeddable (e.g., an embeddable video player),
	// where:
	//   0 = no;
	//   1 = yes.
	Embeddable int `json:"embeddable" bson:"embeddable"`

	// Additional content data. Each Data object represents a different data source.
	Data []Data `json:"data" bson:"data"`

	// Placeholder for exchange-specific extensions to OpenRTB.
	Ext json.RawMessage `json:"ext,omitempty" bson:"ext"`
}

func (c *Content) Copy() Content {

	var producer *Producer
	if c.Producer != nil {
		producer = pointer.Pointer(c.Producer.Copy())
	}

	var categories []ContentCategory
	if len(c.Categories) != 0 {
		categories = make([]ContentCategory, len(c.Categories))
		copy(categories, c.Categories)
	}

	var data []Data
	if len(c.Data) != 0 {
		data = make([]Data, len(c.Data))
		for i := range c.Data {
			data[i] = c.Data[i].Copy()
		}
	}

	var ext []byte
	if len(c.Ext) != 0 {
		ext = make([]byte, len(c.Ext))
		copy(ext, c.Ext)
	}

	return Content{
		ID:                 c.ID,
		Episode:            c.Episode,
		Title:              c.Title,
		Series:             c.Series,
		Season:             c.Season,
		Artist:             c.Artist,
		Genre:              c.Genre,
		Album:              c.Album,
		ISRC:               c.ISRC,
		Producer:           producer,
		URL:                c.URL,
		Categories:         categories,
		ProductionQuality:  c.ProductionQuality,
		VideoQuality:       c.VideoQuality,
		Context:            c.Context,
		ContentRating:      c.ContentRating,
		UserRating:         c.UserRating,
		MediaRating:        c.MediaRating,
		Keywords:           c.Keywords,
		LiveStream:         c.LiveStream,
		SourceRelationship: c.SourceRelationship,
		Length:             c.Length,
		Language:           c.Language,
		Embeddable:         c.Embeddable,
		Data:               data,
		Ext:                ext,
	}
}
