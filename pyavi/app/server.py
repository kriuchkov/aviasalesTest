from flask import Flask, request, make_response, jsonify
from app import storage

import heapq

try:
    import simplejson as json
except ImportError:
    import json


app = Flask(__name__)


def response(dct, code=None):
    return make_response(jsonify(dct), code)


@app.route('/receive', methods=['POST'])
def receive():
    data = request.stream.read()
    if not data:
        return response({'message': 'Тело запроса не найдено', 'code': 401}, 400)
    try:
        storage.load_xml(data.decode("utf-8"))
        return response({'status': 'ok'})
    except Exception as e:
        return response({'message': str(e), 'code': 403}, 400)


@app.route('/itinerary')
def itinerary():
    source = request.args.get('source', None)
    destination = request.args.get('destination', None)
    ret = request.args.get("return", False)

    if not source or not destination:
        return response({'message': 'Точка назначения не найдены', 'code': 402}, 400)

    of_type = {
        'lprice': lambda i: i.price(),      # самый дешевый маршрут
        'sprice': lambda i: -i.price(),     # самый дорогой маршрут
        'ltime': lambda i: i.duration(),    # самый быстрый маршрут
        'stime': lambda i: -i.duration(),   # самый долгий маршрут
        'optimal': lambda i: i.price() * i.duration(),  # оптимальный :D
    }
    get_type = request.args.get("type", None)

    output = storage.get_itinerary(source, destination, ret)

    if get_type and of_type.get(get_type, False):
        h = []
        storage.optimal_itinerary(output, h, of_type.get(get_type))
        output = heapq.heappop(h)

    return app.response_class(
        json.dumps({"status": "ok", "result": output}),
        mimetype='application/json'
    )
