package mapbox

type Coordinate struct {
	Lat float64 `json:"latitude"`
	Lng float64 `json:"longitude"`
}

func (c Coordinate) IsZero() bool {
	return c.Lat == 0 && c.Lng == 0
}

type Geometry struct {
	Coordinates []float64 `json:"coordinates"` // [lng, lat]
}

func (g Geometry) Latitude() float64 {
	return g.Coordinates[0]
}

func (g Geometry) Longitude() float64 {
	return g.Coordinates[1]
}
