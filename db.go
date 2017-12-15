package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/mediocregopher/radix.v2/pool"
)

// Alumno representation
type Alumno struct {
	Id     string `json:"Id"`
	Alumno string `json:"Alumno"`
	Link   string `json:"Link"`
	Time   string `json:"Time"`
}

var db *pool.Pool

// ErrNoAlumno used when no Alumno is found in the database
var ErrNoAlumno = errors.New("Alumno no encontrado")

func init() {
	var err error
	// Establish a pool of 10 connections to the Redis server listening on
	// port 6379 of the local machine.
	db, err = pool.New("tcp", "localhost:6379", 10)
	if err != nil {
		log.Panic(err)
	}
}

// AddAlumno agrega un link de examen
func AddAlumno(alumno Alumno) (bool, string) {
	if existsAlumno(alumno.Id) {
		return false, "El alumno " + alumno.Alumno + " ya entregó su examen."
	}

	resp := db.Cmd("HMSET", "id:"+alumno.Id, "alumno", alumno.Alumno, "link", alumno.Link, "time", alumno.Time)
	if resp.Err != nil {
		log.Fatal(resp.Err)
		return false, "Error al agregar un alumno"
	}

	fmt.Println("Alumno agregado")
	return true, "Alumno agregado"
}

// GetAlumno retrieves Alumno from id
func GetAlumno(id string) (*Alumno, error) {
	if existsAlumno(id) {
		reply, err := db.Cmd("HGETALL", "id:"+id).Map()
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		return populateAlumno(reply), nil
	}
	return nil, ErrNoAlumno
}

func existsAlumno(id string) bool {
	exists, err := db.Cmd("EXISTS", "id:"+id).Int()
	if err != nil {
		log.Fatal(err)
		return false
	} else if exists == 0 {
		return false
	}

	return true
}

func populateAlumno(reply map[string]string) *Alumno {
	alumno := new(Alumno)
	alumno.Alumno = reply["alumno"]
	alumno.Link = reply["link"]
	alumno.Time = reply["time"]

	return alumno
}

// GetAll retrieves all the Alumnos
func GetAll() []Alumno {
	conn, err := db.Get()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Put(conn)

	reply, err := conn.Cmd("KEYS", "*").List()
	if err != nil {
		log.Fatal(err)
	}
	var alumnos []Alumno

	for _, element := range reply {
		alumno, err := conn.Cmd("HGETALL", element).Map()
		if err != nil {
			log.Fatal(err)
		}

		alumnos = append(alumnos, *populateAlumno(alumno))

	}

	return alumnos
}

//PrintAlumno prints Alumno struct to console
func PrintAlumno(a *Alumno) {
	fmt.Println("alumno:" + a.Alumno + ", link:" + a.Link + " , time:" + a.Time)
}
