POST http://localhost:8080/api/procesos
Content-Type: application/json
{
  "nombre": "Proceso test2",
  "prioridad": 1,
  "burst_time": 5,
  "arrival_time": 0,
  "estado": "listo",
  "id_simulacion": 1
}
HTTP 201

POST http://localhost:8080/api/procesos
Content-Type: application/json
{
  "nombre": "Proceso test3",
  "prioridad": 1,
  "burst_time": 5,
  "arrival_time": 0,
  "estado": "listo",
  "id_simulacion": 1
}
HTTP 201

#obtener los procesos
GET http://localhost:8080/api/procesos
HTTP 200


#Edita el proceso con el id=1
PUT http://localhost:8080/api/procesos/2
Content-Type: application/json 
{
  "nombre": "Proceso especial",
  "prioridad": 1,
  "burst_time": 5,
  "arrival_time": 0,
  "estado": "listo",
  "id_simulacion": 1
}
HTTP 200

#obtener el proceso con el id = 1
GET http://localhost:8080/api/procesos/2
HTTP 200


#Elimina el proceso con el id=1
DELETE http://localhost:8080/api/procesos/2
HTTP 204

#obtener todos los procesos de la tabla procesos
GET http://localhost:8080/api/procesos
HTTP 200
