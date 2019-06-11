package bumble

import "math"

const (
	earthRadius = 3958.8
	degToRad    = math.Pi / 180
)

type Location struct {
	Name        string
	Lat         float64
	Lon         float64
	CountryCode string
}

// Distance returns the distance (in miles) between two
// locations on Earth.
func (l *Location) Distance(lat, lon float64) float64 {
	x1, y1, z1 := latLonToXYZ(lat, lon)
	x2, y2, z2 := latLonToXYZ(l.Lat, l.Lon)
	cosTheta := x1*x2 + y1*y2 + z1*z2
	return earthRadius * math.Acos(cosTheta)
}

func latLonToXYZ(lat, lon float64) (x, y, z float64) {
	y = math.Sin(lat * degToRad)
	x = math.Sin(lon*degToRad) * math.Cos(lat*degToRad)
	z = math.Sqrt(1 - x*x - y*y)
	return
}
