package skus

func CalculateLuhn(num int) int {
	var sum int
	double := true

	for i := num; i > 0; i /= 10 {
		d := (i % 10)
		if double {
			d *= 2
			if d >= 10 {
				d -= 9
			}
		}
		sum += d
		double = !double
	}

	return (10 - (sum % 10)) % 10
}

func VerifyLuhn(num int, check_digit int) bool {
	return check_digit == CalculateLuhn(num/10)
}
