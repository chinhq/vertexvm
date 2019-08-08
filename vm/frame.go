package vm

import (
	"log"

	"github.com/go-interpreter/wagon/wasm"
)

// Frame or call frame holds the relevant execution information of a function
type Frame struct {
	fn             *wasm.Function
	ip             int
	basePointer    int
	baseBlockIndex int
}

// NewFrame initialize a call frame for a given function fn
func NewFrame(fn *wasm.Function, basePointer int, baseBlockIndex int) *Frame {
	f := &Frame{
		fn:             fn,
		ip:             -1,
		basePointer:    basePointer,
		baseBlockIndex: baseBlockIndex,
	}
	return f
}

func (frame *Frame) readLEB(maxbit uint32, hasSign bool) int64 {
	ins := frame.instructions()
	var (
		shift  uint32
		bitcnt uint32
		cur    int64
		result int64
		sign   int64 = -1
	)
	for i := frame.ip + 1; i < len(ins); i++ {
		cur = int64(ins[i])
		result |= (cur & 0x7f) << shift
		shift += 7
		sign <<= 7
		bitcnt++
		if cur&0x80 == 0 {
			break
		}
		if bitcnt > (maxbit+7-1)/7 {
			log.Fatal("Unsigned LEB at byte overflow")
		}
	}
	if hasSign && ((sign>>1)&result) != 0 {
		result |= sign
	}
	frame.ip += int(bitcnt)
	return result
}

func (frame *Frame) instructions() []byte {
	return frame.fn.Body.Code
}

func (frame *Frame) hasEnded() bool {
	return frame.ip == len(frame.instructions())-1
}
