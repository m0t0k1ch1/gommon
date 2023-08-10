package bigutil

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/m0t0k1ch1/gommon/internal/testutil"
)

var (
	maxUint256 = new(big.Int).SetBytes([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
)

func TestScan(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			Name         string
			Input        any
			ErrorMessage string
		}{
			{
				"nil",
				nil,
				"src must not be nil",
			},
			{
				"int",
				0,
				"the type of src must be []byte",
			},
			{
				"string",
				"string",
				"the type of src must be []byte",
			},
			{
				"empty []byte",
				[]byte(nil),
				"the length of src must be greater than 0",
			},
			{
				"overlength []byte",
				new(big.Int).Add(maxUint256, big.NewInt(1)).Bytes(),
				"the length of src must be 32 or less",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.Name, func(t *testing.T) {
				var x Int
				err := x.Scan(tc.Input)
				if err == nil {
					t.Errorf("err must not be nil")
					return
				}
				testutil.Equal(t, tc.ErrorMessage, err.Error())
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			Name   string
			Input  []byte
			Output string
		}{
			{
				"0x0",
				[]byte{0},
				"0x0",
			},
			{
				"0x" + maxUint256.Text(16),
				maxUint256.Bytes(),
				"0x" + maxUint256.Text(16),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.Name, func(t *testing.T) {
				var x Int
				if err := x.Scan(tc.Input); err != nil {
					t.Fatal(err)
				}
				testutil.Equal(t, tc.Output, x.String())
			})
		}
	})
}

func TestMarshalJSON(t *testing.T) {
	t.Run("success", func(*testing.T) {
		tcs := []struct {
			Name     string
			Input    Int
			Output10 string
			Output16 string
		}{
			{
				"0x0",
				NewInt(big.NewInt(0)),
				`"0"`,
				`"0x0"`,
			},
			{
				"0x" + maxUint256.Text(16),
				NewInt(maxUint256),
				`"` + maxUint256.Text(10) + `"`,
				`"0x` + maxUint256.Text(16) + `"`,
			},
		}

		for _, tc := range tcs {
			t.Run(tc.Name, func(t *testing.T) {
				tc.Input.SetBaseTo10()
				{
					b, err := json.Marshal(tc.Input)
					if err != nil {
						t.Fatal(err)
					}

					testutil.Equal(t, tc.Output10, string(b))
				}

				tc.Input.SetBaseTo16()
				{
					b, err := json.Marshal(tc.Input)
					if err != nil {
						t.Fatal(err)
					}

					testutil.Equal(t, tc.Output16, string(b))
				}
			})
		}
	})
}

func TestUnmarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			Name  string
			Input []byte
		}{
			{
				`"0"`,
				[]byte(`"0"`),
			},
			{
				`"0x0"`,
				[]byte(`"0x0"`),
			},
			{
				`"` + maxUint256.Text(10) + `"`,
				[]byte(`"` + maxUint256.Text(10) + `"`),
			},
			{
				`"0x` + maxUint256.Text(16) + `"`,
				[]byte(`"0x` + maxUint256.Text(16) + `"`),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.Name, func(t *testing.T) {
				var x Int
				if err := json.Unmarshal(tc.Input, &x); err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, string(tc.Input), `"`+x.String()+`"`)
			})
		}
	})
}
