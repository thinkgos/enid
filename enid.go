package enid

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

const (
	// defaultEpoch is set to the enid epoch of Dec 01 2024 05:06:07 UTC in milliseconds
	defaultEpoch        int64 = 1733029567666
	base32EncodeCharset       = "234567abcdefghijklmnopqrstuvwxyz"
	base58EncodeCharset       = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
)

var (
	// ErrBase58IllegalChar is returned by ParseBase58 when given an invalid []byte
	ErrBase58IllegalChar = errors.New("illegal base58 char")
	// ErrBase32IllegalChar is returned by ParseBase32 when given an invalid []byte
	ErrBase32IllegalChar = errors.New("illegal base32 char")
)

var base32DecodeMap [256]byte
var base58DecodeMap [256]byte

// A JSONSyntaxError is returned from UnmarshalJSON if an invalid Id is provided.
type JSONSyntaxError struct{ original []byte }

func (j JSONSyntaxError) Error() string {
	return fmt.Sprintf("invalid enid Id %q", string(j.original))
}

// Create maps for decoding Base58/Base32. this speeds up the process tremendously.
func init() {
	for i := 0; i < len(base58DecodeMap); i++ {
		base58DecodeMap[i] = 0xFF
	}
	for i := 0; i < len(base58EncodeCharset); i++ {
		base58DecodeMap[base58EncodeCharset[i]] = byte(i)
	}

	for i := 0; i < len(base32DecodeMap); i++ {
		base32DecodeMap[i] = 0xFF
	}
	for i := 0; i < len(base32EncodeCharset); i++ {
		base32DecodeMap[base32EncodeCharset[i]] = byte(i)
	}
}

// A Enid struct holds the basic information needed for a enid generator.
type Enid struct {
	mu    sync.Mutex
	epoch time.Time
	time  int64
	node  int64
	step  int64

	nodeBits      uint8
	stepBits      uint8
	timeShift     uint8
	stepShift     uint8
	nodeMax       int64
	nodeMask      int64
	stepMask      int64
	enableEntropy bool
	entropy       func(n int) int
}

// An Id is a custom type used for a enid Id.  This is used so we can attach methods onto the Id.
type Id int64

type Option func(*Enid)

// WithNode customize this to set a different node for your application.
func WithNode(node int64) Option {
	return func(e *Enid) {
		e.node = node
	}
}

// WithEntropy customize this to set rand node value.
func WithEntropy(f func(n int) int) Option {
	return func(e *Enid) {
		e.enableEntropy = true
		e.entropy = f
	}
}

// WithEpoch customize this to set a different epoch(UTC in milliseconds) for your application.
func WithEpoch(epoch int64) Option {
	return func(e *Enid) {
		curTime := time.Now()
		// add time.Duration to curTime to make sure we use the monotonic clock if available
		e.epoch = curTime.Add(time.UnixMilli(epoch).Sub(curTime))
	}
}

// WithNodeStepBits customize this to set Node/Step.
// We have a total 20 bits to share between Node/Step
func WithNodeStepBits(nodeBits, stepBits uint8) Option {
	return func(e *Enid) {
		e.nodeBits = nodeBits
		e.stepBits = stepBits
	}
}

// MustNew is a convenience function equivalent to New that panics on failure
// instead of returning an error.
func MustNew(opts ...Option) *Enid {
	e, err := New(opts...)
	if err != nil {
		panic(err)
	}
	return e
}

// New returns a new enid node that can be used to generate enid Ids
func New(opts ...Option) (*Enid, error) {
	n := &Enid{
		node:     0,
		nodeBits: 8,
		stepBits: 12,
	}
	WithEpoch(defaultEpoch)(n)
	for _, f := range opts {
		f(n)
	}

	if n.nodeBits+n.stepBits > 20 {
		return nil, errors.New("we have a total 20 bits to share between Node/Step")
	}

	n.timeShift = n.nodeBits + n.stepBits
	n.stepShift = n.nodeBits
	n.nodeMax = -1 ^ (-1 << n.nodeBits)
	n.nodeMask = n.nodeMax << n.stepBits
	n.stepMask = -1 ^ (-1 << n.stepBits)
	if n.node < 0 || n.node > n.nodeMax {
		return nil, errors.New("node number must be between 0 and " + strconv.FormatInt(n.nodeMax, 10))
	}
	return n, nil
}

// Next creates and returns a unique enid Id
// To help guarantee uniqueness
// - Make sure your system is keeping accurate system time
// - Make sure you never have multiple nodes running with the same node Id
func (d *Enid) Next() Id {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Since(d.epoch).Milliseconds()
	if now == d.time {
		d.step = (d.step + 1) & d.stepMask
		if d.step == 0 {
			for now <= d.time {
				now = time.Since(d.epoch).Milliseconds()
			}
		}
	} else {
		d.step = 0
	}
	d.time = now
	node := d.node
	if d.enableEntropy {
		node = int64(d.entropy(int(d.nodeMax)))
	}
	r := Id((now)<<d.timeShift | (d.step << d.stepShift) | (node))
	return r
}

// Int64 returns an int64 of the enid Id
func (d Id) Int64() int64 { return int64(d) }

// ParseInt64 converts an int64 into a enid Id
func ParseInt64(id int64) Id { return Id(id) }

// String returns a string of the enid Id
func (d Id) String() string { return strconv.FormatInt(int64(d), 10) }

// ParseString converts a string into a enid Id
func ParseString(id string) (Id, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	return Id(i), err
}

// MarshalJSON returns a json byte array string of the enid Id.
func (d Id) MarshalJSON() ([]byte, error) {
	buff := make([]byte, 0, 22)
	buff = append(buff, '"')
	buff = strconv.AppendInt(buff, int64(d), 10)
	buff = append(buff, '"')
	return buff, nil
}

// UnmarshalJSON converts a json byte array of a enid Id into an Id type.
func (d *Id) UnmarshalJSON(b []byte) error {
	if len(b) < 3 || b[0] != '"' || b[len(b)-1] != '"' {
		return JSONSyntaxError{b}
	}

	i, err := strconv.ParseInt(string(b[1:len(b)-1]), 10, 64)
	if err != nil {
		return err
	}

	*d = Id(i)
	return nil
}
