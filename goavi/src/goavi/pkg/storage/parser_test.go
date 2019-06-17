package storage

import (
	"encoding/xml"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var xmlFile string = `
<AirFareSearchResponse RequestTime="28-09-2015 20:30:19" ResponseTime="28-09-2015 20:30:26">
   <RequestId>123ABCD</RequestId>
   <PricedItineraries>
      <Flights>
         <OnwardPricedItinerary>
            <Flights>
               <Flight>
                  <Carrier id="AI">AirIndia</Carrier>
                  <FlightNumber>996</FlightNumber>
                  <Source>DXB</Source>
                  <Destination>DEL</Destination>
                  <DepartureTimeStamp>2018-10-27T0005</DepartureTimeStamp>
                  <ArrivalTimeStamp>2018-10-27T0445</ArrivalTimeStamp>
                  <Class>G</Class>
                  <NumberOfStops>0</NumberOfStops>
                  <FareBasis>2820303decf751-5511-447a-aeb1-810a6b10ad7d@@$255_DXB_DEL_996_9_00:05_$255_DEL_BKK_332_9_13:25__A2_1_1</FareBasis>
                  <WarningText />
                  <TicketType>E</TicketType>
               </Flight>
               <Flight>
                  <Carrier id="AI">AirIndia</Carrier>
                  <FlightNumber>332</FlightNumber>
                  <Source>DEL</Source>
                  <Destination>BKK</Destination>
                  <DepartureTimeStamp>2018-10-27T1325</DepartureTimeStamp>
                  <ArrivalTimeStamp>2018-10-27T1920</ArrivalTimeStamp>
                  <Class>G</Class>
                  <NumberOfStops>0</NumberOfStops>
                  <FareBasis>2820303decf751-5511-447a-aeb1-810a6b10ad7d@@$255_DXB_DEL_996_9_00:05_$255_DEL_BKK_332_9_13:25__A2_1_1</FareBasis>
                  <WarningText />
                  <TicketType>E</TicketType>
               </Flight>
            </Flights>
         </OnwardPricedItinerary>
         <Pricing currency="SGD">
            <ServiceCharges type="SingleAdult" ChargeType="BaseFare">167.00</ServiceCharges>
            <ServiceCharges type="SingleAdult" ChargeType="AirlineTaxes">215.70</ServiceCharges>
            <ServiceCharges type="SingleAdult" ChargeType="TotalAmount">382.70</ServiceCharges>
            <ServiceCharges type="SingleChild" ChargeType="BaseFare">129.00</ServiceCharges>
            <ServiceCharges type="SingleChild" ChargeType="AirlineTaxes">215.70</ServiceCharges>
            <ServiceCharges type="SingleChild" ChargeType="TotalAmount">344.70</ServiceCharges>
            <ServiceCharges type="SingleInfant" ChargeType="BaseFare">20.00</ServiceCharges>
            <ServiceCharges type="SingleInfant" ChargeType="TotalAmount">20.00</ServiceCharges>
         </Pricing>
      </Flights>
      <Flights>
         <OnwardPricedItinerary>
            <Flights>
               <Flight>
                  <Carrier id="CZ">China Southern Airlines</Carrier>
                  <FlightNumber>384</FlightNumber>
                  <Source>DXB</Source>
                  <Destination>CAN</Destination>
                  <DepartureTimeStamp>2018-10-27T0140</DepartureTimeStamp>
                  <ArrivalTimeStamp>2018-10-27T1225</ArrivalTimeStamp>
                  <Class>T</Class>
                  <NumberOfStops>0</NumberOfStops>
                  <FareBasis>2820303decf751-5511-447a-aeb1-810a6b10ad7d@@$255_DXB_CAN_384_107_01:40_$255_CAN_BKK_3035_107_14:50__A2_1_1</FareBasis>
                  <WarningText />
                  <TicketType>E</TicketType>
               </Flight>
               <Flight>
                  <Carrier id="CZ">China Southern Airlines</Carrier>
                  <FlightNumber>3035</FlightNumber>
                  <Source>CAN</Source>
                  <Destination>BKK</Destination>
                  <DepartureTimeStamp>2018-10-27T1450</DepartureTimeStamp>
                  <ArrivalTimeStamp>2018-10-27T1710</ArrivalTimeStamp>
                  <Class>Y</Class>
                  <NumberOfStops>0</NumberOfStops>
                  <FareBasis>2820303decf751-5511-447a-aeb1-810a6b10ad7d@@$255_DXB_CAN_384_107_01:40_$255_CAN_BKK_3035_107_14:50__A2_1_1</FareBasis>
                  <WarningText />
                  <TicketType>E</TicketType>
               </Flight>
            </Flights>
         </OnwardPricedItinerary>
         <Pricing currency="SGD">
            <ServiceCharges type="SingleAdult" ChargeType="BaseFare">233.00</ServiceCharges>
            <ServiceCharges type="SingleAdult" ChargeType="AirlineTaxes">152.40</ServiceCharges>
            <ServiceCharges type="SingleAdult" ChargeType="TotalAmount">385.40</ServiceCharges>
            <ServiceCharges type="SingleChild" ChargeType="BaseFare">233.00</ServiceCharges>
            <ServiceCharges type="SingleChild" ChargeType="AirlineTaxes">132.20</ServiceCharges>
            <ServiceCharges type="SingleChild" ChargeType="TotalAmount">365.20</ServiceCharges>
            <ServiceCharges type="SingleInfant" ChargeType="BaseFare">129.00</ServiceCharges>
            <ServiceCharges type="SingleInfant" ChargeType="AirlineTaxes">11.40</ServiceCharges>
            <ServiceCharges type="SingleInfant" ChargeType="TotalAmount">140.40</ServiceCharges>
         </Pricing>
      </Flights>
   </PricedItineraries>
</AirFareSearchResponse>`

func TestParsingXML(t *testing.T) {
	var search SearchResponse
	xml.Unmarshal([]byte(xmlFile), &search)

	// Проверяем количество элементов
	assert.Equal(t, len(search.Itineraries), 2)

	// Проверяем коррекстрость конкретного маршрута
	itinerary := search.Itineraries[0]
	assert.Equal(t, len(itinerary.Onward), 2)

	// Проверяем коррекстрость перелета
	flight := itinerary.Onward[0]
	assert.Equal(t, flight.Carrier, "AirIndia")
	assert.Equal(t, flight.FlightNumber, "996")
	assert.Equal(t, flight.Source, "DXB")
	assert.Equal(t, flight.Destination, "DEL")
	assert.Equal(t, flight.DepartureTimeStamp.Time, time.Date(2018, 10, 27, 0, 5, 0, 0, time.UTC))
	assert.Equal(t, flight.ArrivalTimeStamp.Time, time.Date(2018, 10, 27, 4, 45, 0, 0, time.UTC))
	assert.Equal(t, flight.Class, "G")

	// Проверяем коррекстрость цен
	price := itinerary.Pricing
	assert.Equal(t, len(price.ServiceCharges), 8)

	charge := price.ServiceCharges[1]
	assert.Equal(t, charge.ChargeType, "AirlineTaxes")
	assert.Equal(t, charge.Cost, 215.7)
}
