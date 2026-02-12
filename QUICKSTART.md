# Inicio Rápido - Coordinador

Guía para tener el proyecto funcionando en 5 minutos.

## ⚡ Setup Rápido

### Opción 1: Con Docker (Recomendado)

```bash
# 1. Clonar el repositorio
git clone <repository-url>
cd coordinador

# 2. Levantar base de datos
docker-compose up -d postgres

# 3. Esperar a que PostgreSQL esté listo (10-15 segundos)
docker-compose logs -f postgres
# Busca: "database system is ready to accept connections"

# Listo! PostgreSQL corriendo en localhost:5432
```

**Credenciales de desarrollo:**
- Host: `localhost`
- Puerto: `5432`
- Usuario: `coordinador`
- Password: `coordinador_dev_2024`
- Database: `coordinador_db`

### Opción 2: Sin Docker

```bash
# 1. Instalar PostgreSQL 15+
# (según tu sistema operativo)

# 2. Crear base de datos
createdb coordinador_db

# 3. Crear usuario
psql -c "CREATE USER coordinador WITH PASSWORD 'coordinador_dev_2024';"
psql -c "GRANT ALL PRIVILEGES ON DATABASE coordinador_db TO coordinador;"
```

## 🚀 Iniciar Backend

```bash
cd backend

# Instalar dependencias
go mod download

# Iniciar servidor
go run cmd/api/main.go

# Deberías ver:
# 🚀 Coordinador API iniciando en http://localhost:8080
```

**Verificar que funciona:**
```bash
curl http://localhost:8080/health

# Respuesta esperada:
# {"status":"ok","service":"coordinador-api","version":"0.1.0"}
```

## 🎨 Iniciar Frontends

### Admin Web (Dashboard Coordinador)

```bash
# En una nueva terminal
cd admin-web

# Instalar dependencias (primera vez)
npm install

# Iniciar en modo desarrollo
npm run dev

# Disponible en: http://localhost:3000
```

### Portal Web (Estudiantes/Profesores/Tutores)

```bash
# En otra terminal
cd portal-web

# Instalar dependencias (primera vez)
npm install

# Iniciar en modo desarrollo
npm run dev

# Disponible en: http://localhost:3001
```

## ✅ Verificación

Deberías tener corriendo:

- ✅ PostgreSQL en `localhost:5432`
- ✅ Backend API en `http://localhost:8080`
- ✅ Admin Web en `http://localhost:3000` (próximamente)
- ✅ Portal Web en `http://localhost:3001` (próximamente)

## 🧪 Probar Endpoints

```bash
# Health check
curl http://localhost:8080/health

# Endpoints placeholder
curl http://localhost:8080/api/v1/students
curl http://localhost:8080/api/v1/courses
curl http://localhost:8080/api/v1/reports
```

## 📝 Próximos Pasos

Ahora que tienes todo corriendo:

1. **Lee la documentación**:
   - [README principal](./README.md)
   - [ADR-001: Arquitectura](./docs/adrs/001-arquitectura-general.md)
   - [Guía de Contribución](./CONTRIBUTING.md)

2. **Explora la estructura**:
   ```bash
   tree -L 2 -I 'node_modules|bin'
   ```

3. **Entiende el flujo de agentes**:
   - [Agente Arquitecto](./docs/agents/agente-arquitecto.md)
   - [Índice de Agentes](./docs/agents/README.md)

4. **Comienza a desarrollar**:
   - Ver issues abiertos
   - Consultar al agente correspondiente
   - Seguir workflow de [CONTRIBUTING.md](./CONTRIBUTING.md)

## 🛠️ Comandos Útiles (con Make)

```bash
# Ver todos los comandos disponibles
make help

# Instalar todas las dependencias
make install

# Iniciar backend
make dev-backend

# Iniciar admin web
make dev-admin

# Iniciar portal web
make dev-portal

# Ejecutar todos los tests
make test

# Build de todo
make build

# Limpiar builds
make clean
```

## 🐛 Troubleshooting

### PostgreSQL no inicia

```bash
# Verificar que el puerto 5432 no está ocupado
lsof -i :5432

# Ver logs de Docker
docker-compose logs postgres

# Reiniciar container
docker-compose restart postgres
```

### Backend no conecta a BD

Verificar variables de entorno en `backend/.env`:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=coordinador
DB_PASSWORD=coordinador_dev_2024
DB_NAME=coordinador_db
DB_SSLMODE=disable
```

### Frontend no conecta a API

Verificar `admin-web/.env.local` y `portal-web/.env.local`:
```env
VITE_API_URL=http://localhost:8080/api/v1
```

### Puerto ocupado

```bash
# Si 8080 está ocupado, cambiar en backend/.env
PORT=8081

# Si 3000 o 3001 están ocupados
# Cambiar en package.json del frontend correspondiente
```

## 📞 Ayuda

Si algo no funciona:

1. Verifica que tienes las versiones correctas:
   ```bash
   go version   # 1.21+
   node -v      # 18+
   psql --version  # 15+
   ```

2. Revisa logs:
   ```bash
   # Backend
   go run cmd/api/main.go

   # PostgreSQL
   docker-compose logs postgres

   # Frontend
   npm run dev
   ```

3. Consulta [CONTRIBUTING.md](./CONTRIBUTING.md) para más detalles

4. Abre un issue describiendo el problema

---

## 🎯 Estado Actual del Proyecto

**Fase**: Setup Inicial ✅

**Completado**:
- ✅ Estructura del proyecto
- ✅ Backend básico con Fiber
- ✅ PostgreSQL con Docker
- ✅ Frontends configurados (esqueleto)
- ✅ Documentación base
- ✅ Agente Arquitecto activo

**Próximo**:
- 🔄 Diseño de modelo de datos (Agente DBA)
- 🔄 Implementación de módulos backend
- 🔄 Desarrollo de componentes frontend
- 🔄 CI/CD pipeline

Ver [README.md](./README.md) para roadmap completo.

---

**¡Listo para desarrollar!** 🚀
