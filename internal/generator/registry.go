package generator

type ModelRegistry struct {
	models map[string]any
}

func NewModelRegistry() *ModelRegistry {
	return &ModelRegistry{
		models: make(map[string]any),
	}
}

func (r *ModelRegistry) Register(name string, model any) {
	r.models[name] = model
}

func (r *ModelRegistry) Get(name string) (any, bool) {
	m, ok := r.models[name]
	return m, ok
}
