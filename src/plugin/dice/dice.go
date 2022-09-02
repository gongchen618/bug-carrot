package dice

import "fmt"

type Dice struct {
}

func (p *Dice) DiceGroupMember(a string) {
	fmt.Println("OK ", a)
}

func CallPlugin() Dice {
	p := Dice{}
	return p
}
