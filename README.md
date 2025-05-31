# Stockify 📈

**Stockify** es una aplicación web diseñada para analizar y mostrar información sobre acciones del mercado de valores. Proporciona una interfaz intuitiva para navegar, buscar y ordenar datos de acciones, y ofrece recomendaciones basadas en el análisis de los datos almacenados para ayudar a los usuarios a tomar decisiones de inversión informadas.



# Contenido

- [✨ Características Principales](#-características-principales)

- [🛠️ Tech Stack](#️-tech-stack)

- [📋 Prerrequisitos](#-prerrequisitos)

- [🚀 Configuración e Instalación](#-configuración-e-instalación)

  - [Usando Docker Compose (Método Recomendado)](#usando-docker-compose-método-recomendado)
    - [1. Clonar el Repositorio](#1-clonar-el-repositorio)
    - [2. Configurar Variables de Entorno](#2-configurar-variables-de-entorno)
    - [3. (Opcional) Scripts de Inicialización de CockroachDB](#3-opcional-scripts-de-inicialización-de-cockroachdb)
    - [4. Levantar los Servicios con Docker Compose](#4-levantar-los-servicios-con-docker-compose)
    - [5. Poblar la Base de Datos](#5-poblar-la-base-de-datos)
    - [6. Acceder a la Aplicación](#6-acceder-a-la-aplicación)

  - [Configuración Manual (Sin Docker Compose)](#configuración-manual-sin-docker-compose)
    - [1. Base de Datos (CockroachDB)](#1-base-de-datos-cockroachdb)
    - [2. Backend (Go)](#2-backend-go)
    - [3. Frontend (Vue)](#3-frontend-vue)

- [📖 Documentación de la API (Backend)](#-documentación-de-la-api-backend)
  - [1. Listar Acciones](#1-listar-acciones)
  - [2. Obtener Detalles de un Evento de Stock por Ticker](#2-obtener-detalles-de-un-evento-de-stock-por-ticker)
  - [3. Obtener Recomendaciones de Acciones](#3-obtener-recomendaciones-de-acciones)

- [🚀 Uso de la Aplicación](#-uso-de-la-aplicación)
- [🧪 Ejecutar Pruebas](#-ejecutar-pruebas)
  - [Backend (Go)](#backend-go)
  - [Frontend (Vue)](#frontend-vue)



## ✨ Características Principales

* **Visualización de datos de stocks**: Muestra una lista de eventos de rating de acciones con información relevante como ticker, compañía, brokerage, ratings y precios objetivo.

* **Búsqueda y ordenamiento**: Permite a los usuarios buscar acciones por ticker o compañía y ordenar la lista por diversos criterios (fecha, ticker, compañía, rating, target).

* **Vista de detalle del stock**: Ofrece una vista detallada para cada evento de rating de una acción.

* **Recomendaciones inteligentes**: Implementa un algoritmo que analiza los datos almacenados para sugerir las acciones más prometedoras para invertir, presentando razones estructuradas para cada recomendación.



## 🛠️ Tech Stack

| Capa            | Tecnologías                                                  |
| --------------- | ------------------------------------------------------------ |
| Backend         | **Go (Golang)** con Chi (enrutador) y GORM (ORM)             |
| Frontend        | **Vue 3** + TypeScript + Pinia (gestión de estado) + Tailwind CSS + Vite |
| Base de Datos   | **CockroachDB**                                              |
| Contenerización | **Docker** & **Docker Compose**                              |



## 📋 Prerrequisitos

Para ejecutar este proyecto localmente, necesitarás tener instalados:

* Git

* Docker Desktop (o Docker Engine + Docker Compose CLI)

* Node.js (v18+ recomendado, para el setup manual del frontend)

* Go (v1.21+ recomendado, para el setup manual del backend)



## 🚀 Configuración e Instalación

### Usando Docker Compose (Método Recomendado)

Esta es la forma más sencilla de levantar toda la aplicación (backend, frontend, base de datos).

1. **Clonar el repositorio:**

   ```bash
   git clone <URL_DEL_REPOSITORIO>
   cd nombre-del-repositorio
   ```

2. **Configurar variables de entorno:**
   Necesitarás crear algunos archivos `.env` con las configuraciones necesarias.

   * **Archivo raíz `.env`** (en la misma carpeta que `compose.yml`):
     Este archivo es leído por `docker-compose` para variables que se usan durante el build o para configurar los servicios.

     ```env
     # ./ (raíz del proyecto)/.env
     VITE_APP_API_BASE_URL=http://localhost:3030/api 
     # Esta URL es la que usará el navegador del cliente para llamar al backend.
     # El backend estará expuesto en el puerto 3030 del host.
     ```

   * **Archivo de entorno del backend (`backend/.env`):**
     Crea un archivo llamado `.env` dentro de la carpeta `backend/`.

     ```env
     # backend/.env
     DATABASE_URL="postgresql://root@db:26257/defaultdb?sslmode=disable" 
     # 'db' es el nombre del servicio de la base de datos en docker-compose.yml

     STOCK_API_TOKEN="TU_TOKEN_BEARER_DE_LA_API_EXTERNA_AQUI" 
     # Reemplaza esto con el token real proporcionado en el assesment.md

     SERVER_PORT="8080" 
     # El servidor Go dentro del contenedor escuchará en el puerto 80.
     # Docker Compose mapeará el puerto 3030 del host a este puerto 80 del contenedor.
     ```

3. **Levantar los servicios con Docker Compose:**
   Desde la raíz de tu proyecto (donde está `compose.yml`):

   ```bash
   docker compose up -d --build
   ```

   * `--build`: Reconstruye las imágenes si ha habido cambios en los Dockerfiles o el código fuente.
   * `-d`: Ejecuta los contenedores en segundo plano (detached mode).

4. **Poblar la base de datos (primera vez o periódicamente):**
   El backend incluye un script CLI para poblar la base de datos desde la API externa. El `entrypoint.sh` del backend Dockerfile ya ejecuta este script automáticamente la primera vez que el contenedor del backend se inicia y la base de datos está vacía.

   Si necesitas ejecutarlo manualmente después (por ejemplo, para actualizar datos, aunque la lógica actual solo puebla si está vacía):

   ```bash
   docker compose exec backend /app/stockify_datasync
   ```

5. **Acceder a la aplicación:**

   * **Frontend (aplicación Vue)**: Abre tu navegador y ve a `http://localhost:3000`
   * **Backend API (Go)**: Los endpoints de la API estarán disponibles en `http://localhost:3030/api` (ej. `http://localhost:3030/api/stocks`)
   * **CockroachDB admin UI**: Puedes acceder a la interfaz de administración de CockroachDB en `http://localhost:8081`

### Configuración Manual (Sin Docker Compose)

Si prefieres no usar Docker Compose, puedes ejecutar cada componente por separado.

1. **Base de datos (CockroachDB):**

   * Puedes iniciar una instancia de CockroachDB usando Docker manualmente:

     ```bash
     docker run -d --name cockroachdb-manual -p 26257:26257 -p 8081:8080 cockroachdb/cockroach:latest start-single-node --insecure
     ```

   * Asegúrate que tu `DATABASE_URL` en `backend/.env` apunte a `postgresql://root@localhost:26257/stockify?sslmode=disable` (o la BD que crees).

2. **Backend (Go):**

   * Navega al directorio `backend/`: `cd backend`

   * Configura las variables de entorno (puedes ponerlas en `backend/.env` y asegurarte que tu `config.Load()` las lea, o exportarlas en tu terminal):

     * `DATABASE_URL` (apuntando a tu instancia de CockroachDB, ej. `postgresql://root@localhost:26257/defaultdb?sslmode=disable`)
     * `STOCK_API_TOKEN` (tu token Bearer de la API externa)
     * `SERVER_PORT` (ej. `8080`)

   * Instala dependencias: `go mod tidy`

   * **Poblar la base de datos (si está vacía):**

     ```bash
     go run ./cmd/datasync/main.go
     ```

   * **Iniciar el servidor API:**

     ```bash
     go run ./cmd/server/main.go
     ```

     El backend estará escuchando en el puerto que definiste (ej. `http://localhost:8080`).

3. **Frontend (Vue):**

   * Navega al directorio `frontend/`: `cd frontend`

   * Configura la variable de entorno (puedes ponerla en `frontend/.env`):

     * `VITE_APP_API_BASE_URL` (apuntando a tu servidor backend Go, ej. `http://localhost:8080/api`)

   * Instala dependencias: `npm install` (o `yarn install` / `pnpm install`)

   * **Iniciar el servidor de desarrollo:**

     ```bash
     npm run dev
     ```

     El frontend estará disponible usualmente en `http://localhost:5173`.



## 📖 Documentación de la API (Backend)

La API sigue un estilo RESTful y todas las rutas están prefijadas con `/api`.

### 1. Listar stocks

* **Endpoint:** `GET /api/stocks`

* **Descripción:** Obtiene una lista paginada de los stocks, con opciones de búsqueda y ordenamiento.

* **Query parameters:**

  * `search` (opcional, string): Término para buscar por ticker o nombre de compañía.
  * `sortBy` (opcional, string): Campo por el cual ordenar (ej. `ticker`, `company`, `time`, `rating_to`, `target_to`). Por defecto `time`.
  * `sortOrder` (opcional, string): Orden (`asc` o `desc`). Por defecto `desc` para `time`.
  * `page` (opcional, int): Número de página (por defecto `1`).
  * `pageSize` (opcional, int): Número de ítems por página (por defecto `10`).

* **Respuesta exitosa (200 OK):**

  ```json
  {
  	"page": 1,
  	"pageSize": 10,
  	"stocks": [
  		{
  			"ID": 1076576506564444161,
  			"CreatedAt": "2025-05-30T09:31:16.387798-05:00",
  			"UpdatedAt": "2025-05-30T09:31:16.387798-05:00",
  			"DeletedAt": null,
  			"ticker": "S",
  			"company": "SentinelOne",
  			"brokerage": "DA Davidson",
  			"action": "target lowered by",
  			"rating_to": "Neutral",
  			"rating_from": "Neutral",
  			"target_to": 17,
  			"target_from": 18,
  			"time": "2025-05-29T19:30:05.861791-05:00"
  		},
  		{
  			"ID": 1076577054025449473,
  			"CreatedAt": "2025-05-30T09:34:03.468581-05:00",
  			"UpdatedAt": "2025-05-30T09:34:03.468581-05:00",
  			"DeletedAt": null,
  			"ticker": "NXRT",
  			"company": "NexPoint Residential Trust",
  			"brokerage": "Truist Financial",
  			"action": "target lowered by",
  			"rating_to": "Hold",
  			"rating_from": "Hold",
  			"target_to": 38,
  			"target_from": 42,
  			"time": "2025-05-29T19:30:05.848666-05:00"
  		},
      	// ... más acciones
      ],
      "totalItems": 100,
      "totalPages": 10
  }
  ```

### 2. Obtener detalles de un stock por ticker

* **Endpoint:** `GET /api/stocks/{ticker}`

* **Descripción:** Obtiene el detalle del stock para un ticker específico.

* **Path parameters:**

  * `ticker` (string, requerido): El símbolo del stock (ej. `AAPL`).

* **Respuesta exitosa (200 OK):**

  ```json
  {
  	"ID": 1075275864031002625,
  	"CreatedAt": "2025-05-25T19:15:51.776458-05:00",
  	"UpdatedAt": "2025-05-25T19:15:51.776458-05:00",
  	"DeletedAt": null,
  	"ticker": "A",
  	"company": "Agilent Technologies",
  	"brokerage": "Jefferies Financial Group",
  	"action": "target lowered by",
  	"rating_to": "Hold",
  	"rating_from": "Hold",
  	"target_to": 116,
  	"target_from": 135,
  	"time": "2025-04-21T19:30:06.089698-05:00"
  }
  ```

* **Respuesta de error (404 Not Found):**

  ```json
  {
    "error": "Stock no encontrado"
  }
  ```

### 3. Recomendaciones de stocks

* **Endpoint:** `GET /api/stocks/recommendations`

* **Descripción:** Obtiene una lista de las acciones recomendadas para invertir, basadas en el algoritmo implementado.

* **Respuesta exitosa (200 OK):**

  ```json
  {
  	"recommendations": [
  		{
  			"ID": 1075276044340330497,
  			"CreatedAt": "2025-05-25T19:16:46.802687-05:00",
  			"UpdatedAt": "2025-05-25T19:16:46.802687-05:00",
  			"DeletedAt": null,
  			"ticker": "FICO",
  			"company": "Fair Isaac",
  			"brokerage": "Needham & Company LLC",
  			"action": "target raised by",
  			"rating_to": "Buy",
  			"rating_from": "Buy",
  			"target_to": 2575,
  			"target_from": 2500,
  			"time": "2025-05-12T19:30:08.130819-05:00",
  			"reasons": [
  				{
  					"type": "POSITIVE_RATING",
  					"details": "Rating positivo: Buy"
  				},
  				{
  					"type": "TARGET_INCREASED",
  					"details": "Precio objetivo aumentado de $2500.00 a $2575.00"
  				},
  				{
  					"type": "RECENT_EVENT",
  					"details": "Evento de rating reciente (últimos 3 meses)."
  				}
  			],
  			"score": 342.5
  		},
      // ... más recomendaciones ...
    ]
  }
  ```



## 🚀 Uso de la Aplicación

1. Asegúrate de que todos los servicios (Docker Compose o manuales) estén corriendo.

2. Abre el frontend en tu navegador (ej. `http://localhost:3000`).

3. Navega por la lista de acciones.

4. Usa la barra de búsqueda para filtrar por ticker o compañía.

5. Usa los selectores para ordenar los resultados.

6. Haz clic en el ticker de una acción o en "Ver" para ir a la página de detalles.

7. Revisa la sección "Top Recomendaciones de Hoy" para ver las sugerencias de inversión.



## 🧪 Ejecutar Pruebas

### Backend (Go)

Navega al directorio `backend/` y ejecuta:

```bash
go test ./... -v
```

### Frontend (Vue)

Navega al directorio `frontend/` y ejecuta:

```bash
npm test
```

<br/>

---
<br/>
<br/>

· Hecho con 💛 :) por [Nath](https://github.com/nathsotomayor) ·

