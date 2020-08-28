package bloom_filter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStandardBloomFilter_Insert(t *testing.T) {
	var bloom = NewBloomFilter(100, 0.01)
	bloom.Insert("item")
	assert.True(t, bloom.Contains("item"))
}

func TestStandardBloomFilter_Contains(t *testing.T) {
	var bloom = NewBloomFilter(100, 0.01)
	assert.True(t, !bloom.Contains("item_1"))
	assert.True(t, !bloom.Contains("item_2"))
	bloom.Insert("item_1")
	assert.True(t, bloom.Contains("item_1"))
}