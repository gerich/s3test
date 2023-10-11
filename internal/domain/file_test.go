package domain

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Табличный тест сделать не получилось из-за необходимости сравнивать значения
func TestParts(t *testing.T) {
	tt := []struct {
		count int
		data  string
		parts []string
	}{
		{
			count: 4,
			data:  "foobarbaztest", // 13
			parts: []string{"foo", "bar", "baz", "test"},
		},
		{
			count: 4,
			data:  "foobarbazte", // 11
			parts: []string{"foo", "bar", "baz", "te"},
		},
		{
			count: 6,
			data:  "foobarbazfoobarbaz", // 18
			parts: []string{"foo", "bar", "baz", "foo", "bar", "baz"},
		},
		{
			count: 3,
			data:  "foobarbazfoobar", // 15
			parts: []string{"fooba", "rbazf", "oobar"},
		},
		{
			count: 6,
			data:  "foobarbazfoobarbazf", // 19
			parts: []string{"foo", "bar", "baz", "foo", "bar", "bazf"},
		},
		{
			count: 6,
			data:  "foobarbazfoobarbazfo", // 20
			parts: []string{"foo", "bar", "baz", "foo", "bar", "bazfo"},
		},
		{
			count: 6,
			data:  "foobarbazfoobarbazfoo", // 21
			parts: []string{"foo", "bar", "baz", "foo", "bar", "bazfoo"},
		},
		{
			count: 6,
			data:  "foobarbazfoobarbazfoob", // 22
			parts: []string{"foob", "arba", "zfoo", "barb", "azfo", "ob"},
		},
		{
			count: 6,
			data:  "foobarbazfoobarbazfooba", // 23
			parts: []string{"foob", "arba", "zfoo", "barb", "azfo", "oba"},
		},
	}

	for _, tc := range tt {
		t.Run(fmt.Sprintf("%d bytes and %d parts", len(tc.data), tc.count), func(t *testing.T) {
			data := []byte(tc.data)
			file := &File{data: data, reader: bytes.NewReader(data)}
			for idx, part := range tc.parts {
				assert.Equal(t, []byte(part), file.Parts(tc.count, idx))
			}
		})
	}
}
