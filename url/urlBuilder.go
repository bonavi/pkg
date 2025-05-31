package url

import (
	"net/url"
	"slices"
	"strings"

	"pkg/errors"
	"pkg/maps"
)

func NewBuilder(
	host string,
	path *string,
	params map[string]string,
	pathParams map[string]string,
	isNeedUnescape bool,
) *Builder {
	return &Builder{
		Host:           host,
		Path:           path,
		Params:         params,
		PathParams:     pathParams,
		IsNeedUnescape: isNeedUnescape,
	}
}

type Builder struct {
	Host           string
	Path           *string
	Params         map[string]string
	PathParams     map[string]string
	IsNeedUnescape bool
}

func (b *Builder) Copy() *Builder {

	params := make(map[string]string, len(b.Params))
	for k, v := range b.Params {
		params[k] = v
	}

	pathParams := make(map[string]string, len(b.PathParams))
	for k, v := range b.PathParams {
		pathParams[k] = v
	}

	path := new(string)

	if b.Path != nil {
		*path = *b.Path
	}

	return &Builder{
		Host:           b.Host,
		Path:           path,
		Params:         params,
		PathParams:     pathParams,
		IsNeedUnescape: b.IsNeedUnescape,
	}
}

func (b *Builder) SetParam(kv ...string) *Builder {
	for i := 0; i < len(kv); i += 2 {
		if i+1 >= len(kv) {
			break
		}
		b.Params[kv[i]] = kv[i+1]
	}

	return b
}

func (b *Builder) GetURL() (string, error) {

	// Разбиваем хост на схему и хост
	u, err := url.Parse(b.Host)
	if err != nil {
		return "", errors.InternalServer.Wrap(err)
	}

	if b.Path != nil {

		path := *b.Path

		// Заменяем теги в пути на значения из pathParams
		for key, value := range b.PathParams {
			path = strings.ReplaceAll(path, ":"+key, value)
		}

		u.Path = path
	}

	// Сортируем query параметры, чтобы всегда была одинаковая ссылка
	keys := maps.Keys(b.Params)
	slices.Sort(keys)

	// Создаем объект URL с параметрами
	q := u.Query()
	for _, key := range keys {
		q.Set(key, b.Params[key])
	}

	// Устанавливаем параметры в URL
	if b.IsNeedUnescape {

		// Установка без экранирования
		u.RawQuery, err = url.QueryUnescape(q.Encode())
		if err != nil {
			return "", errors.InternalServer.Wrap(err)
		}
	} else {

		// Установка с экранированием
		u.RawQuery = q.Encode()
	}

	if u.Host == "" {
		return "", errors.InternalServer.New("Host is empty")
	}

	if u.Scheme == "" {
		return "", errors.InternalServer.New("Scheme is empty")
	}

	if u.Path == "" {
		return "", errors.InternalServer.New("Path is empty")
	}

	return u.String(), nil
}
