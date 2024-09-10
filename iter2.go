package iter2

import (
	"io/fs"
	"iter"
	"sync"
)

// Zip returns an iter.Seq that pairs corresponding elements from iter1 and iter2.
// Iteration stops when either of the seq1 or seq2 stops.
func Zip[T1, T2, Pair any](seq1 iter.Seq[T1], seq2 iter.Seq[T2], pair func(v1 T1, v2 T2) Pair) iter.Seq[Pair] {
	return func(yield func(Pair) bool) {
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
			if !yield(pair(v1, v2)) {
				return
			}
		}
	}
}

// Zip2 returns an iter.Seq2 that pairs corresponding elements from iter1 and iter2.
// Iteration stops when either of the seq1 or seq2 stops.
func Zip2[T1, T2 any](seq1 iter.Seq[T1], seq2 iter.Seq[T2]) iter.Seq2[T1, T2] {
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
// Concat yields the values from seqs without interleaving them.
func Concat[T any](seqs ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, seq := range seqs {
			for v := range seq {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Merge combines seqs into one by merging their values.
// Merge may interleave the values yield by the merged Seq.
// A similar func [Concat] does not interleave values, but
// yields all of each source Seq's values in turn before beginning
// to yield values from the next source Seq.
func Merge[T any](seqs ...iter.Seq[T]) iter.Seq[T] {
	var n = len(seqs)
	if n == 0 {
		return func(yield func(T) bool) {}
	}
	return func(yield func(T) bool) {
		doneR := make(chan struct{}) // done reading
		doneW := make(chan struct{}) // done writing
		ch := make(chan T)
		wg := &sync.WaitGroup{}
		wg.Add(n)
		for _, seq := range seqs {
			go func() {
				defer wg.Done()
				for v := range seq {
					select {
					case ch <- v:
					case <-doneR:
						return
					}
				}
			}()
		}
		go func() {
			wg.Wait()
			close(doneW)
		}()

		for {
			select {
			case <-doneW:
				return
			case v := <-ch:
				if !yield(v) {
					// early stop
					close(doneR)
					wg.Wait()
					return
				}
			}
		}
	}

}

// Map returns an iter.Seq that contains a sequence transformed form seq by func f.
func Map[T1, T2 any](seq iter.Seq[T1], f func(T1) T2) iter.Seq[T2] {
	return func(yield func(T2) bool) {
		for v := range seq {
			if !yield(f(v)) {
				return
			}
		}
	}
}

// Map returns an iter.Seq2 that contains a sequence transformed form seq by func f.
func Map2[K1, V1, K2, V2 any](seq iter.Seq2[K1, V1], f func(K1, V1) (K2, V2)) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k, v := range seq {
			if !yield(f(k, v)) {
				return
			}
		}
	}
}

// Map1To2 returns an iter.Seq2 that contains a sequence transformed form seq by func f.
func Map1To2[T, K, V any](seq iter.Seq[T], f func(v T) (K, V)) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for v := range seq {
			if !yield(f(v)) {
				return
			}
		}
	}
}

// Keys returns an iterator over keys in seq2.
func Keys[K, V any](seq2 iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k, _ := range seq2 {
			if !yield(k) {
				return
			}
		}
	}
}

// Keys returns an iterator over values in seq2.
func Values[K, V any](seq2 iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range seq2 {
			if !yield(v) {
				return
			}
		}
	}
}

// Take returns an iterator that yields the first n values in seq.
// Take panics if n < 0.
func Take[T any](seq iter.Seq[T], n int) iter.Seq[T] {
	if n < 0 {
		panic("negative count")
	} else if n == 0 {
		return func(yield func(T) bool) {}
	}
	return func(yield func(T) bool) {
		var count = 0
		for v := range seq {
			if !yield(v) {
				return
			}
			count++
			if count >= n {
				return
			}
		}
	}
}

// Take returns an iterator that yields the first n values in seq2.
// Take panics if n < 0.
func Take2[K, V any](seq2 iter.Seq2[K, V], n int) iter.Seq2[K, V] {
	if n < 0 {
		panic("negative count")
	} else if n == 0 {
		return func(yield func(K, V) bool) {}
	}
	return func(yield func(K, V) bool) {
		var count = 0
		for k, v := range seq2 {
			if !yield(k, v) {
				return
			}
			count++
			if count >= n {
				return
			}
		}
	}
}

// DirEntry is a file or directory of a file tree.
type DirEntry struct {
	// Path contains the argument to WalkDir as a prefix. That is, if WalkDir is called with root argument "dir"
	// and finds a file named "a" in that directory, the Path of yielded DirEntry is "dir/a".
	Path string
	// Entry is the [fs.DirEntry] for the named path.
	Entry fs.DirEntry

	err error
}

// SkipDir skips the current directory (path if d.IsDir() is true, otherwise path's parent directory).
func (dir *DirEntry) SkipDir() {
	dir.err = fs.SkipDir
}

// SkipAll skips all remaining files and directories.
func (dir *DirEntry) SkipAll() {
	dir.err = fs.SkipAll
}

// WalkDir returns an iterator over the file tree rooted at root.
func WalkDir(fsys fs.FS, root string) iter.Seq2[*DirEntry, error] {
	return func(yield func(*DirEntry, error) bool) {
		fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
			dir := &DirEntry{Path: path, Entry: d}
			if !yield(dir, err) {
				dir.err = fs.SkipAll // early stop. skip all.
			}
			return dir.err
		})
	}

}

// Push creates an iterator whose values are yielded by function calls.
// Calling yield pushes the next value onto the sequence, stopping early if yield returns false.
// Stop ends the iteration. It must be called when the caller has no next value to push.
// It is valid to call stop multiple times. Typically, callers should “defer stop()”.
// It is safe to call yield and stop from multiple goroutines simultaneously.
//
// Push is useful when yielding values out of a loop.
func Push[T any]() (seq iter.Seq[T], yield func(T) bool, stop func()) {
	var ch = make(chan T)
	var doneW = make(chan struct{})
	var doneR = make(chan struct{})
	seq = func(yield func(T) bool) {
		defer close(doneR)
		for {
			select {
			case v := <-ch:
				if !yield(v) {
					return
				}
			case <-doneW:
				return
			}
		}
	}
	yield = func(v T) bool {
		select {
		case ch <- v:
			return true
		case <-doneR:
			return false
		}
	}
	var stopLock sync.Mutex
	stop = func() {
		stopLock.Lock()
		defer stopLock.Unlock()

		select {
		case <-doneW:
			return
		default:
			close(doneW)
		}
	}
	return
}

// Push2 creates an iterator whose values are yielded by function calls.
// Push2 works the same way as [Push], except for the type parameters.
func Push2[K, V any]() (seq2 iter.Seq2[K, V], yield func(K, V) bool, stop func()) {
	type pair struct {
		K K
		V V
	}
	var ch = make(chan pair)
	var doneW = make(chan struct{})
	var doneR = make(chan struct{})
	seq2 = func(yield func(K, V) bool) {
		defer close(doneR)
		for {
			select {
			case pair := <-ch:
				if !yield(pair.K, pair.V) {
					return
				}
			case <-doneW:
				return
			}
		}
	}
	yield = func(k K, v V) bool {
		select {
		case ch <- pair{k, v}:
			return true
		case <-doneR:
			return false
		}
	}
	var stopLock sync.Mutex
	stop = func() {
		stopLock.Lock()
		defer stopLock.Unlock()

		select {
		case <-doneW:
			return
		default:
			close(doneW)
		}
	}
	return
}

// Filter returns an iterator over the sequence of elements in seq that pass the test.
func Filter[T any](seq iter.Seq[T], test func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range seq {
			if !test(v) {
				continue
			}
			if !yield(v) {
				return
			}
		}
	}
}

// Filter returns an iterator over the sequence of elements in seq that pass the test.
func Filter2[K, V any](seq iter.Seq2[K, V], test func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			if !test(k, v) {
				continue
			}
			if !yield(k, v) {
				return
			}
		}
	}
}
