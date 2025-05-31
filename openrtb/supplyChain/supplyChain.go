package supplyChain

type SupplyChain struct {
	// Версия цепочки
	Version string `json:"ver"`
	// Признак завершенности цепочки
	Complete int8 `json:"complete"`
	// Узлы цепочки
	Nodes []Node `json:"nodes"`
}

type Node struct {
	// Идентификатор источника (SSP, Exchange)
	AccountSourceID string `json:"asi"`
	// Идентификатор паблишера
	SellerID string `json:"sid"`
	// Идентификатор запроса
	RequestID string `json:"rid"`
	// Есть ли получатель платежа
	// 1 - означает, что узел участвует в потоке платежей и будет передавать деньги следующему узлу
	// 0 - означает, что узел не участвует в финансовых операциях и не имеет получателя платежа
	HasPayee int8 `json:"hp"`
}
