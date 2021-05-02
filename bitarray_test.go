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

func TestBitArray_None(t *testing.T) {
	b := New(19 * 19)
	b.Inverse()

	b.None()
	for i := 0; i < b.Size(); i++ {
		if b.Is(i) {
			t.Errorf("%d pos should not be 1", i)
		}
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
		b.Set(pos)
	}
	b.ShiftLeft(12)
	for pos := 0; pos < len(positions); pos++ {
		if !b.Is(pos + 12) {
			t.Errorf("Error in position: %d", pos+12)
		}
	}

	b.ShiftRight(12)
	for pos := 0; pos < len(positions); pos++ {
		if !b.Is(pos) {
			t.Errorf("Error in position: %d", pos)
		}
	}
}

func TestBitArray_ShiftLeft1(t *testing.T) {
	actual := New(8).Set(0).ShiftLeft(1)
	expected := New(8).Set(1)

	if !expected.Equal(actual) {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestBitArray_ShiftRight1(t *testing.T) {
	actual := New(8).Set(1).ShiftRight(1)
	expected := New(8).Set(0)

	if !expected.Equal(actual) {
		t.Errorf("Expected %v got %v", expected, actual)
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

func TestBitArray_DeMorgansLaw(t *testing.T) {
	size := 514
	a := New(size)
	b := New(size)
	leftSide := a.Clone().Or(b).Inverse()
	rightSide := a.Clone().Inverse().And(b.Clone().Inverse())

	if !leftSide.Equal(rightSide) {
		t.Errorf("should be equal")
	}
}

func TestBitArray_XorGate(t *testing.T) {
	size := 851
	a := New(size)
	b := New(size)
	xor := a.Clone().Or(b).And(a.Clone().Inverse().Or(a.Clone().Inverse().Or(b.Clone().Inverse())))

	if !xor.Equal(a.Xor(b)) {
		t.Errorf("should be equal")
	}
}

func TestBitArray_Minus(t *testing.T) {
	size := 851
	a := New(size)
	b := New(size)
	minus := a.Clone().And(b.Inverse())
	if !minus.Equal(a.Minus(b)) {
		t.Errorf("should be equal")
	}
}

func TestBitArray_SetPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	size := 12
	a := New(size)
	a.Set(12)
}

func TestBitArray_GetPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	size := 12
	a := New(size)
	a.Get(12)
}

func TestBitArray_ShifLeftPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	size := 128
	a := New(size)
	a.ShiftLeft(65)
}

func TestBitArray_ShifRightPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	size := 128
	a := New(size)
	a.ShiftRight(65)
}

func TestBitArray_String(t *testing.T) {
	size := 65
	a := New(size).Inverse().String()
	s :=
		"1111111111111111111111111111111111111111111111111111111111111111\n" +
			"0000000000000000000000000000000000000000000000000000000000000001\n"
	if a != s {
		t.Errorf("%s", a)
	}
}

func TestBitArray_NotEqualSize(t *testing.T) {
	a := New(2)
	b := New(1)

	if a.Equal(b) {
		t.Errorf("Should not be equal. Have different size")
	}
}

func TestBitArray_NotEqual(t *testing.T) {
	a := New(2).Set(1)
	b := New(2)

	if a.Equal(b) {
		t.Errorf("Should not be equal")
	}
}
