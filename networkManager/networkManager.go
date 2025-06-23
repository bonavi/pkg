package networkManager

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"pkg/errors"
)

type AuthManager interface {
	GetToken(ctx context.Context, endpointURL string) (string, error)
}

type SendRequestReq struct {
	ServerURL         string
	AuthHeaderName    string
	URL               string
	Method            string
	Headers           map[string]string
	WithAuthorization bool
	QueryParams       map[string]string
	Body              []byte
}

type SendRequestRes struct {
	Body       []byte
	StatusCode int
	Headers    map[string]string
	Cookies    []*http.Cookie
}

type NetworkManager struct {
	authManager AuthManager
	client      *http.Client
}

func NewNetworkManager(
	authManager AuthManager,
	client *http.Client,
) *NetworkManager {
	return &NetworkManager{
		authManager: authManager,
		client:      client,
	}
}

func (n *NetworkManager) SendRequest(ctx context.Context, req SendRequestReq) (res SendRequestRes, err error) {

	// Создаем HTTP-запрос на логин
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, req.URL, bytes.NewReader(req.Body))
	if err != nil {
		return res, err
	}

	var token string

	// Если требуется авторизация и есть менеджер авторизации
	if req.WithAuthorization && n.authManager != nil {

		// Получаем токен
		if token, err = n.authManager.GetToken(ctx, req.ServerURL); err != nil {
			return res, errors.Default.Wrap(err)
		}

		// Добавляем токен в заголовки запроса
		httpReq.Header.Add(req.AuthHeaderName, token)
	}

	for key, value := range req.Headers {
		httpReq.Header.Add(key, value)
	}

	// Отправляем запрос на сервер
	resp, err := n.client.Do(httpReq)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, errors.Default.Wrap(err)
	}

	switch resp.StatusCode {
	case http.StatusOK:

		// Получаем заголовки ответа
		headers := make(map[string]string)
		for key, values := range resp.Header {
			headers[key] = values[0]
		}

		// Возвращаем тело ответа
		return SendRequestRes{
			Body:       body,
			StatusCode: resp.StatusCode,
			Headers:    make(map[string]string),
			Cookies:    resp.Cookies(),
		}, nil

	default:
		return res, errors.Default.New("unexpected status code").WithParams(
			"statusCode", resp.StatusCode,
			"body", string(body),
		)
	}
}
