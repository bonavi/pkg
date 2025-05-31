package ord

import "pkg/errors"

type ORDContract struct {

	// Идентификатор договора в ОРД, указанной в поле ado_id
	ID string `json:"id" validate:"required"`

	// Регистронезависимый ключ-идентификатор ОРД из списка “Список имен ОРД”
	AdoID ORDName `json:"ado_id" validate:"required"`

	// Единый идентификатор договора, выданный ЕРИР
	UnifiedID *string `json:"unified_id" validate:"required"`

	// Тип договора из списка “Список типов договоров”
	Type ContractType `json:"type" validate:"required"`

	// Номер договора
	Number string `json:"number" validate:"required"`

	// Предмет договора из списка “Список предметов договора”. Не обязательное для доп. соглашений
	SubjectType ContractSubjectType `json:"subject_type" validate:"required"`

	// Осуществляемые посредником действия из списка “Список осуществляемых посредником действий”. Обязательное только для посреднических договоров
	ActionType ContractActionType `json:"action_type"`

	// Дата подписания договора в формате YYYY-MM-DD
	SignDate string `json:"sign_date" validate:"required"`

	// Флаг, выставляемый в true только для договоров найма агентства для поиска клиентов паблишеру. Подробнее в описании ЕРИР. Указывается только в посреднических договорах
	AgentActingForPublisher bool `json:"agent_acting_for_publisher"`

	// Включен ли НДС в стоимость, указанную в договоре. С версии 2.2 игнорируется и принимается за true во всех случаях согласно указаниям ЕРИР
	VatIncluded bool `json:"vat_included"`

	// Сумма договора, дробное число в строке. Не может быть равна 0 для договоров с типом “Посреднический” и доп. соглашений к посредническим договорам.
	Amount string `json:"amount"`

	// Идентификатор родительского договора из списка nroa.parent_contracts. Обязателен только для доп. соглашений
	ParentContractID *string `json:"parent_contract_id"`
}

func (o *ORDContract) Validate() error {

	if o.ID == "" {
		return errors.BadRequest.New("ID is required")
	}

	if o.AdoID == "" {
		return errors.BadRequest.New("AdoID is required")
	}

	if o.UnifiedID == nil {
		return errors.BadRequest.New("UnifiedID is required")
	}

	if o.Type == "" {
		return errors.BadRequest.New("Type is required")
	}

	if o.Number == "" {
		return errors.BadRequest.New("Number is required")
	}

	if o.SubjectType == "" {
		return errors.BadRequest.New("SubjectType is required")
	}

	if o.SignDate == "" {
		return errors.BadRequest.New("SignDate is required")
	}

	if o.Amount == "" {
		return errors.BadRequest.New("Amount is required")
	}

	return nil
}

func (o *ORDContract) Copy() ORDContract {
	return ORDContract{
		ID:                      o.ID,
		AdoID:                   o.AdoID,
		UnifiedID:               o.UnifiedID,
		Type:                    o.Type,
		Number:                  o.Number,
		SubjectType:             o.SubjectType,
		ActionType:              o.ActionType,
		SignDate:                o.SignDate,
		AgentActingForPublisher: o.AgentActingForPublisher,
		VatIncluded:             o.VatIncluded,
		Amount:                  o.Amount,
		ParentContractID:        o.ParentContractID,
	}
}
