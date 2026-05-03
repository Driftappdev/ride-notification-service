package domain

type Priority int

const (
	PriorityLow Priority = iota + 1
	PriorityNormal
	PriorityHigh
	PriorityCritical
)

func (p Priority) Weight() int {
	switch p {
	case PriorityCritical:
		return 1000
	case PriorityHigh:
		return 500
	case PriorityNormal:
		return 100
	case PriorityLow:
		return 10
	default:
		return 100
	}
}
