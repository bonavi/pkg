package ord

type ORDName string

const (
	ORDNameAmberdata ORDName = "amberdata"
	ORDNameMTS       ORDName = "mts"
	ORDNameVK        ORDName = "vk"
	ORDNameOZON      ORDName = "ozon"
	ORDNameBeeline   ORDName = "beeline"
	ORDNameSber      ORDName = "sber"
	ORDNameYandex    ORDName = "yandex"
)

type OrganizationType string

const (
	OrganizationTypeLegalEntity        OrganizationType = "ul"  // Юридическое лицо РФ
	OrganizationTypePhysical           OrganizationType = "fl"  // Физическое лицо РФ
	OrganizationTypeIndividual         OrganizationType = "ip"  // Индивидуальный предприниматель
	OrganizationTypeForeignLegalEntity OrganizationType = "ful" // Иностранное юридическое лицо
	OrganizationTypeForeignPhysical    OrganizationType = "ffl" // Иностранное физическое лицо
)

type ContractType string

const (
	ContractTypeContract             ContractType = "contract"              // Договор оказания услуг
	ContractTypeIntermediaryContract ContractType = "intermediary-contract" // Посреднический договор
	ContractTypeAdditionalAgreement  ContractType = "additional-agreement"  // Дополнительное соглашение
)

type ContractSubjectType string

const (
	ContractSubjectTypeDistribution    ContractSubjectType = "distribution"     // Договор на распространение рекламы
	ContractSubjectTypeRepresentation  ContractSubjectType = "representation"   // Представительство
	ContractSubjectTypeMediation       ContractSubjectType = "mediation"        // Посредничество
	ContractSubjectTypeOrgDistribution ContractSubjectType = "org-distribution" // Договор на организацию распространения рекламы
	ContractSubjectTypeOther           ContractSubjectType = "other"            // Другое
)

type ContractActionType string

const (
	ContractActionTypeDistribution ContractActionType = "distribution" // Действия в целях распространения рекламы
	ContractActionTypeConclude     ContractActionType = "conclude"     // Заключение договоров
	ContractActionTypeCommercial   ContractActionType = "commercial"   // Коммерческое представительство
	ContractActionTypeOther        ContractActionType = "other"        // Другое
)
