package openrtb

import (
	"pkg/openrtb/stableID"
	"pkg/pointer"
)

// Human user of the device; audience for advertising.
type User struct {
	// Exchange-specific ID for the user.
	//
	// At least one of id or buyeruid is recommended.
	ID string `json:"id,omitempty" bson:"id"`

	// Buyer-specific ID for the user as mapped by the exchange for the buyer.
	//
	// At least one of buyeruid or id is recommended.
	BuyerUID string `json:"buyeruid,omitempty" bson:"buyeruid"`

	// Year of birth as a 4-digit integer.
	YearOfBirth int `json:"yob,omitempty" bson:"yob"`

	// Gender, where:
	//   M = male;
	//   F = female;
	//   O = known to be other (i.e., omitted is unknown).
	Gender string `json:"gender,omitempty" bson:"gender"`

	// Comma separated list of keywords, interests, or intent.
	//
	// FIXME: keywords can be a string or an array strings.
	Keywords string `json:"keywords,omitempty" bson:"keywords"`

	// Optional feature to pass bidder data that was set in the exchange’s cookie.
	// The string must be in base85 cookie safe characters and be in any format.
	// Proper JSON encoding must be used to include “escaped” quotation marks.
	CustomData string `json:"customdata,omitempty" bson:"customdata"`

	// Location of the user’s home base defined by a Geo object. This is not necessarily
	// their current location.
	Geo *Geo `json:"geo,omitempty" bson:"geo"`

	// Additional user data. Each Data object represents a different data source.
	Data []Data `json:"data,omitempty" bson:"data"`

	// Placeholder for exchange-specific extensions to OpenRTB.
	Ext *UserExt `json:"ext,omitempty" bson:"ext"`
}

type UserExt struct {
	StableID *stableID.StableID `json:"stable_id,omitempty" bson:"stable_id"`
}

func (u *UserExt) copy() UserExt {

	var stableID *stableID.StableID
	if u.StableID != nil {
		stableID = pointer.Pointer(u.StableID.Copy())
	}

	return UserExt{
		StableID: stableID,
	}
}

func (u *User) Copy() User {

	var geo *Geo
	if u.Geo != nil {
		geo = pointer.Pointer(u.Geo.Copy())
	}

	var data []Data
	if len(u.Data) != 0 {
		data = make([]Data, len(u.Data))
		for i := range u.Data {
			data[i] = u.Data[i].Copy()
		}
	}

	var ext *UserExt
	if u.Ext != nil {
		ext = pointer.Pointer(u.Ext.copy())
	}

	return User{
		ID:          u.ID,
		BuyerUID:    u.BuyerUID,
		YearOfBirth: u.YearOfBirth,
		Gender:      u.Gender,
		Keywords:    u.Keywords,
		CustomData:  u.CustomData,
		Geo:         geo,
		Data:        data,
		Ext:         ext,
	}
}
