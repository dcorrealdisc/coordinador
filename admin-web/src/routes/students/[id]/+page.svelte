<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { getStudent, deleteStudent } from '$lib/api/students';
	import { goto } from '$app/navigation';
	import type { Student } from '$lib/api/types';

	let student: Student | null = null;
	let loading = true;
	let error = '';

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

	const genderLabels: Record<string, string> = {
		M: 'Masculino',
		F: 'Femenino'
	};

	function formatDate(iso: string | undefined): string {
		if (!iso) return '—';
		return new Date(iso).toLocaleDateString('es-CO', {
			year: 'numeric',
			month: 'long',
			day: 'numeric'
		});
	}

	async function handleDelete() {
		if (!student) return;
		if (!confirm(`¿Eliminar a ${student.first_names} ${student.last_names}?`)) return;
		try {
			await deleteStudent(student.id);
			await goto('/students');
		} catch (e) {
			alert(e instanceof Error ? e.message : 'Error al eliminar');
		}
	}

	onMount(async () => {
		try {
			student = await getStudent($page.params.id);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Error al cargar estudiante';
		} finally {
			loading = false;
		}
	});
</script>

<svelte:head>
	<title>{student ? `${student.first_names} ${student.last_names}` : 'Estudiante'} - Coordinador</title>
</svelte:head>

<div class="p-8 max-w-3xl">
	<div class="mb-6">
		<a href="/students" class="text-blue-600 hover:underline text-sm">&larr; Volver a estudiantes</a>
	</div>

	{#if loading}
		<div class="bg-white rounded-lg shadow p-12 text-center text-gray-500">
			Cargando...
		</div>
	{:else if error}
		<div class="bg-white rounded-lg shadow p-12 text-center">
			<p class="text-red-600 mb-4">{error}</p>
			<a href="/students" class="text-blue-600 hover:underline text-sm">Volver al listado</a>
		</div>
	{:else if student}
		<div class="bg-white rounded-lg shadow">
			<!-- Header -->
			<div class="px-6 py-5 border-b flex items-center justify-between">
				<div>
					<h1 class="text-2xl font-bold text-gray-900">
						{student.first_names} {student.last_names}
					</h1>
					<div class="flex items-center gap-3 mt-1">
						{#if student.student_code}
							<span class="text-sm font-mono text-gray-500">{student.student_code}</span>
						{/if}
						<span class="inline-block px-2 py-0.5 text-xs font-medium rounded-full {statusColors[student.status] ?? ''}">
							{statusLabels[student.status] ?? student.status}
						</span>
					</div>
				</div>
				<div class="flex gap-2">
					<a
						href="/students/{student.id}/edit"
						class="border border-gray-300 text-gray-700 px-4 py-2 rounded-lg text-sm font-medium hover:bg-gray-50 transition-colors"
					>
						Editar
					</a>
					<button
						on:click={handleDelete}
						class="border border-red-300 text-red-600 px-4 py-2 rounded-lg text-sm font-medium hover:bg-red-50 transition-colors"
					>
						Eliminar
					</button>
				</div>
			</div>

			<!-- Details -->
			<div class="px-6 py-5 space-y-6">
				<!-- Personal -->
				<section>
					<h2 class="text-sm font-semibold text-gray-500 uppercase tracking-wide mb-3">Datos personales</h2>
					<dl class="grid grid-cols-2 gap-x-6 gap-y-4">
						<div>
							<dt class="text-xs text-gray-500">Nombres</dt>
							<dd class="text-sm text-gray-900 mt-0.5">{student.first_names}</dd>
						</div>
						<div>
							<dt class="text-xs text-gray-500">Apellidos</dt>
							<dd class="text-sm text-gray-900 mt-0.5">{student.last_names}</dd>
						</div>
						<div>
							<dt class="text-xs text-gray-500">Documento</dt>
							<dd class="text-sm text-gray-900 mt-0.5">{student.document_id ?? '—'}</dd>
						</div>
						<div>
							<dt class="text-xs text-gray-500">Fecha de nacimiento</dt>
							<dd class="text-sm text-gray-900 mt-0.5">{formatDate(student.birth_date)}</dd>
						</div>
						<div>
							<dt class="text-xs text-gray-500">Genero</dt>
							<dd class="text-sm text-gray-900 mt-0.5">{student.gender ? genderLabels[student.gender] ?? student.gender : '—'}</dd>
						</div>
					</dl>
				</section>

				<!-- Contact -->
				<section>
					<h2 class="text-sm font-semibold text-gray-500 uppercase tracking-wide mb-3">Contacto</h2>
					<dl class="grid grid-cols-2 gap-x-6 gap-y-4">
						<div>
							<dt class="text-xs text-gray-500">Emails</dt>
							<dd class="text-sm text-gray-900 mt-0.5">
								{#if student.emails && student.emails.length > 0}
									{#each student.emails as email}
										<div>{email}</div>
									{/each}
								{:else}
									—
								{/if}
							</dd>
						</div>
						<div>
							<dt class="text-xs text-gray-500">Telefonos</dt>
							<dd class="text-sm text-gray-900 mt-0.5">
								{#if student.phones && student.phones.length > 0}
									{#each student.phones as phone}
										<div>{phone}</div>
									{/each}
								{:else}
									—
								{/if}
							</dd>
						</div>
					</dl>
				</section>

				<!-- Academic -->
				<section>
					<h2 class="text-sm font-semibold text-gray-500 uppercase tracking-wide mb-3">Academico</h2>
					<dl class="grid grid-cols-2 gap-x-6 gap-y-4">
						<div>
							<dt class="text-xs text-gray-500">Codigo de estudiante</dt>
							<dd class="text-sm font-mono text-gray-900 mt-0.5">{student.student_code ?? '—'}</dd>
						</div>
						<div>
							<dt class="text-xs text-gray-500">Estado</dt>
							<dd class="text-sm text-gray-900 mt-0.5">
								<span class="inline-block px-2 py-0.5 text-xs font-medium rounded-full {statusColors[student.status] ?? ''}">
									{statusLabels[student.status] ?? student.status}
								</span>
							</dd>
						</div>
						<div>
							<dt class="text-xs text-gray-500">Cohorte</dt>
							<dd class="text-sm text-gray-900 mt-0.5">{student.cohort}</dd>
						</div>
						<div>
							<dt class="text-xs text-gray-500">Fecha de ingreso</dt>
							<dd class="text-sm text-gray-900 mt-0.5">{formatDate(student.enrollment_date)}</dd>
						</div>
						{#if student.graduation_date}
							<div>
								<dt class="text-xs text-gray-500">Fecha de graduacion</dt>
								<dd class="text-sm text-gray-900 mt-0.5">{formatDate(student.graduation_date)}</dd>
							</div>
						{/if}
					</dl>
				</section>

				<!-- Audit -->
				<section class="border-t pt-4">
					<dl class="grid grid-cols-2 gap-x-6 gap-y-2 text-xs text-gray-400">
						<div>Creado: {formatDate(student.created_at)}</div>
						<div>Actualizado: {formatDate(student.updated_at)}</div>
					</dl>
				</section>
			</div>
		</div>
	{/if}
</div>
