// Typed access to the Go backend. Inside Wails, calls go through the generated bindings.
// In a plain browser (`npm run dev` without Wails) they fall back to stubs, so the UI can be
// built and screenshot-tested without the desktop shell.
import * as Go from './wailsjs/go/main/App';
import * as Rt from './wailsjs/runtime/runtime';
import type { main } from './wailsjs/go/models';

export type AppInfo = main.AppInfo;
export type Settings = main.Settings;
export type UpdateInfo = main.UpdateInfo;

export const inWails = typeof window !== 'undefined' && 'go' in window;

const defaults: Settings = { accent: 'magenta', glowSize: 1, scanlines: false, autoUpdate: true, lastUpdateCheck: '' } as Settings;
const noUpdate = { current: 'dev', latest: 'dev', available: false, canInstall: false, why: '', notes: '', pageUrl: '' } as UpdateInfo;
const browserInfo = {
	name: 'wails-app', version: 'dev', goVersion: '—', os: 'browser', arch: '—', hostname: 'localhost', cpus: 0
} as AppInfo;

export const backend = {
	appInfo: (): Promise<AppInfo> => (inWails ? Go.AppInfo() : Promise.resolve(browserInfo)),
	ping: (msg: string): Promise<string> =>
		inWails ? Go.Ping(msg) : Promise.resolve(`> ACK "${msg}" (browser stub, no Go backend)`),
	getSettings: (): Promise<Settings> => (inWails ? Go.GetSettings() : Promise.resolve({ ...defaults })),
	saveSettings: (s: Settings): Promise<void> => (inWails ? Go.SaveSettings(s) : Promise.resolve()),
	settingsPath: (): Promise<string> => (inWails ? Go.SettingsPath() : Promise.resolve('(browser: not saved)')),
	checkForUpdate: (): Promise<UpdateInfo> => (inWails ? Go.CheckForUpdate() : Promise.resolve({ ...noUpdate })),
	installUpdate: (): Promise<void> => (inWails ? Go.InstallUpdate() : Promise.reject(new Error('Updates need the desktop app.'))),
	restartForUpdate: (): Promise<void> => (inWails ? Go.RestartForUpdate() : Promise.resolve()),
	openURL: (url: string): void => {
		if (inWails) Go.OpenURL(url);
		else window.open(url, '_blank', 'noopener');
	},
	/** Go raises "update:available" after its once-a-day startup check finds a release. */
	onUpdateAvailable: (fn: (info: UpdateInfo) => void): (() => void) => (inWails ? Rt.EventsOn('update:available', fn) : () => {})
};

export const win = {
	minimise: () => inWails && Rt.WindowMinimise(),
	toggleMaximise: () => inWails && Rt.WindowToggleMaximise(),
	quit: () => inWails && Rt.Quit()
};
