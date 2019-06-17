package storage

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	storage := NewStorage()

	absPath, _ := filepath.Abs("../../../../../fixtures/one3.xml")
	fileBytes, err := ioutil.ReadFile(absPath)
	if err != nil {
		t.Errorf("Не смогли прочитать файл %v", err)
	}
	err = storage.LoadXML(fileBytes)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	absPath, _ = filepath.Abs("../../../../../fixtures/oneOW.xml")
	fileBytes, err = ioutil.ReadFile(absPath)
	if err != nil {
		t.Errorf("Не смогли прочитать файл %v", err)
	}
	err = storage.LoadXML(fileBytes)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	itinerary := storage.GetItinerary("DXB", "BKK", false)
	assert.Greater(t, len(itinerary), 0)

	maxTimeQueue := NewTimeQueueMax()
	storage.OptimalItinerary(itinerary, maxTimeQueue)
	maxTime := maxTimeQueue.PopOrdered()
	assert.Equal(t, maxTime.Onward[0].FlightNumber, "756")

}
