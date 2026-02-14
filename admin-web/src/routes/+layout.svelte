<script lang="ts">
	import '../app.css';
	import { page } from '$app/stores';

	const navItems = [
		{ href: '/', label: 'Dashboard', icon: '□' },
		{ href: '/students', label: 'Estudiantes', icon: '◇' }
	];

	function isActive(href: string, pathname: string): boolean {
		if (href === '/') return pathname === '/';
		return pathname.startsWith(href);
	}
</script>

<div class="flex h-screen bg-gray-100">
	<!-- Sidebar -->
	<aside class="w-64 bg-gray-900 text-white flex flex-col">
		<div class="p-6 border-b border-gray-700">
			<h1 class="text-xl font-bold">Coordinador</h1>
			<p class="text-gray-400 text-sm mt-1">Panel Administrativo</p>
		</div>

		<nav class="flex-1 p-4 space-y-1">
			{#each navItems as item}
				<a
					href={item.href}
					class="flex items-center gap-3 px-4 py-3 rounded-lg text-sm transition-colors
						{isActive(item.href, $page.url.pathname)
						? 'bg-blue-600 text-white'
						: 'text-gray-300 hover:bg-gray-800 hover:text-white'}"
				>
					<span class="text-lg">{item.icon}</span>
					{item.label}
				</a>
			{/each}
		</nav>

		<div class="p-4 border-t border-gray-700 text-xs text-gray-500">
			v0.1.0
		</div>
	</aside>

	<!-- Main content -->
	<main class="flex-1 overflow-auto">
		<slot />
	</main>
</div>
