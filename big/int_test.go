package big

import (
	"encoding/json"
	"testing"

	"github.com/m0t0k1ch1/gommon/internal/testutils"
)

func TestMarshalStringInt(t *testing.T) {
	cases := []struct {
		s    string
		want string
	}{{
		s:    "0",
		want: `{"id":"0"}`,
	}, {
		s:    "2083236893",
		want: `{"id":"2083236893"}`,
	}}

	for _, c := range cases {
		t.Run(c.s, func(t *testing.T) {
			si, err := NewStringInt(c.s)
			if err != nil {
				t.Fatal(err)
			}

			obj := struct {
				ID *StringInt `json:"id"`
			}{
				ID: si,
			}

			b, err := json.Marshal(obj)
			if err != nil {
				t.Fatal(err)
			}

			testutils.Equal(t, c.want, string(b))
		})
	}
}

func TestUnmarshalStringInt(t *testing.T) {
	cases := []struct {
		s    string
		want string
	}{{
		s:    `{"id":"0"}`,
		want: "0",
	}, {
		s:    `{"id":"2083236893"}`,
		want: "2083236893",
	}, {
		s:    `{"id":"0x0"}`,
		want: "0",
	}, {
		s:    `{"id":"0x7c2bac1d"}`,
		want: "2083236893",
	}}

	for _, c := range cases {
		t.Run(c.want, func(t *testing.T) {
			var obj struct {
				ID *StringInt `json:"id"`
			}
			if err := json.Unmarshal([]byte(c.s), &obj); err != nil {
				t.Fatal(err)
			}

			testutils.Equal(t, c.want, obj.ID.String())
		})
	}
}

func TestScanStringInt(t *testing.T) {
	cases := []struct {
		b    []byte
		want string
	}{{
		b:    []byte{},
		want: "0",
	}, {
		b:    []byte{124, 43, 172, 29},
		want: "2083236893",
	}}

	for _, c := range cases {
		t.Run(c.want, func(t *testing.T) {
			var si StringInt
			if err := si.Scan(c.b); err != nil {
				t.Fatal(err)
			}

			testutils.Equal(t, c.want, si.String())
		})
	}
}

func TestStringIntValue(t *testing.T) {
	cases := []struct {
		s    string
		want []byte
	}{{
		s:    "0",
		want: []byte{},
	}, {
		s:    "2083236893",
		want: []byte{124, 43, 172, 29},
	}}

	for _, c := range cases {
		t.Run(c.s, func(t *testing.T) {
			si, err := NewStringInt(c.s)
			if err != nil {
				t.Fatal(err)
			}

			v, err := si.Value()
			if err != nil {
				t.Fatal(err)
			}

			testutils.Equal(t, c.want, v)
		})
	}
}

func TestMarshalHexInt(t *testing.T) {
	cases := []struct {
		s    string
		want string
	}{{
		s:    "0x0",
		want: `{"id":"0x0"}`,
	}, {
		s:    "0x4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b",
		want: `{"id":"0x4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b"}`,
	}}

	for _, c := range cases {
		t.Run(c.s, func(t *testing.T) {
			hi, err := NewHexInt(c.s)
			if err != nil {
				t.Fatal(err)
			}

			obj := struct {
				ID *HexInt `json:"id"`
			}{
				ID: hi,
			}

			b, err := json.Marshal(obj)
			if err != nil {
				t.Fatal(err)
			}

			testutils.Equal(t, c.want, string(b))
		})
	}
}

func TestUnmarshalHexInt(t *testing.T) {
	cases := []struct {
		s    string
		want string
	}{{
		s:    `{"id":"0x0"}`,
		want: "0x0",
	}, {
		s:    `{"id":"0x4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b"}`,
		want: "0x4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b",
	}}

	for _, c := range cases {
		t.Run(c.want, func(t *testing.T) {
			var obj struct {
				ID *HexInt `json:"id"`
			}
			if err := json.Unmarshal([]byte(c.s), &obj); err != nil {
				t.Fatal(err)
			}

			testutils.Equal(t, c.want, obj.ID.String())
		})
	}
}

func TestScanHexInt(t *testing.T) {
	cases := []struct {
		b    []byte
		want string
	}{{
		b:    []byte{},
		want: "0x0",
	}, {
		b:    []byte{74, 94, 30, 75, 170, 184, 159, 58, 50, 81, 138, 136, 195, 27, 200, 127, 97, 143, 118, 103, 62, 44, 199, 122, 178, 18, 123, 122, 253, 237, 163, 59},
		want: "0x4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b",
	}}

	for _, c := range cases {
		t.Run(c.want, func(t *testing.T) {
			var hi HexInt
			if err := hi.Scan(c.b); err != nil {
				t.Fatal(err)
			}

			testutils.Equal(t, c.want, hi.String())
		})
	}
}

func TestHexIntValue(t *testing.T) {
	cases := []struct {
		s    string
		want []byte
	}{{
		s:    "0x0",
		want: []byte{},
	}, {
		s:    "0x4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b",
		want: []byte{74, 94, 30, 75, 170, 184, 159, 58, 50, 81, 138, 136, 195, 27, 200, 127, 97, 143, 118, 103, 62, 44, 199, 122, 178, 18, 123, 122, 253, 237, 163, 59},
	}}

	for _, c := range cases {
		t.Run(c.s, func(t *testing.T) {
			hi, err := NewHexInt(c.s)
			if err != nil {
				t.Fatal(err)
			}

			v, err := hi.Value()
			if err != nil {
				t.Fatal(err)
			}

			testutils.Equal(t, c.want, v)
		})
	}
}
