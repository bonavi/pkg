package chain

type link[T any] interface {
	Run(T) error
	setNext(link[T])
	Apply(T) error
}

type Link[T any] struct {
	next link[T]
}

func (r *Link[T]) Run(i T) error {
	if r.next == nil {
		return nil
	}
	// log.Debug(fmt.Sprintf("%T", r.next)) // Для проверки вызываемого звена

	if err := r.next.Apply(i); err != nil {
		return err
	}

	return r.next.Run(i)
}

func (r *Link[T]) setNext(next link[T]) { //nolint:unused
	r.next = next
}

func SetArrange[T any](links ...link[T]) *Link[T] {

	// Делаем стартовый элемент
	initialLink := &Link[T]{
		next: nil,
	}

	if len(links) == 0 {
		return initialLink
	}

	// Добавляем первое звено как следующий элемент стартового
	initialLink.next = links[0]

	// Проходимся по каждому звену с первого индекса
	for i := 1; i < len(links); i++ {

		// К предыдущему звену добавляем текущее
		links[i-1].setNext(links[i])
	}

	return initialLink
}
