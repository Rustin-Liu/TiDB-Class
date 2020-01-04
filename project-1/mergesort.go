package project_1

import "sync"

func MergeSort(src []int64) {
	semaphore := make(chan struct{}, 5)
	parallelMergeSort(src, 0, int64(len(src)), semaphore)
}

func parallelMergeSort(src []int64, lo int64, hi int64, semaphore chan struct{}) {
	if hi-lo < 2 {
		return
	}

	mi := (lo + hi) / 2

	wg := sync.WaitGroup{}
	wg.Add(2)

	select {
	case semaphore <- struct{}{}:
		go func() {
			parallelMergeSort(src, lo, mi, semaphore)
			<-semaphore
			wg.Done()
		}()
	default:
		mergeSort(src, lo, mi)
		wg.Done()
	}

	select {
	case semaphore <- struct{}{}:
		go func() {
			parallelMergeSort(src, mi, hi, semaphore)
			<-semaphore
			wg.Done()
		}()
	default:
		mergeSort(src, mi, hi)
		wg.Done()
	}

	wg.Wait()
	merge(src, lo, mi, hi)
}

func mergeSort(src []int64, lo int64, hi int64) {
	if hi-lo < 2 {
		return
	}
	mi := (lo + hi) / 2
	mergeSort(src, lo, mi)

	mergeSort(src, mi, hi)
	if src[mi-1] > src[mi] {
		merge(src, lo, mi, hi)
	}
}

func merge(src []int64, lo int64, mi int64, hi int64) {
	a := src[lo:]
	var lb = mi - lo
	var i, j, k int64

	b := make([]int64, lb)
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
