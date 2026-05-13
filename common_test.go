package main

/* run with: watch go test ./...  */

import(
	"testing"
)


func TestEmpty(t *testing.T) {

	are := func(ß func(int16, int16) int16, expected, left, right int16, msg string) {
		actual := ß(left, right)
		if expected != actual {
			t.Fatalf("%s: %d!=%d.", msg, expected, actual)
		}
	}

	are(Minimum, 1, 1, 5, "min 1, 1 first")
	are(Minimum, 1, 5, 1, "min 1, 5 first")
	are(Minimum, -1, -1, 5, "min -1, -1 first")
	are(Minimum, -1, 5, -1, "min -1, 5 first")
	are(Minimum, -5, -1, -5, "min -5, all negative")
	are(Minimum, 0, 0, 0, "min 0, everything")

	are(Maximum, 5, 1, 5, "max 5, 1 first")
	are(Maximum, 5, 5, 1, "max 5, 5 first")
	are(Maximum, 5, -1, 5, "max 5, -1 first")
	are(Maximum, 5, 5, -1, "max 5, 5 first, -1")
	are(Maximum, -1, -1, -5, "max -1, all negative")
	are(Maximum, 0, 0, 0, "max 0, everything")

}


