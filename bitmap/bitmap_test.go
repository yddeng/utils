package bitmap

import "testing"

func TestNewBitmap(t *testing.T) {
	bm := New(8)

	t.Log(bm.Cap(), bm.Len())

	t.Log(bm.Set(0), bm.Set(6), bm.Set(6), bm.Len())

	t.Log(bm.Set(8), bm.Set(9), bm.Len())

	t.Log(bm.Set(2), bm.Clear(2), bm.Clear(3), bm.Len())

	t.Log(bm.String())

	bm2 := New(14)

	t.Log(bm.Cap(), bm2.Len())
	t.Log(bm2.Set(4), bm2.Set(7), bm2.Len(), bm2.String())

	bm.Copy(bm2)

	t.Log(bm.Cap(), bm.Len(), bm.String())
}
