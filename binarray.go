package binarray

import (
	"fmt"
	"log"
	"math/bits"
)

type BinArray struct {
	size   int
	board1 uint64 // holds first max 64 bits
	board2 uint64 // holds the rest of the bits
}

var allMap map[int]*BinArray

func initAll(num int) *BinArray {
	all := BinArray{num, 0, 0}
	for pos := 0; pos < num; pos++ {
		all.Set(pos)
	}
	return &all
}

func NewBinArray(num int) (*BinArray, error) {
	if num > 128 {
		return nil, fmt.Errorf("Size is not supported: %d", num)
	}
	var b BinArray
	b.board1 = 0
	b.board2 = 0
	b.size = num
	if allMap == nil {
		allMap = make(map[int]*BinArray)
	}
	_, ok := allMap[num]
	if !ok {
		allMap[num] = initAll(num)
	}
	return &b, nil
}

func (b *BinArray) Equal(w *BinArray) bool {
	return b.board1 == w.board1 && b.board2 == w.board2
}

func (b *BinArray) Size() int {
	return b.size
}

func (b *BinArray) Clone() *BinArray {
	var clone, _ = NewBinArray(b.size)
	clone.board1 = b.board1
	clone.board2 = b.board2
	return clone
}

func (b *BinArray) None() *BinArray {
	b.board1 = 0
	b.board2 = 0
	return b
}

func (b *BinArray) All() *BinArray {
	return allMap[b.size].Clone()
}

func (b *BinArray) Count() int {
	return bits.OnesCount64(b.board1) + bits.OnesCount64(b.board2)
}

func (b *BinArray) Set(pos int) *BinArray {
	if pos < 64 {
		b.board1 = b.board1 | (1 << pos)
	} else {
		b.board2 = b.board2 | (1 << (pos - 64))
	}
	return b
}

func (b *BinArray) Get(pos int) uint64 {
	if pos < 64 {
		return (b.board1 & (1 << pos)) >> pos
	} else {
		pos -= 64
		return (b.board2 & (1 << pos)) >> pos
	}
}

func (b *BinArray) Is(pos int) bool {
	return b.Get(pos) == 1
}

// boolean operators work on all bits.
// Also on the ones outside the board area
func (b *BinArray) And(w *BinArray) *BinArray {
	b.board1 = b.board1 & w.board1
	b.board2 = b.board2 & w.board2
	return b
}

func (b *BinArray) Or(w *BinArray) *BinArray {
	b.board1 = b.board1 | w.board1
	b.board2 = b.board2 | w.board2
	return b
}

func (b *BinArray) Xor(w *BinArray) *BinArray {
	b.board1 = b.board1 ^ w.board1
	b.board2 = b.board2 ^ w.board2
	return b
}

func (b *BinArray) Inverse() *BinArray {
	b.board1 = ^b.board1
	b.board2 = ^b.board2
	return b
}

func (b *BinArray) Minus(w *BinArray) *BinArray {
	b.board1 = b.board1 & ^w.board1
	b.board2 = b.board2 & ^w.board2
	return b
}

// shifting
func (b *BinArray) ShiftLeft(times int) *BinArray {
	if times > 64 {
		log.Fatal("ShiftLeft not implemented for greater than 64")
	}
	num := b.board1 >> (64 - times)
	b.board1 = b.board1 << times
	b.board2 = b.board2 << times
	b.board2 = b.board2 | num
	return b
}

func (b *BinArray) ShiftRight(times int) *BinArray {
	if times > 64 {
		log.Fatal("ShiftRight not implemented for greater than 64")
	}
	num := b.board2 << (64 - times)
	b.board1 = b.board1 >> times
	b.board2 = b.board2 >> times
	b.board1 = b.board1 | num
	return b
}

func (b *BinArray) String() string {
	return fmt.Sprintf("%064b\n%064b\n", b.board1, b.board2)
}
