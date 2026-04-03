
package simulador

import (
	db "scheduler_os/backend/db/sqlc"
    "database/sql"
    "context"
)


func Simulacion(newQueue []db.Proceso, elegir func([]*db.Proceso, *db.Proceso) (*db.Proceso, []*db.Proceso), queries *db.Queries) (int, []db.Proceso) {
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

		// acá se ejecuta la funcion elegir que me hayan pasado por parámetro, dependiendo el algoritmo va a cambiar.
		runningProceso, readyQueue = elegir(readyQueue, runningProceso) 


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

func FirstComeFirstServe(ready []*db.Proceso, running *db.Proceso) (*db.Proceso, []*db.Proceso){ //NOTA: los primeros parentesis son los parametros de entrada, los segundos los de salida

        //Tomo el primer elemento de la lista de ready y lo paso a running si está disponible para ejecutar
		if running != nil {
			return running, ready
		}
		if len(ready) > 0 {
			running = ready[0]
			//contextSwitches++
			running.Estado = "running"
			//updateProcesoEstado(readyQueue[0].ID, "running") //actualizo en la DB el estado del proceso
			ready = ready[1:] //elimino el primer elemento de la lista de ready, la ready queue se volvio la ready queue 
			// empezando desde el segundo elemento hasta el final por eso 1:. En go la longitud se indica con min:max
		}
		return running, ready
}


func buscaMenor(ready []*db.Proceso) int {
	var posMenor int = 0
	for i:= 0; i < len(ready); i++ {
		if (ready[i].BurstTime < ready[posMenor].BurstTime) {
			posMenor = i
		}
	}
	return posMenor
}
func ShortestJobFirst(ready []*db.Proceso, running *db.Proceso) (*db.Proceso, []*db.Proceso) {
    //sin desalojo
	if running != nil  || len(ready) == 0 {
		return running, ready
	}

    /*Como no todos los procesos tienen arrival 1 lo mejor es buscar linealmente el valor menor en el arreglo*/
    posMin := buscaMenor(ready)
    running = ready[posMin]
    running.Estado = "running"
    ready = append(ready[:posMin], ready[posMin+1:]...) //elimina el elemento minimo de ready. 
    return running, ready
	
}
func ShortestJobFirstPreemptive(ready []*db.Proceso, running *db.Proceso) (*db.Proceso, []*db.Proceso) {
    if len(ready) > 0{ //como es con desalojo, no importa si hay o no un proceso en ejecución en el scheduler, se desaloja si llega alguno con menos burst time. 
        posMin := buscaMenor(ready)
        if running != nil {
            if running.BurstTime > ready[posMin].BurstTime {
                running.Estado = "ready"
                ready = append(ready, running) //el proceso que estaba en ejecución vuelve a la fila de ready
            } else {
                return running, ready
            }
        }
           
        running = ready[posMin]
        running.Estado = "running"
        ready = append(ready[:posMin], ready[posMin+1:]...) //elimina el elemento minimo de ready.
        
    }
    return running, ready
}
func RoundRobin(ready []*db.Proceso, running *db.Proceso, quantum int) (*db.Proceso, []*db.Proceso) {
    //TO-DO
    return running, ready
}
func buscarProcesoDeMayorPrioridad(ready []*db.Proceso) int {
	var posMayor int = 0
	for i:= 0; i < len(ready); i++ {
		if (ready[i].Prioridad < ready[posMayor].Prioridad) {
			posMayor = i
		}
	}
	return posMayor
}
func Priority(ready []*db.Proceso, running *db.Proceso) (*db.Proceso, []*db.Proceso) {
    //Recordar que cuánto más pequeño sea el numero más prioritario es. 
    //TO-DO: IMPLEMENTAR EL AGING DE LOS PROCESOS CON PRIORIDAD, ES DECIR QUE SE INCREMENTA LA PRIORIDAD DE LOS PROCESOS PARA DARLE LUGAR A OTROS y no hambrear xD
    if running != nil || len(ready) == 0{
        return running, ready
    }

   posMayor := buscarProcesoDeMayorPrioridad(ready)
   running = ready[posMayor]
   running.Estado = "running"
   ready = append(ready[:posMayor], ready[posMayor+1:]...) //elimina el elemento de mayor prioridad de ready.
   return running, ready
}
