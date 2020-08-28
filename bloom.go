package bloom_filter

import (
	"github.com/golang-collections/go-datastructures/bitarray"
	"hash"
	"hash/crc64"
	"hash/fnv"
	"math"
)

type StandardBloomFilter struct {
	bitmap   bitarray.BitArray
	optimalM uint64
	optimalK uint64
	hashes   []hash.Hash64
}

func NewBloomFilter(itemsCount uint64, fpRate float64) *StandardBloomFilter {
	var (
		optimalM uint64
		optimalK uint64
		hashes   = []hash.Hash64{crc64.New(crc64.MakeTable(crc64.ISO)), fnv.New64()}
	)
	optimalM = getBitMapSize(itemsCount, fpRate)
	optimalK = getOptimalK(fpRate)
	return &StandardBloomFilter{
		bitmap:   bitarray.NewBitArray(optimalM),
		optimalM: optimalM,
		optimalK: optimalK,
		hashes:   hashes,
	}
}

func (s *StandardBloomFilter) Insert(item string) {
	var (
		h1 uint64
		h2 uint64
		ki uint64
	)
	for ki = 0; ki < s.optimalK; ki++ {
		var index = s.getIndex(h1, h2, ki)
		_ = s.bitmap.SetBit(index)
	}
}

func (s *StandardBloomFilter) Contains(item string) bool {
	var (
		h1 uint64
		h2 uint64
		ki uint64
	)
	for ki = 0; ki < s.optimalK; ki++ {
		var (
			err error
			ok  bool
		)
		var index = s.getIndex(h1, h2, ki)
		ok, err = s.bitmap.GetBit(index)
		if err != nil || !ok {
			return false
		}
	}
	return true
}

func (s *StandardBloomFilter) hashKernel(item string) (uint64, uint64) {
	var (
		hash1 hash.Hash64
		hash2 hash.Hash64
		h1    uint64
		h2    uint64
	)
	hash1 = s.hashes[0]
	hash2 = s.hashes[1]
	defer hash1.Reset()
	defer hash2.Reset()
	_ = hash1.Sum([]byte(item))
	_ = hash2.Sum([]byte(item))
	h1 = hash1.Sum64()
	h2 = hash2.Sum64()
	return h1, h2
}

func (s *StandardBloomFilter) getIndex(h1 uint64, h2 uint64, ki uint64) uint64 {
	return ((h1 + ki) * h2) % s.optimalM
}

func getBitMapSize(itemsCount uint64, fpRate float64) uint64 {
	var ln22 float64 = math.Ln2 * math.Ln2
	var size = uint64(-1*float64(itemsCount)*math.Log(fpRate)/ln22) + 1
	return size
}

func getOptimalK(fpRate float64) uint64 {
	var k = uint64(-1*math.Log(fpRate)/math.Ln2) + 1
	return k
}
