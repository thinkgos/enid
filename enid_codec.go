package enid

import (
	"encoding/base64"
	"encoding/binary"
	"strconv"
	"unsafe"
)

// Base2 returns a string base2 of the enid ID
func (d Id) Base2() string { return strconv.FormatInt(int64(d), 2) }

// ParseBase2 converts a Base2 string into a enid ID
func ParseBase2(id string) (Id, error) {
	i, err := strconv.ParseInt(id, 2, 64)
	return Id(i), err
}

// Base32 uses the z-base-32 character set but encodes and decodes similar
// to base58, allowing it to create an even smaller result string.
// NOTE: There are many different base32 implementations so becareful when
// doing any interoperation.
func (d Id) Base32() string {
	if d < 32 {
		return string(base32EncodeCharset[d])
	}
	b := make([]byte, 0, 12)
	for d >= 32 {
		b = append(b, base32EncodeCharset[d%32])
		d /= 32
	}
	b = append(b, base32EncodeCharset[d])
	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// ParseBase32 parses a base32 []byte into a enid ID
// NOTE: There are many different base32 implementations so becareful when
// doing any interoperation.
func ParseBase32(b []byte) (Id, error) {
	var id int64

	for i := range b {
		if base32DecodeMap[b[i]] == 0xFF {
			return -1, ErrBase32IllegalChar
		}
		id = id*32 + int64(base32DecodeMap[b[i]])
	}
	return Id(id), nil
}

// Base36 returns a base36 string of the enid ID
func (d Id) Base36() string { return strconv.FormatInt(int64(d), 36) }

// ParseBase36 converts a Base36 string into a enid ID
func ParseBase36(id string) (Id, error) {
	i, err := strconv.ParseInt(id, 36, 64)
	return Id(i), err
}

// Base58 returns a base58 string of the enid ID
func (d Id) Base58() string {
	if d < 58 {
		return string(base58EncodeCharset[d])
	}

	b := make([]byte, 0, 11)
	for d >= 58 {
		b = append(b, base58EncodeCharset[d%58])
		d /= 58
	}
	b = append(b, base58EncodeCharset[d])
	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// ParseBase58 parses a base58 []byte into a enid ID
func ParseBase58(b []byte) (Id, error) {
	var id int64

	for i := range b {
		if base58DecodeMap[b[i]] == 0xFF {
			return -1, ErrBase58IllegalChar
		}
		id = id*58 + int64(base58DecodeMap[b[i]])
	}
	return Id(id), nil
}

// Base64 returns a base64 string of the enid ID
func (d Id) Base64() string { return base64.StdEncoding.EncodeToString(d.Bytes()) }

// ParseBase64 converts a base64 string into a enid ID
func ParseBase64(id string) (Id, error) {
	b, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return -1, err
	}
	return ParseBytes(b)
}

// Bytes returns a byte slice of the enid ID
func (d Id) Bytes() []byte {
	s := d.String()
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// ParseBytes converts a byte slice into a enid ID
func ParseBytes(id []byte) (Id, error) {
	i, err := strconv.ParseInt(string(id), 10, 64)
	return Id(i), err
}

// IntBytes returns an array of bytes of the enid ID, encoded as a
// big endian integer.
func (d Id) IntBytes() [8]byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(d))
	return b
}

// ParseIntBytes converts an array of bytes encoded as big endian integer as
// a enid ID
func ParseIntBytes(id [8]byte) Id { return Id(int64(binary.BigEndian.Uint64(id[:]))) }
