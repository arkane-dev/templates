<script lang="ts">
	import '@cyberpunk-apps/neondeck/styles.css';
	import favicon from '#lib/assets/favicon.svg';
	import { page } from '$app/state';
	import { AppShell, Button, Tag } from '@cyberpunk-apps/neondeck';
	import { site } from '#lib';
	import type { LayoutProps } from './$types';

	let { children }: LayoutProps = $props();

	const nav = $derived(site.nav.map((n) => ({ ...n, active: page.url.pathname === n.href })));
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<AppShell brand={site.brand} railCaption="{site.name} // v{site.version}" {nav}>
	{#snippet actions()}
		<Button size="sm" arrow href="/about">Contact</Button>
	{/snippet}
	{#snippet status()}
		<span><Tag tone="success" dot>online</Tag></span>
		<span>{page.url.pathname}</span>
		<span style="margin-left:auto">v{site.version} // root@{site.name}:~#</span>
	{/snippet}

	{@render children()}
</AppShell>
