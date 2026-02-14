<script lang="ts">
	import { goto } from '$app/navigation';
	import { createStudent } from '$lib/api/students';

	let fullName = '';
	let documentId = '';
	let birthDate = '';
	let email = '';
	let phone = '';
	let studentCode = '';
	let status = 'active';
	let cohort = '';
	let enrollmentDate = '';

	let saving = false;
	let error = '';

	// Colombia UUID hardcoded until catalogs API is ready
	const COLOMBIA_ID = '1bc87fb3-2bd9-47b6-930b-eba0d6a36bd8';

	async function handleSubmit() {
		error = '';
		saving = true;
		try {
			await createStudent({
				full_name: fullName,
				document_id: documentId || undefined,
				birth_date: birthDate,
				country_origin_id: COLOMBIA_ID,
				emails: [email],
				phones: phone ? [phone] : undefined,
				student_code: studentCode || undefined,
				status,
				cohort,
				enrollment_date: enrollmentDate
			});
			await goto('/students');
		} catch (e) {
			error = e instanceof Error ? e.message : 'Error al crear estudiante';
		} finally {
			saving = false;
		}
	}
</script>

<svelte:head>
	<title>Nuevo Estudiante - Coordinador</title>
</svelte:head>

<div class="p-8 max-w-2xl">
	<div class="mb-6">
		<a href="/students" class="text-blue-600 hover:underline text-sm">&larr; Volver a estudiantes</a>
		<h1 class="text-2xl font-bold text-gray-900 mt-2">Nuevo Estudiante</h1>
	</div>

	{#if error}
		<div class="bg-red-50 border border-red-200 text-red-700 rounded-lg p-4 mb-6 text-sm">
			{error}
		</div>
	{/if}

	<form on:submit|preventDefault={handleSubmit} class="bg-white rounded-lg shadow p-6 space-y-5">
		<div>
			<label for="fullName" class="block text-sm font-medium text-gray-700 mb-1">
				Nombre completo *
			</label>
			<input
				id="fullName"
				type="text"
				bind:value={fullName}
				required
				minlength="3"
				class="w-full border border-gray-300 rounded-lg px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
			/>
		</div>

		<div class="grid grid-cols-2 gap-4">
			<div>
				<label for="documentId" class="block text-sm font-medium text-gray-700 mb-1">
					Documento de identidad
				</label>
				<input
					id="documentId"
					type="text"
					bind:value={documentId}
					placeholder="CC-12345678"
					class="w-full border border-gray-300 rounded-lg px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
				/>
			</div>
			<div>
				<label for="birthDate" class="block text-sm font-medium text-gray-700 mb-1">
					Fecha de nacimiento *
				</label>
				<input
					id="birthDate"
					type="date"
					bind:value={birthDate}
					required
					class="w-full border border-gray-300 rounded-lg px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
				/>
			</div>
		</div>

		<div class="grid grid-cols-2 gap-4">
			<div>
				<label for="email" class="block text-sm font-medium text-gray-700 mb-1">
					Email *
				</label>
				<input
					id="email"
					type="email"
					bind:value={email}
					required
					class="w-full border border-gray-300 rounded-lg px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
				/>
			</div>
			<div>
				<label for="phone" class="block text-sm font-medium text-gray-700 mb-1">
					Telefono
				</label>
				<input
					id="phone"
					type="text"
					bind:value={phone}
					placeholder="+57 300 123 4567"
					class="w-full border border-gray-300 rounded-lg px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
				/>
			</div>
		</div>

		<div class="grid grid-cols-2 gap-4">
			<div>
				<label for="studentCode" class="block text-sm font-medium text-gray-700 mb-1">
					Codigo de estudiante
				</label>
				<input
					id="studentCode"
					type="text"
					bind:value={studentCode}
					placeholder="202620190"
					maxlength="9"
					pattern="[0-9]{4}[12][0-9]{4}"
					title="Formato: 4 digitos año + semestre (1 o 2) + 4 digitos secuencia"
					class="w-full border border-gray-300 rounded-lg px-4 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-blue-500"
				/>
				<p class="text-xs text-gray-400 mt-1">Formato: YYYYS#### (ej: 202620190)</p>
			</div>
			<div>
				<label for="cohort" class="block text-sm font-medium text-gray-700 mb-1">
					Cohorte *
				</label>
				<input
					id="cohort"
					type="text"
					bind:value={cohort}
					required
					placeholder="2024-1"
					class="w-full border border-gray-300 rounded-lg px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
				/>
			</div>
		</div>

		<div class="grid grid-cols-2 gap-4">
			<div>
				<label for="enrollmentDate" class="block text-sm font-medium text-gray-700 mb-1">
					Fecha de ingreso *
				</label>
				<input
					id="enrollmentDate"
					type="date"
					bind:value={enrollmentDate}
					required
					class="w-full border border-gray-300 rounded-lg px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
				/>
			</div>
			<div>
				<label for="status" class="block text-sm font-medium text-gray-700 mb-1">
					Estado *
				</label>
				<select
					id="status"
					bind:value={status}
					class="w-full border border-gray-300 rounded-lg px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
				>
					<option value="active">Activo</option>
					<option value="graduated">Graduado</option>
					<option value="withdrawn">Retirado</option>
					<option value="suspended">Suspendido</option>
				</select>
			</div>
		</div>

		<div class="flex gap-3 pt-4">
			<button
				type="submit"
				disabled={saving}
				class="bg-blue-600 text-white px-6 py-2 rounded-lg text-sm font-medium hover:bg-blue-700 transition-colors disabled:opacity-50"
			>
				{saving ? 'Guardando...' : 'Crear Estudiante'}
			</button>
			<a
				href="/students"
				class="border border-gray-300 text-gray-700 px-6 py-2 rounded-lg text-sm font-medium hover:bg-gray-50 transition-colors"
			>
				Cancelar
			</a>
		</div>
	</form>
</div>
