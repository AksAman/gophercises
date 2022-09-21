package hnclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	testIDs  []int  = []int{1, 2, 3}
	testID   int    = 1
	testUser string = "test_user"
)

func areSlicesEqual[Type comparable](a, b []Type) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func setupAndTear() (serverURL string, tearDown func()) {
	mux := http.NewServeMux()

	topStoriesHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err := json.NewEncoder(w).Encode(testIDs)
		if err != nil {
			panic(err)
		}
	}

	itemHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		rawItem := map[string]any{
			"by":          testUser,
			"descendants": 71,
			"id":          testID,
			"kids":        []int{8952, 9224},
			"score":       111,
			"time":        1175714200,
			"title":       "My YC app: Dropbox - Throw away your USB drive",
			"type":        "story",
			"url":         "http://www.getdropbox.com/u/2/screencast.html",
		}

		err := json.NewEncoder(w).Encode(rawItem)
		if err != nil {
			panic(err)
		}
	}

	mux.HandleFunc("/topstories.json", topStoriesHandler)
	mux.HandleFunc(fmt.Sprintf("/item/%d.json", testID), itemHandler)

	server := httptest.NewServer(mux)

	serverURL = server.URL
	tearDown = func() {
		server.Close()
	}
	return serverURL, tearDown
}

func compare(t *testing.T, got, want any) {
	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestClient_GetTopItems(t *testing.T) {
	baseURL, tearDown := setupAndTear()
	defer tearDown()

	hnClient := Client{
		baseURL: baseURL,
	}

	ids, err := hnClient.GetTopItems()
	if err != nil {
		t.Errorf("hnClient.GetTopItems() returned an error: %v", err)
	}
	if !areSlicesEqual(ids, testIDs) {
		t.Errorf("got: %v, want: %v", ids, testIDs)
	}
}

func TestClient_GetItem(t *testing.T) {
	baseURL, tearDown := setupAndTear()
	defer tearDown()

	hnClient := Client{
		baseURL: baseURL,
	}

	item, err := hnClient.GetItem(testID)
	if err != nil {
		t.Errorf("hnClient.GetItem(%d) returned an error: %v", testID, err)
	}

	if item.By != testUser {
		t.Errorf("item.By got: %v, want: %v", item.By, testUser)
	}

	if item.ID != int64(testID) {
		t.Errorf("item.ID got: %v, want: %v", item.ID, testID)
	}

	if item.Type != "story" {
		t.Errorf("item.Type got: %v, want: %v", item.Type, "story")
	}
}

func TestClient_defaultify(t *testing.T) {
	client := Client{}
	client.defaultify()

	compare(t, client.baseURL, BASE_URL)
}
