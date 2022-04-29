package param

type Food struct {
	Name        string `bson:"name"`
	Address     string `bson:"address"`
	Description string `bson:"description"`
	Recommender int64  `bson:"recommender"`
}
