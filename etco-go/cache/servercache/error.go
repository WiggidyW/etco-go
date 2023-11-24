package servercache

type ErrServerUnlock struct {
	err error
}

func (e ErrServerUnlock) Unwrap() error { return e.err }
func (e ErrServerUnlock) Error() string {
	return "ErrServerUnlock: " + e.err.Error()
}

type ErrServerObtainLock struct {
	err error
}

func (e ErrServerObtainLock) Unwrap() error { return e.err }
func (e ErrServerObtainLock) Error() string {
	return "ErrServerLock: " + e.err.Error()
}

type ErrServerRefreshLock struct {
	err error
}

func (e ErrServerRefreshLock) Unwrap() error { return e.err }
func (e ErrServerRefreshLock) Error() string {
	return "ErrServerRefreshLock: " + e.err.Error()
}

type ErrServerGet struct {
	err error
}

func (e ErrServerGet) Unwrap() error { return e.err }
func (e ErrServerGet) Error() string {
	return "ErrServerGet: " + e.err.Error()
}

type ErrServerSet struct {
	err error
}

func (e ErrServerSet) Unwrap() error { return e.err }
func (e ErrServerSet) Error() string {
	return "ErrServerSet: " + e.err.Error()
}

type ErrServerDel struct {
	err error
}

func (e ErrServerDel) Unwrap() error { return e.err }
func (e ErrServerDel) Error() string {
	return "ErrServerDel: " + e.err.Error()
}
