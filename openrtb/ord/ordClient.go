package ord

import "pkg/errors"

type ORDClient struct {

	// ИНН юр. или физ. лица
	INN string `json:"inn"`

	// ОПФ и юридическое наименование
	Name string `json:"name"`

	// Тип организации из списка “Список типов организации”
	LegalForm OrganizationType `json:"legal_form"`

	// Номер телефона. Заполняется для иностранных физ- и юрлиц в соответствии с требованиями ЕРИР
	MobilePhone string `json:"mobile_phone"`

	// Номер электронного средства платежа. Заполняется для иностранных физ- и юрлиц в соответствии с требованиями ЕРИР
	EpayNumber string `json:"epay_number"`

	// Регистрационный номер, либо его аналог. Заполняется для иностранных физ- и юрлиц в соответствии с требованиями ЕРИР
	RegNumber string `json:"reg_number"`

	// Номер налогоплательщика либо его аналог в стране регистрации. Заполняется для иностранных физ- и юрлиц в соответствии с требованиями ЕРИР
	AlternativeINN string `json:"alternative_inn"`

	// Код страны регистрации юрлица в соответствии с ОКСМ. Заполняется для иностранных физ- и юрлиц в соответствии с требованиями ЕРИР
	OKSMNumber string `json:"oksm_number"`
}

func (o *ORDClient) Validate() error {

	if o.INN == "" {
		return errors.BadRequest.New("INN is required")
	}

	if o.Name == "" {
		return errors.BadRequest.New("Name is required")
	}

	if o.LegalForm == "" {
		return errors.BadRequest.New("LegalForm is required")
	}

	return nil
}

func (o *ORDClient) Copy() ORDClient {
	return ORDClient{
		INN:            o.INN,
		Name:           o.Name,
		LegalForm:      o.LegalForm,
		MobilePhone:    o.MobilePhone,
		EpayNumber:     o.EpayNumber,
		RegNumber:      o.RegNumber,
		AlternativeINN: o.AlternativeINN,
		OKSMNumber:     o.OKSMNumber,
	}
}
