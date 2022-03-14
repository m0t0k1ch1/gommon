package big

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

type StringInt struct {
	Int *big.Int
}

func NewStringInt(s string) (*StringInt, error) {
	si := new(StringInt)
	if err := si.SetString(s); err != nil {
		return nil, err
	}

	return si, nil
}

func (si *StringInt) SetString(s string) error {
	x, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return errors.New("failed to set string")
	}

	si.Int = x

	return nil
}

func (si *StringInt) String() string {
	return si.Int.Text(10)
}

func (si *StringInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(si.String())
}

func (si *StringInt) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	si.SetString(s)

	return nil
}

func (si *StringInt) Scan(v interface{}) error {
	if v == nil {
		return nil
	}

	b, ok := v.([]byte)
	if !ok {
		return fmt.Errorf("converting %T to StringInt is unsupported", v)
	}

	si.Int = new(big.Int).SetBytes(b)

	return nil
}

func (si *StringInt) Value() (driver.Value, error) {
	return si.Int.Bytes(), nil
}

type HexInt struct {
	Int *big.Int
}

func NewHexInt(s string) (*HexInt, error) {
	hi := new(HexInt)
	if err := hi.SetString(s); err != nil {
		return nil, err
	}

	return hi, nil
}

func (hi *HexInt) SetString(s string) error {
	if !strings.HasPrefix(s, "0x") {
		return errors.New("invalid hex string")
	}

	x, ok := new(big.Int).SetString(s, 0)
	if !ok {
		return errors.New("failed to set string")
	}

	hi.Int = x

	return nil
}

func (hi *HexInt) String() string {
	s := hi.Int.Text(16)
	if len(s) == 0 {
		return "0x0"
	}

	return "0x" + s
}

func (hi *HexInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(hi.String())
}

func (hi *HexInt) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	hi.SetString(s)

	return nil
}

func (hi *HexInt) Scan(v interface{}) error {
	if v == nil {
		return nil
	}

	b, ok := v.([]byte)
	if !ok {
		return fmt.Errorf("converting %T to HexInt is unsupported", v)
	}

	hi.Int = new(big.Int).SetBytes(b)

	return nil
}

func (hi *HexInt) Value() (driver.Value, error) {
	return hi.Int.Bytes(), nil
}
