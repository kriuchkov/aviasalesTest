
export FLASK_APP=$(CURDIR)/pyavi/main.py
pyrun:
	(cd pyavi && pipenv run flask run --port=8080)

flake8:
	(cd pyavi && pipenv run flake8 $(CURDIR)/pyavi/)

pytest:
	(cd pyavi && pipenv run pytest -vs --tb=short)


export GOPATH := $(CURDIR)/goavi
gorun:
	(go install goavi/cmd/goavi && $(CURDIR)/goavi/bin/goavi)

gotest:
	(cd goavi && go test ./src/goavi/... -v)

GODIRS = goavi/src/goavi/cmd/goavi goavi/src/goavi/pkg/
fmt:
	gofmt -l -w $(GODIRS)

check:
	gofmt -l $(GODIRS)

imports:
	goimports -l -w $(GODIRS)


# monkey-testing
upload:
	curl -s 127.0.0.1:8080/receive -X POST -d @fixtures/oneOW.xml
	curl -s 127.0.0.1:8080/receive -X POST -d @fixtures/one3.xml

get:
	curl -s "127.0.0.1:8080/itinerary?source=DXB&destination=BKK&type=stime" | jq