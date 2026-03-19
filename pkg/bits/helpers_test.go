package bits

import "testing"

func TestAt(t *testing.T) {
	t.Run("invalid index", func(t *testing.T) {
		defer func() {
			err := recover()
			if err == nil {
				t.Error("failed to panic on invalid index")
			}
		}()

		At(0, -1)

		t.Error("failed to panic on invalid index")
	})

	t.Run("index ing", func(t *testing.T) {
		for i := range 7 {
			var b byte = 1 << i

			bit := At(b, i)
			if bit != 1 {
				t.Errorf("failed to return correct bit at idx %d: got %d but wanted 1", i, bit)
			}
		}
	})
}

func TestFlip(t *testing.T) {
	t.Run("invalid index", func(t *testing.T) {
		defer func() {
			err := recover()
			if err == nil {
				t.Error("failed to panic on invalid index")
			}
		}()

		Flip(0, -1)

		t.Error("failed to panic on invalid index")
	})

	t.Run("flip to 0", func(t *testing.T) {
		for i := range 7 {
			var b byte = 1 << i

			bit := Flip(b, i)
			if bit != 0 {
				t.Errorf("failed to return correct bit at idx %d: got %d but wanted 0", i, bit)
			}
		}
	})

	t.Run("flip to 1", func(t *testing.T) {
		for i := range 7 {
			bit := Flip(0, i)
			if bit != 1<<i {
				t.Errorf("failed to return correct bit at idx %d: got %d but wanted %d", i, bit, 1<<i)
			}
		}
	})
}
