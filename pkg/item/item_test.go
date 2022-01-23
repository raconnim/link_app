package item

import "testing"

func TestGenerateShortLink(t *testing.T) {
	for i := 0; i < 10; i++ {
		str := GenerateShortLink()
		if len(str) != 10 {
			t.Errorf("expected length = 10, got %v", len(str))
			return
		}
	}
}
