package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"strconv"
	_ "github.com/lib/pq"
	db "scheduler_os/db/sqlc"


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

	mux.HandleFunc("/api/procesos", func(w http.ResponseWriter, r *http.Request) {

		handleProcesos(w, r)
	})

	mux.HandleFunc("/api/procesos/", func(w http.ResponseWriter, r *http.Request) {

		handleProcesos(w, r)
	})


	mux.Handle("/", http.FileServer(http.Dir("./frontend")))

	log.Fatal(http.ListenAndServe(":8080", enableCORS(mux)))

}
func handleProcesos(w http.ResponseWriter, r *http.Request) {
    // 💡 DIAGNÓSTICO: Esto te mostrará qué método está viendo realmente el servidor.
    log.Printf("Entre al handleProcesos. Método de la solicitud: %s", r.Method)
    
    path := strings.TrimPrefix(r.URL.Path, "/api/procesos")
    path = strings.Trim(path, "/")

    switch r.Method {
    case http.MethodOptions: 
        // Si el método es OPTIONS, ya se han puesto los encabezados CORS en enableCORS,
        // solo necesitamos devolver un 200 OK y terminar.
        w.WriteHeader(http.StatusOK)
        return
        
    case "GET":
        if path == "" {
            listProcesos(w, r)
            return
        }
        id, err := strconv.Atoi(path)
        if err != nil {
            http.Error(w, "ID inválido", http.StatusBadRequest)
            return
        }
        getProceso(w, r, int32(id))
        return

    case "POST":
        if path != "" {
            http.Error(w, "POST no admite ID en la URL", http.StatusBadRequest)
            return
        }
        createProceso(w, r)
        return

    case "PUT":
        if path == "" {
            http.Error(w, "Falta ID", http.StatusBadRequest)
            return
        }
        id, err := strconv.Atoi(path)
        if err != nil {
            http.Error(w, "ID inválido", http.StatusBadRequest)
            return
        }
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


	writeJSON(w, procesos, http.StatusOK)
}
func enableCORS(next http.Handler) http.Handler {
    allowedOrigins := map[string]bool{
        "http://127.0.0.1:5500": true,
        "http://localhost:5500": true,
    }

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Obtener el origen de la solicitud
        origin := r.Header.Get("Origin")

        // Si el origen está en nuestra lista blanca, lo devolvemos como permitido
        if allowedOrigins[origin] {
            w.Header().Set("Access-Control-Allow-Origin", origin)
        } else {
            // Si el origen no está permitido (ej: una IP externa), puedes omitir el header
            // o usar una política de denegación más explícita.
            // Para localhost/127.0.0.1, esta rama no debería ejecutarse.
        }

        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        w.Header().Set("Access-Control-Max-Age", "86400")

        next.ServeHTTP(w, r)
    })
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
