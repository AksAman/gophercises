package controllers

import (
	"fmt"
	"net/http"
)

func ProblematicFunc() {
	panic(fmt.Errorf("some Error"))
}

func FakeError(w http.ResponseWriter, r *http.Request) {
	defer func() {
		fmt.Println("FakeError")
	}()
	ProblematicFunc()
}

func FakeErrorAfter(w http.ResponseWriter, r *http.Request) {
	defer func() {
		fmt.Println("FakeErrorAfter")
	}()

	fmt.Fprintf(w, "Hello")
	ProblematicFunc()

}
