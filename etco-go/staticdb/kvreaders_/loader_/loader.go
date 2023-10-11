package loader_

type Loader[R any] interface {
	Load() (R, error)
}
