package consumer

type Approach string

const (
	MultipleFilterSubjects = Approach("multiple-filter-subjects")
	ManyConsumers          = Approach("many-consumers")
	WildCard               = Approach("wild-card")
)

var approaches []Approach

func init() {
	approaches = []Approach{MultipleFilterSubjects, ManyConsumers, WildCard}
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
	panic("unknown approach")
}
