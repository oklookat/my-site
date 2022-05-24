import hljs from 'highlight.js';
import '$lib/assets/highlight.scss';
import { marked } from 'marked';
import type { Config } from '@oklookat/jmarkd';

/** get config for text editor */
export function getEditorConfig(container: HTMLDivElement, data?: string): Config {
	const parser = getParser()
	return {
		container: container,
		placeholder: 'Hello.',
		input: data,
		toolbar: {
			elements: {
				config: {
					preview: {
						parse: (data: string) => {
							return parser(data);
						}
					}
				}
			}
		}
	};
}

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
