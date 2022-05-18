import type { Page } from '$lib_elven/types';
import type { FileTypeSelector } from '../../tools/extension';

/** request param */
export enum Start {
	newest = 'newest',
	oldest = 'oldest'
}

/** request param */
export enum By {
	created = 'created'
}

/** request params */
export type Params = {
	page?: Page;
	start?: Start;
	by?: By;

	/** find files with this extensions. Format like: 'jpg,gif,png' */
	extensions?: string;

	/** search by filename */
	filename?: string;
};

/** file */
export type File = {
	id: string;
	user_id: string;
	hash: string;
	path: string;
	name: string;
	original_name: string;
	extension: string;
	size: number;
	created_at: string;
	updated_at: string;

	///////////////////// not in model

	/** short original name like: 'longlongnamehellowo...' */
	original_name_short?: string;

	/** like path, but with absolute URL to static */
	pathConverted?: URL;

	/** readable size like: '1.1 MB' */
	sizeConverted?: string;

	/** readable created data like: '11 seconds ago' */
	createdAtConverted?: string;

	/** readable extension */
	extensionsSelector?: FileTypeSelector;
};
