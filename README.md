# os_scheduler
Este proyecto se realizó para la cátedra de programación web de la carrera de ingeniería de sistemas utilizando el lenguaje Go para el backend, PostgreSQL para la base de datos, Javascript, CSS, HTML, HTMX, JSON y TEMPL para el frontend y las tecnologías de Make, Air, sqlc y HURL para automatización de procesos. 

**ESTRUCTURA ACTUAL DEL PROYECTO:**
myapp/\
│
├──index.html\
├──backend\
│     ├──main.go\
│     ├──Makefile\
│     ├──request.hurl\
│     ├──sqlc.yaml\
│     ├──go.sum\
│     ├──go.mod\
│     ├──docker-compose.yml\
│     ├──air.toml\
│     ├──descripcion.txt\
│     ├──db\
│     │   ├──queries\
│     │   │   └──queries.sql\
│     │   └──schema
│     │          └──schema.sql\
│     └─sqlc\
│         ├─db.go\   
│         ├─models.go\
│         └─queries.sql.go\
└──frontend\
    ├──README.md\
    ├──css\
    │   └──styles.css\
    ├──js\
    │   └──main.js\
    └── pages\
        ├──algoritmo.html\ 
        └──estadisticas.html   <!--Esta pagina podría ser un template de templ creo-->\

**PARA EJECUTAR EL PROYECTO**
*Para correr el backend*
- cd backend
- Correr el comando make all

*Para correr el frontend* 
- Utilizar live server de vsc

*Para correr los testeos del hurl*
- cd backend
- Correr el comando make test 

Alumna: Julieta Watts\ 
Materia: Programación web\ 
Carrera: Ingenieria de sistemas\ 
Año: 2025\ 