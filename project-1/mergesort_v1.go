package project_1

func MergeSortV1(src []int64) {
	mergeSortV1(src, 0, int64(len(src)))
}
func mergeSortV1(src []int64, lo int64, hi int64) {
	if hi-lo < 2 {
		return
	}
	var mi = (lo + hi) / 2
	mergeSortV1(src, lo, mi)
	mergeSortV1(src, mi, hi)
	if src[mi-1] > src[mi] {
		mergeV1(src, lo, mi, hi)
	}
}

func mergeV1(src []int64, lo int64, mi int64, hi int64) {
	a := src[lo:]
	var lb = mi - lo
	b := make([]int64, lb)
	var i, j, k int64
	for i = 0; i < lb; i++ {
		b[i] = a[i]
	}
	var lc = hi - mi
	c := src[mi:]
	for i, j, k = 0, 0, 0; (j < lb) || (k < lc); {
		if (j < lb) && (!(k < lc) || (b[j] <= c[k])) {
			a[i] = b[j]
			i++
			j++
		}
		if (k < lc) && (!(j < lb) || c[k] < b[j]) {
			a[i] = c[k]
			i++
			k++
		}
	}

}