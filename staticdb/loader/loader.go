package loader

type Loader[R any] interface {
	Load() (R, error)
}
