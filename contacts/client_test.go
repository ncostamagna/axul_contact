package contacts

import "testing"

func BenchmarkHttpSameValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetTemplateHTTP(2)
	}
}

func BenchmarkGrpcSameValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetTemplate(2)
	}
}

func BenchmarkHttpManyValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetTemplateHTTP(uint(i))
	}
}

func BenchmarkGrpcManyValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetTemplate(uint(i))
	}
}
