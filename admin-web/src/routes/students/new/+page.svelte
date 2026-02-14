<script lang="ts">
	import { goto } from '$app/navigation';
	import { createStudent } from '$lib/api/students';

	let firstNames = '';
	let lastNames = '';
	let documentId = '';
	let birthDate = '';
	let gender = '';
	let email = '';
	let phone = '';
	let nationalityCountryId = '';
	let residenceCountryId = '';
	let residenceCityId = '';
	let companyId = '';
	let jobTitleCategoryId = '';
	let professionId = '';
	let studentCode = '';
	let status = 'active';
	let cohort = '';
	let enrollmentDate = '';

	let saving = false;
	let error = '';

	// UUIDs reales de la tabla countries
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

	// UUIDs reales de la tabla job_title_categories
	const JOB_TITLE_CATEGORIES: Record<string, string> = {
		'bd18e609-2d03-48a6-b03d-8a3a17a8e88a': 'Analista',
		'98c706ed-0b4d-44c4-9fcd-3b5beb8796fc': 'Consultor/a',
		'128f0bba-0787-4bcb-a57f-ea193e348e5a': 'Coordinador/a',
		'7c31223f-d8c1-4920-849b-dbaf1f7adf00': 'Desarrollador/a',
		'ab6ba022-4d05-44d5-915e-76c5c62aa4d3': 'Director/a',
		'42b747f9-f52a-40ff-93f6-50360921f76f': 'Diseñador/a',
		'88eaa761-180a-4b52-8868-1b742b743eb5': 'Docente',
		'acdda65b-0fe6-4742-b5d0-ef8ba2f53242': 'Especialista',
		'd19fbeb8-5aa8-4076-9b12-df50f3092cd3': 'Gerente',
		'0364559b-e953-446c-9d8b-e74e902141cf': 'Investigador/a',
		'df34c268-3db5-4ec1-84c5-6e0336d15873': 'Jefe de área',
		'ea7a8c1e-267f-4d72-b81a-30ab671e040e': 'Líder técnico',
		'dc82ce9e-748f-4757-8936-c20115991cc9': 'Otro',
		'a0d1644c-efb5-4642-931e-a407501143ba': 'Profesional independiente',
		'f1e13d75-7ecd-470d-a529-c9b76468a066': 'Subgerente',
	};

	// UUIDs reales de la tabla professions
	const PROFESSIONS: Record<string, string> = {
		'c614fc71-54d6-4eb0-a3a0-f62e2cc7be16': 'Administración de Empresas',
		'9fb31981-a492-40c2-975a-7c343800138c': 'Ciencias de la Computación',
		'721ca450-20c5-47ec-b8ae-3fe0bf1838cb': 'Comunicación Social',
		'8393aafa-bf98-49f1-b9b9-760f385a75ee': 'Contaduría Pública',
		'9699d359-0669-4cd6-9ee9-a7297157b372': 'Derecho',
		'1a8a2166-0c20-444a-8916-1c865984f60c': 'Diseño Industrial',
		'02c00df4-18b1-44b7-90ff-b8214e672ce0': 'Economía',
		'ae45f91b-a74f-462d-ad8b-3ef4ad888d7f': 'Estadística',
		'94ce7930-8b10-45ff-b108-522dfe3f49ca': 'Física',
		'5fa844fc-4ff3-4f31-9ee6-593d6227f6f9': 'Ingeniería Civil',
		'304f98ea-f3e0-4658-9beb-d45068839810': 'Ingeniería Electrónica',
		'4f24dae5-31f0-4123-b325-8f41cc8250c3': 'Ingeniería Industrial',
		'69ee0e79-82de-44b0-a356-e112a4281821': 'Ingeniería Mecánica',
		'134e2792-64e5-45a1-9780-cb4f9132c98d': 'Ingeniería de Sistemas',
		'c7aad62c-7568-4e11-af5c-7aec705218f7': 'Ingeniería de Software',
		'30d8a8d5-4904-402b-81a7-a4e5495c6298': 'Ingeniería de Telecomunicaciones',
		'0fc902a4-eeb0-40f5-bca6-94c9401b611e': 'Matemáticas',
		'33a952eb-88bd-4634-9120-a68c39054fbd': 'Otra',
		'd940817f-384f-4d39-9129-0c19d3f97688': 'Psicología',
	};

	async function handleSubmit() {
		error = '';
		saving = true;
		try {
			await createStudent({
				first_names: firstNames,
				last_names: lastNames,
				document_id: documentId || undefined,
				birth_date: birthDate,
				gender: gender ? gender as 'M' | 'F' : undefined,
				nationality_country_id: nationalityCountryId,
				residence_country_id: residenceCountryId,
				residence_city_id: residenceCityId || undefined,
				emails: [email],
				phones: phone ? [phone] : undefined,
				company_id: companyId || undefined,
				job_title_category_id: jobTitleCategoryId || undefined,
				profession_id: professionId || undefined,
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
		<div class="grid grid-cols-2 gap-4">
			<div>
				<label for="firstNames" class="block text-sm font-medium text-gray-700 mb-1">
					Nombres *
				</label>
				<input
					id="firstNames"
					type="text"
					bind:value={firstNames}
					required
					minlength="2"
					class="w-full border border-gray-300 rounded-lg px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
				/>
			</div>
			<div>
				<label for="lastNames" class="block text-sm font-medium text-gray-700 mb-1">
					Apellidos *
				</label>
				<input
					id="lastNames"
					type="text"
					bind:value={lastNames}
					required
					minlength="2"
					class="w-full border border-gray-300 rounded-lg px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
				/>
			</div>
		</div>

		<div class="grid grid-cols-3 gap-4">
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
			<div>
				<label for="gender" class="block text-sm font-medium text-gray-700 mb-1">
					Genero
				</label>
				<select
					id="gender"
					bind:value={gender}
					class="w-full border border-gray-300 rounded-lg px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
				>
					<option value="">Seleccionar...</option>
					<option value="M">Masculino</option>
					<option value="F">Femenino</option>
				</select>
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
				<label for="professionId" class="block text-sm font-medium text-gray-700 mb-1">
					Profesion
				</label>
				<select
					id="professionId"
					bind:value={professionId}
					class="w-full border border-gray-300 rounded-lg px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
				>
					<option value="">Seleccionar...</option>
					{#each Object.entries(PROFESSIONS) as [id, name]}
						<option value={id}>{name}</option>
					{/each}
				</select>
			</div>
			<div>
				<label for="jobTitleCategoryId" class="block text-sm font-medium text-gray-700 mb-1">
					Categoria de cargo
				</label>
				<select
					id="jobTitleCategoryId"
					bind:value={jobTitleCategoryId}
					class="w-full border border-gray-300 rounded-lg px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
				>
					<option value="">Seleccionar...</option>
					{#each Object.entries(JOB_TITLE_CATEGORIES) as [id, name]}
						<option value={id}>{name}</option>
					{/each}
				</select>
			</div>
		</div>

		<div>
			<label for="companyId" class="block text-sm font-medium text-gray-700 mb-1">
				Empresa
			</label>
			<input
				id="companyId"
				type="text"
				bind:value={companyId}
				placeholder="UUID de empresa"
				class="w-full border border-gray-300 rounded-lg px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
			/>
			<p class="text-xs text-gray-400 mt-1">UUID hasta que exista selector de empresas</p>
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
