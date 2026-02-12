# Coordinador - Sistema de Gestión Académica de Maestría

Sistema integral para la coordinación y gestión académica de programas de maestría, diseñado para facilitar el seguimiento de estudiantes, planificación de cursos, asignación de recursos y generación de reportes analíticos.

## 🎯 Objetivos del Proyecto

1. **Funcional**: Sistema completo de gestión académica
2. **Aprendizaje**: 
   - Desarrollo basado en agentes especializados de IA
   - Arquitectura modular y escalable
   - Stack moderno (Go/Fiber + Svelte)
   - PostgreSQL con optimizaciones para reportes
   - CI/CD y containerización

## 🏗️ Arquitectura

### Monorepo Modular
- **Backend**: API REST en Go/Fiber
- **Admin Web**: Dashboard administrativo en Svelte
- **Portal Web**: Portal para estudiantes/profesores/tutores en Svelte
- **Shared**: Tipos y contratos compartidos

### Patrón CQRS Light
- **Write Path**: Operaciones transaccionales (CRUD)
- **Read Path**: Consultas optimizadas con vistas materializadas

Ver [ADR-001](./docs/adrs/001-arquitectura-general.md) para decisiones arquitectónicas detalladas.

## 📦 Estructura del Proyecto

```
coordinador/
├── backend/              # API en Go/Fiber
│   ├── cmd/api/         # Entry point de la aplicación
│   ├── internal/        # Código privado del backend
│   │   ├── students/    # Módulo de estudiantes
│   │   ├── courses/     # Módulo de cursos
│   │   ├── planning/    # Módulo de planificación académica
│   │   ├── reports/     # Módulo de reportes y analítica
│   │   ├── tutors/      # Módulo de tutores/monitores
│   │   ├── auth/        # Módulo de autenticación
│   │   └── shared/      # Código compartido entre módulos
│   ├── migrations/      # Migraciones de base de datos
│   └── pkg/             # Código público reutilizable
├── admin-web/           # Dashboard administrativo (Svelte)
├── portal-web/          # Portal de usuarios (Svelte)
├── shared/              # Recursos compartidos entre proyectos
│   └── types/          # Definiciones de tipos
├── docs/                # Documentación
│   ├── adrs/           # Architectural Decision Records
│   ├── domain/         # Modelos de dominio
│   └── agents/         # Guías de agentes especializados
└── scripts/             # Scripts de utilidad

```

## 🎭 Agentes Especializados

Este proyecto utiliza agentes de IA especializados para diferentes aspectos del desarrollo:

- **Agente Arquitecto**: Decisiones de diseño y arquitectura
- **Agente Go/Backend**: Implementación del backend en Go
- **Agente DBA**: Diseño de base de datos y optimizaciones
- **Agente Svelte**: Desarrollo de interfaces de usuario
- **Agente DevOps**: CI/CD, containerización y despliegue

Ver [/docs/agents](./docs/agents/) para guías detalladas de cada agente.

## 🚀 Inicio Rápido

### Prerrequisitos
- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Docker & Docker Compose (opcional)
- Git

### Clonar el Proyecto

```bash
# Con SSH (recomendado)
git clone git@github.com:dcorreal/coordinador.git
cd coordinador

# Con HTTPS
git clone https://github.com/dcorreal/coordinador.git
cd coordinador

# Backend
cd backend
go mod download
go run cmd/api/main.go

# Admin Web
cd ../admin-web
npm install
npm run dev

# Portal Web
cd ../portal-web
npm install
npm run dev
```

### Con Docker Compose

```bash
docker-compose up
```

## 📊 Dominio del Negocio

### Entidades Principales

1. **Estudiantes**: Gestión de estudiantes (activos, graduados, históricos)
2. **Cursos**: Catálogo de cursos (obligatorios/electivos, créditos, prerrequisitos)
3. **Pensum**: Estructura del programa de maestría
4. **Programación**: Oferta de cursos por período académico
5. **Inscripciones**: Relación estudiante-curso-período
6. **Calificaciones**: Resultados académicos
7. **Profesores**: Gestión de docentes y asignaciones
8. **Tutores/Monitores**: Pool de apoyo académico

### Funcionalidades Clave

#### Para Coordinadores (Admin Web)
- Carga y gestión de estudiantes (nuevos, activos, graduados)
- Diseño y mantenimiento del pensum
- Programación de períodos académicos
- Asignación de profesores y tutores
- Reportes y analítica:
  - Hoja de vida de estudiantes
  - Desempeño académico
  - Proyecciones de inscripciones
  - Tasas de graduación y deserción
  - Análisis de cohortes

#### Para Estudiantes (Portal Web)
- Consulta de hoja de vida académica
- Visualización de cursos disponibles
- Seguimiento de progreso en el programa

#### Para Profesores (Portal Web)
- Visualización de cursos asignados
- Selección de tutores/monitores
- Gestión de calificaciones

#### Para Tutores (Portal Web)
- Consulta de cursos donde apoyan
- Registro de calificaciones

## 🛠️ Stack Tecnológico

### Backend
- **Lenguaje**: Go 1.21+
- **Framework**: Fiber (HTTP framework)
- **Base de datos**: PostgreSQL 15
- **ORM**: GORM / sqlx
- **Migraciones**: golang-migrate

### Frontend
- **Framework**: Svelte + SvelteKit
- **UI**: TailwindCSS
- **Estado**: Svelte Stores
- **HTTP Client**: Fetch API / Axios

### DevOps
- **Containerización**: Docker
- **Orquestación**: Docker Compose (desarrollo)
- **CI/CD**: GitHub Actions
- **Testing**: Go testing + Vitest (Svelte)

## 📚 Documentación

- [ADRs](./docs/adrs/): Decisiones arquitectónicas
- [Modelo de Dominio](./docs/domain/): Entidades y reglas de negocio
- [Guías de Agentes](./docs/agents/): Instrucciones para agentes especializados
- [API Docs](./docs/api/): Documentación de endpoints (pendiente)

## 🤝 Desarrollo Basado en Agentes

Este proyecto sigue un enfoque de desarrollo asistido por agentes de IA. Cada agente tiene:
- Contexto específico de su dominio
- Memoria de decisiones previas (vía ADRs)
- Expertise en su área técnica
- Responsabilidades claramente definidas

### Workflow de Desarrollo

1. **Arquitecto**: Define decisiones de alto nivel (ADR)
2. **DBA**: Diseña modelo de datos y optimizaciones
3. **Backend**: Implementa lógica de negocio y APIs
4. **Frontend**: Construye interfaces de usuario
5. **DevOps**: Automatiza deployment y operaciones

## 📝 Licencia

[Definir licencia]

## 👤 Autor

Dario Correal

---

## 📚 Recursos Adicionales

- [QUICKSTART.md](./QUICKSTART.md) - Inicio en 5 minutos
- [CONTRIBUTING.md](./CONTRIBUTING.md) - Guía de contribución
- [PROJECT_STATUS.md](./PROJECT_STATUS.md) - Estado actual del proyecto
- [GITHUB_SETUP.md](./GITHUB_SETUP.md) - Configurar Git y GitHub
- [GIT_REFERENCE.md](./GIT_REFERENCE.md) - Referencia rápida de Git

---

**Nota**: Este es un proyecto en desarrollo activo. La documentación se actualizará continuamente.
