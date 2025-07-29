# Ticketera Web (Go + MongoDB)

Aplicación web para generar tickets personalizados, exportarlos a DOCX, administrar textos y logo, con arquitectura hexagonal y MongoDB.

## Requisitos
- Go 1.22+
- MongoDB en ejecución (por defecto en `mongodb://localhost:27017`)

## Instalación
1. Clona el repositorio y entra a la carpeta del proyecto.
2. Instala las dependencias:
   ```sh
   go mod tidy
   ```
3. Asegúrate de tener un archivo de plantilla `ticket_template.docx` en la carpeta `web/`.

## Compilación
```sh
cd cmd/ticketera
go build -o ticketera
```

## Ejecución
Desde la raíz del proyecto:
```sh
./cmd/ticketera/ticketera
```
El servidor se inicia en `http://localhost:8080`

## Uso
- Accede a `http://localhost:8080/` para generar tickets.
- Accede a `http://localhost:8080/admin` para administrar textos y logo.

### Ejemplo para probar
1. Sube textos y logo desde `/admin`.
2. Ve a `/`, completa el formulario y genera un ticket.
3. El archivo `ticket.docx` se descargará con los datos y logo personalizados.

## Personalización de plantilla DOCX
- Edita `web/ticket_template.docx` en Word.
- Usa campos como `{{TITLE}}`, `{{CONTENT}}`, `{{LOGO}}` y cualquier otro texto entre llaves dobles para que sean reemplazados dinámicamente.

### Ejemplo de plantilla
Puedes crear un archivo en Word con el siguiente contenido:

```
Ticket de Atención

Título: {{TITLE}}
Contenido: {{CONTENT}}
Texto personalizado: {{MI_TEXTO}}

Logo: {{LOGO}}
```

Luego, desde `/admin`, agrega el texto con clave `MI_TEXTO` y su valor.

## Docker

### Build
```sh
docker build -t ticketera .
```

### Run
```sh
docker run --rm -p 8080:8080 -v $(pwd)/web:/app/web ticketera
```

> Asegúrate de tener MongoDB accesible desde el contenedor (puedes usar Docker Compose o exponer el puerto 27017 de tu MongoDB local).

---
Desarrollado con Go, MongoDB y Bootstrap.
# ticketera
