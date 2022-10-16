package param

type Ballot struct {
	Title          string         `bson:"title" json:"title"`
	OfferedOptions []string       `bson:"offered_options" json:"offered_options"`
	TargetMember   []BallotMember `bson:"target_member" json:"target_member"`
}

type BallotMember struct {
	Info   MusterPerson `bson:"info" json:"info"`
	Option string       `bson:"option" json:"option"`
}
