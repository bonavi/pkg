run.env: ## start service with environments
	@godotenv -f ./environments/$(env) go run ./cmd/events-service/...

# Используем, чтобы прогнать код линтером (следим, чтобы локальный golangci был всегда последней версии)
lint:
	golangci-lint run -v

# Используем, чтобы обновить моки для тестов
mockery:
	find . -type f -name 'mock_*' -exec rm {} +
	mockery

# Используем для получения тест-кавераджа для каждого файла и проекта в целом
test-coverage-number: mockery
	go test -v -coverprofile=profile.cov ./cmd...
	go tool cover -func profile.cov
	rm profile.cov

# Используем для получения графического тест-кавераджа. Просматривая код в браузере мы можем видеть сколько тесткейсов
# Проходятся по одному и тому же участку кода
test-coverage-html: mockery
	go test -v -coverprofile=profile.cov ./cmd...
	go tool cover -html profile.cov
	rm profile.cov

# Используем для прогона тестов текущего сервиса
test: mockery
	go test -race ./...

# Перед деплоем дергаем эту команду, чтобы проверить, что код готов к сливанию с другой веткой
# Если этого не делать, завалится пайплайн и все равно придется править :)
deploy-check: test lint

# Чтобы обновить сериализатор json структуры лога
easyjson:
	easyjson log/jsonHandler.go
