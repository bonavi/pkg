package macros

import (
	"net/url"
	"pkg/openrtb"
	"strings"
)

// ExpandPriceAndCurrency заменяет макросы цены и валюты в строке.
// Сырые токены получают сырое значение, URL-encoded токены — закодированное значение.
func ExpandPriceAndCurrency(s, price, currency string) string {
	repl := strings.NewReplacer(
		// Raw tokens
		openrtb.AuctionPriceMacros, price,
		"${AUCTION_PRICE}", price,
		openrtb.AuctionCurrencyMacros, currency,
		"${AUCTION_CURRENCY}", currency,

		// URL-encoded tokens
		"%24%7BAUCTION_PRICE%7D", url.QueryEscape(price),
		"%24%7BAUCTION_CURRENCY%7D", url.QueryEscape(currency),
	)
	return repl.Replace(s)
}

// ExpandClickURL заменяет макрос ${CLICK_URL} в adm на переданный URL.
// Заменяет как raw, так и URL-encoded вариант макроса.
func ExpandClickURL(adm, clickURL string) string {
	repl := strings.NewReplacer(
		openrtb.AuctionClickURLMacros, clickURL,
		"%24%7BCLICK_URL%7D", url.QueryEscape(clickURL),
	)
	return repl.Replace(adm)
}

// ExpandAuctionMacros заменяет все OpenRTB макросы аукциона в строке.
// AUCTION_PRICE, AUCTION_CURRENCY, AUCTION_IMP_ID, AUCTION_SEAT_ID, AUCTION_AD_ID
func ExpandAuctionMacros(s, price, currency string) string {
	repl := strings.NewReplacer(
		openrtb.AuctionPriceMacros, price,
		openrtb.AuctionCurrencyMacros, currency,

		"%24%7BAUCTION_PRICE%7D", url.QueryEscape(price),
		"%24%7BAUCTION_CURRENCY%7D", url.QueryEscape(currency),
	)
	return repl.Replace(s)
}
