package common

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func createTransporter() Transporter {
	return NewTransporter(5)
}

func TestDatacenterFetch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Got request ", r.URL.Path)
		if strings.Contains(r.URL.Path, "/fax/send") {
			w.WriteHeader(http.StatusOK)
			return
		}
		if strings.Contains(r.URL.Path, "/sms/send") {
			w.WriteHeader(http.StatusOK)
			return
		}
		// For example, respond with a 404 status
		w.WriteHeader(http.StatusInternalServerError)
		t.Fatalf("unexpected URL path: %s", r.URL.Path)

		fmt.Fprint(w, "Message accepted")
	}))
	defer server.Close()

	ts := createTransporter()
	server1 := server.URL + "/fax/send"
	server2 := server.URL + "/sms/send"

	servers := []string{server1, server2}

	res := ts.DoDatacenterFetch(servers, "", "", nil, "", http.MethodGet)
	if res[0].StatusCode != 200 && res[1].StatusCode != 200 {
		t.Errorf("invalid status code returned by one of the responses")
	}
}

func TestDatacenterFetchWith404(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Got request ", r.URL.Path)
		if requestCount == 1 {
			w.WriteHeader(http.StatusOK)
			return
		}
		// For example, respond with a 404 status
		w.WriteHeader(http.StatusNotFound)
		requestCount += 1
		return
	}))
	defer server.Close()

	ts := createTransporter()

	servers := []string{server.URL + "/abc", server.URL + "/def"}
	res := ts.DoDatacenterFetch(servers, "", "", nil, "", http.MethodGet)
	if res[0].StatusCode != 400 && res[1].StatusCode != 200 {
		t.Errorf("invalid status code returned by one of the responses")
	}
}
