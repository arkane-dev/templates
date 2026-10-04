<!-- Home: the NEONDECK website anatomy. Hero split, numbered bilingual sections, hazard footer. -->
<script lang="ts">
	import { Button, SectionHeader, Panel, MoonScroll, HazardStripe, GlitchText } from '@cyberpunk-apps/neondeck';
	import { site } from '#lib';

	const features = [
		['Fast', 'Prerendered to static HTML. No server needed.'],
		['On-brand', 'NEONDECK tokens and components only.'],
		['Bilingual', 'Hanzi with pinyin or English beside it.']
	];
</script>

<svelte:head><title>{site.title}</title></svelte:head>

<section class="hero nd-grid-bg">
	<div class="hero-text">
		<p class="nd-meta">&gt; RENDERING... 87%</p>
		<h1 class="hero-title"><GlitchText text="HELLO" /><br /><span class="nd-neon">WORLD</span></h1>
		<p class="tagline nd-neon-2">{site.tagline}</p>
		<p class="lede">Replace this with what the project does, in two short sentences.</p>
		<div class="row">
			<Button size="lg" arrow href="/about" sub="read more">Get started</Button>
			<Button size="lg" variant="ghost" href="#features">Features</Button>
		</div>
	</div>
	<div class="hero-visual nd-brackets nd-fog">
		<span class="scn nd-meta">/01</span>
		<MoonScroll text={site.hero.zh} label={site.hero.meaning} caption={site.hero.pinyin} size="clamp(15rem, 28vw, 25rem)" />
		<span class="scn2 nd-meta">//{site.name.toUpperCase()}_01</span>
	</div>
</section>

<div class="content">
	<section id="features">
		<SectionHeader index="02" zh="特点" title="Features" />
		<div class="features">
			{#each features as [t, d], i (t)}
				<div class="ft">
					<span class="nd-index">0{i + 1}</span>
					<h4>{t}</h4>
					<p>{d}</p>
				</div>
			{/each}
		</div>
	</section>

	<section>
		<SectionHeader index="03" zh="内容" title="Content" meta="panels · 3 across" />
		<div class="grid">
			{#each ['Alpha', 'Beta', 'Gamma'] as name, i (name)}
				<Panel title={name} index={`0${i + 4}`} meta="DRAFT">
					<p>Panel body. Swap in real content.</p>
					<Button size="sm" variant="outline" arrow>Open</Button>
				</Panel>
			{/each}
		</div>
	</section>
</div>

<HazardStripe text="> access granted" />

<style>
	.hero { display: grid; grid-template-columns: minmax(0, 1.1fr) minmax(0, 1fr); min-height: 70vh; border-bottom: 1px solid var(--nd-line); }
	.hero-text { padding: var(--nd-space-16) var(--nd-gutter); border-right: 1px solid var(--nd-line); background: linear-gradient(90deg, var(--nd-bg) 60%, transparent); }
	.hero-title { font-size: var(--nd-text-hero); line-height: 0.92; margin-bottom: var(--nd-space-6); }
	.tagline { font-family: var(--nd-font-mono); }
	.lede { color: var(--nd-text-dim); max-width: 46ch; }
	.row { display: flex; flex-wrap: wrap; gap: var(--nd-space-4); align-items: center; }
	.hero-visual { --nd-bracket-color: var(--nd-accent); position: relative; display: grid; place-items: center; margin: var(--nd-space-8); overflow: hidden; }
	.scn { position: absolute; top: var(--nd-space-4); left: var(--nd-space-4); color: var(--nd-accent); }
	.scn2 { position: absolute; bottom: var(--nd-space-4); right: var(--nd-space-4); }
	.content { max-width: var(--nd-content-max); margin: 0 auto; padding: var(--nd-space-16) var(--nd-gutter); display: grid; grid-template-columns: minmax(0, 1fr); gap: var(--nd-space-24); }
	.features { display: grid; grid-template-columns: repeat(auto-fit, minmax(14rem, 1fr)); border-block: 1px solid var(--nd-line); }
	.ft { padding: var(--nd-space-6) var(--nd-space-5); border-right: 1px solid var(--nd-line); }
	.ft:last-child { border-right: 0; }
	.ft h4 { margin: var(--nd-space-2) 0; color: var(--nd-accent-2); font-family: var(--nd-font-ui); letter-spacing: 0.1em; }
	.ft p { color: var(--nd-text-dim); font-size: var(--nd-text-sm); margin: 0; }
	.grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: var(--nd-space-5); }
	@media (max-width: 960px) {
		.hero { grid-template-columns: 1fr; }
		.hero-visual { min-height: 30rem; }
		.hero-visual :global(.nd-moonscroll) { --moon: min(16rem, 62vw) !important; } /* the stage is 1.25× the moon: 62vw keeps it inside the panel down to 320px */
		.grid { grid-template-columns: 1fr; }
	}
</style>
