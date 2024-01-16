package consumer

import "fmt"

type Approach string

const (
	MultipleFilterSubjects = Approach("multiple-filter-subjects")
	ManyConsumers          = Approach("many-consumers")
	Wildcard               = Approach("wildcard")
)

var approaches []Approach

func init() {
	approaches = []Approach{MultipleFilterSubjects, ManyConsumers, Wildcard}
}

func Approaches() []Approach {
	return approaches
}

func NewApproach(name string) Approach {
	for _, a := range approaches {
		if string(a) == name {
			return a
		}
	}
	panic(fmt.Sprintf("unknown approach: '%s'", name))
}
