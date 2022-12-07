package rs

import (
	"github.com/klauspost/reedsolomon"
	"io"
)

type encoder struct {
	writers []io.Writer
	enc     reedsolomon.Encoder
	cache   []byte
}

func NewEncoder(writers []io.Writer) *encoder {
	enc, _ := reedsolomon.New(DATA_SHARDS, PARITY_SHARDS) //数据片数量 和 校验片数量
	return &encoder{writers, enc, nil}
}

func (e *encoder) Write(p []byte) (n int, err error) {
	length := len(p)
	current := 0
	for length != 0 {
		// 剩余多少没有填进去
		remainder := BLOCK_SIZE - len(e.cache)
		// 剩余的够 这次填充
		if remainder > length {
			remainder = length // 剩余 length
		}
		e.cache = append(e.cache, p[current:current+remainder]...)
		if len(e.cache) == BLOCK_SIZE {
			e.Flush()
		}
		//p 填入了多少
		current += remainder
		// 还有多少没有填
		length -= remainder
	}
	return len(p), nil
}

func (e *encoder) Flush() {
	if len(e.cache) == 0 {
		return
	}
	shards, _ := e.enc.Split(e.cache)
	e.enc.Encode(shards)
	for i := range shards {
		e.writers[i].Write(shards[i])
	}
	e.cache = e.cache[:0]
}
