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

-- name: UpdateProcess :exec
UPDATE procesos
SET nombre=$2,
    prioridad=$3,
    burst_time=$4,
    arrival_time=$5, 
    estado=$6
WHERE id = $1;

-- name: DeleteProcess :exec
DELETE FROM procesos
WHERE id= $1;