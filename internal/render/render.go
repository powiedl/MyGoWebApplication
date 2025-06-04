package render

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/justinas/nosurf"
	"github.com/powiedl/myGoWebApplication/internal/config"
	"github.com/powiedl/myGoWebApplication/internal/models"
)

var functions = template.FuncMap{
	"humanReadableDate":HumanReadableDate,
}

// HumanReadableDate returns a time value in the YYYY MM DD format
func HumanReadableDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func AddDefaultData(td *models.TemplateData,r *http.Request) *models.TemplateData{
	td.Success= app.Session.PopString(r.Context(),"success")
	td.Error= app.Session.PopString(r.Context(),"error")
	td.Warning= app.Session.PopString(r.Context(),"warning")
	td.CSRFToken=nosurf.Token(r)
	if app.Session.Exists(r.Context(),"user_id") {
		td.IsAuthenticated = 1
	}
	return td
}

var app *config.AppConfig // pointer to a config.AppConfig (the app config populated in main)

// NewRenderer sets the config for the template
func NewRenderer(a *config.AppConfig) {
	app = a
}

// #region 4-32 bis 4-34 statischer Cache
// meine Version
func Template(w http.ResponseWriter,r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template
	if app.UseCache {
		// get the template cache from app config
		tc = app.TemplateCache
	} else {
		// create a new template cache in every render (so no cache is used)
		tc, _ = CreateTemplateCache(app)
	}
	// get the right template from cache
	t, ok := tc[tmpl]
	if !ok {
		return errors.New("template not found in cache for some reason")
	} else {
		//log.Println("using template ",t)
	}

	// store result in a buffer and double-check if it is a valid value
	buf := new(bytes.Buffer) // creates a buffer of bytes

	td = AddDefaultData(td,r) // add default data to the td

	//log.Println("td.StringMap",td.StringMap)
	//log.Println("td.CSRFToken",td.CSRFToken)
	//tdJson,_ := json.Marshal(td)
	//log.Println("td",td)
	//log.Println("tdJson",string(tdJson))
	err := t.Execute(buf, td) // tries to render the template
	if err != nil {
		log.Println("Tried to render the template, but got this error:",err)
		return err
	}
	//fmt.Println(buf)
	// render that template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println("send the result of the rendering of the template to the client, but got this error:",err)
		return err
	}

	return nil
}

// CreateTemplateCache creates a map and stores the templates in for caching 
func CreateTemplateCache(app *config.AppConfig) (map[string]*template.Template,error) {
	//log.Println("createTemplateCache...")
	theCache := map[string]*template.Template{}

	// get all available files *-page.template.html from folder ../../templates
	pages,err := filepath.Glob(app.Basedir + "templates/*-page.template.html")
	log.Println("CreateTemplateCache, #pages",len(pages))
	if err != nil {
		return theCache,err
	}

	// range through the slice of *-page.template.html
	for _,page := range pages {
		name := filepath.Base(page)
		//log.Println("  Processing page:",name)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			log.Printf("  Processing page: '%s' -   ERROR: %s",name,err)
			return theCache,nil
		}
		matches, err := filepath.Glob(app.Basedir + "templates/*-layout.template.html") // get all layout files
		if err != nil {
			return theCache, err
		}

		if len(matches)>0 {
			// at least one base layout was found - append the layouts to the currently processed page
			ts, err = ts.ParseGlob(app.Basedir + "templates/*-layout.template.html")
			if err != nil {
				log.Printf("  Processing templates for page: '%s' -   ERROR: %s",name,err)
				return theCache,nil
			}	
		}
		
		theCache[name] = ts
	}
	return theCache,nil
}
// #endregion

// #region 4-30 ohne cache
/*
func RenderTemplate(w http.ResponseWriter,tmpl string) {
	parsedTemplate, err := template.ParseFiles("../../templates/" + tmpl,"../../templates/base-layout.template.html")
	if err!=nil {
		fmt.Println("Error parsing template:",err)
		return
	}
	err = parsedTemplate.Execute(w,nil)
	if err != nil {
		fmt.Println("error executing parsed template:",err)
	}
}
*/
// #endregion

// #region 4-31 dynamischer Cache
/*
var tc = make(map[string]*template.Template) // map, wo der key ein string ist und der value ein Pointer auf ein template.Template (das wird von template.ParseFiles zurückgeliefert)
func RenderTemplate(w http.ResponseWriter, t string) {
	var tmpl *template.Template
	var err error
	log.Println("RenderTemplate for ",t)
	// check to see if we already have the template in our cache
	_, inMap := tc[t] // es gibt ein Element mit dem Key t in tc
	if !inMap {
		// need to create the template
		err = createTemplateCache(t)
		if err != nil {
			log.Println(err)
		}
	} else {
		// we have the template in the cache
		log.Println("using the cached template")
	}

	tmpl = tc[t]
	err = tmpl.Execute(w,nil)
	if err != nil {
		log.Println(err)
	}
}

func createTemplateCache(t string) error {
	log.Println("Create a new cache entry for",t)
	templates := []string{ // slice of string
		fmt.Sprintf("../../templates/%s",t),
		"../../templates/base-layout.template.html",
	}
	
	// parse the template
	tmpl, err := template.ParseFiles(templates...) // spread out the slice to all strings inside templates
	if err != nil {
		return err
	}

	// add template to cache
	tc[t] = tmpl

	// return without an error
	return nil
}
*/
// #endregion