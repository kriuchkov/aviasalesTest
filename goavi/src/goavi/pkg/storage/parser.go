package storage

import (
	"encoding/xml"
	"time"
)

// FlightDate: кастомная дата для XML файла
type FlightDate struct {
	time.Time
}

func (c *FlightDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse("2006-01-02T1504", v)
	if err != nil {
		return nil
	}
	*c = FlightDate{parse}
	return nil
}

// SearchResponse: поисковой запрос
type SearchResponse struct {
	RequestID   string      `xml:"RequestId"`
	Itineraries []Itinerary `xml:"PricedItineraries>Flights"`
}

// Itinerary: структура, хранит в себе маршруты перелетов и стоимость.
type Itinerary struct {
	Onward  []Flight `xml:"OnwardPricedItinerary>Flights>Flight"`
	Return  []Flight `xml:"ReturnPricedItinerary>Flights>Flight"`
	Pricing *Price   `xml:"Pricing"`
}

// SourceExist: проверяем есть ли в маршруте начальная точка
func (i *Itinerary) IsSource(point string) bool {
	if len(i.Onward) > 0 && i.Onward[0].Source == point {
		return true
	}
	return false

}

// IsDestination: проверяем есть ли в маршруте конечная точка назначения
// при условии, что в массиве данных точка не может быть транзитной.
func (i *Itinerary) IsDestination(point string) bool {
	var exist bool
	for k := range i.Onward {
		if i.Onward[k].Destination == point {
			exist = true
			break
		}
	}
	// Поищем конечную точку в перелетах обратно
	if !exist {
		for k := range i.Return {
			if i.Return[k].Destination == point {
				exist = true
				break
			}
		}
	}
	return exist
}

// Duration: время перелета без учета времени на ожидания
func (i *Itinerary) Duration() int64 {
	var cost int64
	for k := range i.Onward {
		arrival := i.Onward[k].ArrivalTimeStamp.Time
		cost = cost + int64(arrival.Sub(i.Onward[k].DepartureTimeStamp.Time))
	}
	return cost
}

// PriceInt64: стоимость перелета без конвертации валюты
func (i *Itinerary) PriceInt64() int64 {
	if i.Pricing != nil {
		for s := range i.Pricing.ServiceCharges {
			if i.Pricing.ServiceCharges[s].ChargeType == "TotalAmount" {
				return int64(i.Pricing.ServiceCharges[s].Cost * 100)
			}
		}
	}
	return 0
}

type Flight struct {
	Carrier            string
	Class              string
	FlightNumber       string
	Source             string
	Destination        string
	DepartureTimeStamp FlightDate
	ArrivalTimeStamp   FlightDate
	TicketType         string
	NumberOfStops      string
}

type Price struct {
	Currency       string   `xml:"currency,attr"`
	ServiceCharges []Charge `xml:"ServiceCharges"`
}

type Charge struct {
	ChargeType string  `xml:"ChargeType,attr"`
	Type       string  `xml:"type,attr"`
	Cost       float64 `xml:",chardata"`
}
