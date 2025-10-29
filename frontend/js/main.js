document.addEventListener("DOMContentLoaded", iniciar);

function iniciar () {
    const base_url = "http://localhost:8080/api/procesos"
    listarProcesos();

    const btnForm = document.getElementById("btn-editar-o-agregar");
    const modal = document.getElementById("modal");

    document.getElementById("btn-add").addEventListener("click", () => {
        modal.showModal(); //muestra el modal con el formulario. ¿por qué showModal y no show(): https://webinista.com/demos/dialog-element-tutorial/index.html (para bloquear el resto de la pagina)
        btnForm.innerHTML = "Agregar recurso";

    }); 
    document.getElementById("btn-close-form-popUp").addEventListener("click", () => {
        modal.close();   
    })

    const form = document.getElementById("form-editar-agregar");
    form.addEventListener("submit", (e)  => {
            e.preventDefault(); //previene el comportamiento por defecto (?)
            if(btnForm.textContent === "Agregar recurso") {
                crearProceso()
                modal.close();
            }
            else {
                editarProceso(form.dataset.editingItem);
                modal.close();
            }
        }
    );
    async function obtenerProceso(id) {
        try {
            const response = await fetch(base_url + '/' + id, {
                method: 'GET',
                headers: {'content-type':'application/json'}})
            ;
            if (response.ok) {
                const proceso = await response.json();      
                
            }            
        } catch (error) {
            console.log(error);
        }
    }
    async function listarProcesos() {
        try {
            const response = await fetch(base_url, {
                method: 'GET',
                headers: {'content-type':'application/json'}})
            ;
            if (response.ok) {
                const procesos = await response.json();      
                mostrarProcesos(procesos);
            }            
        } catch (error) {
            console.log(error);
        }

    }
    async function editarProceso(id){
        const formData = new FormData(form);
        
        const nombre = formData.get('nombre');
        const prioridad = parseInt(formData.get('prioridad'));
        const burstTime = parseInt(formData.get('burst_time'));
        const arrivalTime = parseInt(formData.get('arrival_time')) ;

        const nuevoProceso = {
            "nombre": nombre,
            "prioridad": prioridad,
            "burst_time": burstTime,
            "arrival_time": arrivalTime,
            "estado": "new"
        }


        try {
            const response = await fetch(base_url +  '/' + id, {
                method: 'PUT',
                headers: {'content-type':'application/json'},
                body: JSON.stringify(nuevoProceso)
                
            })
            ;
            if (response.ok) {
                const proceso = await response.json();   
                listarProcesos();

                
            }            
        } catch (error) {
            console.log(error);
        }
    }
    async function eliminarProceso(id){
        try {
            const response = await fetch(base_url + '/' + id, {method: 'DELETE'} );
            if (response.ok) {
                listarProcesos();
                
            }            
        } catch (error) {
            console.log(error);
        }
    }
    async function crearProceso(){
        const formData = new FormData(form);

        const nombre = formData.get('nombre');
        const prioridad = parseInt(formData.get('prioridad'));
        const burstTime = parseInt(formData.get('burst_time'));
        const arrivalTime = parseInt(formData.get('arrival_time')) ;
        const estado = "new";

        console.log(arrivalTime);
        const nuevoProceso = {
            "nombre": nombre,
            "prioridad": prioridad,
            "burst_time": burstTime,
            "arrival_time": arrivalTime,
            "estado": estado
        }


        try {
            const response = await fetch(base_url, {
                method: 'POST',
                headers: {'content-type':'application/json'},
                body: JSON.stringify(nuevoProceso)
            })
            ;
            if (response.ok) {
                const proceso = await response.json();      
                console.log(proceso);
                listarProcesos();
                
            }            
        } catch (error) {
            console.log(error);
        }
    }

    function mostrarProceso(proceso){

    }
    function mostrarProcesos(procesos) {
        const container = document.getElementById('listado');
        container.innerHTML = ''; //debo limpiar el html interno, porque si no se va a duplicar. 
        
        //ojo que puede devolver lista vacia y llegar acá y no se debe romper. 
        if (procesos.length == 0) {
            console.log("Esta vacio");
            container.appendChild(document.createElement('h3').textContent = "No hay procesos");

        } else {
             procesos.forEach(proceso => {
                const card = document.createElement('article');

                const title = document.createElement('h2');
                title.textContent = "Nombre " + proceso.nombre;

                const list = document.createElement('ul');

                const prioridad = document.createElement('li');
                prioridad.textContent = "Prioridad " + proceso.prioridad;

                const burst_time = document.createElement('li');
                burst_time.textContent = "Burst time " + proceso.burst_time;
                            
                const arrival_time = document.createElement('li');
                arrival_time.textContent = "Arrival Time " + proceso.arrival_time;

                const estado = document.createElement('li');
                estado.textContent = "Estado " + proceso.estado;

                list.appendChild(prioridad);
                list.appendChild(burst_time);
                list.appendChild(arrival_time);
                list.appendChild(estado);

                card.appendChild(title);
                card.appendChild(list);

                const btnDelete = document.createElement('button');
                btnDelete.innerHTML = "Eliminar";
                btnDelete.dataset.id = proceso.id; 
                btnDelete.addEventListener("click", (event) => {
                        eliminarProceso(proceso.id); 
                    }
                );

                let btnModify = document.createElement("button");
                btnModify.innerHTML = "Editar";
                btnModify.dataset.id = proceso.id; 
                btnModify.addEventListener("click", () => {
                    modal.showModal(); 


                    document.getElementById("nombre").value = proceso.nombre;
                    document.getElementById("prioridad").value = proceso.prioridad;
                    document.getElementById("arrival_time").value = proceso.arrival_time;

                    document.getElementById("burst_time").value = proceso.burst_time;

                    btnForm.innerHTML = "Actualizar cambios"; 
                });
                card.appendChild(btnDelete);
                card.appendChild(btnModify);

                form.dataset.editingItem = proceso.id; //le agrego al formulario en el dataset el id del elemento que quiero modificar. 
                
                container.appendChild(card);

                card.classList.add('card');
                list.classList.add('card-footer');

            });
            
        }
    }
}