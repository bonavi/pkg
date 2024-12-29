package url

import (
	"net/url"
	"slices"
	"strings"

	"pkg/errors"
	"pkg/maps"
)

func BuildURL(
	host string,
	path string,
	params map[string]string,
	pathParams map[string]string,
	isNeedUnescape bool,
) (string, error) {

	// Разбиваем хост на схему и хост
	u, err := url.Parse(host)
	if err != nil {
		return "", errors.InternalServer.Wrap(err)
	}

	// Заменяем теги в пути на значения из pathParams
	for key, value := range pathParams {
		path = strings.ReplaceAll(path, ":"+key, value)
	}

	u.Path = path

	// Сортируем query параметры, чтобы всегда была одинаковая ссылка
	keys := maps.Keys(params)
	slices.Sort(keys)

	// Создаем объект URL с параметрами
	q := u.Query()
	for _, key := range keys {
		q.Set(key, params[key])
	}

	// Устанавливаем параметры в URL
	if isNeedUnescape {

		// Установка без экранирования
		u.RawQuery, err = url.QueryUnescape(q.Encode())
		if err != nil {
			return "", errors.InternalServer.Wrap(err)
		}
	} else {

		// Установка с экранированием
		u.RawQuery = q.Encode()
	}

	return u.String(), nil
}
