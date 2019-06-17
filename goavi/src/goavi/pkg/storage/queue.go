package storage

import (
	"sort"
)

type Item struct {
	value *Itinerary
	cost  int64
}

type queueBase []*Item

func (pq queueBase) Len() int {
	return len(pq)
}

// Swap: меняем места
func (pq queueBase) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

// setTimeCost: устанавливаем стоимость приоритета по времени
func (pq queueBase) setTimeCost(i int) {
	if pq[i].cost == 0 {
		pq[i].cost = pq[i].value.Duration()
	}
}

// setPriceCost: устанавливаем стоимость приоритета по цене
func (pq queueBase) setPriceCost(i int) {
	if pq[i].cost == 0 {
		pq[i].cost = pq[i].value.PriceInt64()
	}
}

func (pq *queueBase) Push(i *Itinerary) {
	*pq = append(*pq, &Item{value: i})
}

func (pq *queueBase) Pop() *Itinerary {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item.value
}

type StorageList interface {
	PushOrdered(*Itinerary)
	PopOrdered() *Itinerary
	Len() int
}

type queueInterface interface {
	sort.Interface
	Push(x *Itinerary)
	Pop() *Itinerary
}

type queueWrapper struct {
	queueInterface
}

func (pq *queueWrapper) PushOrdered(v *Itinerary) {
	pq.Push(v)
	pq.up(pq.Len() - 1)
}

func (pq *queueWrapper) PopOrdered() *Itinerary {
	n := pq.Len() - 1
	pq.Swap(0, n)
	pq.down(0, n)
	return pq.Pop()
}

// container/heap
func (pq *queueWrapper) fix(i int) {
	if !pq.down(i, pq.Len()) {
		pq.up(i)
	}
}

func (pq *queueWrapper) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !pq.Less(j, i) {
			break
		}
		pq.Swap(i, j)
		j = i
	}
}

func (pq *queueWrapper) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 {
			break
		}
		j := j1
		if j2 := j1 + 1; j2 < n && pq.Less(j2, j1) {
			j = j2
		}
		if !pq.Less(j, i) {
			break
		}
		pq.Swap(i, j)
		i = j
	}
	return i > i0
}

// самый долгий маршрут
type TimeQueueMax struct {
	queueBase
}

// сравниваем значения
func (pq TimeQueueMax) Less(i, j int) bool {
	pq.setTimeCost(i)
	pq.setTimeCost(j)
	return pq.queueBase[i].cost > pq.queueBase[j].cost
}

func NewTimeQueueMax() StorageList {
	l := &queueWrapper{new(TimeQueueMax)}
	n := l.Len()
	for i := n/2 - 1; i >= 0; i-- {
		l.down(i, n)
	}
	return l
}

// самый быстрый маршрут
type TimeQueueMin struct {
	queueBase
}

func (pq TimeQueueMin) Less(i, j int) bool {
	pq.setTimeCost(i)
	pq.setTimeCost(j)
	return pq.queueBase[i].cost < pq.queueBase[j].cost
}

func NewTimeQueueMin() StorageList {
	l := &queueWrapper{new(TimeQueueMin)}
	n := l.Len()
	for i := n/2 - 1; i >= 0; i-- {
		l.down(i, n)
	}
	return l
}

// самый дорогой маршрут
type PriceQueueMax struct {
	queueBase
}

func (pq PriceQueueMax) Less(i, j int) bool {
	pq.setPriceCost(i)
	pq.setPriceCost(j)
	return pq.queueBase[i].cost > pq.queueBase[j].cost
}

func NewPriceQueueMax() StorageList {
	l := &queueWrapper{new(PriceQueueMax)}
	n := l.Len()
	for i := n/2 - 1; i >= 0; i-- {
		l.down(i, n)
	}
	return l
}

// самый дешевый маршрут
type PriceQueueMin struct {
	queueBase
}

func (pq PriceQueueMin) Less(i, j int) bool {
	pq.setPriceCost(i)
	pq.setPriceCost(j)
	return pq.queueBase[i].cost < pq.queueBase[j].cost
}

func NewPriceQueueMin() StorageList {
	l := &queueWrapper{new(PriceQueueMin)}
	n := l.Len()
	for i := n/2 - 1; i >= 0; i-- {
		l.down(i, n)
	}
	return l
}

// оптимальный маршрут :D
type OptimalQueue struct {
	queueBase
}

func (pq OptimalQueue) Less(i, j int) bool {
	if pq.queueBase[i].cost == 0 {
		pq.queueBase[i].cost = (pq.queueBase[i].value.PriceInt64() * 2) - pq.queueBase[i].value.Duration()
	}
	if pq.queueBase[j].cost == 0 {
		pq.queueBase[j].cost = (pq.queueBase[j].value.PriceInt64() * 2) - pq.queueBase[j].value.Duration()
	}
	return pq.queueBase[i].cost < pq.queueBase[j].cost
}

func NewOptimalQueue() StorageList {
	l := &queueWrapper{new(OptimalQueue)}
	n := l.Len()
	for i := n/2 - 1; i >= 0; i-- {
		l.down(i, n)
	}
	return l
}
