CREATE TABLE proceso (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nombre VARCHAR(40) NOT NULL,
    prioridad INT NOT NULL,
    burst_time INT NOT NULL, 
    arrival_time INT NOT NULL,
    estado VARCHAR(10) NOT NULL, 
    id_simulacion INT NOT NULL

);
CREATE TABLE simulacion (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nombre VARCHAR(40) NOT NULL,
    process_time INT NOT NULL,
    context_switches INT NOT NULL,
    dispatch_latency INT NOT NULL,
    average_turnaround_time INT NOT NULL,
    average_waiting_time INT NOT NULL,
    average_throughput INT NOT NULL,
    algoritmo VARCHAR(40) NOT NULL,
    quantum INT NULL,
    prioridad INT NULL
)

