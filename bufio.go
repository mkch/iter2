package iter2

import (
	"bufio"
	"iter"
)

// Tokens returns an iterator over tokens scanned by the scanner.
func Tokens(scanner *bufio.Scanner) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for scanner.Scan() {
			if !yield(scanner.Bytes()) {
				return
			}
		}
	}
}

// TokenStrings returns an iterator over token strings scanned by the scanner.
func TokenStrings(scanner *bufio.Scanner) iter.Seq[string] {
	return func(yield func(string) bool) {
		for scanner.Scan() {
			if !yield(scanner.Text()) {
				return
			}
		}
	}
}
