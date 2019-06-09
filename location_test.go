package bumble

import (
	"math"
	"testing"
)

const (
	hkToNZ = 5845.0
	hkToCA = 7061.0
)

func TestLocationDistance(t *testing.T) {
	hongKong := &Location{Name: "Hong Kong", Lat: 22.3193, Lon: 114.1694}
	newZealand := &Location{Name: "New Zealand", Lat: -40.9006, Lon: 174.8860}
	california := &Location{Name: "California", Lat: 36.7783, Lon: -119.4179}

	actualDist := hongKong.Distance(newZealand.Lat, newZealand.Lon)
	if math.Abs(actualDist-hkToNZ) > 10 {
		t.Errorf("expected HK to NZ distance %f but got %f", hkToNZ, actualDist)
	}
	actualDist = hongKong.Distance(california.Lat, california.Lon)
	if math.Abs(actualDist-hkToCA) > 10 {
		t.Errorf("expected HK to CA distance %f but got %f", hkToCA, actualDist)
	}
}
