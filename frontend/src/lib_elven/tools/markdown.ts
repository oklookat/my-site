import hljs from 'highlight.js';
import '$lib/assets/highlight.scss';
import { marked } from 'marked';

/** get parser to parse markdown */
export function getParser(): (data: string) => string {
	marked.setOptions({
		renderer: new marked.Renderer(),
		highlight: function (code, lang) {
			const language = hljs.getLanguage(lang) ? lang : 'plaintext';
			return hljs.highlight(code, { language }).value;
		},
		// highlight.js css expects a top-level 'hljs' class.
		langPrefix: 'hljs language-',
		pedantic: false,
		gfm: true,
		breaks: false,
		sanitize: false,
		smartLists: true,
		smartypants: false,
		xhtml: false
	});

	return (data: string): string => {
		return marked.parse(data);
	};
}
