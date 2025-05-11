package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/powiedl/myGoWebApplication/internal/driver"
	"github.com/powiedl/myGoWebApplication/internal/models"
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
	// {"","/reservation","POST",[]postData{
	// 	{key:"startingDate",value:"2025-05-01"},
	// 	{key:"endingDate",value:"2025-05-03"},
	// },http.StatusOK,},
	// {"","/reservation-json","POST",[]postData{
	// 	{key:"start",value:"2025-05-01"},
	// 	{key:"end",value:"2025-05-03"},
	// },http.StatusOK,},
	// {"","/make-reservation","POST",[]postData{
	// 	{key:"full_name",value:"Test Name"},
	// 	{key:"email",value:"test@name.local"},
	// 	{key: "phone",value:"1324"},
	// },http.StatusOK,},
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

// TestRepository MakeReservation tests the MakeReservation GET request handler
func TestRepository_MakeReservation(t *testing.T) {
	reservation := models.Reservation{
		BungalowID: 1,
		Bungalow: models.Bungalow{
			ID: 1,
			BungalowName: "The Solitude Shack",
		},
	}

	req,_ := http.NewRequest("GET","/make-reservation",nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// rr bedeutet "request recorder" und ist ein initialisierter Response Recorder für http Requests aus dem httptest Package
	// damit kann man einen Client "faken" - und somit einen gültigen Request/Response Zyklus während eines Tests zur Verfügung stellen
	rr := httptest.NewRecorder()
	session.Put(ctx,"reservation",reservation)
	//log.Println(reservation)

	// der zu testende Handler wird als Variable gespeichert
	handler := http.HandlerFunc(Repo.MakeReservation)

	// hier der Aufruf des zu testenden Handlers - mit dem Request Recorder als writer und dem req (mit der Session!) als Reques
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusOK {
		t.Errorf("Handler MakeReservation failed: unexpected response code: got %d, expected %d",rr.Code,http.StatusOK)
	}

	// test case without a reservation in the session
	req,_ = http.NewRequest("GET","/make-reservation",nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Handler MakeReservation failed: unexpected response code: got %d, expected %d",rr.Code,http.StatusTemporaryRedirect)
	}

	// test error returned from database query (invalid BungalowId)
	req,_ = http.NewRequest("GET","/make-reservation",nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.BungalowID = 99
	session.Put(ctx,"reservation",reservation)	

	handler.ServeHTTP(rr,req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Handler MakeReservation failed: unexpected response code: got %d, expected %d",rr.Code,http.StatusTemporaryRedirect)
	}
}

// TestRepository_PostMakeReservation tests the PostMakeReservation handler
func TestRepository_PostMakeReservation(t *testing.T) {
	// case #1: reservation works fine
	postedData := url.Values{}
	postedData.Add("full_name", "Peter Griffin")
	postedData.Add("email", "peter@griffin.family")
	postedData.Add("phone", "1234567890")

	// data to put in session
	layout := "2006-01-02"
	sd, _ := time.Parse(layout, "2037-01-01")
	ed, _ := time.Parse(layout, "2037-01-02")
	bungalowId, _ := strconv.Atoi("1")

	reservation := models.Reservation{
		StartDate:  sd,
		EndDate:    ed,
		BungalowID: bungalowId,
		Bungalow: models.Bungalow{
			BungalowName: "some bungalow name for tests",
		},
	}

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostMakeReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// case #2: missing post body
	// create request
	req, _ = http.NewRequest("POST", "/make-reservation", nil)

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler returned wrong response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// case #3: missing session data
	postedData = url.Values{}
	postedData.Add("full_name", "Peter Griffin")
	postedData.Add("email", "peter@griffin.family")
	postedData.Add("phone", "1234567890")

	// data to put in session
	layout = "2006-01-02"
	sd, _ = time.Parse(layout, "2037-01-01")
	ed, _ = time.Parse(layout, "2037-01-02")
	bungalowId, _ = strconv.Atoi("1")

	reservation = models.Reservation{
		StartDate:  sd,
		EndDate:    ed,
		BungalowID: bungalowId,
		Bungalow: models.Bungalow{
			BungalowName: "some bungalow name for tests",
		},
	}

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler returned wrong response code for missing session data: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// case #4: invalid/insufficient data
	postedData = url.Values{}
	postedData.Add("full_name", "P")
	postedData.Add("email", "peter@griffin.family")
	postedData.Add("phone", "1234567890")

	// data to put in session
	layout = "2006-01-02"
	sd, _ = time.Parse(layout, "2037-01-01")
	ed, _ = time.Parse(layout, "2037-01-02")
	bungalowId, _ = strconv.Atoi("1")

	reservation = models.Reservation{
		StartDate:  sd,
		EndDate:    ed,
		BungalowID: bungalowId,
		Bungalow: models.Bungalow{
			BungalowName: "some bungalow name for tests",
		},
	}

	// create request
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("PostMakeReservation handler returned wrong response code invalid/insufficient data: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// case #5:  failure inserting reservation into database
	postedData = url.Values{}
	postedData.Add("full_name", "Peter Griffin")
	postedData.Add("email", "peter@griffin.family")
	postedData.Add("phone", "1234567890")

	// data to put in session
	layout = "2006-01-02"
	sd, _ = time.Parse(layout, "2037-01-01")
	ed, _ = time.Parse(layout, "2037-01-02")
	bungalowId, _ = strconv.Atoi("99")

	reservation = models.Reservation{
		StartDate:  sd,
		EndDate:    ed,
		BungalowID: bungalowId,
		Bungalow: models.Bungalow{
			BungalowName: "some bungalow name for tests",
		},
	}

	// create request
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler failed when trying to inserting a reservation into the database: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// case #6: failure to inserting restriction into database
	postedData = url.Values{}
	postedData.Add("full_name", "Peter Griffin")
	postedData.Add("email", "peter@griffin.family")
	postedData.Add("phone", "1234567890")

	// data to put in session
	layout = "2006-01-02"
	sd, _ = time.Parse(layout, "2037-01-01")
	ed, _ = time.Parse(layout, "2037-01-02")
	bungalowId, _ = strconv.Atoi("999")

	reservation = models.Reservation{
		StartDate:  sd,
		EndDate:    ed,
		BungalowID: bungalowId,
		Bungalow: models.Bungalow{
			BungalowName: "some bungalow name for tests",
		},
	}

	// create request
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))

	// get the context
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler failed when trying to inserting a reservation into the database: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

// #region falscher Test für PostMakeReservation
/*
// TestRepository_PostMakeReservation tests the PostMakeReservation handler
func TestRepository_PostMakeReservation(t *testing.T) {
	// OK Case
	
	// build request body by concattening a long string - not the prefered method
	reqBody := "start_date=2037-01-01"
	reqBody = fmt.Sprintf("%s&%s",reqBody,"end_date=2037-01-04")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"full_name=Peter Griffin")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"email=peter@Griffin.io")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"phone=1234565")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"bungalow_id=1")
	req, _ := http.NewRequest("POST","/make-reservation",strings.NewReader(reqBody))

	// build request body by using url.Values - the prefered method
	postedData := url.Values{}
	postedData.Add("start_date","2037-01-01")
	postedData.Add("end_date","2037-01-04")
	postedData.Add("full_name","Peter Griffin")
	postedData.Add("email","peter@Griffin.io")
	postedData.Add("phone","1234565")
	postedData.Add("bungalow_id","1")
	req, _ = http.NewRequest("POST","/make-reservation",strings.NewReader(postedData.Encode()))

	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Handler PostMakeReservation failed: unexpected response code: got %d, expected %d",rr.Code,http.StatusSeeOther)
	}

	// missing body
	req, _ = http.NewRequest("POST","/make-reservation",nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Handler PostMakeReservation failed (OK test): unexpected response code: got %d, expected %d",rr.Code,http.StatusSeeOther)
	}

	// missing start_date
	reqBody = "end_date=2037-01-01"
	reqBody = fmt.Sprintf("%s&%s",reqBody,"full_name=Peter Griffin")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"email=peter@Griffin.io")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"phone=1234565")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"bungalow_id=1")

  req, _ = http.NewRequest("POST","/make-reservation",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Handler PostMakeReservation failed (missing start_date): unexpected response code: got %d, expected %d",rr.Code,http.StatusTemporaryRedirect)
	}

	// missing end_date
	reqBody = "start_date=2037-01-01"
	reqBody = fmt.Sprintf("%s&%s",reqBody,"full_name=Peter Griffin")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"email=peter@Griffin.io")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"phone=1234565")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"bungalow_id=1")

  req, _ = http.NewRequest("POST","/make-reservation",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Handler PostMakeReservation failed (missing end_date): unexpected response code: got %d, expected %d",rr.Code,http.StatusTemporaryRedirect)
	}

	// invalid bungalow_id
	reqBody = "start_date=2037-01-01"
	reqBody = fmt.Sprintf("%s&%s",reqBody,"end_date=2037-01-04")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"full_name=Peter Griffin")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"email=peter@Griffin.io")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"phone=1234565")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"bungalow_id=invalid")

  req, _ = http.NewRequest("POST","/make-reservation",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Handler PostMakeReservation failed (invalid bungalow_id): unexpected response code: got %d, expected %d",rr.Code,http.StatusSeeOther)
	}

	// invalid/insufficient data
	reqBody = "start_date=2037-01-01"
	reqBody = fmt.Sprintf("%s&%s",reqBody,"end_date=2037-01-04")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"full_name=P")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"email=peter@Griffin.io")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"phone=1234565")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"bungalow_id=1")

  req, _ = http.NewRequest("POST","/make-reservation",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Handler PostMakeReservation failed (invalid full_name): unexpected response code: got %d, expected %d",rr.Code,http.StatusSeeOther)
	}

	// failure inserting reservation into database
	reqBody = "start_date=2037-01-01"
	reqBody = fmt.Sprintf("%s&%s",reqBody,"end_date=2037-01-04")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"full_name=Peter")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"email=peter@Griffin.io")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"phone=1234565")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"bungalow_id=99")

  req, _ = http.NewRequest("POST","/make-reservation",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Handler PostMakeReservation failed (expected failure at inserting reservation into database, bungalow_id=99): unexpected response code: got %d, expected %d",rr.Code,http.StatusTemporaryRedirect)
	}

	// failure inserting restriction into database
	reqBody = "start_date=2037-01-01"
	reqBody = fmt.Sprintf("%s&%s",reqBody,"end_date=2037-01-04")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"full_name=Peter")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"email=peter@Griffin.io")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"phone=1234565")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"bungalow_id=999")

  req, _ = http.NewRequest("POST","/make-reservation",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Handler PostMakeReservation failed (expected failure at inserting restriction into database, bungalow_id=999): unexpected response code: got %d, expected %d",rr.Code,http.StatusTemporaryRedirect)
	}
}
*/
// #endregion

// TestRespository_ReservationJSON tests the MakeReservation get request handler
func TestRepository_ReservationJSON(t *testing.T) {
	// OK Case
	reqBody := "start=2037-01-01"
	reqBody = fmt.Sprintf("%s&%s",reqBody,"end=2037-01-04")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"bungalow_id=1")

	req, _ := http.NewRequest("POST","/reservation-json",strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.ReservationJSON)
	handler.ServeHTTP(rr,req)

	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()),&j)
	if err != nil {
		t.Error("Handler ReservationJSON failed (expected correct data, but can't parse json), got '",rr.Body.String(),"'")
	}

	// empty request body
	log.Println("empty request body")
	req, _ = http.NewRequest("POST","/reservation-json",nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ReservationJSON)
	handler.ServeHTTP(rr,req)
	
	err = json.Unmarshal(rr.Body.Bytes(), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// start date < 2036-12-31, expecting availability
	if j.OK || j.Message != "Internal server error" {
		t.Error("got availability with empty request body")
	}
	
	// invalid bungalow_id
	reqBody = "start=2037-01-01"
	reqBody = fmt.Sprintf("%s&%s",reqBody,"end=2037-01-04")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"bungalow_id=invalid")

	req, _ = http.NewRequest("POST","/reservation-json",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ReservationJSON)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Handler ReservationJSON failed (invalid bungalow_id), got %d, expected %d",rr.Code,http.StatusTemporaryRedirect)
	}
	
	// invalid start_date
	reqBody = "start=invalid"
	reqBody = fmt.Sprintf("%s&%s",reqBody,"end=2037-01-04")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"bungalow_id=1")

	req, _ = http.NewRequest("POST","/reservation-json",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ReservationJSON)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Handler ReservationJSON failed (invalid start), got %d, expected %d",rr.Code,http.StatusTemporaryRedirect)
	}

	// invalid end_date
	reqBody = "start=2037-01-01"
	reqBody = fmt.Sprintf("%s&%s",reqBody,"end=invalid")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"bungalow_id=1")

	req, _ = http.NewRequest("POST","/reservation-json",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ReservationJSON)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Handler ReservationJSON failed (invalid end), got %d, expected %d",rr.Code,http.StatusTemporaryRedirect)
	}

	// bungalow is available
	reqBody = "start=2036-11-30"
	reqBody = fmt.Sprintf("%s&%s",reqBody,"end=2036-12-02")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"bungalow_id=1")

	req, _ = http.NewRequest("POST","/reservation-json",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ReservationJSON)
	handler.ServeHTTP(rr,req)

	err = json.Unmarshal([]byte(rr.Body.String()),&j)
	if err != nil {
		t.Errorf("Handler ReservationJSON failed (bungalow available), but got an invalid json")
	}
	if j.OK != true {
		t.Errorf("Handler ReservationJSON failed (bungalow available), got %v, expected %v",j.OK,true)
	}

	// bungalow is not available
	reqBody = "start=2039-11-30"
	reqBody = fmt.Sprintf("%s&%s",reqBody,"end=2039-12-02")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"bungalow_id=1")

	req, _ = http.NewRequest("POST","/reservation-json",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ReservationJSON)
	handler.ServeHTTP(rr,req)

	err = json.Unmarshal([]byte(rr.Body.String()),&j)
	if err != nil {
		t.Errorf("Handler ReservationJSON failed (bungalow available), but got an invalid json")
	}
	if j.OK != false {
		t.Errorf("Handler ReservationJSON failed (bungalow available), got %v, expected %v",j.OK,false)
	}

		// simulate error in DB function
	reqBody = "start=2038-01-01"
	reqBody = fmt.Sprintf("%s&%s",reqBody,"end=2039-12-02")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"bungalow_id=1")

	req, _ = http.NewRequest("POST","/reservation-json",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ReservationJSON)
	handler.ServeHTTP(rr,req)

	err = json.Unmarshal([]byte(rr.Body.String()),&j)
	if err != nil {
		t.Errorf("Handler ReservationJSON failed (bungalow available), but got an invalid json")
	}
	if j.OK != false {
		t.Errorf("Handler ReservationJSON failed (bungalow available), got %v, expected %v",j.OK,false)
	}

}

// TestRepository_ChooseBungalow tests the ChooseBungalow handler
func TestRepository_ChooseBungalow(t *testing.T) {
	reservation := models.Reservation{
		BungalowID: 1,
		Bungalow: models.Bungalow{
			ID: 1,
			BungalowName: "The Solitude Shack",
		},
	}
	// OK case
	req,_ := http.NewRequest("GET","/choose-bungalow/1",nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx,"reservation",reservation)
	
	handler := http.HandlerFunc(Repo.ChooseBungalow)
	handler.ServeHTTP(rr,req)
	
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Handler ChooseBungalow (correct request) failed, got %d, expected %d",rr.Code,http.StatusSeeOther)
	}
}

// TestRepository_BookBungalow tests the BookBungalow handler
// example for query parameters ...
func TestRepository_BookBungalow(t *testing.T) {
	reservation := models.Reservation{
		BungalowID: 1,
		Bungalow: models.Bungalow{
			ID:1,
			BungalowName: "The Solitude Shack",
		},
	}

	req,_ := http.NewRequest("GET","/book-bungalow?s=2036-11-01&e=2036-12-02&id=1",nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx,"reservation",reservation)

	handler := http.HandlerFunc(Repo.BookBungalow)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Handler BookBungalow (correct request), got %d, expected %d",rr.Code,http.StatusSeeOther)
	}
}

// TestNewRepo tests the NewRepo function using the reflect package
func TestNewRepo(t *testing.T) {
  var db driver.DB
	testRepo := NewRepo(&app,&db)

	// reflect.TypeOf(testRepo) is used to determine the runtime type of testRepo - often it is used to further determine the type in case of an empty interface is required
	if reflect.TypeOf(testRepo).String() != "*handlers.Repository" {
		t.Errorf("Did not get correct type from NewRepo: got %s, expected *Repository",reflect.TypeOf(testRepo).String())
	}
}

func getCtx(req *http.Request) context.Context {
	ctx,err := session.Load(req.Context(),req.Header.Get("X-Session"))
	if err != nil {
		log.Println("getCtx for '",req.URL,"' - error:",err)
	}
	return ctx
}