// Auto-update state, shared by the banner in the layout and the Updates panel in Settings.
import { backend, type UpdateInfo } from './backend.js';

type Phase = 'idle' | 'checking' | 'current' | 'available' | 'installing' | 'ready' | 'error';

export const update = $state({
	phase: 'idle' as Phase,
	info: null as UpdateInfo | null,
	error: '',
	dismissed: false // "Later" hides the banner until the next start
});

function failed(e: unknown) {
	update.phase = 'error';
	update.error = e instanceof Error ? e.message : String(e);
}

export async function checkForUpdate() {
	update.phase = 'checking';
	update.error = '';
	try {
		const info = await backend.checkForUpdate();
		update.info = info;
		update.phase = info.available ? 'available' : 'current';
		update.dismissed = false;
	} catch (e) {
		failed(e);
	}
}

/** Download, verify and swap in the new version. It runs after a restart. */
export async function installUpdate() {
	update.phase = 'installing';
	update.error = '';
	try {
		await backend.installUpdate();
		update.phase = 'ready';
	} catch (e) {
		failed(e);
	}
}

export function listenForUpdates(): () => void {
	return backend.onUpdateAvailable((info) => {
		update.info = info;
		update.phase = 'available';
	});
}
