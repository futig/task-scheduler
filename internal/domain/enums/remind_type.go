package enums

import "fmt"

type RemindType int

const (
	Start RemindType = iota
	End
)

func (r RemindType) ParseToFutureVerb() (string, error) {
	if r == Start {
		return "начнется", nil
	} else if r == End {
		return "закончится", nil
	}
	return "", fmt.Errorf("типа напоминаний %d не существует", r)
}

func (r RemindType) ParseToPresentVerb() (string, error) {
	if r == Start {
		return "началась", nil
	} else if r == End {
		return "закончилась", nil
	}
	return "", fmt.Errorf("типа напоминаний %d не существует", r)
}
