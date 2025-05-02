package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST","/an-url",nil)
	f := New(r.PostForm)

	isValid := f.Valid()
	if !isValid {
		t.Error("expected valid, got invalid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST","/an-url",nil)
	form := New(r.PostForm)
	form.Required("req1","req2")
	
	if form.Valid() {
		t.Error("required fields missing: expected invalid, got valid")
	}

	postedData := url.Values{}
	postedData.Add("req1","req1 value");
	postedData.Add("req2","req2 value");
	postedData.Add("req3","req3 value");

	r, _ = http.NewRequest("POST","/an-url",nil)
	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("req1","req2")

	if !form.Valid() {
		t.Error("required fields are all present: expected valid, got invalid")
	}
}

func TestForm_Has(t *testing.T) {
	
	postedData := url.Values{}
	postedData.Add("req1","req1 value");
	f := New(postedData)
	has :=f.Has("startDate")
	if has {
		t.Error("check existence of a field: expected false, but got true")
	}
		
	has = f.Has("req1")
	if !has {
		t.Error("check existence of a field: expected true, but got false")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("y","abc")
	form :=New(postedData)
	form.MinLength("x",10)
	if form.Valid() {
		t.Error("check minimum length: for a non existing field the form accepts a minimum length")
	}

	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("check minimum length: should report an error, but does not")
	}

	form = New(postedData)
	form.MinLength("y",100)
	if form.Valid() {
		t.Error("check minimum length: returns true, but given length=3, required length=100")
	}

	form = New(postedData)
	form.MinLength("y",1)
	if !form.Valid() {
		t.Error("check minimum length: returns false, but given length=3, required length=1")
	}

	isError = form.Errors.Get("y")
	if isError != "" {
		t.Error("check minimum length: returns true, but it should return false (because given length=3 is greater than required length=1)")
	}

}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("check email: form shows valid email for non-existent field")
	}
	form = New(postedData)
	postedData.Add("email2","email@email.email")
	form.IsEmail("email2")
	if !form.Valid() {
		t.Error("check email: got false, but we provided a valid email")
	}

	form = New(postedData)
	postedData.Add("wrong_email","wrong@")
	form.IsEmail("wrong_email")
	if form.Valid() {
		t.Error("check email: got true, but we provided an invalid email")
	}
}
