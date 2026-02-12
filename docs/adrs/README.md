# Architectural Decision Records (ADRs)

Este directorio contiene todos los registros de decisiones arquitectónicas del proyecto Coordinador.

## ¿Qué es un ADR?

Un ADR (Architectural Decision Record) es un documento que captura una decisión arquitectónica importante junto con su contexto y consecuencias. Sirve como memoria histórica del proyecto.

## Índice de ADRs

| ID | Título | Estado | Fecha | Agente |
|----|--------|--------|-------|--------|
| [001](./001-arquitectura-general.md) | Arquitectura General del Sistema | Aceptado | 2026-02-12 | Arquitecto |

## Estados Posibles

- **Propuesto**: Decisión en discusión
- **Aceptado**: Decisión aprobada e implementada
- **Rechazado**: Decisión evaluada pero no adoptada
- **Obsoleto**: Decisión reemplazada por otra posterior

## Proceso para Crear un ADR

1. Copiar el template del [Agente Arquitecto](../agents/agente-arquitecto.md)
2. Asignar el siguiente número secuencial (XXX)
3. Completar todas las secciones del template
4. Actualizar este índice
5. Commit con mensaje: `docs: ADR-XXX [título corto]`

## Convenciones

- **Numeración**: Secuencial de 3 dígitos (001, 002, 003, ...)
- **Nombre de archivo**: `XXX-titulo-kebab-case.md`
- **No borrar**: Los ADRs rechazados u obsoletos se mantienen para historia
- **Referencias**: Usar links relativos entre ADRs relacionados

## Por Tema

### Arquitectura y Diseño
- [ADR-001](./001-arquitectura-general.md): Arquitectura General del Sistema

### Base de Datos
_Pendiente_

### Frontend
_Pendiente_

### Backend
_Pendiente_

### DevOps
_Pendiente_

---

**Última actualización**: 2026-02-12
