package braillify

import (
	"testing"
)

const brailleBaseCodepoint rune = '\u2800'

func TestBrailleMap(t *testing.T) {
	// From wikipedia - https://en.wikipedia.org/wiki/Braille_Patterns
	binaryMap := map[int]rune{
		0: 0x01,
		1: 0x02,
		2: 0x04,
		3: 0x40,
		4: 0x08,
		5: 0x10,
		6: 0x20,
		7: 0x80,
	}
	for i := 0; i < 256; i++ {
		have := brailleMap[uint8(i)]
		want := brailleBaseCodepoint
		for j := 0; j < 8; j++ {
			if (i & (1 << uint(j))) != 0 {
				want += binaryMap[j]
			}
		}
		if want != have {
			t.Fatalf("brailleMap contains '\\u%04x' for '0x%02x' want '\\u%04x'", have, i, want)
		}
	}
}
