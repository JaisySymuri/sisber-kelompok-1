package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "qwer1234"
	dbName := "Pegawai"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName+"?parseTime=true")
	if err != nil {
		fmt.Println("Can't connect to DB, continue anyway")
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func getCurrentDirectoryName() (string, error) {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Extract the directory's name using filepath.Base
	dirName := filepath.Base(cwd)
	return dirName, nil
}

// Helper function to scan all rows into an Employee struct
func scanEmployee(row *sql.Rows) (Employee, error) {
	var emp Employee
	err := row.Scan(&emp.ID_Pegawai, &emp.NIK, &emp.Nama, &emp.Username, &emp.Password, &emp.Alamat, &emp.Tempat_Lahir, &emp.Tanggal_Lahir, &emp.No_HP, &emp.Pekerjaan, &emp.Gender, &emp.updated_at, &emp.created_at)
	return emp, err
}

// Helper function to scan a single row into an Employee struct
func scanEmployeeRow(row *sql.Row) (Employee, error) {
	var emp Employee
	err := row.Scan(&emp.ID_Pegawai, &emp.NIK, &emp.Nama, &emp.Username, &emp.Password, &emp.Alamat, &emp.Tempat_Lahir, &emp.Tanggal_Lahir, &emp.No_HP, &emp.Pekerjaan, &emp.Gender, &emp.updated_at, &emp.created_at)
	return emp, err
}

func parseEmployeeInput(r *http.Request) (Employee, error) {
	var input Employee
	var err error

	input.NIK = r.FormValue("NIK")
	input.Nama = r.FormValue("Nama")
	input.Username = r.FormValue("Username")
	input.Password = r.FormValue("Password")
	input.Alamat = r.FormValue("Alamat")
	input.Tempat_Lahir = r.FormValue("Tempat_Lahir")

	tglLahirStr := r.FormValue("Tanggal_Lahir")
	input.Tanggal_Lahir, err = time.Parse("2006-01-02", tglLahirStr)
	if err != nil {
		return Employee{}, err
	}

	input.No_HP = r.FormValue("No_HP")

	input.Pekerjaan = r.FormValue("Pekerjaan")
	input.Gender = r.FormValue("Gender")
	input.ID_Pegawai = r.FormValue("uid")

	return input, nil
}


func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM Employee9998 ORDER BY ID_Pegawai DESC")
	if err != nil {
		panic(err.Error())
	}
	res := []Employee{}
	for selDB.Next() {
		emp, err := scanEmployee(selDB)
		if err != nil {
			panic(err.Error())
		}
		res = append(res, emp)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}


func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nID_Pegawai := r.URL.Query().Get("ID_Pegawai")
	row := db.QueryRow("SELECT * FROM Employee9998 WHERE ID_Pegawai=?", nID_Pegawai)
	
	emp, err := scanEmployeeRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		panic(err.Error())
	}	
	tmpl.ExecuteTemplate(w, "Show", emp)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nID_Pegawai := r.URL.Query().Get("ID_Pegawai")
	row := db.QueryRow("SELECT * FROM Employee9998 WHERE ID_Pegawai=?", nID_Pegawai)

	emp, err := scanEmployeeRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		panic(err.Error())
	}		
	tmpl.ExecuteTemplate(w, "Edit", emp)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {	
	db := dbConn()
	if r.Method == "POST" {
		input, err := parseEmployeeInput(r)
		if err != nil {
			http.Error(w, "Invalid input data", http.StatusBadRequest)
			return
		}

		insForm, err := db.Prepare("INSERT INTO Employee9998(NIK, Nama, Username, Password, Alamat, Tempat_Lahir, Tanggal_Lahir, No_HP, Pekerjaan, Gender) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(input.NIK, input.Nama, input.Username, input.Password, input.Alamat, input.Tempat_Lahir, input.Tanggal_Lahir, input.No_HP, input.Pekerjaan, input.Gender)
		log.Println("INSERT: NIK:", input.NIK, "| Nama:", input.Nama, "| Username:", input.Username, "| Password:", input.Password, "| Alamat:", input.Alamat, "| Tempat_Lahir:", input.Tempat_Lahir, "| Tanggal_Lahir:", input.Tanggal_Lahir, "| No_HP:", input.No_HP, "| Pekerjaan:", input.Pekerjaan, "| Gender:", input.Gender)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		input, err := parseEmployeeInput(r)
		if err != nil {
			http.Error(w, "Invalid input data", http.StatusBadRequest)
			return
		}
		
		insForm, err := db.Prepare("UPDATE Employee9998 SET NIK=?, Nama=?, Username=?, Password=?, Alamat=?, Tempat_Lahir=?, Tanggal_Lahir=?, No_HP=?, Pekerjaan=?, Gender=? WHERE ID_Pegawai=?")
		if err != nil {
			panic(err.Error())
		}
		_, err = insForm.Exec(input.NIK, input.Nama, input.Username, input.Password, input.Alamat, input.Tempat_Lahir, input.Tanggal_Lahir, input.No_HP, input.Pekerjaan, input.Gender, input.ID_Pegawai)
		if err != nil {
			panic(err.Error())
		}
		log.Println("UPDATE: NIK:", input.NIK, "| Nama:", input.Nama, "| Username:", input.Username, "| Password:", input.Password, "| Alamat:", input.Alamat, "| Tempat_Lahir:", input.Tempat_Lahir, "| Tanggal_Lahir:", input.Tanggal_Lahir, "| No_HP:", input.No_HP, "| Pekerjaan:", input.Pekerjaan, "| Gender:", input.Gender, "|ID_PEGAWAI:", input.ID_Pegawai)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	emp := r.URL.Query().Get("ID_Pegawai")
	delForm, err := db.Prepare("DELETE FROM Employee9998 WHERE ID_Pegawai=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

// GET /employees - Get all employees
func GetEmployees(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	db := dbConn()
	selDB, err := db.Query("SELECT * FROM Employee9998 ORDER BY ID_Pegawai DESC")
	if err != nil {
		panic(err.Error())
	}
	res := []Employee{}
	for selDB.Next() {
		emp, err := scanEmployee(selDB)
		if err != nil {
			panic(err.Error())
		}
		res = append(res, emp)
	}	
	defer db.Close()	

	// Convert employees slice to JSON and send it in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// GET /employees/{id} - Get an employee by ID_Pegawai
func GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the ID_Pegawai from the URL path
	idPegawai := r.URL.Path[len("/employees/"):]
	if idPegawai == "" {
		http.Error(w, "Missing ID_Pegawai parameter", http.StatusBadRequest)
		return
	}

	db := dbConn()
	defer db.Close()

	// Query the database to get the employee with the specified ID_Pegawai
	row := db.QueryRow("SELECT * FROM Employee9998 WHERE ID_Pegawai = ?", idPegawai)

	emp, err := scanEmployeeRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Employee not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Send the employee data as JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(emp); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}

func (e *Employee) UnmarshalJSON(b []byte) error {
	type EmployeeAlias Employee // Create an alias to avoid recursive call of UnmarshalJSON
	var emp EmployeeAlias

	err := json.Unmarshal(b, &emp)
	if err != nil {
		return err
	}

	e.ID_Pegawai = emp.ID_Pegawai
	e.NIK = emp.NIK
	e.Nama = emp.Nama
	e.Username = emp.Username
	e.Password = emp.Password
	e.Alamat = emp.Alamat
	e.Tempat_Lahir = emp.Tempat_Lahir
	e.No_HP = emp.No_HP
	e.Pekerjaan = emp.Pekerjaan
	e.Gender = emp.Gender
	e.updated_at = emp.updated_at
	e.created_at = emp.created_at

	return nil
}

func (j *JSONEmployee) ToEmployee() (*Employee, error) {
	parsedTime, err := time.Parse("2006-01-02", j.Tanggal_Lahir)
	if err != nil {
		return nil, err
	}

	return &Employee{
		ID_Pegawai:    j.ID_Pegawai,
		NIK:           j.NIK,
		Nama:          j.Nama,
		Username:      j.Username,
		Password:      j.Password,
		Alamat:        j.Alamat,
		Tempat_Lahir:  j.Tempat_Lahir,
		Tanggal_Lahir: parsedTime,
		No_HP:         j.No_HP,
		Pekerjaan:     j.Pekerjaan,
		Gender:        j.Gender,
	}, nil
}

func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	db := dbConn()
	defer db.Close()

	// Parse the JSON request body into a JSONEmployee object
	var jsonEmp JSONEmployee
	err := json.NewDecoder(r.Body).Decode(&jsonEmp)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert the JSONEmployee to Employee
	emp, err := jsonEmp.ToEmployee()
	if err != nil {
		fmt.Println("Error converting to Employee:", err)
		http.Error(w, "Invalid date format for Tanggal_Lahir", http.StatusBadRequest)
		return
	}

	fmt.Println("Received Employee data:", emp)

	// Perform the database insertion
	_, err = db.Exec("INSERT INTO Employee9998 (NIK, Nama, Username, Password, Alamat, Tempat_Lahir, Tanggal_Lahir, No_HP, Pekerjaan, Gender, updated_at, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		emp.NIK, emp.Nama, emp.Username, emp.Password, emp.Alamat, emp.Tempat_Lahir, emp.Tanggal_Lahir, emp.No_HP, emp.Pekerjaan, emp.Gender, time.Now(), time.Now())
	if err != nil {
		http.Error(w, "Failed to insert employee", http.StatusInternalServerError)
		return
	}

	// Send a success response
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Employee created successfully")
}

func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	
	// Get the ID from the URL path variable
	vars := mux.Vars(r)
	id := vars["id"]

	// Parse the JSON payload from the request body into a JSONEmployee struct
	var empData JSONEmployee
	err := json.NewDecoder(r.Body).Decode(&empData)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Convert the Tanggal_Lahir string to time.Time format
	tanggalLahir, err := time.Parse("2006-01-02", empData.Tanggal_Lahir)
	if err != nil {
		http.Error(w, "Invalid date format for Tanggal_Lahir", http.StatusBadRequest)
		return
	}

	// Create a new Employee struct with the parsed data
	employee := Employee{
		ID_Pegawai:    id,
		NIK:           empData.NIK,
		Nama:          empData.Nama,
		Username:      empData.Username,
		Password:      empData.Password,
		Alamat:        empData.Alamat,
		Tempat_Lahir:  empData.Tempat_Lahir,
		Tanggal_Lahir: tanggalLahir,
		No_HP:         empData.No_HP,
		Pekerjaan:     empData.Pekerjaan,
		Gender:        empData.Gender,
	}


	db := dbConn()
	defer db.Close()

	// Execute the SQL UPDATE query
	stmt, err := db.Prepare("UPDATE Employee9998 SET NIK=?, Nama=?, Username=?, Password=?, Alamat=?, Tempat_Lahir=?, Tanggal_Lahir=?, No_HP=?, Pekerjaan=?, Gender=?, updated_at=NOW() WHERE ID_Pegawai=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		employee.NIK,
		employee.Nama,
		employee.Username,
		employee.Password,
		employee.Alamat,
		employee.Tempat_Lahir,
		employee.Tanggal_Lahir,
		employee.No_HP,
		employee.Pekerjaan,
		employee.Gender,
		employee.ID_Pegawai,
	)
	if err != nil {
		// Handle MySQL-specific error (e.g., duplicate entry)
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			http.Error(w, fmt.Sprintf("MySQL Error: %s", mysqlErr.Message), http.StatusInternalServerError)
			return
		}

		// Handle generic database error
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	// Send a success response to the client
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Employee with ID_Pegawai=%s updated successfully", id)
}

// DeleteEmployee is the handler function for deleting an employee from the database
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the URL path variable
	vars := mux.Vars(r)
	id := vars["id"]

	
	db := dbConn()
	defer db.Close()

	// Execute the SQL DELETE query
	stmt, err := db.Prepare("DELETE FROM Employee9998 WHERE ID_Pegawai=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		// Handle MySQL-specific error (e.g., foreign key constraint violation)
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			http.Error(w, fmt.Sprintf("MySQL Error: %s", mysqlErr.Message), http.StatusInternalServerError)
			return
		}

		// Handle generic database error
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	// Send a success response to the client
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Employee with ID_Pegawai=%s deleted successfully", id)
}




func main() {
	dirName, err := getCurrentDirectoryName()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	// Print the directory's name
	fmt.Println("Directory name:", dirName)
	fmt.Println("Server started on: http://localhost:9998")
	r := mux.NewRouter()

	// Register your other handlers here
	r.HandleFunc("/", Index)
	r.HandleFunc("/show", Show)
	r.HandleFunc("/new", New)
	r.HandleFunc("/edit", Edit)
	r.HandleFunc("/insert", Insert)
	r.HandleFunc("/update", Update)
	r.HandleFunc("/delete", Delete)
	r.HandleFunc("/employees", GetEmployees)
	r.HandleFunc("/employees/{id}", GetEmployeeByID)
	r.HandleFunc("/create", CreateEmployee)

	// Register the UpdateEmployee handler for the "/updateemp/{id}" route
	r.HandleFunc("/updateemp/{id}", UpdateEmployee).Methods("PUT")
	r.HandleFunc("/deleteemp/{id}", DeleteEmployee).Methods("DELETE")

	http.Handle("/", r)
	http.ListenAndServe(":9998", nil)
}


type Employee struct {
	ID_Pegawai    string
	NIK           string
	Nama          string
	Username      string
	Password      string
	Alamat        string
	Tempat_Lahir  string
	Tanggal_Lahir time.Time 
	No_HP         string
	Pekerjaan     string
	Gender        string
	updated_at time.Time
	created_at time.Time
}

// JSONEmployee is a struct used only for JSON decoding
type JSONEmployee struct {
	ID_Pegawai    string
	NIK           string
	Nama          string
	Username      string
	Password      string
	Alamat        string
	Tempat_Lahir  string
	Tanggal_Lahir string `json:"Tanggal_Lahir"`
	No_HP         string
	Pekerjaan     string
	Gender        string
}
