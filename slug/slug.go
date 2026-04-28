package slug

import (
	"regexp"
	"strings"

	"pkg/errors"
)

const (
	MaxSlugLength = 40 // Максимальная длина slug
	MinSlugLength = 3  // Минимальная длина slug
)

var (
	// Разрешены только буквы, цифры, дефис и подчеркивание
	SlugRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9-_]*[a-z0-9]$`)
)

// Normalize нормализует slug
func Normalize(slug string) string {
	// Убираем пробелы в начале и конце
	slug = strings.TrimSpace(slug)

	// Приводим к нижнему регистру
	slug = strings.ToLower(slug)

	// Заменяем все спецсимволы на тире
	slug = regexp.MustCompile(`[^a-z0-9\s-]`).ReplaceAllString(slug, "-")

	// Заменяем множественные пробелы на тире
	slug = regexp.MustCompile(`\s+`).ReplaceAllString(slug, "-")

	// Убираем множественные тире
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")

	// Убираем тире в начале и конце
	slug = strings.Trim(slug, "-")

	return slug
}

// Validate проверяет корректность slug
func Validate(slug string) error {
	// Нормализуем slug
	slug = Normalize(slug)

	// Проверка на пустое значение
	if slug == "" {
		return errors.Default.New("slug не может быть пустым")
	}

	// Проверка длины
	if len(slug) > MaxSlugLength {
		return errors.Default.New("длина slug превышает максимально допустимую").
			WithParams("max_length", MaxSlugLength)
	}

	if len(slug) < MinSlugLength {
		return errors.Default.New("длина slug меньше минимально допустимой").
			WithParams("min_length", MinSlugLength)
	}

	// Проверка формата с помощью регулярного выражения
	if !SlugRegex.MatchString(slug) {
		return errors.Default.New("неверный формат slug: разрешены только строчные буквы, цифры, дефисы и подчеркивания, должен начинаться и заканчиваться буквой или цифрой")
	}

	return nil
}
