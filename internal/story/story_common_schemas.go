package story

type DrawingPoint struct {
	X float32 `bson:"x" json:"x"`
	Y float32 `bson:"y" json:"y"`
}

type CaertsianCoordinateColor struct {
	Start DrawingPoint `bson:"start" json:"start"`
	End   DrawingPoint `bson:"end"   json:"end"`
	Color string       `bson:"color" json:"color"`
}
