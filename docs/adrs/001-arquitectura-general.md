# ADR-001: Arquitectura General del Sistema

**Estado**: Aceptado  
**Fecha**: 2026-02-12  
**Decisor**: Dario Correal  
**Agente**: Arquitecto

## Contexto y Problema

Se requiere desarrollar un sistema de gestión académica para programas de maestría que permita:

1. **Gestión de estudiantes**: Matrícula, seguimiento académico, historial completo (activos, graduados, desertores)
2. **Planificación académica**: Diseño de pensum, programación de períodos, asignación de recursos
3. **Gestión de cursos**: Catálogo, prerrequisitos, créditos, oferta por período
4. **Gestión de personal**: Profesores, tutores/monitores, asignaciones
5. **Reportes y analítica**: Desempeño individual, reportes consolidados, proyecciones, análisis de cohortes

### Escala del Sistema
- **Estudiantes**: 100-500 activos
- **Usuarios concurrentes**: <100
- **Volumen de datos**: Medio (años de histórico académico)

### Objetivos de Aprendizaje
- Desarrollo basado en agentes especializados de IA
- Arquitectura moderna y escalable
- Stack Go/Fiber (backend) + Svelte (frontend)
- PostgreSQL con optimizaciones para reportes
- CI/CD y containerización

## Decisión

### 1. Arquitectura: Monolito Modular

Adoptamos un **monolito modular** en lugar de microservicios.

**Justificación**:
- Escala del sistema (500 estudiantes) no justifica complejidad de microservicios
- Facilita desarrollo inicial y refactoring
- Transacciones ACID más simples
- Un único punto de deployment
- Módulos bien definidos permiten evolución futura a microservicios si es necesario

**Módulos principales**:
```
internal/
├── students/      # Gestión de estudiantes
├── courses/       # Catálogo y gestión de cursos
├── planning/      # Planificación académica (pensum, programación)
├── reports/       # Reportes y analítica
├── tutors/        # Gestión de tutores/monitores
├── auth/          # Autenticación y autorización
└── shared/        # Código compartido entre módulos
```

### 2. Patrón CQRS Light

Separación lógica entre operaciones de **escritura** (transaccionales) y **lectura** (reportes):

```
┌──────────────────┐      ┌──────────────────┐
│   Write Path     │      │    Read Path     │
│  (Transaccional) │──────│   (Reportes)     │
│                  │ sync │                  │
└──────────────────┘      └──────────────────┘
         │                         │
         ▼                         ▼
    PostgreSQL              Vistas Materializadas
                           + Índices optimizados
```

**Justificación**:
- El sistema tiene dos patrones de uso muy distintos:
  - **Escritura**: CRUD de estudiantes, inscripciones, calificaciones (baja frecuencia)
  - **Lectura**: Reportes complejos, analítica, dashboards (alta frecuencia, queries pesados)
- Vistas materializadas permiten pre-computar reportes complejos
- Mantiene simplicidad de un solo servicio
- Optimización donde realmente se necesita

### 3. Organización del Código: Monorepo

Estructura:
```
coordinador/
├── backend/        # Go/Fiber API
├── admin-web/      # Dashboard administrativo (Svelte)
├── portal-web/     # Portal usuarios (Svelte)
├── shared/         # Tipos compartidos
└── docs/           # Documentación centralizada
```

**Justificación**:
- Versionado sincronizado (un tag = release completo)
- Compartir tipos y contratos entre proyectos
- CI/CD simplificado (un pipeline)
- Refactoring más fácil en cambios transversales
- Mejor para proyecto educativo (visibilidad completa)

### 4. Dos Aplicaciones Frontend Separadas

**Admin Web** (para coordinadores):
- Dashboard complejo con reportes y analítica
- Carga masiva de datos
- Gestión de pensum y programación
- Asignación de recursos
- UX enfocada en productividad

**Portal Web** (para estudiantes/profesores/tutores):
- Interfaces simples y específicas por rol
- Consulta de información personal
- Registro de calificaciones (profesores/tutores)
- UX enfocada en simplicidad y acceso móvil

**Justificación**:
- **Bundle size** optimizado (código de admin no se carga en portal)
- **Seguridad**: Diferentes niveles de autenticación
- **UX diferenciada**: Diseños completamente distintos sin condicionales
- **Deploy independiente**: Actualizar uno sin afectar el otro
- **Mantenibilidad**: Código más limpio y enfocado

### 5. Stack Tecnológico

#### Backend
- **Lenguaje**: Go 1.21+
- **Framework**: Fiber v2
- **Base de datos**: PostgreSQL 15
- **ORM**: GORM (considerando sqlx para queries complejas)

**Justificación Go/Fiber**:
- Rendimiento excepcional para APIs REST
- Type-safe desde el compilador
- Excelente para concurrencia (goroutines)
- Fiber: sintaxis similar a Express, muy rápido
- Binary único para deployment
- Alineado con objetivos de aprendizaje del desarrollador

#### Frontend
- **Framework**: Svelte + SvelteKit
- **Styling**: TailwindCSS
- **Estado**: Svelte Stores

**Justificación Svelte**:
- Liviano y rápido (compilado, no virtual DOM)
- Curva de aprendizaje suave
- Excelente DX (Developer Experience)
- Alineado con objetivos de aprendizaje del desarrollador

#### Base de Datos
- **RDBMS**: PostgreSQL 15
- **Estrategia de reportes**: Vistas materializadas + índices

**Justificación PostgreSQL**:
- Ideal para datos relacionales (estudiantes, cursos, inscripciones)
- Vistas materializadas para reportes complejos
- Excelente soporte para JSON (flexibilidad futura)
- ACID compliant
- Maduro y confiable

### 6. Desarrollo Basado en Agentes

Agentes especializados:
1. **Agente Arquitecto**: Decisiones de diseño, ADRs
2. **Agente Go/Backend**: Implementación backend
3. **Agente DBA**: Diseño DB, optimizaciones
4. **Agente Svelte**: Frontend development
5. **Agente DevOps**: CI/CD, containerización

**Justificación**:
- Expertise especializada en cada área
- Memoria contextual vía ADRs y documentación
- Consistencia en decisiones técnicas
- Objetivo educativo del proyecto

## Consecuencias

### Positivas
- ✅ Arquitectura pragmática para la escala del sistema
- ✅ Desarrollo más rápido inicialmente
- ✅ Menor complejidad operacional
- ✅ Stack moderno y de alto rendimiento
- ✅ Optimización específica para reportes (CQRS light)
- ✅ Dos frontends enfocados en sus usuarios
- ✅ Monorepo facilita gestión del proyecto

### Negativas
- ⚠️ Monolito puede crecer si no se mantiene disciplina modular
- ⚠️ CQRS requiere mantener sincronización de vistas materializadas
- ⚠️ Dos frontends = más código a mantener

### Riesgos y Mitigaciones

**Riesgo**: Monolito se vuelve difícil de mantener
- **Mitigación**: Arquitectura modular estricta, boundaries claros entre módulos

**Riesgo**: Vistas materializadas desactualizadas
- **Mitigación**: Triggers o jobs programados para refresh, monitoring de freshness

**Riesgo**: Duplicación de código entre frontends
- **Mitigación**: Shared components library, tipos compartidos en `/shared`

## Alternativas Consideradas

### Microservicios desde día 1
- **Pros**: Máxima modularidad, escalabilidad independiente
- **Contras**: Overkill para la escala, complejidad operacional alta, transacciones distribuidas complejas
- **Decisión**: Rechazado - La escala no justifica la complejidad

### Monolito sin CQRS
- **Pros**: Más simple
- **Contras**: No optimiza para los reportes pesados que son caso de uso crítico
- **Decisión**: Rechazado - Los reportes justifican la separación read/write

### Stack TypeScript (NestJS + Next.js)
- **Pros**: Un solo lenguaje, ecosystem maduro
- **Contras**: No alineado con objetivos de aprendizaje de Go
- **Decisión**: Rechazado - Prioridad a aprendizaje de Go

### Un solo frontend con roles
- **Pros**: Menos código
- **Contras**: Bundle más grande, UX comprometida, menos seguro
- **Decisión**: Rechazado - Dos audiencias muy distintas justifican separación

## Referencias

- [Go Fiber Documentation](https://docs.gofiber.io/)
- [Svelte Documentation](https://svelte.dev/)
- [PostgreSQL Materialized Views](https://www.postgresql.org/docs/current/rules-materializedviews.html)
- [CQRS Pattern](https://martinfowler.com/bliki/CQRS.html)
- [Modular Monolith](https://www.kamilgrzybek.com/design/modular-monolith-primer/)

## Notas

Esta arquitectura está diseñada para evolucionar. Si la escala crece significativamente (>1000 estudiantes activos, >500 usuarios concurrentes), podemos:
1. Extraer módulo de reportes como servicio separado
2. Implementar caché distribuido (Redis)
3. Separar base de datos read/write físicamente
4. Evolucionar a microservicios selectivamente

La arquitectura modular facilita esta evolución sin reescritura completa.
