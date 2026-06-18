package openrtb

import (
	"pkg/openrtb/supplyChain"
	"pkg/pointer"
)

// Request source details on post-auction decisioning (e.g., header bidding).
type Source struct {
	// Entity responsible for the final impression sale decision, where:
	//    0 = exchange;
	//    1 = upstream source.
	//
	// Recommended.
	FinalSaleDecision int `json:"fd,omitempty" bson:"fd"`

	// Transaction ID that must be common across all participants in this bid request
	// (e.g., potentially multiple exchanges).
	//
	// Recommended.
	TransactionID string `json:"tid,omitempty" bson:"tid"`

	// Payment ID chain string containing embedded syntax described in the TAG Payment
	// ID Protocol v1.0.
	//
	// Recommended.
	PaymentChain string `json:"pchain,omitempty" bson:"pchain"`

	// Placeholder for exchange-specific extensions to OpenRTB.
	Ext *SourceExt `json:"ext,omitempty" bson:"ext"`
}

type SourceExt struct {
	Schain *supplyChain.SupplyChain `json:"schain,omitempty" bson:"schain"`
}

func (s *SourceExt) copy() SourceExt {
	var schain *supplyChain.SupplyChain
	if s.Schain != nil {
		schain = pointer.Pointer(s.Schain.Copy())
	}

	return SourceExt{
		Schain: schain,
	}
}

func (s *Source) Copy() Source {

	var ext *SourceExt
	if s.Ext != nil {
		ext = pointer.Pointer(s.Ext.copy())
	}

	return Source{
		FinalSaleDecision: s.FinalSaleDecision,
		TransactionID:     s.TransactionID,
		PaymentChain:      s.PaymentChain,
		Ext:               ext,
	}
}
