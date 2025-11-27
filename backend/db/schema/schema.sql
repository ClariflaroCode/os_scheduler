CREATE TABLE proceso (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nombre VARCHAR(40) NOT NULL,
    prioridad INT NOT NULL,
    burst_time INT NOT NULL, 
    arrival_time INT NOT NULL,
    waiting_time INT DEFAULT 0,
    completion_time INT DEFAULT 0,
    estado VARCHAR(10) NOT NULL, 
    id_simulacion INT NOT NULL

);
CREATE TABLE simulacion (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nombre VARCHAR(40) NOT NULL,
    process_time INT NOT NULL,
    context_switches INT NOT NULL,
    dispatch_latency INT NOT NULL,
    average_turnaround_time DOUBLE PRECISION NOT NULL,
    average_waiting_time DOUBLE PRECISION NOT NULL,
    average_throughput DOUBLE PRECISION NOT NULL,
    algoritmo VARCHAR(40) NOT NULL,
    quantum INT NULL,
    prioridad INT NULL
);


