package enid

import "math/rand/v2"

var defaultEnid = MustNew(WithEntropy(rand.IntN))

// Next creates and returns a unique enid Id, use the default enid generator,
// which use random number.
// NOTE: Make sure your system is keeping accurate system time
func Next() Id { return defaultEnid.Next() }
