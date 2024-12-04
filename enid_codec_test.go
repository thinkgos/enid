package enid

import "testing"

func Test_Base2(t *testing.T) {
	node, err := New()
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	oID := node.Next()
	i := oID.Base2()

	pID, err := ParseBase2(i)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}
	if pID != oID {
		t.Fatalf("pID %v != oID %v", pID, oID)
	}

	ms := `111101111111101110110101100101001000000000000000000000000000`
	_, err = ParseBase2(ms)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}

	ms = `1112316766490855473152`
	_, err = ParseBase2(ms)
	if err == nil {
		t.Fatalf("no error parsing %s", ms)
	}
}

func Test_Base32(t *testing.T) {
	node, err := New()
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	for i := 0; i < 100; i++ {
		sf := node.Next()
		b32i := sf.Base32()
		psf, err := ParseBase32([]byte(b32i))
		if err != nil {
			t.Fatal(err)
		}
		if sf != psf {
			t.Fatal("Parsed does not match String.")
		}
	}
}

func Test_Base36(t *testing.T) {
	node, err := New()
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	oID := node.Next()
	i := oID.Base36()

	pID, err := ParseBase36(i)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}
	if pID != oID {
		t.Fatalf("pID %v != oID %v", pID, oID)
	}

	ms := `8hgmw4blvlkw`
	_, err = ParseBase36(ms)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}

	ms = `68h5gmw443blv2lk1w`
	_, err = ParseBase36(ms)
	if err == nil {
		t.Fatalf("no error parsing, %s", err)
	}
}

func Test_Base58(t *testing.T) {
	node, err := New()
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	for i := 0; i < 10; i++ {
		sf := node.Next()
		b58 := sf.Base58()
		psf, err := ParseBase58([]byte(b58))
		if err != nil {
			t.Fatal(err)
		}
		if sf != psf {
			t.Fatal("Parsed does not match String.")
		}
	}
}

func Test_Base64(t *testing.T) {
	node, err := New()
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	oID := node.Next()
	i := oID.Base64()

	pID, err := ParseBase64(i)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}
	if pID != oID {
		t.Fatalf("pID %v != oID %v", pID, oID)
	}

	ms := `MTExNjgxOTQ5NDY2MDk5NzEyMA==`
	_, err = ParseBase64(ms)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}

	ms = `MTExNjgxOTQ5NDY2MDk5NzEyMA`
	_, err = ParseBase64(ms)
	if err == nil {
		t.Fatalf("no error parsing, %s", err)
	}
}

func Test_Bytes(t *testing.T) {
	node, err := New()
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	oID := node.Next()
	i := oID.Bytes()

	pID, err := ParseBytes(i)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}
	if pID != oID {
		t.Fatalf("pID %v != oID %v", pID, oID)
	}

	ms := []byte{0x31, 0x31, 0x31, 0x36, 0x38, 0x32, 0x31, 0x36, 0x37, 0x39, 0x35, 0x37, 0x30, 0x34, 0x31, 0x39, 0x37, 0x31, 0x32}
	_, err = ParseBytes(ms)
	if err != nil {
		t.Fatalf("error parsing, %#v", err)
	}

	ms = []byte{0xFF, 0xFF, 0xFF, 0x31, 0x31, 0x31, 0x36, 0x38, 0x32, 0x31, 0x36, 0x37, 0x39, 0x35, 0x37, 0x30, 0x34, 0x31, 0x39, 0x37, 0x31, 0x32}
	_, err = ParseBytes(ms)
	if err == nil {
		t.Fatalf("no error parsing, %#v", err)
	}
}

func Test_IntBytes(t *testing.T) {
	node, err := New()
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	oID := node.Next()
	i := oID.IntBytes()

	pID := ParseIntBytes(i)
	if pID != oID {
		t.Fatalf("pID %v != oID %v", pID, oID)
	}

	ms := [8]uint8{0xf, 0x7f, 0xc0, 0xfc, 0x2f, 0x80, 0x0, 0x0}
	mi := int64(1116823421972381696)
	pID = ParseIntBytes(ms)
	if pID.Int64() != mi {
		t.Fatalf("pID %v != mi %v", pID.Int64(), mi)
	}
}

func Test_ParseBase32(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    Id
		wantErr bool
	}{
		{
			name:    "ok",
			arg:     "3bodfmr6k2222",
			want:    1427970479175499776,
			wantErr: false,
		},
		{
			name:    "capital case is invalid encoding",
			arg:     "B8WJM1ZROYYYY",
			want:    -1,
			wantErr: true,
		},
		{
			name:    "l is not allowed",
			arg:     "b8wjm1zroyyyl",
			want:    -1,
			wantErr: true,
		},
		{
			name:    "v is not allowed",
			arg:     "b8wjm1zroyyyv",
			want:    -1,
			wantErr: true,
		},
		{
			name:    "2 is not allowed",
			arg:     "b8wjm1zroyyy2",
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseBase32([]byte(tt.arg))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseBase32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseBase32() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ParseBase58(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    Id
		wantErr bool
	}{
		{
			name:    "ok",
			arg:     "4jgmnx8Js8A",
			want:    1428076403798048768,
			wantErr: false,
		},
		{
			name:    "0 not allowed",
			arg:     "0jgmnx8Js8A",
			want:    -1,
			wantErr: true,
		},
		{
			name:    "I not allowed",
			arg:     "Ijgmnx8Js8A",
			want:    -1,
			wantErr: true,
		},
		{
			name:    "O not allowed",
			arg:     "Ojgmnx8Js8A",
			want:    -1,
			wantErr: true,
		},
		{
			name:    "l not allowed",
			arg:     "ljgmnx8Js8A",
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseBase58([]byte(tt.arg))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseBase58() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseBase58() got = %v, want %v", got, tt.want)
			}
		})
	}
}
