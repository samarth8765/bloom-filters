package bloom

import (
	"hash"
	"math/rand"

	"github.com/spaolacci/murmur3"
)

const DEFAULT_HASH_FUNCTIONS = 2
const DEFAULT_BYTE_SIZE = 10

type BloomFilter struct {
	size     uint64
	bitArray []byte
	numHash  uint8
	hashFns  []hash.Hash64
}

func NewBloomFilter(size uint64, numHash uint8) *BloomFilter {
	if size == 0 {
		size = DEFAULT_BYTE_SIZE
	}

	if numHash == 0 {
		numHash = DEFAULT_HASH_FUNCTIONS
	}

	hashFns := make([]hash.Hash64, numHash)
	for i := 0; i < int(numHash); i++ {
		hashFns[i] = murmur3.New64WithSeed(rand.Uint32())
	}

	return &BloomFilter{
		size:     size,
		bitArray: make([]byte, size),
		numHash:  numHash,
		hashFns:  hashFns,
	}
}

func (bloom *BloomFilter) Add(key []byte) error {
	positions, err := bloom.getPositions(key)
	if err != nil {
		return err
	}

	for _, pos := range positions {
		setBit(bloom.bitArray, pos)
	}

	return nil
}

func (bloom *BloomFilter) Check(key []byte) (bool, []uint64, error) {
	positions, err := bloom.getPositions(key)
	if err != nil {
		return false, nil, err
	}

	for _, pos := range positions {
		if getBit(bloom.bitArray, pos) == 0 {
			return false, positions, nil
		}
	}

	return true, positions, nil
}

func (bloom *BloomFilter) getPositions(key []byte) ([]uint64, error) {
	positions := make([]uint64, bloom.numHash)
	for i, hashFn := range bloom.hashFns {
		hashFn.Reset()
		_, err := hashFn.Write(key)
		if err != nil {
			return nil, err
		}
		hash := hashFn.Sum64() % (bloom.size * 8)
		positions[i] = hash

	}
	return positions, nil
}
