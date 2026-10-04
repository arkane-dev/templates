<!-- Dashboard: the NEONDECK app anatomy. Panel grid, readouts, a Go round trip. -->
<script lang="ts">
	import { onMount } from 'svelte';
	import { SectionHeader, Panel, Readout, Meter, Button, Input, Tag } from '@cyberpunk-apps/neondeck';
	import { backend, type AppInfo } from '#lib';

	let info = $state<AppInfo | null>(null);
	let msg = $state('hello');
	let reply = $state('');
	let busy = $state(false);

	onMount(async () => (info = await backend.appInfo()));

	async function ping() {
		busy = true;
		reply = await backend.ping(msg);
		busy = false;
	}
</script>

<div class="content">
	<SectionHeader index="01" zh="控制台" title="Dashboard" meta={info?.hostname ?? ''} />
	<div class="dash">
		<Panel title="System" index="02" meta={info?.os ?? '…'} class="span-2">
			<div class="readouts">
				<Readout label="CPUs" value={info?.cpus ?? '—'} neon />
				<Readout label="Arch" value={info?.arch ?? '—'} />
				<Readout label="Version" value={info?.version ?? '—'} />
			</div>
			<Meter label="Example load" value={42} meta="replace with real data" />
		</Panel>

		<Panel title="Backend" index="03" accent="cyan" active>
			<p>Round trip to Go through the generated bindings.</p>
			<Input label="Message" bind:value={msg} />
			<div class="row">
				<Button arrow onclick={ping} disabled={busy}>Ping Go</Button>
				{#if reply}<Tag tone="success">ok</Tag>{/if}
			</div>
			{#if reply}<pre>{reply}</pre>{/if}
		</Panel>
	</div>
</div>

<style>
	.content { padding: var(--nd-space-10) var(--nd-gutter); }
	.dash { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: var(--nd-space-5); }
	.dash :global(.span-2) { grid-column: span 2; }
	.readouts { display: grid; grid-template-columns: repeat(3, 1fr); gap: var(--nd-space-6); margin-bottom: var(--nd-space-6); }
	.row { display: flex; align-items: center; gap: var(--nd-space-3); margin-top: var(--nd-space-4); }
	pre { margin-top: var(--nd-space-4); white-space: pre-wrap; }
	@media (max-width: 1100px) {
		.dash { grid-template-columns: 1fr; }
		.dash :global(.span-2) { grid-column: auto; }
	}
</style>
