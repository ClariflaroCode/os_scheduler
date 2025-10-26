CREATE TABLE procesos (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nombre VARCHAR(40) NOT NULL,
    prioridad INT NOT NULL,
    burst_time INT NOT NULL, 
    arrival_time INT NOT NULL,
    estado VARCHAR(10) NOT NULL, 
    id_simulacion INT NOT NULL

)


