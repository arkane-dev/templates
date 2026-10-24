<!-- Shown when a new release is out: from the startup check or a manual one. -->
<script lang="ts">
	import { Button } from '@cyberpunk-apps/neondeck';
	import { backend, update, installUpdate } from '#lib';

	const show = $derived(!update.dismissed && ['available', 'installing', 'ready'].includes(update.phase) && update.info);
</script>

{#if show && update.info}
	<aside class="banner" role="status" aria-live="polite">
		{#if update.phase === 'ready'}
			<p><strong>Version {update.info.latest} is installed.</strong> Restart to use it.</p>
			<div class="actions">
				<Button size="sm" onclick={() => backend.restartForUpdate()}>Restart now</Button>
				<Button size="sm" variant="ghost" onclick={() => (update.dismissed = true)}>Later</Button>
			</div>
		{:else}
			<p><strong>Version {update.info.latest} is out.</strong> You have {update.info.current}.</p>
			<div class="actions">
				{#if update.info.canInstall}
					<Button size="sm" onclick={installUpdate} disabled={update.phase === 'installing'}>
						{update.phase === 'installing' ? 'Updating…' : 'Update'}
					</Button>
				{:else}
					<Button size="sm" arrow onclick={() => backend.openURL(update.info!.pageUrl)}>Download</Button>
				{/if}
				<Button size="sm" variant="outline" onclick={() => backend.openURL(update.info!.pageUrl)}>What's new</Button>
				<Button size="sm" variant="ghost" onclick={() => (update.dismissed = true)}>Later</Button>
			</div>
		{/if}
	</aside>
{/if}

<style>
	.banner {
		position: fixed;
		right: var(--nd-space-6);
		bottom: calc(var(--nd-statusbar-h) + var(--nd-space-4));
		z-index: var(--nd-z-toast);
		max-width: 26rem;
		padding: var(--nd-space-4) var(--nd-space-5);
		background: var(--nd-surface-3);
		border: 1px solid var(--nd-line-strong);
		border-left: 2px solid var(--nd-accent);
		display: grid;
		gap: var(--nd-space-3);
	}
	p { margin: 0; }
	.actions { display: flex; flex-wrap: wrap; gap: var(--nd-space-2); }
</style>
