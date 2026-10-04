// Shared, reactive user preferences. Loaded once from Go, applied to the document, saved on change.
import { backend, type Settings } from './backend';

export const prefs = $state<{ value: Settings | null; path: string }>({ value: null, path: '' });

export function apply(s: Settings) {
	const root = document.documentElement;
	root.dataset.ndAccent = s.accent;
	root.style.setProperty('--nd-glow-size', String(s.glowSize));
}

export async function loadPrefs() {
	prefs.value = await backend.getSettings();
	prefs.path = await backend.settingsPath();
	apply(prefs.value);
}

export async function savePrefs() {
	if (!prefs.value) return;
	apply(prefs.value);
	await backend.saveSettings(prefs.value);
}
