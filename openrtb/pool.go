package openrtb

import (
	"pkg/pointer"
	"sync"
)

var bidRequestPool = sync.Pool{
	New: func() any { return new(BidRequest) },
}

// WithBidRequest берёт *BidRequest из пула, копирует в него original, выполняет fn и возвращает объект в пул.
// Использовать только для синхронного кода внутри fn — не запускать горутины, которые переживут вызов fn.
// Для горутин использовать AcquireBidRequest и ReleaseBidRequest.
func WithBidRequest(original *BidRequest, fn func(*BidRequest)) {
	r := bidRequestPool.Get().(*BidRequest)
	r.fillFrom(original)
	defer func() {
		r.clearForPool()
		bidRequestPool.Put(r)
	}()
	fn(r)
}

// AcquireBidRequest берёт *BidRequest из пула и копирует в него original.
// Вызывающий обязан вызвать ReleaseBidRequest после завершения работы с объектом.
func AcquireBidRequest(original *BidRequest) *BidRequest {
	r := bidRequestPool.Get().(*BidRequest)
	r.fillFrom(original)
	return r
}

// ReleaseBidRequest очищает *BidRequest и возвращает объект в пул.
func ReleaseBidRequest(r *BidRequest) {
	r.clearForPool()
	bidRequestPool.Put(r)
}

// fillFrom копирует все поля из src в r, переиспользуя ёмкость существующих срезов.
func (r *BidRequest) fillFrom(src *BidRequest) {

	// Скалярные поля
	r.ID = src.ID
	r.Test = src.Test
	r.AuctionType = src.AuctionType
	r.TimeMax = src.TimeMax
	r.AllImpressions = src.AllImpressions

	// Импрешены — переиспользуем backing array если ёмкость достаточна
	r.Impressions = r.Impressions[:0]
	for i := range src.Impressions {
		r.Impressions = append(r.Impressions, src.Impressions[i].Copy())
	}

	// Указатели на вложенные объекты
	if src.App != nil {
		r.App = pointer.Pointer(src.App.Copy())
	} else {
		r.App = nil
	}

	if src.Site != nil {
		r.Site = pointer.Pointer(src.Site.Copy())
	} else {
		r.Site = nil
	}

	if src.Device != nil {
		r.Device = pointer.Pointer(src.Device.Copy())
	} else {
		r.Device = nil
	}

	if src.User != nil {
		r.User = pointer.Pointer(src.User.Copy())
	} else {
		r.User = nil
	}

	if src.Source != nil {
		r.Source = pointer.Pointer(src.Source.Copy())
	} else {
		r.Source = nil
	}

	if src.Regulations != nil {
		r.Regulations = pointer.Pointer(src.Regulations.Copy())
	} else {
		r.Regulations = nil
	}

	if src.Ext != nil {
		r.Ext = pointer.Pointer(src.Ext.copy())
	} else {
		r.Ext = nil
	}

	// Срезы строк — переиспользуем backing array
	r.Seats = append(r.Seats[:0], src.Seats...)
	r.BlockedSeats = append(r.BlockedSeats[:0], src.BlockedSeats...)
	r.Currencies = append(r.Currencies[:0], src.Currencies...)
	r.Languages = append(r.Languages[:0], src.Languages...)
	r.BlockedAdvDomains = append(r.BlockedAdvDomains[:0], src.BlockedAdvDomains...)
	r.BlockedApps = append(r.BlockedApps[:0], src.BlockedApps...)

	// BlockedCategories — переиспользуем backing array
	r.BlockedCategories = r.BlockedCategories[:0]
	r.BlockedCategories = append(r.BlockedCategories, src.BlockedCategories...)
}

// clearForPool обнуляет элементы срезов для освобождения внутренних указателей,
// сохраняет backing arrays срезов для переиспользования, обнуляет указатели и скалярные поля.
func (r *BidRequest) clearForPool() {

	// Обнуляем элементы импрешенов, чтобы освободить вложенные указатели (Banner, Video и т.д.)
	for i := range r.Impressions {
		r.Impressions[i] = Impression{}
	}
	r.Impressions = r.Impressions[:0]

	// Обнуляем указатели на вложенные объекты
	r.App = nil
	r.Site = nil
	r.Device = nil
	r.User = nil
	r.Source = nil
	r.Regulations = nil
	r.Ext = nil

	// Сбрасываем длину срезов строк, сохраняя backing arrays
	r.Seats = r.Seats[:0]
	r.BlockedSeats = r.BlockedSeats[:0]
	r.Currencies = r.Currencies[:0]
	r.Languages = r.Languages[:0]
	r.BlockedCategories = r.BlockedCategories[:0]
	r.BlockedAdvDomains = r.BlockedAdvDomains[:0]
	r.BlockedApps = r.BlockedApps[:0]

	// Обнуляем скалярные поля
	r.ID = ""
	r.Test = 0
	r.AuctionType = 0
	r.TimeMax = 0
	r.AllImpressions = 0
}