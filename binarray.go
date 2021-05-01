package binarray

import (
	"fmt"
	"log"
	"math/bits"
)

type BinArray struct {
	size   int
	blocks []uint64
}

// NewBinArray returns a *BinArray possible to hold num bits.
func NewBinArray(num int) *BinArray {
	b := newBinArray(num)
	if cachedAllValues == nil {
		cachedAllValues = make(map[int]*BinArray)
	}
	_, ok := cachedAllValues[num]
	if !ok {
		cachedAllValues[num] = allValueFor(num)
	}
	return b
}

func newBinArray(num int) *BinArray {
	var b BinArray
	b.size = num
	numBlocks := b.size/64 + 1
	b.blocks = make([]uint64, numBlocks)
	return &b
}

var cachedAllValues map[int]*BinArray

func allValueFor(num int) *BinArray {
	all := newBinArray(num)
	for pos := 0; pos < num; pos++ {
		all.Set(pos)
	}
	return all
}

// Equal is a deep equality comparison.
func (b *BinArray) Equal(w *BinArray) bool {
	for i := range b.blocks {
		if b.blocks[i] != w.blocks[i] {
			return false
		}
	}
	return true
}

// Size returns the total number of bits.
func (b *BinArray) Size() int {
	return b.size
}

// Clone creates a copy of the original bin matrix.
// It does not change the receiver.
func (b *BinArray) Clone() *BinArray {
	var clone = newBinArray(b.size)
	for i := range b.blocks {
		clone.blocks[i] = b.blocks[i]
	}
	return clone
}

// None sets all bits to 0.
func (b *BinArray) None() *BinArray {
	for i := range b.blocks {
		b.blocks[i] = 0
	}
	return b
}

// All returns a bit matrix with all bits set to 1.
// It does not change the receiver.
func (b *BinArray) All() *BinArray {
	return cachedAllValues[b.size].Clone()
}

// Count returns the number of 1 bits in the matrix.
// It does not change the receiver.
func (b *BinArray) Count() int {
	count := 0
	for i := range b.blocks {
		count += bits.OnesCount64(b.blocks[i])
	}
	return count
}

// Set changes the pos bit to 1 and returns the matrix.
func (b *BinArray) Set(pos int) *BinArray {
	index, position := indexPos(pos)
	b.blocks[index] |= 1 << position
	return b
}

// Get returns the pos value of the bin matrix.
// It does not change the receiver.
func (b *BinArray) Get(pos int) uint64 {
	index, position := indexPos(pos)
	return (b.blocks[index] & (1 << position)) >> position
}

// Is checks if the pos bit of the bin matrix is 1.
// It does not change the receiver.
func (b *BinArray) Is(pos int) bool {
	return b.Get(pos) == 1
}

// And performs the and boolean operator on all bits.
// Also on the ones outside the size area.
func (b *BinArray) And(w *BinArray) *BinArray {
	for i := range b.blocks {
		b.blocks[i] &= w.blocks[i]
	}
	return b
}

// Or performs the or boolean operator on all bits.
// Also on the ones outside the size area.
func (b *BinArray) Or(w *BinArray) *BinArray {
	for i := range b.blocks {
		b.blocks[i] |= w.blocks[i]
	}
	return b
}

// Xor performs the xor boolean operator on all bits.
// Also on the ones outside the size area.
func (b *BinArray) Xor(w *BinArray) *BinArray {
	for i := range b.blocks {
		b.blocks[i] ^= w.blocks[i]
	}
	return b
}

// Inverse inverses all bits.
// Also the ones outside the size area.
func (b *BinArray) Inverse() *BinArray {
	for i := range b.blocks {
		b.blocks[i] = ^b.blocks[i]
	}
	return b
}

// Minus performs the minus boolean operator on all bits.
// Also on the ones outside the size area.
func (b *BinArray) Minus(w *BinArray) *BinArray {
	for i := range b.blocks {
		b.blocks[i] = b.blocks[i] & ^w.blocks[i]
	}
	return b
}

// ShiftLeft shifts all bits times to the left.
// The parameter times cannot be greater than 64.
func (b *BinArray) ShiftLeft(times int) *BinArray {
	if times > 64 {
		log.Fatal("ShiftLeft not implemented for greater than 64")
	}
	for i := 0; i < len(b.blocks)-1; i++ {
		tmp := b.blocks[i] >> (64 - times)
		b.blocks[i] <<= times
		b.blocks[i+1] <<= times
		b.blocks[i+1] |= tmp
	}
	return b
}

// ShiftLeft shifts all bits times to the right.
// The parameter times cannot be greater than 64.
func (b *BinArray) ShiftRight(times int) *BinArray {
	if times > 64 {
		log.Fatal("ShiftRight not implemented for greater than 64")
	}
	for i := 0; i < len(b.blocks)-1; i++ {
		tmp := b.blocks[i+1] << (64 - times)
		b.blocks[i] >>= times
		b.blocks[i+1] >>= times
		b.blocks[i] |= tmp
	}
	return b
}

// String returns a string of all bits
// organized in 64 bits
func (b *BinArray) String() string {
	str := ""
	for i := range b.blocks {
		str += fmt.Sprintf("%064b\n", b.blocks[i])
	}
	return str
}

// indexPos is a helper function that returns
// the index of the blocks and the position within the block
// corresponding to the pos argument.
func indexPos(pos int) (int, int) {
	index := pos / 64
	position := pos - index*64
	return index, position
}
