// Typed access to the Go backend. Inside Wails, calls go through the generated bindings.
// In a plain browser (`npm run dev` without Wails) they fall back to stubs, so the UI can be
// built and screenshot-tested without the desktop shell.
import * as Go from './wailsjs/go/main/App';
import * as Rt from './wailsjs/runtime/runtime';
import type { main } from './wailsjs/go/models';

export type AppInfo = main.AppInfo;
export type Settings = main.Settings;

export const inWails = typeof window !== 'undefined' && 'go' in window;

const defaults: Settings = { accent: 'magenta', glowSize: 1, scanlines: false } as Settings;
const browserInfo = {
	name: 'wails-app', version: 'dev', goVersion: '—', os: 'browser', arch: '—', hostname: 'localhost', cpus: 0
} as AppInfo;

export const backend = {
	appInfo: (): Promise<AppInfo> => (inWails ? Go.AppInfo() : Promise.resolve(browserInfo)),
	ping: (msg: string): Promise<string> =>
		inWails ? Go.Ping(msg) : Promise.resolve(`> ACK "${msg}" (browser stub, no Go backend)`),
	getSettings: (): Promise<Settings> => (inWails ? Go.GetSettings() : Promise.resolve({ ...defaults })),
	saveSettings: (s: Settings): Promise<void> => (inWails ? Go.SaveSettings(s) : Promise.resolve()),
	settingsPath: (): Promise<string> => (inWails ? Go.SettingsPath() : Promise.resolve('(browser: not saved)'))
};

export const win = {
	minimise: () => inWails && Rt.WindowMinimise(),
	toggleMaximise: () => inWails && Rt.WindowToggleMaximise(),
	quit: () => inWails && Rt.Quit()
};
