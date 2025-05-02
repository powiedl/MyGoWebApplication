package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/powiedl/myGoWebApplication/internal/config"
	"github.com/powiedl/myGoWebApplication/internal/models"
)

var testApp config.AppConfig 
var session *scs.SessionManager

type myWriter struct{}

// #region wir definieren jetzt alle Methoden, die es braucht, damit etwas vom Typ http.ResponseWriter ist
func (tw *myWriter)Header() http.Header {
	var h http.Header
	return h
}

func (tw *myWriter)Write(b []byte) (int,error) {
	length :=len(b)
	return length,nil
}

func (tw *myWriter)WriteHeader(i int){
	// wir schreiben keinen Header, aber das Interface verlangt, dass diese Methode exisitert, also legen wir eine leere Methode an
}
// #endregion

func TestMain(m *testing.M) {
	// Data to be available in the session
	gob.Register(models.Reservation{})

	// don't forget to change to tru in Production
	testApp.Basedir = "../../"

	infoLog := log.New(os.Stdout,"[INFO]\t",log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog
	errorLog := log.New(os.Stdout,"[ERROR]\t",log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false
	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}


