package annotation

type Annotation interface {
	isAnnotation()
}

// TODO: Generics
type Class struct {
	Name string
}

func (c *Class) isAnnotation() {}
