<script lang="ts">
	import '@cyberpunk-apps/neondeck/styles.css';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { AppShell, Tag } from '@cyberpunk-apps/neondeck';
	import { app, backend, inWails, prefs, loadPrefs, type AppInfo } from '#lib';
	import WindowControls from '#lib/WindowControls.svelte';
	import type { LayoutProps } from './$types';

	let { children }: LayoutProps = $props();
	let info = $state<AppInfo | null>(null);

	onMount(async () => {
		info = await backend.appInfo();
		await loadPrefs();
	});

	// With the hash router, page.url.hash holds the route ("#/settings"); "" means home.
	const current = $derived(page.url.hash || '#/');
</script>

<AppShell brand={app.brand} home="#/" scanlines={prefs.value?.scanlines ?? false}>
	{#snippet actions()}
		<WindowControls />
	{/snippet}

	{#snippet sidebar()}
		<nav class="side" aria-label="App">
			{#each app.nav as item (item.href)}
				<a href={item.href} class:active={current === item.href} aria-current={current === item.href ? 'page' : undefined}>
					<span class="zh" lang="zh-Hans">{item.zh}</span>{item.label}
				</a>
			{/each}
		</nav>
	{/snippet}

	{#snippet status()}
		<span><Tag tone={inWails ? 'success' : 'warning'} dot>{inWails ? 'native' : 'browser'}</Tag></span>
		<span>{info ? `${info.os}/${info.arch}` : '…'}</span>
		<span>{info?.goVersion ?? ''}</span>
		<span style="margin-left:auto">v{info?.version ?? '…'} // root@{app.name}:~#</span>
	{/snippet}

	{@render children()}
</AppShell>

<style>
	/* The top bar is the window's title bar: drag it to move the window. */
	:global(.nd-shell > .topbar) { --wails-draggable: drag; }
	:global(.nd-shell > .topbar a, .nd-shell > .topbar button) { --wails-draggable: no-drag; }

	.side { display: grid; gap: 2px; }
	.side a {
		display: flex;
		align-items: baseline;
		gap: var(--nd-space-3);
		padding: var(--nd-space-2) var(--nd-space-3);
		border-left: 2px solid transparent;
		color: var(--nd-text-dim);
		font-family: var(--nd-font-ui);
		font-weight: 600;
		letter-spacing: var(--nd-tracking-label);
		text-decoration: none;
		text-transform: uppercase;
	}
	.side a:hover { background: var(--nd-surface-2); color: var(--nd-text); text-shadow: none; }
	.side a.active { border-left-color: var(--nd-accent); background: var(--nd-accent-tint); color: var(--nd-accent); }
	.zh { font-family: var(--nd-font-cjk); font-weight: 900; letter-spacing: 0.04em; }
</style>
