package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func connectionDB() (connection *sql.DB) {
	Driver := "mysql"
	User := "root"
	Password := ""
	Name := "sistema"

	connection, err := sql.Open(Driver, User+":"+Password+"@tcp(127.0.0.1)/"+Name)
	if err != nil {
		panic(err.Error())
	}
	return connection
}

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {
	http.HandleFunc("/", Start)
	http.HandleFunc("/crear", crear)
	http.HandleFunc("/insertar", insertar)

	http.HandleFunc("/borrar", Borrar)

	http.HandleFunc("/editar", Editar)

	http.HandleFunc("/actualizar", Actualizar)

	log.Println("Servidor Iniciado")

	http.ListenAndServe(":8080", nil)
}

type Empleado struct {
	Id     int
	Nombre string
	Correo string
}

func Start(w http.ResponseWriter, r *http.Request) {

	connectionStarted := connectionDB()

	readData, err := connectionStarted.Query("SELECT * FROM empleados")

	if err != nil {
		panic(err.Error())
	}
	empleado := Empleado{}
	arregloEmpleado := []Empleado{}

	for readData.Next() {
		var id int
		var nombre, correo string
		err = readData.Scan(&id, &nombre, &correo)

		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

		arregloEmpleado = append(arregloEmpleado, empleado)

	}
	fmt.Println(arregloEmpleado)

	//fmt.Fprintf(w, "Primer goLang")
	plantillas.ExecuteTemplate(w, "inicio", arregloEmpleado)
}
func crear(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Primer goLang")
	plantillas.ExecuteTemplate(w, "crear", nil)
}
func Borrar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")
	connectionStarted := connectionDB()

	deleteData, err := connectionStarted.Prepare("DELETE FROM empleados WHERE id=?")

	if err != nil {
		panic(err.Error())
	}
	deleteData.Exec(idEmpleado)

	http.Redirect(w, r, "/", 301)
}
func insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		connectionStarted := connectionDB()

		insertData, err := connectionStarted.Prepare("INSERT INTO empleados(nombre, correo) VALUES (?,?)")

		if err != nil {
			panic(err.Error())
		}
		insertData.Exec(nombre, correo)

		http.Redirect(w, r, "/", 301)
	}
}
func Editar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")

	connectionStarted := connectionDB()

	selecteData, err := connectionStarted.Query("SELECT * FROM empleados WHERE id=?", idEmpleado)

	empleado := Empleado{}
	for selecteData.Next() {
		var id int
		var nombre, correo string
		err = selecteData.Scan(&id, &nombre, &correo)

		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo
	}
	fmt.Println(empleado)
	plantillas.ExecuteTemplate(w, "editar", empleado)
}
func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		connectionStarted := connectionDB()

		updateData, err := connectionStarted.Prepare("UPDATE empleados SET nombre=?, correo=? WHERE id=?")

		if err != nil {
			panic(err.Error())
		}
		updateData.Exec(nombre, correo, id)

		http.Redirect(w, r, "/", 301)
	}
}
