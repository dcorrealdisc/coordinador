<script lang="ts">
	import { goto } from '$app/navigation';
	import { createStudent } from '$lib/api/students';

	let fullName = '';
	let documentId = '';
	let birthDate = '';
	let email = '';
	let phone = '';
	let nationalityCountryId = '';
	let residenceCountryId = '';
	let residenceCityId = '';
	let studentCode = '';
	let status = 'active';
	let cohort = '';
	let enrollmentDate = '';

	let saving = false;
	let error = '';

	// UUIDs reales de la tabla countries (hardcoded hasta que exista el API de catálogos)
	const COUNTRIES: Record<string, string> = {
		'f5f3d79f-6a18-4ccf-99a5-bae8d8d99b3d': 'Argentina',
		'4af2b16d-04c7-4c2c-8918-7d036151367b': 'Bolivia',
		'4b851a3b-b4d8-49eb-9a0a-22f2cba47550': 'Brasil',
		'8ba4e777-f638-4adc-8f06-d8388f180606': 'Canada',
		'177809c9-fb44-48bb-a8d6-91ee7a040ada': 'Chile',
		'1bc87fb3-2bd9-47b6-930b-eba0d6a36bd8': 'Colombia',
		'ece29466-fa59-4eea-b922-db72b08fc4e0': 'Ecuador',
		'f7e2cde3-ad75-4ab7-8c70-61b75330f8b6': 'España',
		'88aa7905-d4e3-4e47-8c4c-e5a86db5aac9': 'Estados Unidos',
		'8955d0a7-6378-4c66-ac5b-8b795890fffc': 'Mexico',
		'4c5ad867-9d1d-4bca-8fea-a3d01f2dc3dc': 'Paraguay',
		'ea7ca018-8bf7-4bf2-9e52-7ea343fc4cc6': 'Peru',
		'2875ed1a-1f25-459c-8c0e-28b2314f587b': 'Uruguay',
		'd80250ae-34bd-4bd9-af8a-0732c62b4b02': 'Venezuela',
	};

	async function handleSubmit() {
		error = '';
		saving = true;
		try {
			await createStudent({
				full_name: fullName,
				document_id: documentId || undefined,
				birth_date: birthDate,
				nationality_country_id: nationalityCountryId,
				residence_country_id: residenceCountryId,
				residence_city_id: residenceCityId || undefined,
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

		<div class="grid grid-cols-3 gap-4">
			<div>
				<label for="nationalityCountryId" class="block text-sm font-medium text-gray-700 mb-1">
					Nacionalidad *
				</label>
				<select
					id="nationalityCountryId"
					bind:value={nationalityCountryId}
					required
					class="w-full border border-gray-300 rounded-lg px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
				>
					<option value="">Seleccionar...</option>
					{#each Object.entries(COUNTRIES) as [id, name]}
						<option value={id}>{name}</option>
					{/each}
				</select>
			</div>
			<div>
				<label for="residenceCountryId" class="block text-sm font-medium text-gray-700 mb-1">
					Pais de residencia *
				</label>
				<select
					id="residenceCountryId"
					bind:value={residenceCountryId}
					required
					class="w-full border border-gray-300 rounded-lg px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
				>
					<option value="">Seleccionar...</option>
					{#each Object.entries(COUNTRIES) as [id, name]}
						<option value={id}>{name}</option>
					{/each}
				</select>
			</div>
			<div>
				<label for="residenceCityId" class="block text-sm font-medium text-gray-700 mb-1">
					Ciudad de residencia
				</label>
				<input
					id="residenceCityId"
					type="text"
					bind:value={residenceCityId}
					placeholder="UUID de ciudad"
					class="w-full border border-gray-300 rounded-lg px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
				/>
				<p class="text-xs text-gray-400 mt-1">UUID hasta que exista selector de ciudades</p>
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
