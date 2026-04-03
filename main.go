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
    simulador "scheduler_os/simulador"
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

    mux.HandleFunc("/procesos", handleProcesos)
    mux.HandleFunc("/procesos/", handleProcesos)
    mux.HandleFunc("/simulaciones", handleSimulaciones)
    mux.HandleFunc("/simulaciones/", handleSimulaciones)


	mux.HandleFunc("/", handleHome)
    
    mux.HandleFunc("/agregar", showCreateForm)
    mux.HandleFunc("/algoritmo", showAlgoritmoForm)
    mux.HandleFunc("/ejecutar", ejecutarSimulacion)
    mux.HandleFunc("/algoritmo/opciones", OpcionesAlgoritmoHandler) //Para que muestre el input quantum si se selecciona RR
    mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./backend/static"))))
    err = http.ListenAndServe(":8080", mux)
    if err != nil {
        fmt.Println("Error al iniciar el servidor:", err)
    }


}
func showAlgoritmoForm(w http.ResponseWriter, r *http.Request) {
    views.AgregarSimulacion().Render(context.Background(), w)

}
func crearSimulacion(w http.ResponseWriter, r *http.Request) {

    if err := r.ParseForm(); err != nil {
        http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
        return
    }

    var s db.CreateSimulacionParams

    // Campos del formulario
    s.Nombre = r.FormValue("nombre")
    algoritmo := r.FormValue("algoritmo")

    switch algoritmo {

    case "first_come_first_serve":
        s.Algoritmo = "FCFS"
        s.Quantum = sql.NullInt32{Valid: false}
    case "shortest_job_first":
        s.Algoritmo = "SJF"
        s.Quantum = sql.NullInt32{Valid: false}
    case "shortest_job_first_preemptive":
        s.Algoritmo = "SJF_PREEMPTIVE"
        s.Quantum = sql.NullInt32{Valid: false}
    case "round_robin":
        s.Algoritmo = "RR"

        quantumStr := r.FormValue("quantum")
        quantumInt, err := strconv.Atoi(quantumStr)
        if err != nil {
            http.Error(w, "Quantum inválido", http.StatusBadRequest)
            return
        }

        s.Quantum = sql.NullInt32{
            Int32: int32(quantumInt),
            Valid: true,
        }

    case "priority":
        s.Algoritmo = "PRIORITY"
        prioridadStr := r.FormValue("prioridad")
        prioridadInt, err := strconv.Atoi(prioridadStr)
        if err != nil {
            http.Error(w, "Prioridad inválida", http.StatusBadRequest)
            return
        }

        s.Prioridad = sql.NullInt32{
            Int32: int32(prioridadInt),
            Valid: true,
        }
    default:
        http.Error(w, "Algoritmo no soportado", http.StatusBadRequest)
        return
    }

    // Valores por defecto 
    s.ProcessTime = 0
    s.ContextSwitches = 0
    s.DispatchLatency = 0
    s.AverageTurnaroundTime = 0
    s.AverageWaitingTime = 0
    s.AverageThroughput = 0

    _, err := queries.CreateSimulacion(context.Background(), s)
    if err != nil {
        log.Println("Error DB:", err)
        http.Error(w, "Error al guardar simulación", http.StatusInternalServerError)
        return
    }
    //TO-DO: ver de hacer que el listar procesos directamente siempre llame al view y muestre los procesos 
    // y el grafo y el listar procesos por estado consulte directamente a las queues.
    /*procesos, err := queries.ListProcess(context.Background())
    if err != nil {
        http.Error(w, "Error al cargar procesos", http.StatusInternalServerError)
        return
    }*/
    lastSimulacion, err := queries.GetLastSimulacion(context.Background())
    if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return

        
    }
    procesos, err := queries.GetProcessesBySimulacion(context.Background(), lastSimulacion.ID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    views.HomeView(procesos).Render(context.Background(), w)

}


func showCreateForm(w http.ResponseWriter, r *http.Request) {
    views.AgregarForm().Render(context.Background(), w)
}
func MostrarSimulaciones(w http.ResponseWriter, r *http.Request) {
    simulaciones, err := queries.ListSimulacion(context.Background())
    if err != nil {
        http.Error(w, "Error al cargar simulaciones: "+err.Error(), http.StatusInternalServerError)
        return
    }
    var procesos []db.Proceso
    procesos, err = queries.ListProcess(context.Background())
    if err != nil {
        http.Error(w, "Error al cargar procesos: "+err.Error(), http.StatusInternalServerError)
        return
    }

    views.ListarSimulaciones(simulaciones, procesos).Render(context.Background(), w)
}
func handleHome(w http.ResponseWriter, r *http.Request) {
    /*
    procesos, err := queries.ListProcess(context.Background())
    if err != nil {
        http.Error(w, "Error al cargar procesos: "+err.Error(), http.StatusInternalServerError)
        return
    }*/
    var procesos = []db.Proceso{}
    lastSimulacion, err := queries.GetLastSimulacion(context.Background())
    if err == nil {
        procesos, err = queries.GetProcessesBySimulacion(context.Background(), lastSimulacion.ID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }
/*
    if procesos == nil {
        procesos = []db.Proceso{}
    }
*/
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

func handleSimulaciones(w http.ResponseWriter, r *http.Request) {

    path := strings.TrimPrefix(r.URL.Path, "/simulaciones")
    pathParts := strings.Split(path, "/")

    switch r.Method {            
    case "GET":
        if (len(pathParts) == 1 && pathParts[0] != "") {
            id := pathParts[0]
            if id != "" {
                id, err := strconv.Atoi(id)
                if err != nil {
                    http.Error(w, "ID inválido", http.StatusBadRequest)
                    return
                }
                
                var procesos []db.Proceso
                simulacion, err := queries.GetSimulacion(context.Background(), int32(id))
                if err != nil {
                    http.Error(w, "Error al cargar la simulacion", http.StatusInternalServerError)
                    return
                }
                procesos, err = queries.GetProcessesBySimulacion(context.Background(), int32(id))
                if err != nil {
                    http.Error(w, "Error al cargar los procesos de la simulacion", http.StatusInternalServerError)
                    return
                }
                views.EstadisticasView(simulacion, procesos).Render(context.Background(), w)
                return
            }
        }
        MostrarSimulaciones(w, r)
        return

    case "POST":
        crearSimulacion(w, r)
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

	//procesos, err := queries.ListProcess(context.Background())
        var procesos = []db.Proceso{}
    lastSimulacion, err := queries.GetLastSimulacion(context.Background())
    if err == nil {
        procesos, err = queries.GetProcessesBySimulacion(context.Background(), lastSimulacion.ID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }

    views.HomeView(procesos).Render(context.Background(), w)
    return

}


func createProceso(w http.ResponseWriter, r *http.Request) {
    var err error
    
    if err = r.ParseForm(); err != nil {
        http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
        return
    }
    
    ultima, err := queries.GetLastSimulacion(context.Background())
    if err != nil {
        http.Error(w, "No hay simulación activa", http.StatusBadRequest)
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

    p.Nombre = r.FormValue("nombre")
    p.Estado = r.FormValue("estado") // Viene del input hidden

    p.Prioridad, err = toInt32("prioridad")
    if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

    p.BurstTime, err = toInt32("burst_time")
    if err != nil { http.Error(w, "El campo 'Burst Time' es inválido.", http.StatusBadRequest); return }

    p.ArrivalTime, err = toInt32("arrival_time")
    if err != nil { http.Error(w, "El campo 'Arrival time' es inválido.", http.StatusBadRequest); return }
    
    p.IDSimulacion = ultima.ID
    
    _, err = queries.CreateProcess(context.Background(), p)
    if err != nil {
        log.Println("Error DB:", err)
        http.Error(w, "Error al guardar simulación", http.StatusInternalServerError)
        return
    }

    listProcesos(w, r)
    return
    //http.Redirect(w, r, "/", http.StatusSeeOther) Elimina la redireccion 
    
}
func deleteProceso(w http.ResponseWriter, r *http.Request, id int32) {
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

// Ejecutar simulacion
func ejecutarSimulacion(w http.ResponseWriter, r *http.Request) {

  
    
    //Busco los parametros de la ultima simulacion creada por el usuario 

    //TO-DO crear la tabla simulacion, crear una query que devuelva la ultima simulacion creada, otra que cree una simulacion nueva y otra query que devuelva todas las simulaciones. 
    simulacion, err := queries.ListSimulacion(context.Background())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if len(simulacion) == 0 {
        http.Error(w, "No hay simulaciones cargadas", http.StatusBadRequest)
        return
    }
    ultima := simulacion[len(simulacion)-1]
    
    newQueue, err := queries.GetProcessesBySimulacion(context.Background(), ultima.ID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    var terminatedQueue []db.Proceso
    var totalTime int

    //Aquí iria la llamada a la funcion que ejecuta la simulacion con los procesos new 
    switch ultima.Algoritmo {
    case "FCFS":
        log.Println("Ejecutando FCFS")
        totalTime, terminatedQueue = simulador.Simulacion(newQueue, simulador.FirstComeFirstServe, queries)
        calcularEstadisticasSimulacion(terminatedQueue, "FCFS", totalTime)
        
    case "SJF":
        log.Println("Ejecutando SJF")
        totalTime, terminatedQueue = simulador.Simulacion(newQueue, simulador.ShortestJobFirst, queries)
        calcularEstadisticasSimulacion(terminatedQueue, "SJF", totalTime)
    case "SJF_PREEMPTIVE":
        log.Println("Ejecutando SJF con desalojo")
        totalTime, terminatedQueue = simulador.Simulacion(newQueue, simulador.ShortestJobFirstPreemptive, queries)
        calcularEstadisticasSimulacion(terminatedQueue, "SJF_PREEMPTIVE", totalTime)
    case "RR":
        log.Println("Ejecutando RR con quantum:", simulacion[len(simulacion)-1].Quantum)
        //rr(newQueue, simulacion[len(simulacion)-1].Quantum)
    default:
        http.Error(w, "Algoritmo no soportado", http.StatusBadRequest)
        return
    }


    //TO-DO: actualizar la tabla simulacion con las estadisticas de la simulacion que se acaba de ejecutar

    fmt.Fprintln(w, "Simulación ejecutada")
}
func OpcionesAlgoritmoHandler(w http.ResponseWriter, r *http.Request) {
    algoritmo := r.URL.Query().Get("algoritmo")

    if algoritmo == "round_robin" {
        views.InputQuantum().Render(r.Context(), w)
        return
    } else if algoritmo == "priority" {
        views.InputPrioridad().Render(r.Context(), w)
        return
    } else {
        // No se necesita ningún campo adicional para otros algoritmos
        w.WriteHeader(http.StatusNoContent) // 204 No Content
        return
    }
}

func calcularEstadisticasSimulacion(terminatedQueue []db.Proceso, algoritmo string, totalTime int) {
    //Funcion que calcula las estadisticas de la simulacion y las guarda en la DB
    //TO-DO calcular avg waiting time, avg turnaround time, avg response time, cpu utilization, throughput
    //y guardar en la tabla estadisticas_simulacion
    turnaround_time := float64(0)
    waiting_time := float64(0)
    contextSwitches := 0
    
    for _, p := range terminatedQueue {
        turnaround_time += float64(p.CompletionTime.Int32 - p.ArrivalTime)
        waiting_time += float64(p.WaitingTime.Int32)
        contextSwitches++;
    }
    averageTurnaroundTime := turnaround_time / float64(len(terminatedQueue))
    averageWaitingTime := waiting_time / float64(len(terminatedQueue))
    averageThroughput := float64(len(terminatedQueue)) / float64(totalTime)


    ultima, err := queries.GetLastSimulacion(context.Background())
    if err != nil {
        log.Println("Error al obtener la última simulación:", err)
        return
    }
    err = queries.UpdateSimulacion(context.Background(), db.UpdateSimulacionParams{
        ID:              ultima.ID,
        Algoritmo:       algoritmo,
        AverageWaitingTime:     averageWaitingTime,
        AverageTurnaroundTime:  averageTurnaroundTime,
        AverageThroughput: averageThroughput,
        ProcessTime:    int32(totalTime),
        ContextSwitches: int32(contextSwitches),
    })
    if err != nil {
        log.Println("Error al actualizar la simulación:", err)
        return
    }
    log.Printf("Simulación %d actualizada con estadísticas.", ultima.ID)
    return
}
