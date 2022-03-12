package big

import (
	"encoding/json"
	"errors"
	"math/big"
	"strings"
)

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
