package maps

import "pkg/errors"

func Values[K comparable, V any](mapa map[K]V) []V {
	slice := make([]V, 0, len(mapa))
	for _, v := range mapa {
		slice = append(slice, v)
	}
	return slice
}

func Keys[K comparable, V any](mapa map[K]V) []K {
	slice := make([]K, 0, len(mapa))
	for k := range mapa {
		slice = append(slice, k)
	}
	return slice
}

func Join[K comparable, V any](leftMap, rightMap map[K]V) map[K]V {

	joinMap := make(map[K]V, len(leftMap)+len(rightMap))

	for k, v := range leftMap {
		joinMap[k] = v
	}

	for k, v := range rightMap {
		joinMap[k] = v
	}

	return joinMap
}

func Revert[K comparable, V comparable](mapa map[K]V) (map[V]K, error) {
	revertMap := make(map[V]K, len(mapa))

	for k, v := range mapa {
		if _, ok := revertMap[v]; ok {
			return nil, errors.BadRequest.New("value is not unique")
		}
		revertMap[v] = k
	}

	return revertMap, nil
}
