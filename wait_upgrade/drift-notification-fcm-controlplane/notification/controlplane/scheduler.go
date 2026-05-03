package controlplane

import (
	"container/heap"
	"time"

	"github.com/driftappdev/libpackage/notification/domain"
)

type ScheduledItem struct {
	Notification domain.Notification
	Score        int64
	index        int
}

type PriorityQueue []*ScheduledItem

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	if pq[i].Score == pq[j].Score {
		return pq[i].Notification.Priority.Weight() > pq[j].Notification.Priority.Weight()
	}
	return pq[i].Score < pq[j].Score
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index, pq[j].index = i, j
}
func (pq *PriorityQueue) Push(x any) {
	item := x.(*ScheduledItem)
	item.index = len(*pq)
	*pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[:n-1]
	return item
}

func Score(n domain.Notification) int64 {
	now := time.Now().UTC()
	due := now
	if n.ScheduledAt != nil { due = n.ScheduledAt.UTC() }

	// Lower score = earlier execution.
	// Priority subtracts weight so critical jumps ahead when due time is close.
	return due.UnixMilli() - int64(n.Priority.Weight())
}

func NewQueue(items []domain.Notification) PriorityQueue {
	pq := PriorityQueue{}
	heap.Init(&pq)
	for _, n := range items {
		heap.Push(&pq, &ScheduledItem{Notification: n, Score: Score(n)})
	}
	return pq
}
