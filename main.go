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
    procesos, err := queries.ListProcess(context.Background())
    if err != nil {
        http.Error(w, "Error al cargar procesos", http.StatusInternalServerError)
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

    views.ListarSimulaciones(simulaciones).Render(context.Background(), w)
}
func handleHome(w http.ResponseWriter, r *http.Request) {
    
    procesos, err := queries.ListProcess(context.Background())
    if err != nil {
        http.Error(w, "Error al cargar procesos: "+err.Error(), http.StatusInternalServerError)
        return
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
	procesos, err := queries.ListProcess(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
/*
	if procesos == nil {
		procesos = []db.Proceso{}
	}*/

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

    //Las simulaciones son de términos cortos/tiempos cortos o short terms. 



    //obtengo los procesos en estado "new" que son los creados por el usuario
    newQueue, err := queries.ListProcessByEstado(context.Background(), "new")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
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

    var terminatedQueue []db.Proceso
    var totalTime int

    //Aquí iria la llamada a la funcion que ejecuta la simulacion con los procesos new 
    switch ultima.Algoritmo {
    case "FCFS":
        log.Println("Ejecutando FCFS")
        totalTime, terminatedQueue = fcfs(newQueue)
        calcularEstadisticasSimulacion(terminatedQueue, "FCFS", totalTime)
    case "SJF":
        log.Println("Ejecutando SJF")
        //sjf(newQueue)
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
func fcfs(newQueue []db.Proceso) (int, []db.Proceso){ //NOTA: los primeros parentesis son los parametros de entrada, los segundos los de salida
    //Funcion que ejecuta el algoritmo FCFS
    //Acá recibo la cola new y la voy procesando, miro los arrival time, se simulan los ciclos de reloj, y cuando un proceso llega a su arrival time se cambia su estado a ready
    clock := 0 //Empiezo el ciclo de reloj en 0

    //Mientras que haya procesos en new o ready o running o waiting debo seguir simulando. La simulacion termina cuando todos los procesos estan en terminated. 
    readyQueue := []*db.Proceso{} //Debe ser puntero porque con esto voy a modificar los procesos en la db. 
    runningProceso := (*db.Proceso)(nil) //Puntero a proceso en estado running, inicialmente nil
    waitingQueue := []db.Proceso{}
    terminatedQueue := []db.Proceso{}



    //Mientras que haya procesos en new o ready o running o waiting debo seguir simulando. La simulacion termina cuando todos los procesos estan en terminated.
    for len(newQueue) > 0 || len(readyQueue) > 0 || runningProceso != nil || len(waitingQueue) > 0 {
        for i := 0; i < len(newQueue); {
            if newQueue[i].ArrivalTime == int32(clock) {
                p := newQueue[i]
                readyQueue = append(readyQueue, &p)
                newQueue = append(newQueue[:i], newQueue[i+1:]...)
                continue   // esto es para que no incremente el i, ya que al eliminar un elemento, el siguiente elemento se mueve a la posición actual. 
            }
            i++
        }


        //Tomo el primer elemento de la lista de ready y lo paso a running si está disponible para ejecutar
        if len(readyQueue) > 0 && runningProceso == nil {
            runningProceso = readyQueue[0]
            runningProceso.Estado = "running"
            //updateProcesoEstado(readyQueue[0].ID, "running") //actualizo en la DB el estado del proceso
            readyQueue = readyQueue[1:] //elimino el primer elemento de la lista de ready, la ready queue se volvio la ready queue 
            // empezando desde el segundo elemento hasta el final por eso 1:. En go la longitud se indica con min:max
        }
        // Incremento de waiting time
        for i := range readyQueue {
            readyQueue[i].WaitingTime = sql.NullInt32{
                Int32: readyQueue[i].WaitingTime.Int32 + 1,
                Valid: true,
            }
        }
        if runningProceso != nil {
            runningProceso.BurstTime--

            if runningProceso.BurstTime == 0 {
                //updateProcesoEstado(runningProceso.ID, "terminated") //actualizo en la DB el estado del proceso
                runningProceso.Estado = "terminated"

                runningProceso.CompletionTime = sql.NullInt32{
                    Int32: int32(clock),
                    Valid: true,
                }


                terminatedQueue = append(terminatedQueue, *runningProceso)

                queries.UpdateProcess(context.Background(), db.UpdateProcessParams{
                    ID:             runningProceso.ID,
                    WaitingTime:    runningProceso.WaitingTime,
                    CompletionTime: runningProceso.CompletionTime,
                    Estado:         "terminated",
                })

                runningProceso = nil //libero la CPU
            } 
        }
        clock++

    }
    return clock, terminatedQueue

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

    
    for _, p := range terminatedQueue {
        turnaround_time += float64(p.CompletionTime.Int32 - p.ArrivalTime)
        waiting_time += float64(p.WaitingTime.Int32)
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
    })
    if err != nil {
        log.Println("Error al actualizar la simulación:", err)
        return
    }
    log.Printf("Simulación %d actualizada con estadísticas.", ultima.ID)
    return
}
