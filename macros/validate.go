package macros

import (
	"pkg/openrtb"
	"strings"
)

// ContainsAny возвращает true, если хотя бы один из подстрок встретился в строке
func ContainsAny(s string, subs []string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// ContainsPriceMacros проверяет наличие макросов цены в строке
func ContainsPriceMacros(s string) bool {
	priceMacros := []string{
		openrtb.AuctionPriceMacros,
		"${AUCTION_PRICE}",
		"%24%7BAUCTION_PRICE%7D", // URL-encoded ${AUCTION_PRICE}
	}
	return ContainsAny(s, priceMacros)
}


// ContainsADMURL проверяет наличие макросов ADM_URL в строке
func ContainsADMURL(s string) bool {
	admURLVariants := []string{
		openrtb.ADMURLMacros, // ${ADM_URL}
		"${ADM_URL}",
		"%24%7BADM_URL%7D", // url-encoded ${ADM_URL}
	}
	return ContainsAny(s, admURLVariants)
}
