package hw10programoptimization

import (
	"bytes"
	"testing"
)

const jsonData = `{"Email": "user1@example.com"}
{"Email": "user2@example.com"}
{"Email": "user3@test.com"}
{"Email": "user4@example.com"}
{"Email": "user5@test.com"}`

func BenchmarkGetDomainStat(b *testing.B) {
	data := bytes.NewBufferString(jsonData)
	domain := "example.com"

	for i := 0; i < b.N; i++ {
		_, err := GetDomainStat(data, domain)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		data.Reset()
		data.WriteString(jsonData)
	}
}
