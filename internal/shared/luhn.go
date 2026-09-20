package shared

import "errors"

var ErrInvalidSKUID = errors.New("invalid sku id: must be numeric and satisfy luhn checksum")

func CalculateLuhn(num int64) int64 {
	var sum int64
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

func VerifyLuhn(num int64, check_digit int64) error {
	if check_digit != CalculateLuhn(num/10) {
		return ErrInvalidSKUID
	}
	return nil
}
