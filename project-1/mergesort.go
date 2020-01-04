package project_1

import "sync"

func MergeSort(src []int64) {
	mergeSort(src, 0, int64(len(src)))
}
func mergeSort(src []int64, lo int64, hi int64) {
	if hi-lo < 2 {
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	var mi = (lo + hi) / 2
	go func() {
		defer wg.Done()
		mergeSort(src, lo, mi)
	}()

	mergeSort(src, mi, hi)

	wg.Wait()
	if src[mi-1] > src[mi] {
		merge(src, lo, mi, hi)
	}
}

func merge(src []int64, lo int64, mi int64, hi int64) {
	a := src[lo:]
	var lb = mi - lo
	b := make([]int64, lb)
	var i, j, k int64
	for i = 0; i < lb; i++ {
		b[i] = a[i]
	}
	var lc = hi - mi
	c := src[mi:]
	for i, j, k = 0, 0, 0; j < lb; {
		if k < lc && c[k] < b[j] {
			a[i] = c[k]
			i++
			k++
		}

		if lc <= k || b[j] <= c[k] {
			a[i] = b[j]
			i++
			j++
		}
	}
}
