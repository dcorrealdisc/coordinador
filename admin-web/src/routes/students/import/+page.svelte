<script lang="ts">
	import { importStudents } from '$lib/api/students';
	import type { ImportResult } from '$lib/api/types';

	let file: File | null = null;
	let dragOver = false;
	let uploading = false;
	let error = '';
	let result: ImportResult | null = null;

	function handleFileSelect(e: Event) {
		const input = e.target as HTMLInputElement;
		if (input.files && input.files.length > 0) {
			selectFile(input.files[0]);
		}
	}

	function handleDrop(e: DragEvent) {
		dragOver = false;
		if (e.dataTransfer?.files && e.dataTransfer.files.length > 0) {
			selectFile(e.dataTransfer.files[0]);
		}
	}

	function selectFile(f: File) {
		const ext = f.name.split('.').pop()?.toLowerCase();
		if (ext !== 'csv' && ext !== 'xlsx') {
			error = 'Formato no soportado. Solo se aceptan archivos .csv y .xlsx';
			file = null;
			return;
		}
		error = '';
		result = null;
		file = f;
	}

	function removeFile() {
		file = null;
		result = null;
		error = '';
	}

	async function handleUpload() {
		if (!file) return;
		uploading = true;
		error = '';
		result = null;
		try {
			result = await importStudents(file);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Error al importar archivo';
		} finally {
			uploading = false;
		}
	}
</script>

<svelte:head>
	<title>Importar Estudiantes - Coordinador</title>
</svelte:head>

<div class="p-8 max-w-3xl">
	<div class="mb-6">
		<a href="/students" class="text-blue-600 hover:underline text-sm">&larr; Volver a estudiantes</a>
		<h1 class="text-2xl font-bold text-gray-900 mt-2">Importar Estudiantes</h1>
		<p class="text-gray-500 text-sm mt-1">Carga un archivo CSV o Excel (.xlsx) con los datos de los estudiantes.</p>
	</div>

	{#if error}
		<div class="bg-red-50 border border-red-200 text-red-700 rounded-lg p-4 mb-6 text-sm">
			{error}
		</div>
	{/if}

	<!-- Drop zone -->
	{#if !result}
		<div
			class="bg-white rounded-lg shadow p-6 mb-6"
			role="region"
			aria-label="Zona de carga de archivo"
		>
			<div
				class="border-2 border-dashed rounded-lg p-10 text-center transition-colors
					{dragOver ? 'border-blue-500 bg-blue-50' : 'border-gray-300'}"
				on:dragover|preventDefault={() => (dragOver = true)}
				on:dragleave={() => (dragOver = false)}
				on:drop|preventDefault={handleDrop}
				role="button"
				tabindex="0"
			>
				{#if file}
					<div class="flex items-center justify-center gap-3">
						<span class="text-2xl">{file.name.endsWith('.xlsx') ? '📊' : '📄'}</span>
						<div class="text-left">
							<p class="font-medium text-gray-900">{file.name}</p>
							<p class="text-xs text-gray-500">{(file.size / 1024).toFixed(1)} KB</p>
						</div>
						<button
							on:click={removeFile}
							class="ml-4 text-red-500 hover:text-red-700 text-sm"
						>
							Quitar
						</button>
					</div>
				{:else}
					<p class="text-gray-500 mb-3">Arrastra un archivo aqui o</p>
					<label class="inline-block bg-blue-600 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-blue-700 transition-colors cursor-pointer">
						Seleccionar archivo
						<input
							type="file"
							accept=".csv,.xlsx"
							on:change={handleFileSelect}
							class="hidden"
						/>
					</label>
					<p class="text-xs text-gray-400 mt-3">Formatos aceptados: .csv, .xlsx</p>
				{/if}
			</div>

			{#if file}
				<div class="mt-4 flex gap-3">
					<button
						on:click={handleUpload}
						disabled={uploading}
						class="bg-blue-600 text-white px-6 py-2 rounded-lg text-sm font-medium hover:bg-blue-700 transition-colors disabled:opacity-50"
					>
						{uploading ? 'Importando...' : 'Importar'}
					</button>
					<a
						href="/students"
						class="border border-gray-300 text-gray-700 px-6 py-2 rounded-lg text-sm font-medium hover:bg-gray-50 transition-colors"
					>
						Cancelar
					</a>
				</div>
			{/if}
		</div>
	{/if}

	<!-- Results -->
	{#if result}
		<div class="bg-white rounded-lg shadow p-6 mb-6">
			<h2 class="text-lg font-semibold text-gray-900 mb-4">Resultado de la importacion</h2>

			<div class="grid grid-cols-3 gap-4 mb-6">
				<div class="bg-gray-50 rounded-lg p-4 text-center">
					<p class="text-2xl font-bold text-gray-900">{result.total_rows}</p>
					<p class="text-xs text-gray-500 mt-1">Filas totales</p>
				</div>
				<div class="bg-green-50 rounded-lg p-4 text-center">
					<p class="text-2xl font-bold text-green-700">{result.created}</p>
					<p class="text-xs text-green-600 mt-1">Creados</p>
				</div>
				<div class="rounded-lg p-4 text-center {result.errors.length > 0 ? 'bg-red-50' : 'bg-gray-50'}">
					<p class="text-2xl font-bold {result.errors.length > 0 ? 'text-red-700' : 'text-gray-900'}">{result.errors.length}</p>
					<p class="text-xs {result.errors.length > 0 ? 'text-red-600' : 'text-gray-500'} mt-1">Errores</p>
				</div>
			</div>

			{#if result.errors.length > 0}
				<h3 class="text-sm font-semibold text-gray-700 mb-2">Detalle de errores</h3>
				<div class="overflow-auto max-h-80">
					<table class="w-full text-sm">
						<thead class="bg-gray-50 border-b">
							<tr>
								<th class="text-left px-4 py-2 text-xs font-medium text-gray-500 uppercase">Fila</th>
								<th class="text-left px-4 py-2 text-xs font-medium text-gray-500 uppercase">Campo</th>
								<th class="text-left px-4 py-2 text-xs font-medium text-gray-500 uppercase">Valor</th>
								<th class="text-left px-4 py-2 text-xs font-medium text-gray-500 uppercase">Error</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-gray-200">
							{#each result.errors as err}
								<tr>
									<td class="px-4 py-2 font-mono text-gray-700">{err.row}</td>
									<td class="px-4 py-2 font-mono text-gray-700">{err.field}</td>
									<td class="px-4 py-2 text-gray-600 max-w-[200px] truncate">{err.value || '—'}</td>
									<td class="px-4 py-2 text-red-600">{err.message}</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{/if}

			<div class="mt-6 flex gap-3">
				<a
					href="/students"
					class="bg-blue-600 text-white px-6 py-2 rounded-lg text-sm font-medium hover:bg-blue-700 transition-colors"
				>
					Ver estudiantes
				</a>
				<button
					on:click={() => { result = null; file = null; }}
					class="border border-gray-300 text-gray-700 px-6 py-2 rounded-lg text-sm font-medium hover:bg-gray-50 transition-colors"
				>
					Importar otro archivo
				</button>
			</div>
		</div>
	{/if}

	<!-- Column reference -->
	<div class="bg-white rounded-lg shadow p-6">
		<h2 class="text-sm font-semibold text-gray-700 mb-3">Columnas esperadas</h2>
		<div class="overflow-auto">
			<table class="w-full text-xs">
				<thead class="bg-gray-50 border-b">
					<tr>
						<th class="text-left px-3 py-2 font-medium text-gray-500 uppercase">Columna</th>
						<th class="text-left px-3 py-2 font-medium text-gray-500 uppercase">Requerida</th>
						<th class="text-left px-3 py-2 font-medium text-gray-500 uppercase">Ejemplo</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-200">
					{#each [
						{ col: 'first_names', req: true, ex: 'Juan Carlos' },
						{ col: 'last_names', req: true, ex: 'Perez Lopez' },
						{ col: 'document_id', req: false, ex: 'CC-12345678' },
						{ col: 'birth_date', req: true, ex: '1990-05-15' },
						{ col: 'gender', req: false, ex: 'M o F' },
						{ col: 'email', req: true, ex: 'juan@email.com' },
						{ col: 'phone', req: false, ex: '+57 300 1234567' },
						{ col: 'nationality_country_id', req: true, ex: 'UUID' },
						{ col: 'residence_country_id', req: true, ex: 'UUID' },
						{ col: 'residence_city_id', req: false, ex: 'UUID' },
						{ col: 'company_id', req: false, ex: 'UUID' },
						{ col: 'job_title_category_id', req: false, ex: 'UUID' },
						{ col: 'profession_id', req: false, ex: 'UUID' },
						{ col: 'student_code', req: false, ex: '202620190' },
						{ col: 'status', req: true, ex: 'active' },
						{ col: 'cohort', req: true, ex: '2024-1' },
						{ col: 'enrollment_date', req: true, ex: '2024-02-01' },
					] as item}
						<tr>
							<td class="px-3 py-1.5 font-mono text-gray-700">{item.col}</td>
							<td class="px-3 py-1.5">
								{#if item.req}
									<span class="text-red-600 font-medium">Si</span>
								{:else}
									<span class="text-gray-400">No</span>
								{/if}
							</td>
							<td class="px-3 py-1.5 text-gray-500">{item.ex}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>
</div>
