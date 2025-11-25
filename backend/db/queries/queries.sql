--• Una consulta para crear un nuevo registro (Create...).
--• Una consulta para obtener un registro por su ID (Get...).
--• Una consulta para listar todos los registros (List...).
---• Una consulta para actualizar un registro (Update...).
--• Una consulta para borrar un registro (Delete...).
--Utiliza las anotaciones que requiere sqlc (-- name: ...).




-- name: CreateProcess :one
INSERT INTO procesos (nombre, prioridad, burst_time, arrival_time, estado, id_simulacion)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetProcess :one
SELECT *
FROM procesos
WHERE id = $1;

-- name: ListProcess :many
SELECT *
FROM procesos
ORDER BY id;

-- name: ListProcessByEstado :many
SELECT *
FROM procesos
WHERE estado = $1
ORDER BY id;

-- name: UpdateProcess :exec
UPDATE procesos
SET nombre=$2,
    prioridad=$3,
    burst_time=$4,
    arrival_time=$5, 
    estado=$6,
    id_simulacion=$7
WHERE id = $1;

-- name: DeleteProcess :exec
DELETE FROM procesos
WHERE id= $1;


-- name: CreateSimulacion :one
INSERT INTO simulaciones (nombre, process_time, context_switches, dispatch_latency, average_turnaround_time, average_waiting_time, average_throughput, algoritmo, quantum, prioridad)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;
-- name: GetSimulacion :one
SELECT *
FROM simulaciones
WHERE id = $1;

-- name: ListSimulacion :many
SELECT *
FROM simulaciones
ORDER BY id;

-- name: UpdateSimulacion :exec
UPDATE simulaciones
SET nombre=$2,
    process_time=$3,
    context_switches=$4,
    dispatch_latency=$5,
    average_turnaround_time=$6,
    average_waiting_time=$7,
    average_throughput=$8,
    algoritmo=$9,
    quantum=$10,
    prioridad=$11
WHERE id = $1;

-- name: GetLastSimulacion :one
SELECT *
FROM simulaciones
ORDER BY id DESC
LIMIT 1;