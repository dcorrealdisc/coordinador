<svelte:head>
	<title>Dashboard - Coordinador</title>
</svelte:head>

<div class="p-8">
	<h1 class="text-2xl font-bold text-gray-900">Dashboard</h1>
	<p class="text-gray-600 mt-2">Bienvenido al panel administrativo del sistema Coordinador.</p>

	<div class="grid grid-cols-1 md:grid-cols-3 gap-6 mt-8">
		<a href="/students" class="bg-white rounded-lg shadow p-6 hover:shadow-md transition-shadow">
			<h3 class="text-lg font-semibold text-gray-900">Estudiantes</h3>
			<p class="text-gray-500 text-sm mt-1">Gestionar estudiantes de la maestria</p>
		</a>

		<div class="bg-white rounded-lg shadow p-6 opacity-50">
			<h3 class="text-lg font-semibold text-gray-900">Cursos</h3>
			<p class="text-gray-500 text-sm mt-1">Proximamente</p>
		</div>

		<div class="bg-white rounded-lg shadow p-6 opacity-50">
			<h3 class="text-lg font-semibold text-gray-900">Reportes</h3>
			<p class="text-gray-500 text-sm mt-1">Proximamente</p>
		</div>
	</div>
</div>
