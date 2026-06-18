package openrtb_test

import (
	"testing"

	"pkg/openrtb"
)

func TestSource_Unmarshal(t *testing.T) {
	expected := openrtb.Source{
		FinalSaleDecision: 1,
		TransactionID:     "transaction-id",
		PaymentChain:      "payment-chain",
		Ext:               &openrtb.SourceExt{},
	}

	assertEqualJSON(t, "source", &expected)
}
