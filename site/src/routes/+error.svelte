<!-- Error page: "signal lost". -->
<script lang="ts">
	import { page } from '$app/state';
	import { Button, NeonSign } from '@cyberpunk-apps/neondeck';
</script>

<svelte:head><title>{page.status} · Signal lost</title></svelte:head>

<section class="err nd-grid-bg nd-fog">
	<NeonSign text="信号丢失" label="Signal lost" caption="XIN HAO DIU SHI" tone="red" size="clamp(2.5rem, 6vw, 4rem)" flicker />
	<div class="msg">
		<p class="code nd-mono">ERR_{page.status}</p>
		<h1>Signal lost<span class="nd-cursor"></span></h1>
		<p class="lede">{page.error?.message ?? 'Unknown error'}</p>
		<Button arrow href="/">Return to base</Button>
	</div>
</section>

<style>
	.err { display: flex; flex-wrap: wrap; align-items: center; justify-content: center; gap: var(--nd-space-12); min-height: 70vh; padding: var(--nd-space-12) var(--nd-gutter); }
	.code { color: var(--nd-danger); font-size: var(--nd-text-lg); margin: 0; }
	.lede { color: var(--nd-text-dim); }
</style>
