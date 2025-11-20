# os_scheduler
Este proyecto se realizó para la cátedra de programación web de la carrera de ingeniería de sistemas utilizando el lenguaje Go para el backend, PostgreSQL para la base de datos, Javascript, CSS, HTML, HTMX, JSON y TEMPL para el frontend y las tecnologías de Make, Air, sqlc y HURL para automatización de procesos. 

**ESTRUCTURA ACTUAL DEL PROYECTO:**  
myapp  
│  
├──backend\
│     ├──db\
│     │   ├──queries\
│     │   │   └──queries.sql\
│     │   └──schema  
│     │          └──schema.sql   
│     ├──sqlc\
│     │    ├─db.go   
│     │    ├─models.go\
│     │    └─queries.sql.go\
│     ├──sqlc.yaml\
|     ├──views\ 
|     |     ├─ layout.templ\ 
|     |     ├─ grafo.templ\ 
|     |     ├─ entity_list.templ\ 
|     |     ├─ entity_form.templ\ 
|     |     └─ estadisticas.templ\ 
|     └──css\
│          └──pico.css\
│    
├──main.go\
├──Makefile\
├──request.hurl\
├──go.sum\
├──go.mod\
├──docker-compose.yml\
└── README.md

**PARA EJECUTAR EL PROYECTO**  
*Para correr el backend*
- Correr el comando make all


*Para correr los testeos del hurl*
- Correr el comando make test 

Alumna: Julieta Watts   
Materia: Programación web  
Carrera: Ingenieria de sistemas   
Año: 2025   
