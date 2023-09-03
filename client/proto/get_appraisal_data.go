package proto

type AppraisalWithCharacter[A any] struct {
	Appraisal   *A
	CharacterId int32
}

func (awc AppraisalWithCharacter[A]) Unwrap() (
	appraisal *A,
	characterId int32,
) {
	return awc.Appraisal, awc.CharacterId
}
