package protoerr

import (
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/WiggidyW/etco-go/proto"
)

const (
	OK              proto.ErrorCode = proto.ErrorCode_EC_OK
	SERVER_ERR      proto.ErrorCode = proto.ErrorCode_EC_SERVER_ERROR
	INVALID_REQUEST proto.ErrorCode = proto.ErrorCode_EC_INVALID_REQUEST
	INVALID_MERGE   proto.ErrorCode = proto.ErrorCode_EC_INVALID_MERGE
	BOOTSTRAP_UNSET proto.ErrorCode = proto.ErrorCode_EC_BOOTSTRAP_UNSET
	NOT_FOUND       proto.ErrorCode = proto.ErrorCode_EC_NOT_FOUND
	TOKEN_INVALID   proto.ErrorCode = proto.ErrorCode_EC_TOKEN_INVALID
)

// protoerror.Error

type Error struct {
	Code proto.ErrorCode
	Err  error
}

func New(code proto.ErrorCode, err error) Error {
	return Error{Code: code, Err: err}
}

func MsgNew(code proto.ErrorCode, msg string) Error {
	return Error{Code: code, Err: errors.New(msg)}
}

func (e Error) Unwrap() error { return e.Err }

func (e Error) Error() string {
	if e.Err == nil {
		return e.Code.String()
	} else {
		return e.Code.String() + ": " + e.sanitizedInnerError()
	}
}

// panics if e.Err is nil
func (e Error) sanitizedInnerError() string {
	err := e.Err.Error()
	if utf8.ValidString(err) {
		return err
	}

	var b strings.Builder
	b.Grow(len(err))
	for _, c := range err {
		if c >= 0 && c <= 127 { // valid ASCII
			b.WriteRune(c)
		}
	}
	return b.String()
}

func (e Error) ToProto() *proto.ErrorResponse {
	if e.Err == nil {
		return &proto.ErrorResponse{Code: e.Code}
	} else {
		return &proto.ErrorResponse{
			Code:  e.Code,
			Error: e.sanitizedInnerError(),
		}
	}
}

// error (Go interface)

func ErrToProtoErr(e error) Error {
	var protoError Error
	if errors.As(e, &protoError) {
		return protoError
	} else {
		return Error{Code: SERVER_ERR, Err: e}
	}
}

func ErrToProto(e error) *proto.ErrorResponse {
	if e == nil {
		return nil
	} else {
		return ErrToProtoErr(e).ToProto()
	}
}
