package utils

import (
	"bytes"
	"compress/gzip"
	"io"
	"sync"
)

const (
	MinCompressionSize = 1024
)

var (
	gZipBufferPool = sync.Pool{ // выдаёт готовый *bytes.Buffer. Вместо make([]byte, …)
		New: func() any {
			return new(bytes.Buffer)
		}}
	gZipPool = sync.Pool{ // выдаёт готовый *gzip.Writer
		New: func() any {
			w, _ := gzip.NewWriterLevel(io.Discard, 4) // 4 – среднее значение
			return w
		},
	}
)

// CompressGzip сжимает данные, если пользы мало вернёт исходник
// level==0 - используем дефолт значение 4
func CompressGzip(data []byte, level int) ([]byte, error) {
	if len(data) < MinCompressionSize { // тут сомнения нужно ли делать проверку, может оверхед
		return data, nil
	}
	if level == 0 { // тут сомнения нужны ли нам разные уровни компрессии, пока как "заглушка" для констант
		level = 4
	}

	// берём буфер из пула, очищаем, возвращаем в пул
	buffer := gZipBufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer gZipBufferPool.Put(buffer)

	//	берём *gzip.Writer из пула и привязываем его к буферу через Reset
	gZipWriter := gZipPool.Get().(*gzip.Writer)
	defer gZipPool.Put(gZipWriter)

	gZipWriter.Reset(buffer)

	if _, err := gZipWriter.Write(data); err != nil {
		_ = gZipWriter.Close()
		return nil, err
	}

	if err := gZipWriter.Close(); err != nil {
		return nil, err
	}

	// если выгоды почти нет — вернём исходник, думаю нужно ли перегружать цп у партнеров
	if buffer.Len() >= int(float64(len(data))*0.98) {
		return data, nil
	}

	out := make([]byte, buffer.Len())
	copy(out, buffer.Bytes())
	return out, nil
}
