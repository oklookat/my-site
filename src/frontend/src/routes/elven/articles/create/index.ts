import NetworkArticle from '$lib_elven/network/network_article';
import type { Article } from '$lib_elven/types/articles';
import type { RequestHandlerOutput } from '@sveltejs/kit';
import type { RequestEvent } from '@sveltejs/kit/types/private';

export async function get(event: RequestEvent): Promise<RequestHandlerOutput> {
	/** creating / editing this article */
	let article: Article = {
		title: '',
		content: ''
	};

	const params = event.url.searchParams;
	const isEditMode = params.has('id');
	if (!isEditMode) {
		return {
			body: { article: article }
		};
	}

	try {
		const networkArticle = new NetworkArticle(event.locals.user.token || '');
		const resp = await networkArticle.get(params.get('id')!);
		if(resp.ok) {
			article = await resp.json()
		}
	} catch (err) {}
	return {
		body: { article: article }
	};
}
