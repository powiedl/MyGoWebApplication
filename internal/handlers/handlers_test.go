package handlers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct{
	key   string
	value string
}


// returns string except the first character/rune
func trimLeftChar(s string) string {
	for i := range s {
			if i > 0 {
					return s[i:]
			}
	}
	return s[:0]
}

var allTheHandlerTests = []struct{
	name string             // name of the test
	url string              // url of the handler to test
	method string           // http verb
	params []postData       // parameters of the handler
	expectedStatusCode int  // expected status code
} {
	{"home","/","GET",[]postData{},http.StatusOK,},
	{"","/about","GET",[]postData{},http.StatusOK,},
	{"","/eremite","GET",[]postData{},http.StatusOK,},
	{"","/couple","GET",[]postData{},http.StatusOK,},
	{"","/family","GET",[]postData{},http.StatusOK,},
	{"","/reservation","GET",[]postData{},http.StatusOK,},
	{"","/make-reservation","GET",[]postData{},http.StatusOK,},
	//{"","/reservation-overview","GET",[]postData{},http.StatusOK,},
	{"","/reservation","POST",[]postData{
		{key:"startingDate",value:"2025-05-01"},
		{key:"endingDate",value:"2025-05-03"},
	},http.StatusOK,},
	{"","/reservation-json","POST",[]postData{
		{key:"start",value:"2025-05-01"},
		{key:"end",value:"2025-05-03"},
	},http.StatusOK,},
	{"","/make-reservation","POST",[]postData{
		{key:"full_name",value:"Test Name"},
		{key:"email",value:"test@name.local"},
		{key: "phone",value:"1324"},
	},http.StatusOK,},
	{
		"unknown-route","/does-not-exist","GET",[]postData{},http.StatusNotFound,
	},

}
func TestAllHandlers(t *testing.T) {
	routes := getRoutes()

	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _,test := range allTheHandlerTests {
		var testName=test.name
		if test.name == "" {
			testName = trimLeftChar(test.url)
		} else {
			testName = test.name
		}
		testName = testName + " " + test.method
		t.Logf("Running Test '%s'",testName)
		if test.method == "GET" {
			response, err := testServer.Client().Get(testServer.URL + test.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if response.StatusCode != test.expectedStatusCode {
				t.Errorf("%s: expected status code %d, but got %d",testName,test.expectedStatusCode,response.StatusCode)
			}
		} else if test.method == "POST" {
			values := url.Values{}
			for _,param := range test.params {
				values.Add(param.key,param.value)
			}
			response, err := testServer.Client().PostForm(testServer.URL + test.url,values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			 }
			 if response.StatusCode != test.expectedStatusCode {
				t.Errorf("%s: expected status code %d, but got %d",testName,test.expectedStatusCode,response.StatusCode)
			 }
		} else {
			log.Println("Unknown method:",test.method)
		}
	}
}