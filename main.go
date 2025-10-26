package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"strconv"
	_ "github.com/lib/pq"
	db "ejercicio.com/servidor-go/db/sqlc"

)

var queries *db.Queries


func main() {
	connStr := "postgres://root:root@localhost:5432/myapp?sslmode=disable"

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	err = conn.Ping()
	if err != nil {
		log.Fatalf("La base de datos no responde: %v", err)
	}


	queries = db.New(conn)

	// API 
	http.HandleFunc("/api/procesos", handleProcesos)
	http.HandleFunc("/api/procesos/", handleProcesos)


	http.Handle("/", http.FileServer(http.Dir(".")))

	fmt.Println("Servidor escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func handleProcesos(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entre al handleProcesos")
    path := strings.TrimPrefix(r.URL.Path, "/api/procesos")
    path = strings.Trim(path, "/")
	log.Printf("Path final: '%s'", path)
    switch r.Method {
    case "GET":
        if path == "" {
			log.Printf("Entre al listar procesos")
            listProcesos(w, r)
            return
        }
        id, err := strconv.Atoi(path)
        if err != nil {
            http.Error(w, "ID inválido", http.StatusBadRequest)
            return
        }
		log.Printf("Entre al listar proceso")
        getProceso(w, r, int32(id))
        return

    case "POST":
		log.Printf("Entre al POST")
        if path != "" {
            http.Error(w, "POST no admite ID en la URL", http.StatusBadRequest)
            return
        }
		log.Printf("me VOY AL CREATE")
        createProceso(w, r)
        return

    case "PUT":
        if path == "" {
            http.Error(w, "Falta ID", http.StatusBadRequest)
            return
        }
		log.Printf("Pase el primer if del put")
        id, err := strconv.Atoi(path)
        if err != nil {
            http.Error(w, "ID inválido", http.StatusBadRequest)
            return
        }
		log.Printf("Entrare a updateproceso")
        updateProceso(w, r, int32(id))
        return

    case "DELETE":
        if path == "" {
            http.Error(w, "Falta ID", http.StatusBadRequest)
            return
        }
        id, err := strconv.Atoi(path)
        if err != nil {
            http.Error(w, "ID inválido", http.StatusBadRequest)
            return
        }
        deleteProceso(w, r, int32(id))
        return

    default:
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
    }
}

func listProcesos(w http.ResponseWriter, r *http.Request) {
	procesos, err := queries.ListProcess(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, procesos, http.StatusOK)
}

func getProceso(w http.ResponseWriter, r *http.Request, id int32) {
	proceso, err := queries.GetProcess(context.Background(), id)
	if err == sql.ErrNoRows {
		http.Error(w, "No encontrado", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, proceso, http.StatusOK)
}

func createProceso(w http.ResponseWriter, r *http.Request) {
    var p db.CreateProcessParams
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("Error al decodificar JSON: %v", err)
		http.Error(w, "JSON inválido: " + err.Error(), http.StatusBadRequest)
		return
	}


	if p.Estado == "" {
		http.Error(w, "El campo 'estado' es obligatorio", http.StatusBadRequest)
		return
	}

	newP, err := queries.CreateProcess(context.Background(), p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, newP, http.StatusCreated)
}

func updateProceso(w http.ResponseWriter, r *http.Request, id int32) {
    var p db.UpdateProcessParams
	log.Printf("pase la verga de los params con el update")
    if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
        http.Error(w, "JSON inválido", http.StatusBadRequest)
        return
    }
	log.Printf("esquive el error del decoder en el update")
    p.ID = id

    err := queries.UpdateProcess(context.Background(), p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	log.Printf("esquive el segundo if del update")
    // Recuperar el proceso actualizado para devolverlo como JSON
    proceso, err := queries.GetProcess(context.Background(), id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("✅ Proceso %d actualizado correctamente", id)
    writeJSON(w, proceso, http.StatusOK)
}


func deleteProceso(w http.ResponseWriter, r *http.Request, id int32) {
	err := queries.DeleteProcess(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
