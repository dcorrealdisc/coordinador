# Agente Svelte - Guía de Trabajo

## Rol y Responsabilidades

El Agente Svelte es responsable del desarrollo frontend del sistema Coordinador. Trabaja con dos aplicaciones SvelteKit separadas: el dashboard administrativo (`admin-web`) y el portal de usuarios (`portal-web`).

### Responsabilidades Principales

1. **Desarrollo de interfaces de usuario**
   - Componentes Svelte reutilizables
   - Páginas con layouts consistentes
   - Formularios con validación client-side
   - Tablas con filtros, búsqueda y paginación

2. **Integración con el API**
   - Consumir endpoints REST del backend Go/Fiber
   - Manejo de estados de carga, error y éxito
   - Tipado TypeScript de responses y requests

3. **UX/UI**
   - Diseño responsive con Tailwind CSS
   - Feedback visual al usuario (toasts, loading states)
   - Navegación intuitiva

4. **Testing**
   - Tests de componentes
   - Tests de integración con API

## Contexto del Proyecto

### Stack Frontend

- **SvelteKit** (framework)
- **Svelte 4** (UI library)
- **TypeScript** (tipado)
- **Tailwind CSS 3** (estilos)
- **Axios** (HTTP client)
- **Vite 5** (bundler)

### Dos Aplicaciones Separadas

#### `admin-web` (Dashboard Administrativo)
- **Puerto**: 5173 (default Vite)
- **Usuarios**: Coordinadores del programa de maestría
- **Funciones**: CRUD de estudiantes, cursos, planificación, reportes, gestión de tutores
- **Alcance**: Gestión completa del sistema

#### `portal-web` (Portal de Usuarios)
- **Puerto**: 3001
- **Usuarios**: Estudiantes, profesores, tutores
- **Funciones**: Ver información personal, cursos inscritos, calificaciones, expresar interés en tutorías
- **Alcance**: Vista limitada según rol

### Estado Actual de las Apps

Ambas apps tienen `package.json` y `node_modules` pero **no tienen código fuente aún**. Necesitan el scaffolding inicial de SvelteKit.

### Prioridad de Implementación

**Empezar con `admin-web`** porque:
- Es donde se gestiona todo el sistema
- Los endpoints del API ya están listos para estudiantes
- Permite validar visualmente el backend end-to-end

## API Backend Disponible

### Base URL
```
http://localhost:8080/api/v1
```

### Health Check
```
GET http://localhost:8080/health
Response: { "status": "ok", "service": "coordinador-api", "version": "0.1.0", "database": "ok" }
```

### Módulo de Estudiantes (implementado)

#### Crear Estudiante
```
POST /api/v1/students
Content-Type: application/json

Request Body:
{
  "full_name": "string (required, min 3, max 255)",
  "document_id": "string (optional, max 50)",
  "birth_date": "string (required, format YYYY-MM-DD)",
  "profile_photo_url": "string (optional, URL)",
  "city_origin_id": "string (optional, UUID)",
  "country_origin_id": "string (required, UUID)",
  "emails": ["string (required, min 1, each must be email)"],
  "phones": ["string (optional)"],
  "company_id": "string (optional, UUID)",
  "status": "string (required, one of: active, graduated, withdrawn, suspended)",
  "cohort": "string (required, max 10, ej: '2024-1')",
  "enrollment_date": "string (required, format YYYY-MM-DD)"
}

Response (201):
{
  "success": true,
  "message": "Student created successfully",
  "data": { Student object }
}
```

#### Listar Estudiantes
```
GET /api/v1/students?status=active&cohort=2024-1&search=Juan&country_id=UUID&limit=20&offset=0

Response (200):
{
  "success": true,
  "message": "Students retrieved successfully",
  "data": {
    "items": [ Student objects ],
    "total": 42,
    "limit": 20,
    "offset": 0
  }
}
```

#### Obtener Estudiante por ID
```
GET /api/v1/students/:id

Response (200):
{
  "success": true,
  "message": "Student retrieved successfully",
  "data": { Student object }
}
```

#### Actualizar Estudiante
```
PUT /api/v1/students/:id
Content-Type: application/json

Request Body (todos opcionales):
{
  "full_name": "string",
  "document_id": "string",
  "profile_photo_url": "string",
  "emails": ["string"],
  "phones": ["string"],
  "company_id": "string (UUID)",
  "status": "string"
}

Response (200):
{
  "success": true,
  "message": "Student updated successfully",
  "data": { Student object }
}
```

#### Eliminar Estudiante (Soft Delete)
```
DELETE /api/v1/students/:id

Response (200):
{
  "success": true,
  "message": "Student deleted successfully"
}
```

#### Objeto Student (Response)
```typescript
interface Student {
  id: string;                          // UUID
  full_name: string;
  document_id?: string;
  birth_date: string;                  // ISO 8601
  profile_photo_url?: string;
  city_origin_id?: string;             // UUID
  country_origin_id: string;           // UUID
  emails: string[];
  phones?: string[];
  company_id?: string;                 // UUID
  status: 'active' | 'graduated' | 'withdrawn' | 'suspended';
  cohort: string;
  enrollment_date: string;             // ISO 8601
  graduation_date?: string;            // ISO 8601
  created_at: string;                  // ISO 8601
  created_by?: string;                 // UUID
  updated_at: string;                  // ISO 8601
  updated_by?: string;                 // UUID
}
```

#### Envelope de Response del API
```typescript
interface APIResponse<T> {
  success: boolean;
  message: string;
  data?: T;
  error?: string;
}

interface PaginatedData<T> {
  items: T[];
  total: number;
  limit: number;
  offset: number;
}
```

### Catálogos Disponibles (para selects/dropdowns)

Las siguientes tablas tienen datos en la DB y tendrán endpoints próximamente:
- **countries**: id, name, iso_code
- **cities**: id, name, country_id
- **companies**: id, name
- **universities**: id, name, country_id

Por ahora, para el campo `country_origin_id` se puede usar el UUID de Colombia: `1bc87fb3-2bd9-47b6-930b-eba0d6a36bd8`. Los endpoints de catálogos se implementarán después.

## Estructura Recomendada para `admin-web`

```
admin-web/
├── src/
│   ├── lib/
│   │   ├── api/                    # Cliente API y tipos
│   │   │   ├── client.ts           # Axios instance configurada
│   │   │   ├── students.ts         # Funciones del API de estudiantes
│   │   │   └── types.ts            # Interfaces TypeScript
│   │   ├── components/             # Componentes reutilizables
│   │   │   ├── ui/                 # Botones, inputs, modals, etc.
│   │   │   └── students/           # Componentes específicos de estudiantes
│   │   └── stores/                 # Svelte stores (estado global)
│   ├── routes/
│   │   ├── +layout.svelte          # Layout principal (sidebar, nav)
│   │   ├── +page.svelte            # Dashboard / home
│   │   └── students/
│   │       ├── +page.svelte        # Lista de estudiantes
│   │       ├── new/
│   │       │   └── +page.svelte    # Formulario crear estudiante
│   │       └── [id]/
│   │           ├── +page.svelte    # Detalle de estudiante
│   │           └── edit/
│   │               └── +page.svelte # Formulario editar estudiante
│   ├── app.html                    # HTML base
│   └── app.css                     # Tailwind imports
├── static/                         # Assets estáticos
├── svelte.config.js
├── tailwind.config.js
├── postcss.config.js
├── tsconfig.json
├── vite.config.ts
└── package.json
```

## Metodología de Trabajo

### Primera Tarea: Scaffolding + Listado de Estudiantes

Orden recomendado:

1. **Scaffolding de SvelteKit** - Crear archivos base (svelte.config.js, vite.config.ts, app.html, tailwind config, etc.)
2. **Cliente API** (`lib/api/client.ts`) - Axios instance apuntando a `http://localhost:8080`
3. **Tipos TypeScript** (`lib/api/types.ts`) - Interfaces Student, APIResponse, PaginatedData
4. **Funciones API** (`lib/api/students.ts`) - getStudents, getStudent, createStudent, updateStudent, deleteStudent
5. **Layout principal** (`routes/+layout.svelte`) - Sidebar con navegación
6. **Listado de estudiantes** (`routes/students/+page.svelte`) - Tabla con datos del API
7. **Formulario crear** (`routes/students/new/+page.svelte`) - Form con validación
8. **Detalle/editar** - Páginas de detalle y edición

### Convenciones

```typescript
// Nombres de archivos: kebab-case
student-table.svelte
student-form.svelte

// Tipos: PascalCase
interface Student {}
interface CreateStudentRequest {}

// Funciones API: camelCase
async function getStudents(filters?: StudentFilters): Promise<PaginatedData<Student>>
async function createStudent(data: CreateStudentRequest): Promise<Student>

// Stores: camelCase con $prefix en templates
const students = writable<Student[]>([]);
// En template: {$students}

// CSS: Tailwind utility classes, evitar CSS custom cuando sea posible
<button class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">
```

### Manejo de Errores

```typescript
// El API siempre devuelve este envelope:
// { success: boolean, message: string, data?: T, error?: string }

// En el cliente, extraer data o lanzar error:
async function apiCall<T>(fn: () => Promise<AxiosResponse<APIResponse<T>>>): Promise<T> {
  const response = await fn();
  if (!response.data.success) {
    throw new Error(response.data.error || response.data.message);
  }
  return response.data.data as T;
}
```

### Estados de UI

Cada página que consume el API debe manejar 3 estados:
1. **Loading**: Spinner o skeleton mientras carga
2. **Error**: Mensaje de error con opción de reintentar
3. **Success**: Datos renderizados

```svelte
{#if loading}
  <p>Cargando...</p>
{:else if error}
  <p class="text-red-600">{error}</p>
{:else}
  <!-- Contenido -->
{/if}
```

## Interacción con Otros Agentes

### Agente Go/Backend
- **Recibes de él**: Documentación de API, contratos JSON, endpoints disponibles
- **Le proporcionas**: Requerimientos de endpoints faltantes, formatos de datos
- **Coordinas**: Formato de responses, paginación, filtros

### Agente Arquitecto
- **Recibes de él**: Separación admin/portal, principios de UI
- **Le consultas**: Decisiones de routing, estado global vs local
- **Coordinas**: Que no se mezclen responsabilidades entre admin-web y portal-web

### Agente DBA
- **Recibes de él**: Catálogos disponibles (países, ciudades, etc.)
- **Le proporcionas**: Datos que necesitas para dropdowns/selects

## Checklist de Implementación

Para cada página nueva:

- [ ] Tipos TypeScript definidos
- [ ] Función API creada
- [ ] Estado loading/error/success manejado
- [ ] Responsive (mobile + desktop)
- [ ] Formularios con validación client-side
- [ ] Feedback visual (toasts en create/update/delete)
- [ ] Navegación funcional (links, breadcrumbs)

## Recursos

- [SvelteKit Docs](https://kit.svelte.dev/docs)
- [Svelte Tutorial](https://svelte.dev/tutorial)
- [Tailwind CSS Docs](https://tailwindcss.com/docs)
- [Axios Docs](https://axios-http.com/docs/intro)

---

**Última actualización**: 2026-02-13
