package bloom

import (
	"hash"
	"math/rand"

	"github.com/spaolacci/murmur3"
)

const DEFAULT_HASH_FUNCTIONS = 2

type BloomFilter struct {
	N        uint64
	BitArray []bool
	hashFns  []hash.Hash64
	NumHash  uint8
}

func NewBloomFilter(n uint64, m uint8) *BloomFilter {
	if m <= 0 {
		m = DEFAULT_HASH_FUNCTIONS
	}

	hashFns := make([]hash.Hash64, m)
	for i := 0; i < int(m); i++ {
		hashFns[i] = murmur3.New64WithSeed(rand.Uint32())
	}

	return &BloomFilter{
		N:        n,
		BitArray: make([]bool, n),
		NumHash:  m,
		hashFns:  hashFns,
	}
}

func (bloom *BloomFilter) Add(key []byte) error {
	for _, hashFn := range bloom.hashFns {
		hashFn.Reset()
		_, err := hashFn.Write(key)
		if err != nil {
			return err
		}
		hash := hashFn.Sum64() % (bloom.N)
		bloom.BitArray[hash] = true
	}
	return nil
}

func (bloom *BloomFilter) Check(key string) (bool, []uint64, error) {
	hashIdxs := make([]uint64, bloom.NumHash)
	for i, hashFn := range bloom.hashFns {
		hashFn.Reset()
		_, err := hashFn.Write([]byte(key))
		if err != nil {
			return false, []uint64{}, err
		}
		hash := hashFn.Sum64() % (bloom.N)
		if !bloom.BitArray[hash] {
			return false, []uint64{}, nil
		}

		hashIdxs[i] = hash
	}
	return true, hashIdxs, nil
}
