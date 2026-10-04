<!-- Frameless-window controls (min / max / close), drawn in NEONDECK style. Hidden outside Wails. -->
<script lang="ts">
	import { inWails, win } from './backend';
</script>

{#if inWails}
	<div class="wc">
		<button aria-label="Minimise" onclick={win.minimise}>—</button>
		<button aria-label="Maximise" onclick={win.toggleMaximise}>□</button>
		<button class="close" aria-label="Close" onclick={win.quit}>✕</button>
	</div>
{/if}

<style>
	.wc { display: flex; align-self: stretch; margin-right: calc(var(--nd-space-4) * -1); border-left: 1px solid var(--nd-line); --wails-draggable: no-drag; }
	button {
		width: 2.75rem;
		border: 0;
		background: transparent;
		color: var(--nd-text-dim);
		font-family: var(--nd-font-mono);
		font-size: var(--nd-text-sm);
		cursor: pointer;
	}
	button:hover { background: var(--nd-surface-2); color: var(--nd-text); }
	.close:hover { background: var(--nd-danger); color: var(--nd-text-on-neon); }
</style>
