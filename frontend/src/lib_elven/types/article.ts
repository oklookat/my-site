/** request param */
export enum By {
	created = 'created',
	published = 'published',
	updated = 'updated'
}

/** request params */
export type Params = {
	page?: Page;

	drafts?: boolean;

	newest?: boolean;

	by?: By;

	/** search by title */
	title?: string;
};

/** raw article */
export type RAW = {
	id: string;
	user_id: string;
	cover_id: string | undefined;
	is_published: boolean;
	title: string;
	content: string;
	published_at: string;
	updated_at: string;
	/** available only when GET (JOIN) */
	cover_path?: string;
	/** available only when GET (JOIN) */
	cover_extension?: string;
};
