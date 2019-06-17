from app import server
from app.tests.test_storage import XML
from app import storage

import pytest
import urllib


@pytest.fixture
def app():
    return server.app


def test_server(client):
    response = client.post("/receive")
    assert response.json['code'] == 401

    del storage.Storage[:]
    response = client.post("/receive", data=XML)
    assert response.json['status'] == 'ok'

    urls = (
        ('source', "DXB"),
        ('destination', "BKK"),
    )
    response = client.get("/itinerary?%s" % urllib.parse.urlencode(urls))
    assert len(response.json['result']) == 2

    urls = (
        ('source', "DXB"),
        ('destination', "BKK"),
        ('type', 'stime'),
    )
    response = client.get("/itinerary?%s" % urllib.parse.urlencode(urls))
    assert response.json['result'][0] == -47100.0
