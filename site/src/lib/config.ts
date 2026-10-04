// Project identity. new-project.sh rewrites the `name`, `brand` and `title` lines; edit the rest by hand.
export const site = {
	name: 'template', // kebab-case project folder name
	brand: 'TEMPLATE_', // top-bar brand block; NEONDECK brands end in "_"
	title: 'Template', // <title> suffix
	tagline: 'THE CITY NEVER SLEEPS.',
	hero: { zh: '霓虹都市', pinyin: 'NI HONG DU SHI', meaning: 'Neon City' },
	nav: [
		{ label: 'Home', href: '/' },
		{ label: 'About', href: '/about' }
	],
	version: '0.1.0'
};
