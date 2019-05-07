package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func router() *mux.Router {

	router := mux.NewRouter()
	router.Path("/find").Queries("word", "{word}").HandlerFunc(getFindEndpoint).Methods("GET")
	router.Path("/compare").Queries("word1", "{word1}", "word2", "{word2}").HandlerFunc(getCompareEndpoint).Methods("GET")
	return router

}

func TestFindEndpointWithBadRequest(t *testing.T) {

	request, _ := http.NewRequest("GET", "/find", nil)
	query := request.URL.Query()
	query.Add("word", "")
	request.URL.RawQuery = query.Encode()
	response := httptest.NewRecorder()
	router().ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "should throw error for missing parameter GET find")
}

func TestFindEndpointWithNoMatch(t *testing.T) {

	request, _ := http.NewRequest("GET", "/find", nil)
	query := request.URL.Query()
	query.Add("word", "atc")
	request.URL.RawQuery = query.Encode()
	response := httptest.NewRecorder()
	router().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "should have a sucessful response GET find, but with no match")
}

func TestFindEndpointWithMatch(t *testing.T) {

	request, _ := http.NewRequest("GET", "/find", nil)
	query := request.URL.Query()
	query.Add("word", "reef")
	request.URL.RawQuery = query.Encode()
	response := httptest.NewRecorder()
	router().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "should have a sucessful response GET find, with a match")
}

func TestCompareEndpointWithBadRequest(t *testing.T) {

	request, _ := http.NewRequest("GET", "/compare", nil)
	query := request.URL.Query()
	query.Add("word1", "")
	query.Add("word2", "")
	request.URL.RawQuery = query.Encode()
	response := httptest.NewRecorder()
	router().ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "should throw error for missing parameter GET compare")
}

func TestComapreEndpointWithNoMatch(t *testing.T) {

	request, _ := http.NewRequest("GET", "/compare", nil)
	query := request.URL.Query()
	query.Add("word1", "spray")
	query.Add("word2", "yarps")
	request.URL.RawQuery = query.Encode()
	response := httptest.NewRecorder()
	router().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "should have a sucessful response GET compare, but with no match")
}

func TestCompareEndpointWithMatch(t *testing.T) {

	request, _ := http.NewRequest("GET", "/compare", nil)
	query := request.URL.Query()
	query.Add("word1", "case")
	query.Add("word2", "aces")
	request.URL.RawQuery = query.Encode()
	response := httptest.NewRecorder()
	router().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "should have a sucessful response GET find, with a match")
}
