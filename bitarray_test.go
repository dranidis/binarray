package bitarray

import (
	"testing"
)

func TestBitArray_AllSetCount(t *testing.T) {
	b := New(81)
	count := b.All().Count()
	if count != 81 {
		t.Errorf("Wrong number of 1 bits: %d", b)
	}
}

func TestBitArray_CachedAll(t *testing.T) {
	b1 := New(81)
	count1 := b1.All().Count()
	b2 := New(81)
	count2 := b2.All().Count()
	b3 := New(81)
	count3 := b3.All().Count()
	if count1 != 81 {
		t.Errorf("Wrong number of 1 bits: %d", b1)
	}
	if count2 != 81 {
		t.Errorf("Wrong number of 1 bits: %d", b2)
	}
	if count3 != 81 {
		t.Errorf("Wrong number of 1 bits: %d", b3)
	}
}

func TestBitArray_SetGetLow(t *testing.T) {
	b := New(81)
	for pos := 0; pos < b.size; pos++ {
		if b.Get(pos) != uint64(0) {
			t.Errorf("Not all zero at beginning")
		}
	}
	b.Set(63)
	a := b.Get(63)
	if a != 1 {
		t.Errorf("Not Set/Get: %d", a)
	}
}

func TestBitArray_SetGetHi(t *testing.T) {
	b := New(81)
	for pos := 0; pos < b.size; pos++ {
		if b.Get(pos) != uint64(0) {
			t.Errorf("Not all zero at beginning")
		}
	}
	b.Set(80)
	a := b.Get(80)
	if a != 1 {
		t.Errorf("Not Set/Get: %d", a)
	}
}

func TestBitArray_InverseInversesAllBits(t *testing.T) {
	b := New(81)
	count := b.Inverse().Count()
	if count != 81 {
		t.Errorf("Wrong number of 1 bits: %d", b)
	}
}

func TestBitArray_ShiftLeftRightCancel(t *testing.T) {
	b := New(81)
	positions := []int{43, 61, 63, 64, 65, 80}
	for pos := 0; pos < len(positions); pos++ {
		b = b.Set(pos)
	}
	b = b.ShiftLeft(12).ShiftRight(12)
	for pos := 0; pos < len(positions); pos++ {
		if !b.Is(pos) {
			t.Errorf("Error in position: %d", pos)
		}
	}
}

func TestBitArray_Equal(t *testing.T) {
	b := New(80)

	b.Inverse()
	all := b.All()

	if !b.Equal(all) {
		t.Errorf("Should be equal \n%v \n%v", b, all)
	}

}
