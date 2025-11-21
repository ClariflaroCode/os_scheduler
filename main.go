package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strings"
    "fmt"
	_ "github.com/lib/pq"
	db "scheduler_os/backend/db/sqlc"
    views "scheduler_os/backend/views"
)

var queries *db.Queries

func main() {
    connStr := "postgres://root:root@localhost:5432/myapp?sslmode=disable"

    conn, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Error DB: %v", err)
    }

    if err := conn.Ping(); err != nil {
        log.Fatalf("La base no responde: %v", err)
    }

    queries = db.New(conn)
	mux := http.NewServeMux()

	mux.HandleFunc("/procesos", func(w http.ResponseWriter, r *http.Request) {

		handleProcesos(w, r)
	})


	mux.HandleFunc("/", handleHome)
    mux.HandleFunc("/agregar", showCreateForm)
//    mux.HandleFunc("/estadisticas", showEstadisticas)

  //  mux.HandleFunc("/algoritmo", showAlgoritmoForm)
    mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./backend/static"))))
    err = http.ListenAndServe(":8080", mux)
    if err != nil {
        fmt.Println("Error al iniciar el servidor:", err)
    }


}
func showCreateForm(w http.ResponseWriter, r *http.Request) {
    views.StaticLayout( views.AgregarForm()).Render(context.Background(), w)
}
func handleHome(w http.ResponseWriter, r *http.Request) {
    
    procesos, err := queries.ListProcess(context.Background())
    if err != nil {
        http.Error(w, "Error al cargar procesos: "+err.Error(), http.StatusInternalServerError)
        return
    }

    if procesos == nil {
        procesos = []db.Proceso{}
    }

    views.StaticLayout(views.HomeView(procesos)).Render(context.Background(), w)
}

func handleProcesos(w http.ResponseWriter, r *http.Request) {
    
    path := strings.TrimPrefix(r.URL.Path, "/procesos")
    path = strings.Trim(path, "/")

    switch r.Method {            
    case "GET":
        if path == "" {
            listProcesos(w, r)
            return
        }
        return

    case "POST":
        createProceso(w, r)
        return
        
    default:
        // Si no es OPTIONS y no es ninguno de los métodos permitidos, devuelve 405.
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
    }
}
func listProcesos(w http.ResponseWriter, r *http.Request) {
	procesos, err := queries.ListProcess(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if procesos == nil {
		procesos = []db.Proceso{}
	}

    views.StaticLayout( views.ListarProcesos(procesos)).Render(context.Background(), w)
}


func createProceso(w http.ResponseWriter, r *http.Request) {
    var p db.CreateProcessParams
	


	if p.Estado == "" {
		http.Error(w, "El campo 'estado' es obligatorio", http.StatusBadRequest)
		return
	}

    _, err := queries.CreateProcess(context.Background(), p)  //_, sirve para cuando no vas a usar la variable :D
    if err != nil {
        return
    }
    
}