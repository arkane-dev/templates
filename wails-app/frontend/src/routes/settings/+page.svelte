<!-- Settings: persisted by Go to the OS config dir. Accent, glow and scanlines apply live. -->
<script lang="ts">
	import { SectionHeader, Panel, Button, Tag } from '@cyberpunk-apps/neondeck';
	import { prefs, savePrefs } from '#lib';

	const accents = ['magenta', 'cyan', 'yellow', 'red', 'violet', 'jade', 'gold'];
	let saved = $state(false);

	async function save() {
		await savePrefs();
		saved = true;
		setTimeout(() => (saved = false), 1500);
	}
</script>

<div class="content">
	<SectionHeader index="01" zh="设置" title="Settings" meta={prefs.path} />
	{#if prefs.value}
		<Panel title="Appearance" index="02">
			<fieldset>
				<legend class="nd-label">Accent</legend>
				<div class="accents">
					{#each accents as a (a)}
						<label class="swatch" data-nd-accent={a}>
							<input type="radio" name="accent" value={a} bind:group={prefs.value.accent} onchange={save} />
							<span class="chip"></span>{a}
						</label>
					{/each}
				</div>
			</fieldset>

			<label class="field">
				<span class="nd-label">Glow · {prefs.value.glowSize.toFixed(1)}×</span>
				<input type="range" min="0" max="1.5" step="0.1" bind:value={prefs.value.glowSize} onchange={save} />
			</label>

			<label class="check">
				<input type="checkbox" bind:checked={prefs.value.scanlines} onchange={save} />
				<span class="nd-label">Scanlines overlay</span>
			</label>

			<div class="row">
				<Button variant="outline" onclick={save}>Save now</Button>
				{#if saved}<Tag tone="success" dot>saved</Tag>{/if}
			</div>
		</Panel>
	{:else}
		<p class="nd-meta">&gt; loading settings<span class="nd-cursor"></span></p>
	{/if}
</div>

<style>
	.content { padding: var(--nd-space-10) var(--nd-gutter); max-width: 48rem; }
	fieldset { border: 0; margin: 0 0 var(--nd-space-6); padding: 0; }
	.accents { display: flex; flex-wrap: wrap; gap: var(--nd-space-3); margin-top: var(--nd-space-2); }
	.swatch { display: inline-flex; align-items: center; gap: var(--nd-space-2); padding: var(--nd-space-1) var(--nd-space-3); border: 1px solid var(--nd-line-strong); cursor: pointer; font-family: var(--nd-font-mono); font-size: var(--nd-text-xs); text-transform: uppercase; }
	.swatch:has(input:checked) { border-color: var(--nd-accent); color: var(--nd-accent); }
	.swatch:has(input:focus-visible) { outline: 2px solid var(--nd-focus); outline-offset: 2px; }
	.swatch input { position: absolute; opacity: 0; pointer-events: none; }
	.chip { width: 0.75rem; height: 0.75rem; background: var(--nd-accent); }
	.field { display: grid; gap: var(--nd-space-2); margin-bottom: var(--nd-space-6); max-width: 20rem; }
	input[type='range'] { accent-color: var(--nd-accent); }
	.check { display: flex; align-items: center; gap: var(--nd-space-3); margin-bottom: var(--nd-space-6); }
	.check input { accent-color: var(--nd-accent); width: 1rem; height: 1rem; }
	.row { display: flex; align-items: center; gap: var(--nd-space-3); }
</style>
