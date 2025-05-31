package slices

import (
	"fmt"

	"pkg/errors"
)

// ToMap возвращает map, где ключом является поле структуры, а значением сама структура
// Example:
// AccountGroupsMap := slice.ToMap(_accountGroups, func(ag model.AccountGroup) uint32 { return ag.ID })
func ToMap[K comparable, V any](slice []V, field func(V) K) map[K]V {
	mapBySlice := make(map[K]V, len(slice))
	for _, v := range slice {
		mapBySlice[field(v)] = v
	}
	return mapBySlice
}

func ToMapSlices[K comparable, V any](slice []V, field func(V) K) map[K][]V {
	mapBySlice := make(map[K][]V, len(slice))
	for _, v := range slice {
		mapBySlice[field(v)] = append(mapBySlice[field(v)], v)
	}
	return mapBySlice
}

// GetFields возвращает массив значений полей из массива структур
// Example:
// AccountGroupsIDs := slice.GetFields(_accountGroups, func(ag model.AccountGroup) uint32 { return ag.ID })
func GetFields[K, V any](slice []V, field func(V) K) []K {
	fields := make([]K, 0, len(slice))
	for _, v := range slice {
		fields = append(fields, field(v))
	}
	return fields
}

// Map это синоним GetFields
func Map[K, V any](slice []V, field func(V) K) []K {
	return GetFields(slice, field)
}

// GetUniqueFields возвращает массив уникальных значений полей из массива структур
// Example:
// AccountGroupsIDs := slice.GetUniqueFields(_accountGroups, func(ag model.AccountGroup) uint32 { return ag.ID })
func GetUniqueFields[S any, V comparable](slice []S, field func(S) V) []V {

	fieldsMap := make(map[V]struct{}, len(slice))

	for _, s := range slice {
		fieldsMap[field(s)] = struct{}{}
	}

	fields := make([]V, 0, len(fieldsMap))

	for k := range fieldsMap {
		fields = append(fields, k)
	}

	return fields
}

// In проверяет, содержится ли ХОТЬ ОДНО значение в переданном массиве
// Example:
// if slice.In(1, 1, 2, 3) {
func In[K comparable](value K, slice ...K) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// GetMapValueStruct Возвращает проверочную мапу всех возможных значений поля
// Example:
// AccountGroupsIDs := slice.GetMapValueStruct(_accountGroups, func(ag model.AccountGroup) uint32 { return ag.ID })
func GetMapValueStruct[K comparable, V any](slice []V, field func(V) K) map[K]struct{} {
	m := make(map[K]struct{}, len(slice))
	for _, v := range slice {
		m[field(v)] = struct{}{}
	}
	return m
}

// JoinExclusive возвращает массивы, содержащие только уникальные элементы
// Example:
// leftObjectsExclusive, rightObjectsExclusive := slice.JoinExclusive(leftObjects, rightObjects)
func JoinExclusive[T comparable](leftObjects, rightObjects []T) (leftObjectsExclusive, rightObjectsExclusive []T) {
	leftObjectsMap := GetMapValueStruct(leftObjects, func(v T) T { return v })
	rightObjectsMap := GetMapValueStruct(rightObjects, func(v T) T { return v })

	for _, leftObject := range leftObjects {
		if _, ok := rightObjectsMap[leftObject]; !ok {
			leftObjectsExclusive = append(leftObjectsExclusive, leftObject)
		}
	}

	for _, rightObject := range rightObjects {
		if _, ok := leftObjectsMap[rightObject]; !ok {
			rightObjectsExclusive = append(rightObjectsExclusive, rightObject)
		}
	}

	return leftObjectsExclusive, rightObjectsExclusive
}

// First возвращает первый элемент массива
// Если массив пустой, возвращает nil
// Example:
// firstElement := slice.First([]int{1, 2, 3})
// Возвращает указатель на копию первого элемента массива!
func First[T any](array []T) (value T, err error) {
	// Если массив пустой
	if len(array) == 0 {
		// Возвращаем ошибку
		return value, errors.NotFound.Wrap(ErrSliceIsEmpty).WithParams("Type", fmt.Sprintf("%T", value))
	}
	// Возвращаем первый элемент массива
	return array[0], nil
}

var ErrSliceIsEmpty = errors.New("slice is empty")

// FirstWithError получает массив и ошибку (как правило в таком формате возвращают функции получения массива элементов)
// Если пришедшая ошибка не пустая, просто возвращаем ее
// Если пришедший массив пустой, возвращаем ошибку
// Если массив не пустой, то возвращаем первый элемент этого массива
// Example:
// firstElement, err := slice.FirstWithError(s.GetArray())
func FirstWithError[T any](array []T, initialErr error) (value T, err error) {
	// Если пришедшая ошибка не пустая
	if initialErr != nil {
		// Возвращаем ее
		return value, initialErr
	}
	// Получаем первый элемент массива или ошибку
	return First(array)
}

// Contains проверяет, содержится ли значение в массиве, используя мапу для быстрого поиска.
// Пример использования:
// isPresent := slice.Contains([]int{1, 2, 3}, 2)
func Contains[T comparable](slice []T, value T) bool {

	// Проходимся по каждому элементу
	for _, v := range slice {

		// Если элемент найден
		if v == value {
			return true
		}
	}

	return false
}

// ContainsAll проверяет, содержатся ли все значения target в массиве slice, используя мапу для быстрого поиска.
// Пример использования:
// isPresent := slice.Contains([]int{1, 2, 3}, 1, 2)
func ContainsAll[T comparable](slice []T, target ...T) bool {

	sliceMap := GetMapValueStruct(slice, func(v T) T { return v })

	// Проходимся по каждому элементу
	for _, v := range target {

		// Если элемент не найден
		if _, ok := sliceMap[v]; !ok {
			return false
		}
	}

	return true
}

func ContainsAny[T comparable](slice []T, target ...T) bool {

	sliceMap := GetMapValueStruct(slice, func(v T) T { return v })

	// Проходимся по каждому элементу
	for _, v := range target {

		// Если элемент найден
		if _, ok := sliceMap[v]; ok {
			return true
		}
	}

	return false
}

// Filter возвращает массив, содержащий только те элементы, которые удовлетворяют условию
func Filter[T any](slice []T, filter func(T) bool) (filtered []T) {
	for _, v := range slice {
		if filter(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}
