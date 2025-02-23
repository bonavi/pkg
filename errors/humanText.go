package errors

var humanTextByLevel = map[ErrorType]string{
	BadRequest:     "Введены неверные данные",
	InternalServer: "Произошла непредвиденная ошибка",
	NotFound:       "Данные не найдены",
	Forbidden:      "Доступ запрещен",
	Teapot:         "Разработчик забыл написать текст ошибки",
	BadGateway:     "Произошла ошибка на сервере внешнего сервиса",
	Unauthorized:   "Пользователь не авторизован",
	Timeout:        "Превышено время ожидания",
}

func (et ErrorType) HumanText() string {

	humanText, ok := humanTextByLevel[et]
	if !ok {
		return "Произошла непредвиденная ошибка"
	}

	return humanText
}
