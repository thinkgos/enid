package enid

import "testing"

func Benchmark_ParseBase32(b *testing.B) {
	node, _ := New(WithNode(1))
	sf := node.Next()
	b32i := sf.Base32()

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = ParseBase32([]byte(b32i))
	}
}
func Benchmark_Base32(b *testing.B) {
	node, _ := New(WithNode(1))
	sf := node.Next()

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sf.Base32()
	}
}
func Benchmark_ParseBase58(b *testing.B) {
	node, _ := New(WithNode(1))
	sf := node.Next()
	b58 := sf.Base58()

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = ParseBase58([]byte(b58))
	}
}
func Benchmark_Base58(b *testing.B) {
	node, _ := New(WithNode(1))
	sf := node.Next()

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sf.Base58()
	}
}
func Benchmark_Generate(b *testing.B) {
	node, _ := New(WithNode(1))

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = node.Next()
	}
}

func Benchmark_GenerateMaxSequence(b *testing.B) {
	node, _ := New(WithNode(1), WithNodeStepBits(1, 19))

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = node.Next()
	}
}

func Benchmark_Unmarshal(b *testing.B) {
	// Generate the ID to unmarshal
	node, _ := New(WithNode(1))
	id := node.Next()
	bytes, _ := id.MarshalJSON()

	var id2 Id

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = id2.UnmarshalJSON(bytes)
	}
}

func Benchmark_Marshal(b *testing.B) {
	// Generate the ID to marshal
	node, _ := New(WithNode(1))
	id := node.Next()

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = id.MarshalJSON()
	}
}
