package chain

type Applier[T any] interface {
	Apply(T) error
}

// Run последовательно вызывает Apply у каждого звена, прерываясь при первой ошибке
func Run[T any](steps []Applier[T], input T) error {
	for _, step := range steps {
		if err := step.Apply(input); err != nil {
			return err
		}
	}
	return nil
}
