package bigutil

import (
	"database/sql/driver"
	"math/big"

	eth_hexutil "github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

const (
	defaultBase = 16

	maxByteLength = 32
)

type Int struct {
	x    big.Int
	base int
}

func NewInt(x *big.Int) Int {
	return Int{
		x:    *x,
		base: defaultBase,
	}
}

func (x Int) Base() int {
	return x.base
}

func (x *Int) SetBaseTo10() {
	x.base = 10
}

func (x *Int) SetBaseTo16() {
	x.base = 16
}

func (x Int) Bytes() []byte {
	return x.x.Bytes()
}

func (x Int) String() string {
	switch x.base {
	case 10:
		return x.x.String()
	case 16:
		return eth_hexutil.EncodeBig(&x.x)
	default:
		return eth_hexutil.EncodeBig(&x.x)
	}
}

func (x Int) BigInt() *big.Int {
	return &x.x
}

func (x Int) Value() (driver.Value, error) {
	return x.Bytes(), nil
}

func (x *Int) Scan(src any) error {
	if src == nil {
		return errors.New("src must not be nil")
	}

	b, ok := src.([]byte)
	if !ok {
		return errors.New("the type of src must be []byte")
	}
	if len(b) == 0 {
		return errors.New("the length of src must be greater than 0")
	}
	if len(b) > maxByteLength {
		return errors.Errorf("the length of src must be %d or less", maxByteLength)
	}

	x.x = *new(big.Int).SetBytes(b)
	x.base = defaultBase

	return nil
}

func (x Int) MarshalText() ([]byte, error) {
	switch x.base {
	case 10:
		return []byte(x.x.String()), nil
	case 16:
		return eth_hexutil.Big(x.x).MarshalText()
	default:
		return eth_hexutil.Big(x.x).MarshalText()
	}
}

func (x *Int) UnmarshalText(text []byte) error {
	if len(text) >= 2 && text[0] == '0' && text[1] == 'x' {
		x.SetBaseTo16()
		return (*eth_hexutil.Big)(&x.x).UnmarshalText(text)
	}

	x.SetBaseTo10()
	return x.x.UnmarshalText(text)
}
