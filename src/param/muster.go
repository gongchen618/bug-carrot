package param

type Muster struct {
	Title  string         `bson:"title" json:"title"`
	People []MusterPerson `bson:"people" json:"people"`
}

type MusterPerson struct {
	Name string `bson:"name" json:"name"`
	QQ   int64  `bson:"qq" json:"qq"`
}
