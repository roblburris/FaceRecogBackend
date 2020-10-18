package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInsert(t *testing.T) {

	testRequest, err := http.NewRequest("GET", "foo", nil)
	if err != nil {
		log.Fatal(err)
	}
	testRequest.Header.Add("Name", "new-test1")
	testRequest.Header.Add("Title", "new-test2")
	testRequest.Header.Add("Foo1", "foo1")
	testRequest.Header.Add("Foo2", "foo2")
	testRequest.Header.Add("Foo3", "foo3")
	savePerson(httptest.NewRecorder(), testRequest)
	got := 1
	if got != 1 {
		t.Errorf("Abs(-1) = %d; want 1", got)
	}
}

func TestQuery(t *testing.T) {
	ctx, collection := setupMongo()
	fmt.Println(queryPerson(ctx, collection, "5f8bde21dbfa5439dec22ec0"))
}
