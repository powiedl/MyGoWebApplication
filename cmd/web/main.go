package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/powiedl/myGoWebApplication/internal/config"
	"github.com/powiedl/myGoWebApplication/internal/driver"
	"github.com/powiedl/myGoWebApplication/internal/handlers"
	"github.com/powiedl/myGoWebApplication/internal/helpers"
	"github.com/powiedl/myGoWebApplication/internal/models"
	"github.com/powiedl/myGoWebApplication/internal/render"
)

const portNumber = 8080
var app config.AppConfig
var session *scs.SessionManager

var infoLog *log.Logger
var errorLog *log.Logger

// #region 4-26
/*
// Home is the handler for the home page
func Home(w http.ResponseWriter,r *http.Request) {
	_,_ = fmt.Fprintf(w,"This is the homepage")
}

// About is the handler for the about page
func About(w http.ResponseWriter,r *http.Request) {
	owner,saying := getData()
	sum :=addValues(2,3)
	_, _ = fmt.Fprintf(w, fmt.Sprintf("This is the about page of %s, \n I like to say %s, \nand as a side note, 2+3 is %d",owner,saying,sum))
}

// Divide is the handler for the divide page
func Divide(w http.ResponseWriter, r *http.Request) {
	x:=100.0
	y:=4.0
	result,err := divide(x,y)
	if err != nil {
		_,_ = fmt.Fprintf(w,fmt.Sprintf("ERROR! Cannot divide %d by %d! Error returned:%s\n",x,y,err))
		return
	}
	_,_ = fmt.Fprintf(w,fmt.Sprintf("The result of the division of %.2f by %.2f is %.2f\n",x,y,result))

}

// Divide0 is the handler for the 0divide page - which result in an error
func Divide0(w http.ResponseWriter, r *http.Request) {
	x:=100.0
	y:=0.0
	result,err := divide(x,y)
	if err != nil {
		_,_ = fmt.Fprintf(w,fmt.Sprintf("ERROR! Cannot divide %f by %f! Error returned:%s\n",x,y,err))
		return
	}
	_,_ = fmt.Fprintf(w,fmt.Sprintf("The result of the division of %.2f by %.2f is %.2f\n",x,y,result))

}

// getData returns a name and a saying
func getData() (string,string) {
	o := "Rick Sanchez"
	s := "Wubba Lubba Dup Dup!"
	return o,s
}

// addValues adds two integers and returns the sum
func addValues(x, y int) int {
	return x+y
}

func divide(x,y float64) (float64,error) {
	if y==0 {
		err := errors.New("Divisor is zero!")
		return 0,err
	}
	r := x/y
	return r,nil
}
*/
// #endregion

func main() {
	db,err := run()
	if err != nil {
		log.Fatal("main did not start correctly, EPIC FAIL")
	}
	defer db.SQL.Close()
	defer close(app.MailChan)
	
	// #region 14-133 - ein erstes Testmail
	// from := "rick@c-137.universe"
	// auth := smtp.PlainAuth("",from,"","localhost")
	// err = smtp.SendMail("localhost:1025",auth,from,[]string{"morty@lost-in-roy.game"},[]byte("This is my first email..."))
	// if err != nil {
	// 	log.Println("Error while sending the mail:",err)
	// }
	// #endregion
/*
// #region 4-34
  // Mit 5-38 wird das nicht mehr benötigt, weil das Definieren der Routen in routes.go "übersiedelt" ist (und ein eigenes Routing Package verwendet wird)
	http.HandleFunc("/",handlers.Repo.Home)
	http.HandleFunc("/about",handlers.Repo.About)
	// #endregion
*/

	fmt.Println("Starting E-Mail listener")
	listenForMail()

	// #region 14-134 - ein erstes Testamil mit Go Simple Mail
	/*
	msg := models.MailData{
		From: "sender@mail.local",
		To: "powidl@mail.local",
		Subject: "First Mail via Go simple-mail",
		Content: "This is the content of the first mail. Greetings ...",
	}
	app.MailChan <- msg
  */
	// #endregion
	
	// #region bis inkl 4-33
	/*
	http.HandleFunc("/",handlers.Home)
	http.HandleFunc("/about",handlers.About)
	*/
	// #endregion
	// #region 4-26 inside main
	//http.HandleFunc("/divide",Divide)
	//http.HandleFunc("/0divide",Divide0)
	// #endregion
	
	fmt.Printf("Starting application on port %d\n",portNumber)
	/* // bis 5-38 notwendig
	_ = http.ListenAndServe(fmt.Sprintf(":%d",portNumber),nil)
	*/
	// #region ab 5-38 - eigenen Server verwenden
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d",portNumber),
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalln("ERROR! Cannot start server, details=",err)
	}
}

func run() (*driver.DB,error) {
		// Data to be available in the session
		gob.Register(models.Reservation{})
		gob.Register(models.User{})
		gob.Register(models.Bungalow{})
		gob.Register(models.Restriction{})

		mailChan := make(chan models.MailData) // bidirectional channel for mail
		app.MailChan = mailChan

		// don't forget to change to tru in Production
		app.InProduction = false
		app.Basedir = "../../"

		infoLog = log.New(os.Stdout,"[INFO]\t",log.Ldate|log.Ltime)
		app.InfoLog = infoLog
		errorLog = log.New(os.Stdout,"[ERROR]\t",log.Ldate|log.Ltime|log.Lshortfile)
		app.ErrorLog = errorLog
	
		session = scs.New()
		session.Lifetime = 24 * time.Hour
		session.Cookie.Persist = true
		session.Cookie.SameSite = http.SameSiteLaxMode
		session.Cookie.Secure = app.InProduction
		app.Session = session

		// connect to database
		log.Println("Connecting to database ...")
		db, err := driver.ConnectSQL("host=localhost port=5432 dbname=MyGoWebApplication_development user=postgres password="+config.DbPassword)
		if err != nil {
			log.Fatal("Cannot connect to the database, error=",err)
		}
		log.Println("Successfully connected to the database")
	  
		tc, err := render.CreateTemplateCache(&app)
		if err != nil {
			log.Fatalln("cannot create template cache")
			return nil,err
		}
	
		app.TemplateCache = tc
		app.UseCache = false // disable cache
		//log.Println("app.TemplateCache",app.TemplateCache)
		//log.Println("app.Basedir",app.Basedir)
	
		repo := handlers.NewRepo(&app,db) // create a new repository "based on" app and the database connection pool
		handlers.NewHandlers(repo)
	
		render.NewRenderer(&app) // call render.NewTemplates with the address of the app variable (which means, that the parameter is a pointer)
		
		helpers.NewHelpers(&app)
		// everything went fine, return nil
		return db,nil
}