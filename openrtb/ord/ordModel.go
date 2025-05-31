package ord

import (
	"pkg/errors"
	"pkg/pointer"
)

// Nroa - Информация о каждом креативе для отчетности в ОРД
//
// https://docs.google.com/document/d/1u-1JT9xMZO5-vUdtg6A8wgfa1i0XuMxfNTchFSGVDYk/edit?tab=t.0
type Nroa struct {

	// Признак не подлежащего регистрации креатива.  Например, исследование в виде BLS
	//
	// Используется по договорённости между DSP и SSP
	//
	// Если равен true, остальные поля расширения не заполняются; SSP не учитывает ответ с креативом в отчётности перед ОРД.
	//
	// Значение по умолчанию: false
	Unreportable *bool `json:"unreportable"`

	// Идентификатор креатива в системе ОРД
	Erid string `json:"erid"`

	// Значение 1 означает, что в креатив уже вшита маркировка о рекламе; маркировка со стороны SSP не требуется
	//
	// По умолчанию значение равно 0
	HasNroaMarkup int `json:"has_nroa_markup"` // Признак наличия маркапа в креативе

	// Информация об исполнителе по изначальному договору
	Contractor *ORDClient `json:"contractor"`

	// Информация о заказчике по изначальному договору
	Client *ORDClient `json:"client"`

	// Основная информация изначального договора
	InitialContract *ORDContract `json:"initial_contract"`

	// Список “родительских” договоров
	//
	// Заполняется только если initial_contract имеет тип “дополнительное соглашение”
	//
	// Содержит информацию о договоре (или договорах, если это цепочка), по отношению к которому initial_contract является ДС.
	ParentContracts []ORDContract `json:"parent_contracts"`
}

func (n *Nroa) Copy() Nroa {

	var parentContracts []ORDContract
	if len(n.ParentContracts) != 0 {
		parentContracts = make([]ORDContract, len(n.ParentContracts))
		for i := range n.ParentContracts {
			parentContracts[i] = n.ParentContracts[i].Copy()
		}
	}

	var contractor *ORDClient
	if n.Contractor != nil {
		contractor = pointer.Pointer(n.Contractor.Copy())
	}

	var client *ORDClient
	if n.Client != nil {
		client = pointer.Pointer(n.Client.Copy())
	}

	var initialContract *ORDContract
	if n.InitialContract != nil {
		initialContract = pointer.Pointer(n.InitialContract.Copy())
	}

	return Nroa{
		Unreportable:    n.Unreportable,
		Erid:            n.Erid,
		HasNroaMarkup:   n.HasNroaMarkup,
		Contractor:      contractor,
		Client:          client,
		InitialContract: initialContract,
		ParentContracts: parentContracts,
	}
}

func (n *Nroa) Validate() error {

	if n.Erid == "" {
		return errors.BadRequest.New("Erid is required")
	}

	if n.Contractor == nil {
		return errors.BadRequest.New("Contractor is required")
	}

	if err := n.Contractor.Validate(); err != nil {
		return err
	}

	if n.Client == nil {
		return errors.BadRequest.New("Client is required")
	}

	if err := n.Client.Validate(); err != nil {
		return err
	}

	if n.InitialContract == nil {
		return errors.BadRequest.New("InitialContract is required")
	}

	if err := n.InitialContract.Validate(); err != nil {
		return err
	}

	return nil
}
