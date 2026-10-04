// Project identity. new-project.sh rewrites the `name` and `brand` lines.
export const app = {
	name: 'wails-app', // must match appName in main.go
	brand: 'WAILS_APP_', // top-bar brand block; NEONDECK brands end in "_"
	nav: [
		{ label: 'Dashboard', zh: '控制台', href: '#/' },
		{ label: 'Settings', zh: '设置', href: '#/settings' }
	]
};
