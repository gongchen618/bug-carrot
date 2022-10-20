package param

type Ballot struct {
	Title        string         `bson:"title" json:"title"`
	Remark       string         `bson:"remark" json:"remark"`
	TargetMember []BallotMember `bson:"target_member" json:"target_member"`
}

type BallotMember struct {
	People       PersonWithQQ `bson:"people" json:"people"`
	AnsweredFlag bool         `bson:"answered_flag" json:"answered_flag"`
	Answer       string       `bson:"answer" json:"answer"`
}
