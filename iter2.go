package iter2

import "iter"

// Zip returns an iter.Seq2 that pairs corresponding elements from iter1 and iter2.
// Iteration stops when either of the seq1 or seq2 stops.
func Zip[T1, T2 any](seq1 iter.Seq[T1], seq2 iter.Seq[T2]) iter.Seq2[T1, T2] {
	return func(yield func(T1, T2) bool) {
		next1, stop1 := iter.Pull(seq1)
		defer stop1()
		next2, stop2 := iter.Pull(seq2)
		defer stop2()
		for {
			v1, ok1 := next1()
			if !ok1 {
				return
			}
			v2, ok2 := next2()
			if !ok2 {
				return
			}
			if !yield(v1, v2) {
				return
			}
		}
	}
}

// Concat returns the concation of seqs.
func Concat[T any](seqs ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		var ok bool
		yield2 := func(v T) bool {
			ok = yield(v)
			return ok
		}
		for _, seq := range seqs {
			seq(yield2)
			if !ok {
				return
			}
		}
	}
}
