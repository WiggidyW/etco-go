package util

type Unwrappable2[T1 any, T2 any] interface {
	Unwrap() (T1, T2)
}

type Unwrappable3[T1 any, T2 any, T3 any] interface {
	Unwrap() (T1, T2, T3)
}

func Unwrap2WithErr[
	U Unwrappable2[
		T1,
		T2,
	],
	T1 any,
	T2 any,
](
	u U,
	errIn error,
) (
	t1 T1,
	t2 T2,
	errOut error,
) {
	if errIn != nil {
		return t1, t2, errIn
	} else {
		t1, t2 = u.Unwrap()
		return t1, t2, nil
	}
}

func Unwrap3WithErr[
	U Unwrappable3[
		T1,
		T2,
		T3,
	],
	T1 any,
	T2 any,
	T3 any,
](
	u U,
	errIn error,
) (
	t1 T1,
	t2 T2,
	t3 T3,
	errOut error,
) {
	if errIn != nil {
		return t1, t2, t3, errIn
	} else {
		t1, t2, t3 = u.Unwrap()
		return t1, t2, t3, nil
	}
}
