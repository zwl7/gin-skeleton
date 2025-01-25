package utils

import (
	"strings"
)

const (
	default_alphabet   = "asdfghjklURLEncoderConfigqwertyui"
	default_block_size = uint(24)
	min_length         = 5
	one                = uint64(1)
)

type urlEncoder struct {
	alphabet   string
	block_size uint
}

type URLEncoderConfig struct {
	Alphabet  string
	BlockSize uint
}

func NewURLEncoder(config *URLEncoderConfig) *urlEncoder {
	_alphabet := default_alphabet
	_blockSize := default_block_size
	if config.Alphabet != "" {
		_alphabet = config.Alphabet
	}
	if config.BlockSize != 0 {
		_blockSize = config.BlockSize
	}
	_urlEncoder := &urlEncoder{
		alphabet:   _alphabet,
		block_size: _blockSize,
	}
	return _urlEncoder
}

func getBit(n uint64, pos uint) int {
	if (n & (one << pos)) != 0 {
		return 1
	}
	return 0
}

func (encoder *urlEncoder) encode(n uint64) uint64 {
	for i, j := uint(0), uint(encoder.block_size-1); i < j; i, j = i+1, j-1 {
		if getBit(n, i) != getBit(n, j) {
			n ^= ((one << i) | (one << j))
		}
	}
	return n
}

func (encoder *urlEncoder) enbase(x uint64) string {
	n := uint64(len(encoder.alphabet))
	result := []byte{}
	for {
		ch := encoder.alphabet[x%n]
		result = append(result, ch)
		x = x / n
		if x == 0 && len(result) >= min_length {
			break
		}
	}
	revResult := []byte{}
	for i := len(result) - 1; i >= 0; i-- {
		revResult = append(revResult, result[i])
	}
	return string(revResult)
}

func (encoder *urlEncoder) debase(x string) uint64 {
	n := uint64(len(encoder.alphabet))
	result := uint64(0)
	bits := []byte(x)
	for _, bitValue := range bits {
		result = result*n + uint64(strings.IndexByte(encoder.alphabet, bitValue))
	}
	return result
}

func (encoder *urlEncoder) EncodeURL(n uint64) string {
	return encoder.enbase(encoder.encode(n))
}

func (encoder *urlEncoder) DecodeURL(n string) uint64 {
	return encoder.encode(encoder.debase(n))
}
