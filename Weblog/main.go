package main

import(
	"net/http"
	"os"
	//"io"
	"html/template"
	"log"
//	"fmt"
)


func main() {
	http.HandleFunc("/", makeHandler(handler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

//h1 := func (w http.ResponseWriter,r *http.Request {
//	io.WriteString(w, string)}


type Page struct{
	Body []byte
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, *Page)) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		body, _ := os.ReadFile("test.txt")

		p := &Page{Body: body}
		
		//_, err := o
		//if err != nil{
		//	fmt.Printf("Cant find the file")
		//}
		fn(w, r, p)
	}
}

func handler(w http.ResponseWriter, r *http.Request, p *Page){
	//body, err := os.ReadFile("test.txt")
	//p.Body = body
	//if err != nil {
	//	fmt.Printf("Cant find the file")
	//}
	//str := string(body)
	//io.WriteString(w, str)
	templateRender(w, p)
}

func templateRender(w http.ResponseWriter, p *Page){
	err := templates.ExecuteTemplate(w, "view.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var templates = template.Must(template.ParseFiles("view.html"))

