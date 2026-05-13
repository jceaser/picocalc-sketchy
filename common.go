package main


func Between(low, value, high int16) int16 {
	return Maximum(low, Minimum(value, high))
}

func Maximum(left, right int16) int16 {
	if left > right {
		return left
	}
	return right
}

func Minimum(left, right int16) int16 {
	if left < right {
		return left
	}
	return right
}
