package enid

import (
	"bytes"
	"math/rand/v2"
	"reflect"
	"slices"
	"testing"
)

func Test_PrintAll(t *testing.T) {
	node, err := New(WithEntropy(rand.IntN))
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}
	id := node.Next()

	t.Logf("Int64    : %#v", id.Int64())
	t.Logf("String   : %#v", id.String())
	t.Logf("Base2    : %#v", id.Base2())
	t.Logf("Base32   : %#v", id.Base32())
	t.Logf("Base36   : %#v", id.Base36())
	t.Logf("Base58   : %#v", id.Base58())
	t.Logf("Base64   : %#v", id.Base64())
	t.Logf("Bytes    : %#v", id.Bytes())
	t.Logf("IntBytes : %#v", id.IntBytes())
}

// lazy check if Generate will create duplicate IDs
// would be good to later enhance this with more smarts
func Test_GenerateDuplicateID(t *testing.T) {
	node, _ := New(WithNode(1))

	var x, y Id
	for i := 0; i < 1000000; i++ {
		y = node.Next()
		if x == y {
			t.Errorf("x(%d) & y(%d) are the same", x, y)
		}
		x = y
	}
}

func Test_Order(t *testing.T) {
	node, _ := New(WithEntropy(rand.IntN))
	n := 100000
	bs := make([]int64, 0, n)
	for i := 0; i < n; i++ {
		bs = append(bs, node.Next().Int64())
	}
	if !slices.IsSorted(bs) {
		t.Error("not a order id generate")
	}
}

// I feel like there's probably a better way
func Test_Race(t *testing.T) {
	node, _ := New(WithNode(1))
	go func() {
		for i := 0; i < 1000000000; i++ {
			New(WithNode(1))
		}
	}()

	for i := 0; i < 4000; i++ {
		node.Next()
	}
}

func Test_Int64(t *testing.T) {
	node, err := New()
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	oID := node.Next()
	i := oID.Int64()

	pID := ParseInt64(i)
	if pID != oID {
		t.Fatalf("pID %v != oID %v", pID, oID)
	}

	mi := int64(1116766490855473152)
	pID = ParseInt64(mi)
	if pID.Int64() != mi {
		t.Fatalf("pID %v != mi %v", pID.Int64(), mi)
	}

}

func Test_String(t *testing.T) {
	node, err := New()
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	oID := node.Next()
	si := oID.String()

	pID, err := ParseString(si)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}

	if pID != oID {
		t.Fatalf("pID %v != oID %v", pID, oID)
	}

	ms := `1116766490855473152`
	_, err = ParseString(ms)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}

	ms = `1112316766490855473152`
	_, err = ParseString(ms)
	if err == nil {
		t.Fatalf("no error parsing %s", ms)
	}
}

//******************************************************************************
// Marshall Test Methods

func Test_MarshalJSON(t *testing.T) {
	id := Id(13587)
	expected := "\"13587\""

	bytes, err := id.MarshalJSON()
	if err != nil {
		t.Fatalf("Unexpected error during MarshalJSON")
	}

	if string(bytes) != expected {
		t.Fatalf("Got %s, expected %s", string(bytes), expected)
	}
}

func Test_MarshalsIntBytes(t *testing.T) {
	id := Id(13587).IntBytes()
	expected := []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x35, 0x13}
	if !bytes.Equal(id[:], expected) {
		t.Fatalf("Expected ID to be encoded as %v, got %v", expected, id)
	}
}

func Test_UnmarshalJSON(t *testing.T) {
	tt := []struct {
		json        string
		expectedID  Id
		expectedErr error
	}{
		{`"13587"`, 13587, nil},
		{`1`, 0, JSONSyntaxError{[]byte(`1`)}},
		{`"invalid`, 0, JSONSyntaxError{[]byte(`"invalid`)}},
	}

	for _, tc := range tt {
		var id Id
		err := id.UnmarshalJSON([]byte(tc.json))
		if !reflect.DeepEqual(err, tc.expectedErr) {
			t.Fatalf("Expected to get error '%s' decoding JSON, but got '%s'", tc.expectedErr, err)
		}

		if id != tc.expectedID {
			t.Fatalf("Expected to get ID '%s' decoding JSON, but got '%s'", tc.expectedID, id)
		}
	}
}