from collections import namedtuple
from bs4 import BeautifulSoup, NavigableString
from datetime import datetime

import heapq


# Общее хранилище для данных по маршруту
Storage = list()


class Itinerary(namedtuple('itinerary', ['Onward', 'Return', 'Pricing'], defaults=list())):
    """ Класс с маршрутами и стоимостью """

    def is_source(self, point):
        """
        Проверяем есть ли в маршруте начальная точка
        point: str
        """
        for f in self.Onward:
            if f.Source == point:
                return True
        return False

    def is_destination(self, point):
        """
        Проверяем есть ли в маршруте конечная точка назначения
        point: str
        """
        exist = False
        for f in self.Onward:
            if f.Destination == point:
                exist = True

        if not exist and self.Return:
            for f in self.Onward:
                if f.Destination == point:
                    exist = True
        return exist

    def price(self):
        """ Стоимость маршрута """
        for c in self.Pricing:
            if c.ChargeType == 'TotalAmount' and c.TypeOf == 'SingleAdult':
                return float(c.Value)
        return 0

    def duration(self):
        """ Время перелета без учета времени на ожидания """
        cost = 0
        for f in self.Onward:
            arrival_time = datetime.strptime(f.ArrivalTimeStamp, "%Y-%m-%dT%H%M")
            departure_time = datetime.strptime(f.DepartureTimeStamp, "%Y-%m-%dT%H%M")
            if arrival_time and departure_time:
                cost = cost + (arrival_time - departure_time).total_seconds()
        return cost


Charges = namedtuple('Charges', ['TypeOf', 'ChargeType', 'Value'])


def create_charges(*args):
    return Charges._make(args)


Flight = namedtuple('Flight', ['Carrier', 'FlightNumber', 'Source',
                               'Destination', 'DepartureTimeStamp', 'ArrivalTimeStamp', 'Class',
                               'NumberOfStops', 'FareBasis', 'TicketType'])


def create_flight_from_dom(xml):
    return Flight._make([
        xml.Carrier.get_text(),
        xml.FlightNumber.get_text(),
        xml.Source.get_text(),
        xml.Destination.get_text(),
        xml.DepartureTimeStamp.get_text(),
        xml.ArrivalTimeStamp.get_text(),
        xml.Class.get_text(),
        xml.NumberOfStops.get_text(),
        xml.FareBasis.get_text(),
        xml.TicketType.get_text()
    ])


def get_itinerary(src, dest, ret):
    """ Возвращаем список маршрутов по указанным точкам """
    output = list()
    for itinerary in Storage:
        if itinerary.is_source(src) and itinerary.is_destination(dest):
            if not ret or (ret and len(itinerary.Return) > 0):
                output.append(itinerary)
    return output


def optimal_itinerary(itinerary, heap, cost):
    """ Находим с помощью кучи оптимальный маршрут через лямбду функцию """
    for i in itinerary:
        heapq.heappush(heap, (cost(i), i))
    return heap


def load_xml(xml):
    xml_string = BeautifulSoup(xml, 'xml')
    for itinerary in xml_string.PricedItineraries.children:
        if isinstance(itinerary, NavigableString):
            continue

        ownFlight = list()
        retFlight = list()
        pricing = list()

        # Проверяем есть ли маршрут с перелетами
        own = itinerary.find('OnwardPricedItinerary')
        if own:
            for flight in own.Flights.children:
                if isinstance(flight, NavigableString):
                    continue

                # Добавляем перелет
                ownFlight.append(create_flight_from_dom(flight))

        # Проверяем есть ли маршрут с перелетами
        ret = itinerary.find('ReturnPricedItinerary')
        if ret:
            for flight in ret.Flights.children:
                if isinstance(flight, NavigableString):
                    continue

                # Добавляем перелет
                retFlight.append(create_flight_from_dom(flight))

        price = itinerary.find('Pricing')
        if price:
            for p in price.children:
                if isinstance(p, NavigableString):
                    continue
                # Добавляем цену маршрута
                pricing.append(create_charges(p['type'], p['ChargeType'], p.get_text()))

        itinerary = Itinerary._make([ownFlight, retFlight, pricing])
        Storage.append(itinerary)
    return Storage
