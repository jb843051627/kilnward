package metrics

import "time"

type Point struct {
	At    time.Time
	Value float64
}

func Average(points []Point) float64 {
	if len(points) == 0 {
		return 0
	}
	var total float64
	for _, point := range points {
		total += point.Value
	}
	return total / float64(len(points))
}

func Range(points []Point) (float64, float64) {
	if len(points) == 0 {
		return 0, 0
	}
	min, max := points[0].Value, points[0].Value
	for _, point := range points[1:] {
		if point.Value < min {
			min = point.Value
		}
		if point.Value > max {
			max = point.Value
		}
	}
	return min, max
}
