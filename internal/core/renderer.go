package core

type Renderer interface {
	Render(name string, data any) []byte
}
