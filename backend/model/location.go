package model

type Location struct {
	Longitude float64 `bson:"longitude" json:"longitude"`
	Latitude  float64 `bson:"latitude" json:"latitude"`
	Address   string  `bson:"address" json:"address"`
	TimeStamp int64   `bson:"timestamp" json:"timestamp"`
	Speed     string  `bson:"speed" json:"speed"`
	Datestamp int64   `bson:"datestamp" json:"datestamp"`
}
