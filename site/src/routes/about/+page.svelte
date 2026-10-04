<!-- About: a long-read page using the paper (editorial) surface. -->
<script lang="ts">
	import { SectionHeader, HanziMark, Seal, Button } from '@cyberpunk-apps/neondeck';
	import { site } from '#lib';
</script>

<svelte:head><title>About · {site.title}</title></svelte:head>

<div class="content">
	<SectionHeader index="01" zh="关于" title="About" meta="paper surface" />
	<article class="poster nd-paper">
		<div class="band"><HanziMark text="关于" label="About" tone="sun" size="clamp(3rem, 6vw, 5rem)" /></div>
		<div class="body">
			<div class="seal"><Seal text="霓虹甲板" label="Neon Deck seal" size="clamp(5rem, 9vw, 7rem)" tilt={-2} /></div>
			<p class="nd-label">Long read ─── 001</p>
			<h2 class="title">About</h2>
			<p>Paper blocks are for long reads: cream stock, black ink, big hanzi, one cinnabar seal. Replace this copy.</p>
			<Button arrow href="/">Back home</Button>
		</div>
	</article>
</div>

<style>
	.content { max-width: var(--nd-content-max); margin: 0 auto; padding: var(--nd-space-16) var(--nd-gutter); }
	.poster { display: grid; grid-template-columns: auto 1fr; min-height: 24rem; overflow: hidden; }
	.band { display: grid; place-items: center; padding: var(--nd-space-6); background: var(--nd-ink); }
	.body { position: relative; display: grid; align-content: start; justify-items: start; padding: var(--nd-space-12) var(--nd-space-10); }
	.body > p { max-width: 52ch; }
	.seal { position: absolute; right: var(--nd-space-8); top: var(--nd-space-8); }
	.title { font-size: var(--nd-text-5xl); color: var(--nd-ink); margin: var(--nd-space-2) 0 var(--nd-space-4); }
	@media (max-width: 720px) {
		.poster { grid-template-columns: 1fr; }
		.band :global(.nd-hanzi) { writing-mode: horizontal-tb; }
		.body { padding: var(--nd-space-5) var(--nd-space-5) var(--nd-space-8); }
		.seal { position: static; display: flex; justify-content: flex-end; justify-self: stretch; margin-bottom: var(--nd-space-3); /* in flow on phones: never covers the text */ }
	}
</style>
