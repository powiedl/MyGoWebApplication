package render

import (
	"net/http"
	"testing"

	"github.com/powiedl/myGoWebApplication/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r,err := getSession()
	if err != nil {
		t.Fatal(err)
	}

	session.Put(r.Context(),"flash","a flash message")
	result := AddDefaultData(&td,r)
	if result.Flash != "a flash message" {
		t.Errorf("expected value for key flash '%s', but found '%s' in session","a flash message",result.Flash)
	}
}

func TestRenderTemplate(t *testing.T) {
	tc, err := CreateTemplateCache(app)
	if err != nil {
		t.Error(err)
	}
	app.TemplateCache = tc

	r,err := getSession()
	if err != nil {
		t.Fatal(err)
	}

	var testW myWriter

	err = RenderTemplate(&testW,r,"home-page.template.html",&models.TemplateData{})
	if err != nil {
		t.Error("Writing template to browser failed:",err)
	}

	err = RenderTemplate(&testW,r,"does-not-exist-page.template.html",&models.TemplateData{})
	if err == nil {
		t.Error("Writing a not existing template to browser did NOT fail, but it should fail")
	}
}

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

func TestCreateTemplateCache(t *testing.T) {
	
	_, err := CreateTemplateCache(app)
	if err != nil {
		t.Error(err)
	}
}

func getSession() (*http.Request,error) {
  r,err := http.NewRequest("GET","/",nil)
	if err != nil {
		return nil,err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx,r.Header.Get("X-Session"))
  r = r.WithContext(ctx)
	return r,nil
}