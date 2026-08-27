package validation

import "fmt"

func Temperature(value float64) error {
	if value < -50 || value > 2500 {
		return fmt.Errorf("temperature %.2f outside range", value)
	}
	return nil
}

func Percentage(value float64) error {
	if value < 0 || value > 100 {
		return fmt.Errorf("percentage %.2f outside range", value)
	}
	return nil
}

func Positive(value int) error {
	if value <= 0 {
		return fmt.Errorf("value must be positive")
	}
	return nil
}
