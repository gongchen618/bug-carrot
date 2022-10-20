package param

type Muster struct {
	Title  string         `bson:"title" json:"title"`
	People []PersonWithQQ `bson:"people" json:"people"`
}
