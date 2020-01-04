package project_1

func MergeSort(src []int64) {
	result := make(chan []int64)
	go mergeSort(src, result)

	res := <-result
	copy(src, res)

	close(result)
}

func mergeSort(src []int64, result chan []int64) {
	if len(src) < 2 {
		result <- src
		return
	}

	leftChan := make(chan []int64)
	rightChan := make(chan []int64)
	middle := len(src) / 2

	go mergeSort(src[:middle], leftChan)
	go mergeSort(src[middle:], rightChan)

	leftData := <-leftChan
	rightData := <-rightChan

	close(leftChan)
	close(rightChan)
	result <- merge(leftData, rightData)
}

func merge(left []int64, right []int64) []int64 {
	result := make([]int64, len(left)+len(right))
	leftIndex, rightIndex := 0, 0

	for i := 0; i < cap(result); i++ {
		switch {
		case leftIndex >= len(left):
			result[i] = right[rightIndex]
			rightIndex++
		case rightIndex >= len(right):
			result[i] = left[leftIndex]
			leftIndex++
		case left[leftIndex] < right[rightIndex]:
			result[i] = left[leftIndex]
			leftIndex++
		default:
			result[i] = right[rightIndex]
			rightIndex++
		}
	}

	return result
}
