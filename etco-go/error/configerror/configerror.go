package configerror

import "fmt"

type ErrInvalid struct{ Err error }

func (e ErrInvalid) Unwrap() error { return e.Err }
func (e ErrInvalid) Error() string {
	return e.Err.Error()
}

type ErrMarketInvalid struct {
	Market    string
	ErrString string
}

func (e ErrMarketInvalid) Unwrap() error { return nil }
func (e ErrMarketInvalid) Error() string {
	return fmt.Sprintf(
		"'%s': %s",
		e.Market,
		e.ErrString,
	)
}

type ErrPricingInvalid struct{ ErrString string }

func (e ErrPricingInvalid) Unwrap() error { return nil }
func (e ErrPricingInvalid) Error() string { return e.ErrString }

type ErrBuybackSystemInvalid struct{ Err error }

func (e ErrBuybackSystemInvalid) Unwrap() error { return e.Err }
func (e ErrBuybackSystemInvalid) Error() string { return e.Err.Error() }

type ErrShopLocationInvalid struct{ Err error }

func (e ErrShopLocationInvalid) Unwrap() error { return e.Err }
func (e ErrShopLocationInvalid) Error() string { return e.Err.Error() }

type ErrShopTypeInvalid struct{ Err error }

func (e ErrShopTypeInvalid) Unwrap() error { return e.Err }
func (e ErrShopTypeInvalid) Error() string { return e.Err.Error() }

type ErrBuybackTypeInvalid struct{ Err error }

func (e ErrBuybackTypeInvalid) Unwrap() error { return e.Err }
func (e ErrBuybackTypeInvalid) Error() string { return e.Err.Error() }
