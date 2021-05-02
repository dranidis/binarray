package bitarray

import (
	"fmt"
	"log"
	"math/bits"
)

type BitArray struct {
	size   int
	blocks []uint64
}

const (
	INTSIZE = 64
)

// New returns a *BitArray with the provided size.
func New(size int) *BitArray {
	b := newBitArray(size)
	if cachedAllValues == nil {
		cachedAllValues = make(map[int]*BitArray)
	}
	_, ok := cachedAllValues[size]
	if !ok {
		cachedAllValues[size] = allValueFor(size)
	}
	return b
}

func newBitArray(size int) *BitArray {
	var b BitArray
	b.size = size
	numBlocks := b.size/INTSIZE + 1
	b.blocks = make([]uint64, numBlocks)
	return &b
}

var cachedAllValues map[int]*BitArray

func allValueFor(size int) *BitArray {
	all := newBitArray(size)
	for index := 0; index < size; index++ {
		all.Set(index)
	}
	return all
}

// Equal is a deep equality comparison.
func (b *BitArray) Equal(w *BitArray) bool {
	if b.Size() != w.Size() {
		return false
	}
	b.zeroTheRemainder()
	w.zeroTheRemainder()
	for i := range b.blocks {
		if b.blocks[i] != w.blocks[i] {
			return false
		}
	}
	return true
}

// Size returns the size of the bit array.
func (b *BitArray) Size() int {
	return b.size
}

// Clone creates a copy of the original bit array.
// It does not change the receiver.
func (b *BitArray) Clone() *BitArray {
	var clone = newBitArray(b.size)
	for i := range b.blocks {
		clone.blocks[i] = b.blocks[i]
	}
	return clone
}

// None sets all bits to 0.
func (b *BitArray) None() *BitArray {
	for i := range b.blocks {
		b.blocks[i] = 0
	}
	return b
}

// All returns a bit array with all bits set to 1.
// It does not change the receiver.
func (b *BitArray) All() *BitArray {
	return cachedAllValues[b.size].Clone()
}

// Count returns the number of 1 bits in the array.
// It sets all the bits outside the size to 0.
func (b *BitArray) Count() int {
	b.zeroTheRemainder()
	count := 0
	for i := range b.blocks {
		count += bits.OnesCount64(b.blocks[i])
	}
	return count
}

// Set changes the pos bit to 1 and returns the array.
func (b *BitArray) Set(index int) *BitArray {
	if index >= b.size {
		log.Panic(fmt.Sprintf("BitArray.Set: index out of range: %d", b.size))
	}
	blockIndex, position := indexPos(index)
	b.blocks[blockIndex] |= 1 << position
	return b
}

// Get returns the pos value of the bit array.
// It does not change the receiver.
func (b *BitArray) Get(index int) uint64 {
	if index >= b.size {
		log.Panic(fmt.Sprintf("BitArray.Get: index out of range: %d", b.size))
	}
	blockIndex, position := indexPos(index)
	return (b.blocks[blockIndex] & (1 << position)) >> position
}

// Is checks if the pos bit of the bit array is 1.
// It does not change the receiver.
func (b *BitArray) Is(pos int) bool {
	return b.Get(pos) == 1
}

// And performs the and boolean operator on all bits.
// Also on the ones outside the size area.
func (b *BitArray) And(w *BitArray) *BitArray {
	for i := range b.blocks {
		b.blocks[i] &= w.blocks[i]
	}
	return b
}

// Or performs the or boolean operator on all bits.
// Also on the ones outside the size area.
func (b *BitArray) Or(w *BitArray) *BitArray {
	for i := range b.blocks {
		b.blocks[i] |= w.blocks[i]
	}
	return b
}

// Xor performs the xor boolean operator on all bits.
// Also on the ones outside the size area.
func (b *BitArray) Xor(w *BitArray) *BitArray {
	for i := range b.blocks {
		b.blocks[i] ^= w.blocks[i]
	}
	return b
}

// Inverse inverses all bits.
// Also the ones outside the size area.
func (b *BitArray) Inverse() *BitArray {
	numBLocks := len(b.blocks)
	for i := 0; i < numBLocks; i++ {
		b.blocks[i] = ^b.blocks[i]
	}
	return b
}

// Minus performs the minus boolean operator on all bits.
// Also on the ones outside the size area.
func (b *BitArray) Minus(w *BitArray) *BitArray {
	for i := range b.blocks {
		b.blocks[i] = b.blocks[i] & ^w.blocks[i]
	}
	return b
}

// ShiftLeft shifts all bits times to the left.
// The parameter times cannot be greater than 64.
func (b *BitArray) ShiftLeft(times int) *BitArray {
	if times > INTSIZE {
		log.Panic("ShiftLeft not implemented for greater than 64")
	}
	prevTmp := uint64(0)
	for i := 0; i < len(b.blocks); i++ {
		tmp := b.blocks[i] >> (INTSIZE - times)
		b.blocks[i] <<= times
		b.blocks[i] |= prevTmp
		prevTmp = tmp
	}
	return b
}

// ShiftLeft shifts all bits times to the right.
// The parameter times cannot be greater than 64.
func (b *BitArray) ShiftRight(times int) *BitArray {
	if times > INTSIZE {
		log.Panic("ShiftRight not implemented for greater than 64")
	}
	prevTmp := uint64(0)

	for i := len(b.blocks) - 1; i >= 0; i-- {
		tmp := b.blocks[i] << (INTSIZE - times)
		b.blocks[i] >>= times
		b.blocks[i] |= prevTmp
		prevTmp = tmp
	}
	return b
}

// String returns a string of all bits
// organized in 64 bits
func (b *BitArray) String() string {
	b.zeroTheRemainder()
	str := ""
	for i := range b.blocks {
		str += fmt.Sprintf("%064b\n", b.blocks[i])
	}
	return str
}

// StringBreak returns a string of all bits inside the range
// organized in breakLine number of bits
func (b *BitArray) StringBreak(breakLine int) string {
	b.zeroTheRemainder()
	str := ""
	for i := 0; i < b.size; i++ {
		if i%breakLine == 0 {
			str += "\n"
		}
		str += fmt.Sprintf(" %d", b.Get(i))
	}
	return str
}

// indexPos is a helper function that returns
// the index of the blocks and the position within the block
// corresponding to the pos argument.
func indexPos(pos int) (int, int) {
	index := pos / INTSIZE
	position := pos - index*INTSIZE
	return index, position
}

func (b *BitArray) zeroTheRemainder() {
	numBLocks := len(b.blocks)
	b.blocks[numBLocks-1] &= cachedAllValues[b.size].blocks[numBLocks-1]
}
