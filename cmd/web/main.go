package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/powiedl/myGoWebApplication/pkg/config"
	"github.com/powiedl/myGoWebApplication/pkg/handlers"
	"github.com/powiedl/myGoWebApplication/pkg/render"
)

const portNumber = 8080

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
  var app config.AppConfig
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatalln("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false // disable cache

	repo := handlers.NewRepo(&app) // create a new repository "based on" app
	handlers.NewHandlers(repo)

	render.NewTemplates(&app) // call render.NewTemplates with the address of the app variable (which means, that the parameter is a pointer)
	
	// #region 4-34
	http.HandleFunc("/",handlers.Repo.Home)
	http.HandleFunc("/about",handlers.Repo.About)
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
	_ = http.ListenAndServe(fmt.Sprintf(":%d",portNumber),nil)
}