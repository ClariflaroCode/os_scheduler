package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strings"
    "fmt"
    "strconv"
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

	/*mux.HandleFunc("/procesos", func(w http.ResponseWriter, r *http.Request) {

		handleProcesos(w, r)
	})*/
    mux.HandleFunc("/procesos", handleProcesos)
    mux.HandleFunc("/procesos/", handleProcesos)



	mux.HandleFunc("/", handleHome)
    
    mux.HandleFunc("/agregar", showCreateForm)
    mux.HandleFunc("/estadisticas", showEstadisticas)

    mux.HandleFunc("/algoritmo", showAlgoritmoForm)
    mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./backend/static"))))
    err = http.ListenAndServe(":8080", mux)
    if err != nil {
        fmt.Println("Error al iniciar el servidor:", err)
    }


}
func showAlgoritmoForm(w http.ResponseWriter, r *http.Request) {
    views.AgregarSimulacion().Render(context.Background(), w)
}
func showEstadisticas(w http.ResponseWriter, r *http.Request) {
    views.EstadisticasView().Render(context.Background(), w)
}
func showCreateForm(w http.ResponseWriter, r *http.Request) {
    views.AgregarForm().Render(context.Background(), w)
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
        estado := r.URL.Query().Get("estado")

        if estado != "" {
            getProcesoByEstado(w, r, estado)
            return
        }

        listProcesos(w, r)

        return

    case "POST":
        createProceso(w, r)
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
func getProcesoByEstado(w http.ResponseWriter, r *http.Request, estado string) {
    procesos, err := queries.ListProcessByEstado(context.Background(), estado)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    views.ListarProcesos(procesos, estado).Render(context.Background(), w)
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

    views.HomeView(procesos).Render(context.Background(), w)
    return

}


func createProceso(w http.ResponseWriter, r *http.Request) {
	
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
        return
    }

    toInt32 := func(key string) (int32, error) {
        valStr := r.FormValue(key)
        valInt, err := strconv.ParseInt(valStr, 10, 32) 
        if err != nil {
            return 0, fmt.Errorf("el campo '%s' es inválido o está vacío", key)
        }
        return int32(valInt), nil
    }

    var p db.CreateProcessParams
    var err error

    p.Nombre = r.FormValue("nombre")
    p.Estado = r.FormValue("estado") // Viene del input hidden

    p.Prioridad, err = toInt32("prioridad")
    if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

    p.BurstTime, err = toInt32("burst_time")
    if err != nil { http.Error(w, "El campo 'Burst Time' es inválido.", http.StatusBadRequest); return }

    p.ArrivalTime, err = toInt32("arrival_time")
    if err != nil { http.Error(w, "El campo 'Arrival time' es inválido.", http.StatusBadRequest); return }
    
    _, err = queries.CreateProcess(context.Background(), p)
    
    if err != nil {
        log.Printf("Error al crear proceso en DB: %v", err)
        http.Error(w, "Error interno al guardar el proceso. Ver logs del servidor.", http.StatusInternalServerError)
        return
    }
    listProcesos(w, r)
    return
    //http.Redirect(w, r, "/", http.StatusSeeOther) Elimina la redireccion 
    
}
func deleteProceso (w http.ResponseWriter, r *http.Request, id int32) {
    err := queries.DeleteProcess(context.Background(), id)
    //Si no hubo error SQL → error = nil
    //Si la query falló → error != nil
    if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} 
    log.Println("BORRANDO ID:", id)

    w.WriteHeader(http.StatusOK) 

     
}