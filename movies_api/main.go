package main

import (
	"fmt"
	"strconv"
	"database/sql"
	"io"
	"log"
	"os"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/go-sql-driver/mysql"
	"encoding/json"
)
//used for database
type Director struct{
	ID_D int64 `json:"id_d"`
	FirstName string `json:"firstname"`
	LastName string	`json:"lastname"`
}
// Used for database
type Movie struct{
	ID int64 	`json:"id"`
	Isbn int64	`json:"isbn"`
	Title string	`json:"title"`
	ID_D int64	`json:"id_d"`
}


//type movies []Movie

// Function for Database to golang decleration
func MoviesByID() ([]Movie, error){
	// An movie slice to hold value from returned rows.
	var movie []Movie

	rows, err := db.Query("SELECT * FROM movie")
	if err != nil {
		return nil, fmt.Errorf("MoviesByID have error: %v",err)
	}

	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next(){
		var mov Movie
		if err := rows.Scan(&mov.ID, &mov.Isbn, &mov.Title, &mov.ID_D); err != nil{
			return nil, fmt.Errorf("MoviesByID: %v",err)
		}
		movie = append(movie, mov)
	}
	if err = rows.Err(); err != nil{
		return movie, err
	}
	return movie, nil
}

func DirectorByID() ([]Director, error){
	// like the above decleration
	var director []Director

	rows, err := db.Query("SELECT * FROM director")
	if err != nil {
		return nil, fmt.Errorf("DirectorByID Have error: %v",err)
	}

	defer rows.Close()

	for rows.Next(){
		var dir Director
		if err := rows.Scan(&dir.ID_D, &dir.FirstName, &dir.LastName); err != nil{
			return director, err
		}
		director = append(director, dir)
	}
	if err = rows.Err(); err != nil{
		return director, err
	}
	return director, nil
}

func main(){
	//mux for api
	r := mux.NewRouter()

	Database()

	movies, err := MoviesByID()
	if err != nil{
		log.Fatal(err)
	}
	fmt.Printf("Movies found: %v\n", movies)

	directors, err := DirectorByID()
	if err != nil{
		log.Fatal(err)
	}
	fmt.Printf("Directors found: %v\n", directors)

	

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMoviesByID).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting the server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getMoviesByID(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	//reads the parameters
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	movies, err := MoviesByID()
	if err != nil{
		log.Fatal(err)
	}
	fmt.Printf("Movies found: %v\n", movies)
	
	for i, item := range movies{
		if item.ID == id{
			json.NewEncoder(w).Encode(movies[i])
		}
	}
}

func getMovies(w http.ResponseWriter, r *http.Request){
	//var err error
	movies, err := MoviesByID()
        if err != nil{
                log.Fatal(err)
        }
        fmt.Printf("Movies found: %v\n", movies)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
	// We need a slice for Encode(arg []slice)
}

func deleteMovie(w http.ResponseWriter, r *http.Request){
	//err := error 
	movies, err := MoviesByID()
        if err != nil{
                log.Fatal(err)
        }
        fmt.Printf("Movies found: %v\n", movies)

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil{
		log.Fatal(err)
	}
	for i, item := range movies{
		if item.ID == id{
			movies = append(movies[:i], movies[i+1:]...)
			break
		}
	}

	result, err := db.Prepare("DELETE FROM movie WHERE ID=?")
	if err != nil{
		log.Fatal(err)
	}
	defer result.Close()
	if _, err := result.Exec(id); err != nil{
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(movies)
}

//func makeHandler(fn func(http.ResponseWriter, *http.Request)

func createMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	var movie Movie

	//err := error
	movies, err := MoviesByID()
        if err != nil{
                log.Fatal(err)
        }
        fmt.Printf("Movies found: %v\n", movies)

	body, _ := io.ReadAll(r.Body)
	fmt.Println("Raw request body:", string(body))
	
//	_ = json.NewDecoder(r.Body).Decode(&movie)
//	movies = append(movies, movie)

	if err := json.Unmarshal(body, &movie); err != nil {
        	log.Fatal("JSON Decoding Error:", err)
	}
	fmt.Printf("Attempting to insert: ID=%d, Isbn=%d, Title=%s, ID_D=%d\n", movie.ID, movie.Isbn, movie.Title, movie.ID_D)
	result, err := db.Prepare("INSERT INTO movie(ID, Isbn, Title, ID_D) VALUES (?,?,?,?)")
	if err != nil{
		log.Fatal(err)
	}
	defer result.Close()

	if _, err := result.Exec(movie.ID, movie.Isbn, movie.Title, movie.ID_D); err != nil {
		log.Fatal(err)
	}
	//for _, val := range movies{
	/*result, err := db.Exec("INSERT INTO movie(ID, Isbn, Title) VALUES (?,?,?,?)", movie.ID, movie.Isbn, movie.Title, movie.ID_D)
	if err != nil {
		//return 0, fmt.Errorf("Inserting failed: %v", err)
		fmt.Printf("Error in executing the command INSERT")
	}*/

	/*id, err := result.LastInsertId()
	if err != nil{
		//return 0, fmt.Errorf("Inserting: %v",err)
		fmt.Printf("Failure in command id: %d: %v", id, err)
	}*/
	//return id, nil

	json.NewEncoder(w).Encode(movies)
}

func updateMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	movies, err := MoviesByID()
        if err != nil{
                log.Fatal(err)
        }
        fmt.Printf("Movies found: %v\n", movies)

	var movie Movie
	body, _ := io.ReadAll(r.Body)
	fmt.Println("Raw request body: %v", string(body))

	if err := json.Unmarshal(body, &movie); err != nil {
		log.Fatal("Error Decoding json: ", err)
	}

	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil{
		log.Fatal(err)
	}

	/*for i, item := range movies{
		if item.ID == id{
			movies = append(movies[:i], movies[i+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = id
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movies)
		}else{
			fmt.Printf("item.ID is not equal to id: %v, %v", item.ID, id)
		}
	}*/
	result, err := db.Prepare("UPDATE movie SET ID = ?, Isbn = ?, Title = ?, ID_D = ? WHERE ID = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer result.Close()
	if _, err := result.Exec(movie.ID, movie.Isbn, movie.Title, movie.ID_D, id); err != nil {
		log.Fatal(err)
	}
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movies)
}




var db *sql.DB

func Database(){
	cfg := mysql.Config{
		User: os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net: "tcp",
		Addr: "127.0.0.1:3306",
		DBName: "movies",
		AllowNativePasswords: true,
	}

	//Database handle

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}


