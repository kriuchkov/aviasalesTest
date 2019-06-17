package storage

import (
	"encoding/xml"
)

type Storage struct {
	Data []*Itinerary
}

func NewStorage() *Storage {
	return &Storage{Data: []*Itinerary{}}
}

// LoadXML: Загружаем XML файл в Storage
func (s *Storage) LoadXML(bodyBytes []byte) error {
	var search SearchResponse
	err := xml.Unmarshal(bodyBytes, &search)
	for i := range search.Itineraries {
		itinerary := search.Itineraries[i]
		s.Data = append(s.Data, &itinerary)
	}
	return err
}

// GetItinerary: список маршрутов
func (s *Storage) GetItinerary(src, dest string, ret bool) []*Itinerary {
	var result []*Itinerary
	for i := range s.Data {
		if s.Data[i].IsDestination(dest) && s.Data[i].IsSource(src) {
			if !ret || (ret && s.Data[i].Return != nil) {
				result = append(result, s.Data[i])
			}
		}
	}
	return result
}

// OptimalItinerary: добавляем маршруты в очередь, в зависимости от интерфейса
// очередь будет отсортирована в нужной нам последовательности.
func (s *Storage) OptimalItinerary(itinerary []*Itinerary, queue StorageList) {
	for i := range itinerary {
		queue.PushOrdered(itinerary[i])
	}
}
