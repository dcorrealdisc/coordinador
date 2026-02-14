<script lang="ts">
	import { onMount } from 'svelte';
	import { getStudents, deleteStudent } from '$lib/api/students';
	import type { Student, PaginatedData } from '$lib/api/types';

	let data: PaginatedData<Student> | null = null;
	let loading = true;
	let error = '';

	// Filters
	let search = '';
	let statusFilter = '';
	let searchTimeout: ReturnType<typeof setTimeout>;

	const statusLabels: Record<string, string> = {
		active: 'Activo',
		graduated: 'Graduado',
		withdrawn: 'Retirado',
		suspended: 'Suspendido'
	};

	const statusColors: Record<string, string> = {
		active: 'bg-green-100 text-green-800',
		graduated: 'bg-blue-100 text-blue-800',
		withdrawn: 'bg-gray-100 text-gray-800',
		suspended: 'bg-red-100 text-red-800'
	};

	async function loadStudents() {
		loading = true;
		error = '';
		try {
			data = await getStudents({
				search: search || undefined,
				status: statusFilter || undefined,
				limit: 20,
				offset: 0
			});
		} catch (e) {
			error = e instanceof Error ? e.message : 'Error al cargar estudiantes';
		} finally {
			loading = false;
		}
	}

	function onSearchInput() {
		clearTimeout(searchTimeout);
		searchTimeout = setTimeout(loadStudents, 300);
	}

	async function handleDelete(student: Student) {
		if (!confirm(`Eliminar a ${student.first_names + ' ' + student.last_names}?`)) return;
		try {
			await deleteStudent(student.id);
			await loadStudents();
		} catch (e) {
			alert(e instanceof Error ? e.message : 'Error al eliminar');
		}
	}

	function formatDate(iso: string): string {
		return new Date(iso).toLocaleDateString('es-CO');
	}

	onMount(loadStudents);
</script>

<svelte:head>
	<title>Estudiantes - Coordinador</title>
</svelte:head>

<div class="p-8">
	<!-- Header -->
	<div class="flex items-center justify-between mb-6">
		<div>
			<h1 class="text-2xl font-bold text-gray-900">Estudiantes</h1>
			<p class="text-gray-500 text-sm mt-1">
				{#if data}
					{data.total} estudiante{data.total !== 1 ? 's' : ''} encontrado{data.total !== 1 ? 's' : ''}
				{/if}
			</p>
		</div>
		<div class="flex gap-3">
			<a
				href="/students/import"
				class="border border-gray-300 text-gray-700 px-4 py-2 rounded-lg text-sm font-medium hover:bg-gray-50 transition-colors"
			>
				Importar
			</a>
			<a
				href="/students/new"
				class="bg-blue-600 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-blue-700 transition-colors"
			>
				+ Nuevo Estudiante
			</a>
		</div>
	</div>

	<!-- Filters -->
	<div class="bg-white rounded-lg shadow mb-6 p-4 flex gap-4">
		<input
			type="text"
			placeholder="Buscar por nombre..."
			bind:value={search}
			on:input={onSearchInput}
			class="flex-1 border border-gray-300 rounded-lg px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
		/>
		<select
			bind:value={statusFilter}
			on:change={loadStudents}
			class="border border-gray-300 rounded-lg px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
		>
			<option value="">Todos los estados</option>
			<option value="active">Activo</option>
			<option value="graduated">Graduado</option>
			<option value="withdrawn">Retirado</option>
			<option value="suspended">Suspendido</option>
		</select>
	</div>

	<!-- Content -->
	{#if loading}
		<div class="bg-white rounded-lg shadow p-12 text-center text-gray-500">
			Cargando estudiantes...
		</div>
	{:else if error}
		<div class="bg-white rounded-lg shadow p-12 text-center">
			<p class="text-red-600 mb-4">{error}</p>
			<button
				on:click={loadStudents}
				class="text-blue-600 hover:underline text-sm"
			>
				Reintentar
			</button>
		</div>
	{:else if data && data.items.length === 0}
		<div class="bg-white rounded-lg shadow p-12 text-center text-gray-500">
			No se encontraron estudiantes.
		</div>
	{:else if data}
		<div class="bg-white rounded-lg shadow overflow-hidden">
			<table class="w-full">
				<thead class="bg-gray-50 border-b">
					<tr>
						<th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Código</th>
						<th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Nombre</th>
						<th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Email</th>
						<th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Cohorte</th>
						<th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Estado</th>
						<th class="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Ingreso</th>
						<th class="text-right px-6 py-3 text-xs font-medium text-gray-500 uppercase">Acciones</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-200">
					{#each data.items as student}
						<tr class="hover:bg-gray-50">
							<td class="px-6 py-4 text-sm font-mono text-gray-700">
								{student.student_code ?? '—'}
							</td>
							<td class="px-6 py-4">
								<div class="font-medium text-gray-900">{student.first_names + ' ' + student.last_names}</div>
								{#if student.document_id}
									<div class="text-xs text-gray-500">{student.document_id}</div>
								{/if}
							</td>
							<td class="px-6 py-4 text-sm text-gray-600">
								{student.emails[0]}
							</td>
							<td class="px-6 py-4 text-sm text-gray-600">
								{student.cohort}
							</td>
							<td class="px-6 py-4">
								<span class="inline-block px-2 py-1 text-xs font-medium rounded-full {statusColors[student.status] ?? ''}">
									{statusLabels[student.status] ?? student.status}
								</span>
							</td>
							<td class="px-6 py-4 text-sm text-gray-600">
								{formatDate(student.enrollment_date)}
							</td>
							<td class="px-6 py-4 text-right space-x-2">
								<a
									href="/students/{student.id}"
									class="text-blue-600 hover:underline text-sm"
								>
									Ver
								</a>
								<button
									on:click={() => handleDelete(student)}
									class="text-red-600 hover:underline text-sm"
								>
									Eliminar
								</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>
