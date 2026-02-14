# Agentes Especializados

Este directorio contiene las guías de trabajo para los agentes especializados de IA que apoyan el desarrollo del proyecto Coordinador.

## 🤖 Concepto de Agentes

Los agentes son asistentes de IA especializados en dominios específicos del proyecto. Cada agente tiene:

- **Expertise específico**: Conocimiento profundo en su área
- **Contexto persistente**: Memoria de decisiones previas vía ADRs
- **Responsabilidades claras**: Scope bien definido
- **Autonomía guiada**: Toman decisiones dentro de su dominio

## 📋 Agentes Activos

| Agente | Estado | Responsabilidades Principales | Guía |
|--------|--------|-------------------------------|------|
| Arquitecto | ✅ Activo | Decisiones arquitectónicas, ADRs, coherencia del sistema | [Ver guía](./agente-arquitecto.md) |
| DBA | ✅ Activo | Diseño de base de datos, optimizaciones, migraciones | [Ver guía](./agente-dba.md) |
| Go/Backend | ✅ Activo | Implementación backend, APIs, lógica de negocio | [Ver guía](./agente-go-backend.md) |
| Svelte | ✅ Activo | Desarrollo frontend, componentes, UX | [Ver guía](./agente-svelte.md) |
| DevOps | 📝 Pendiente | CI/CD, containerización, deployment | _Próximamente_ |

## 🔄 Flujo de Trabajo entre Agentes

```
┌──────────────┐
│  Arquitecto  │ ← Toma decisiones de alto nivel
└──────┬───────┘
       │ Define estructura
       ▼
┌──────────────┐
│     DBA      │ ← Diseña modelo de datos
└──────┬───────┘
       │ Crea schema
       ▼
┌──────────────┐
│ Go/Backend   │ ← Implementa lógica de negocio
└──────┬───────┘
       │ Expone APIs
       ▼
┌──────────────┐
│   Svelte     │ ← Construye interfaces
└──────┬───────┘
       │ Requiere deployment
       ▼
┌──────────────┐
│   DevOps     │ ← Automatiza y despliega
└──────────────┘
```

## 📖 Cómo Usar los Agentes

### Para el Desarrollador (Dario)

1. **Identifica el tipo de tarea**:
   - Decisión arquitectónica → Agente Arquitecto
   - Diseño de tabla → Agente DBA
   - Implementar endpoint → Agente Go/Backend
   - Crear componente → Agente Svelte
   - Setup CI/CD → Agente DevOps

2. **Proporciona contexto**:
   - Menciona ADRs relevantes
   - Referencias a código existente
   - Restricciones o requerimientos

3. **Revisa el output**:
   - Verifica adherencia a principios
   - Asegura consistencia con decisiones previas
   - Valida que se documentó apropiadamente

### Para Claude (AI Assistant)

Al asumir el rol de un agente:

1. **Lee tu guía específica** (archivo en este directorio)
2. **Revisa ADRs relevantes** en `/docs/adrs`
3. **Consulta código existente** para mantener consistencia
4. **Documenta decisiones** si corresponde (actualiza ADRs)
5. **Coordina con otros agentes** si la tarea lo requiere

## 🎯 Principios de los Agentes

### Especialización
- Cada agente es experto en su dominio
- No mezclar responsabilidades entre agentes
- Consultar a otros agentes cuando sea necesario

### Consistencia
- Seguir decisiones en ADRs existentes
- Mantener patrones establecidos
- Proponer cambios arquitectónicos cuando algo no calza

### Documentación
- Documentar decisiones importantes
- Actualizar ADRs cuando corresponda
- Mantener guías de agentes actualizadas

### Pragmatismo
- Soluciones simples sobre complejas
- YAGNI (You Aren't Gonna Need It)
- Balancear pureza con productividad

## 📚 Recursos Compartidos

Todos los agentes deben conocer:

- [README del proyecto](../../README.md)
- [ADR-001: Arquitectura General](../adrs/001-arquitectura-general.md)
- [Estructura del proyecto](../../README.md#-estructura-del-proyecto)

## 🆕 Crear un Nuevo Agente

Cuando se requiera un nuevo agente especializado:

1. Crear archivo `agente-[nombre].md` en este directorio
2. Seguir el template del Agente Arquitecto
3. Definir claramente:
   - Rol y responsabilidades
   - Contexto específico de su dominio
   - Metodología de trabajo
   - Interacción con otros agentes
   - Checklist de revisión
4. Actualizar este índice

## Template de Guía de Agente

```markdown
# Agente [Nombre] - Guía de Trabajo

## 🎯 Rol y Responsabilidades
[Definir qué hace este agente]

## 📚 Contexto del Proyecto
[Información específica del dominio]

## 🔧 Metodología de Trabajo
[Cómo trabajar en este dominio]

## 🚨 Señales de Alerta
[Qué cuestionar, cuándo proponer cambios]

## 🔄 Interacción con Otros Agentes
[Cómo coordinar con otros agentes]

## 📝 Checklist de Revisión
[Lista de verificación antes de completar tareas]

## 🎓 Recursos de Referencia
[Links útiles específicos del dominio]
```

---

**Última actualización**: 2026-02-12
