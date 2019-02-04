package curve25519_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"testing"

	fixed "github.com/mmcloughlin/bug25519/fixed"
	"golang.org/x/crypto/curve25519"
)

type Implementation func(*[32]byte, *[32]byte, *[32]byte)

type TestVector struct {
	In     [32]byte
	Base   [32]byte
	Expect [32]byte
}

func TestTestVectorsCurrent(t *testing.T) {
	CheckTestVectors(t, curve25519.ScalarMult)
}

func TestTestVectorsFixed(t *testing.T) {
	CheckTestVectors(t, fixed.ScalarMult)
}

func CheckTestVectors(t *testing.T, mul Implementation) {
	tvs := LoadTestVectors(t, "testdata/testvectors.json")
	failed := 0
	for _, tv := range tvs {
		var got, in, base [32]byte
		copy(in[:], tv.In[:])
		copy(base[:], tv.Base[:])
		mul(&got, &in, &base)
		if !bytes.Equal(got[:], tv.Expect[:]) {
			t.Logf("    in = %x", tv.In)
			t.Logf("  base = %x", tv.Base)
			t.Logf("   got = %x", got)
			t.Logf("expect = %x", tv.Expect)
			t.Fail()
			failed++
		}
	}
	t.Logf("failed %d of %d", failed, len(tvs))
}

func LoadTestVectors(t *testing.T, filename string) []TestVector {
	t.Helper()

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	var raw []struct {
		InHex     string `json:"in"`
		BaseHex   string `json:"base"`
		ExpectHex string `json:"expect"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}

	var tvs []TestVector
	for _, r := range raw {
		tvs = append(tvs, TestVector{
			In:     DecodeHex32(t, r.InHex),
			Base:   DecodeHex32(t, r.BaseHex),
			Expect: DecodeHex32(t, r.ExpectHex),
		})
	}

	return tvs
}

func DecodeHex32(t *testing.T, h string) [32]byte {
	t.Helper()
	b, err := hex.DecodeString(h)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 32 {
		t.Fatal("expected length 32")
	}
	var b32 [32]byte
	copy(b32[:], b)
	return b32
}
