# 📝 Comandos Git - Referencia Rápida

## 🚀 Setup Inicial (Solo una vez)

```bash
# Opción 1: Usar script automático
./setup-git.sh

# Opción 2: Manual
git config user.name "Dario Correal"
git config user.email "tu@email.com"
git add .
git commit -m "feat: setup inicial del proyecto"
```

## 🔗 Conectar con GitHub (Solo una vez)

```bash
# Después de crear el repo en GitHub
git remote add origin git@github.com:TU_USUARIO/coordinador.git
git push -u origin main
```

## 📅 Workflow Diario

```bash
# 1. Antes de empezar a trabajar
git pull origin main

# 2. Ver qué cambió
git status

# 3. Agregar cambios
git add .                    # Todo
git add archivo.go           # Un archivo
git add backend/             # Una carpeta

# 4. Hacer commit
git commit -m "feat: descripción del cambio"

# 5. Subir a GitHub
git push origin main
```

## 📋 Tipos de Commits (Conventional Commits)

```bash
git commit -m "feat: nueva funcionalidad"
git commit -m "fix: corrección de bug"
git commit -m "docs: actualizar documentación"
git commit -m "refactor: restructurar código"
git commit -m "test: agregar tests"
git commit -m "chore: actualizar dependencias"
git commit -m "perf: mejorar performance"
git commit -m "style: formateo de código"
```

## 🌿 Trabajar con Branches

```bash
# Ver branches
git branch

# Crear nueva branch
git checkout -b feature/nueva-feature

# Cambiar de branch
git checkout main
git checkout develop

# Subir branch a GitHub
git push -u origin feature/nueva-feature

# Eliminar branch local
git branch -d feature/nombre

# Eliminar branch remota
git push origin --delete feature/nombre
```

## 🔄 Sincronizar entre Máquinas

```bash
# En el portátil (terminar el día)
git add .
git commit -m "feat: avance del día"
git push origin main

# En la oficina (empezar el día)
git pull origin main
# ¡Continúa trabajando!

# En la oficina (terminar el día)
git add .
git commit -m "feat: continuar desarrollo"
git push origin main

# En el portátil (al día siguiente)
git pull origin main
# ¡Todo sincronizado!
```

## 📊 Ver Historial

```bash
# Ver commits
git log

# Ver commits en una línea
git log --oneline

# Ver últimos 5 commits
git log -5 --oneline

# Ver cambios de un archivo
git log --follow archivo.go

# Ver quién cambió qué
git blame archivo.go
```

## ↩️ Deshacer Cambios

```bash
# Deshacer cambios en archivo (antes de add)
git checkout -- archivo.go

# Quitar archivo del staging (después de add)
git reset HEAD archivo.go

# Deshacer último commit (mantiene cambios)
git reset --soft HEAD~1

# Deshacer último commit (descarta cambios) ⚠️
git reset --hard HEAD~1

# Revertir commit creando uno nuevo
git revert HEAD
```

## 🔍 Buscar y Comparar

```bash
# Buscar en código
git grep "texto a buscar"

# Ver cambios no commiteados
git diff

# Ver cambios en staging
git diff --staged

# Comparar branches
git diff main develop

# Ver archivos modificados
git status -s
```

## 🏷️ Tags y Versiones

```bash
# Crear tag
git tag -a v0.1.0 -m "Release inicial"

# Listar tags
git tag

# Subir tag a GitHub
git push origin v0.1.0

# Subir todos los tags
git push origin --tags

# Eliminar tag local
git tag -d v0.1.0

# Eliminar tag remoto
git push origin :refs/tags/v0.1.0
```

## 🔧 Configuración

```bash
# Ver configuración
git config --list

# Configurar editor
git config --global core.editor "code"

# Configurar colores
git config --global color.ui true

# Configurar alias
git config --global alias.st status
git config --global alias.co checkout
git config --global alias.br branch
git config --global alias.ci commit

# Usar los alias
git st      # en vez de git status
git co main # en vez de git checkout main
```

## 🆘 Resolver Conflictos

```bash
# Si git pull genera conflictos:

# 1. Ver archivos en conflicto
git status

# 2. Editar archivos y resolver conflictos
#    Busca markers: <<<<<<<, =======, >>>>>>>

# 3. Agregar archivos resueltos
git add archivo-resuelto.go

# 4. Completar merge
git commit -m "fix: resolver conflictos de merge"

# 5. Subir
git push origin main
```

## 🔐 SSH (Recomendado)

```bash
# Generar clave SSH
ssh-keygen -t ed25519 -C "tu@email.com"

# Ver clave pública (para agregar en GitHub)
cat ~/.ssh/id_ed25519.pub

# Probar conexión
ssh -T git@github.com

# Cambiar de HTTPS a SSH
git remote set-url origin git@github.com:USUARIO/coordinador.git
```

## 📦 Stash (Guardar temporalmente)

```bash
# Guardar cambios sin commit
git stash

# Ver stashes guardados
git stash list

# Aplicar último stash
git stash pop

# Aplicar stash específico
git stash apply stash@{0}

# Eliminar stash
git stash drop stash@{0}

# Limpiar todos los stashes
git stash clear
```

## 🔍 Información Útil

```bash
# Ver remote URLs
git remote -v

# Ver qué branch estás siguiendo
git branch -vv

# Ver todos los branches (local y remoto)
git branch -a

# Ver tamaño del repositorio
git count-objects -vH

# Limpiar archivos no rastreados
git clean -fd
```

## 📝 .gitignore

```bash
# Ignorar archivo después de haberlo commiteado
git rm --cached archivo.log
echo "archivo.log" >> .gitignore
git add .gitignore
git commit -m "chore: ignorar archivo.log"

# Ver archivos ignorados
git status --ignored
```

## 🚨 Emergencias

```bash
# Olvidé hacer pull antes de hacer cambios
git stash
git pull origin main
git stash pop

# Hice commit en branch equivocada
git log  # copiar hash del commit
git checkout branch-correcta
git cherry-pick HASH_DEL_COMMIT

# Subí algo que no debía
git revert HEAD
git push origin main

# Necesito volver a un commit anterior
git log --oneline
git checkout HASH_COMMIT  # solo ver
git checkout -b nueva-branch HASH_COMMIT  # crear branch desde ahí
```

## 💡 Tips

1. **Pull frecuente**: `git pull` antes de empezar a trabajar
2. **Commits pequeños**: Commits frecuentes y focalizados
3. **Mensajes claros**: Usa Conventional Commits
4. **Revisa antes**: `git status` y `git diff` antes de commit
5. **Branches**: Usa branches para experimentar

## 🎯 Workflow Recomendado

```bash
# Inicio del día
git checkout main
git pull origin main
git checkout -b feature/nueva-feature

# Durante desarrollo
# ... trabajar ...
git add .
git commit -m "feat: progreso de feature"

# Varios commits durante el día...

# Fin del día
git push -u origin feature/nueva-feature

# Al terminar la feature
# En GitHub: Create Pull Request
# Después de merge:
git checkout main
git pull origin main
git branch -d feature/nueva-feature
```

---

**Ver también**:
- [GITHUB_SETUP.md](./GITHUB_SETUP.md) - Setup completo
- [CONTRIBUTING.md](./CONTRIBUTING.md) - Guía de contribución
- [Git Book](https://git-scm.com/book/es/v2) - Documentación completa
