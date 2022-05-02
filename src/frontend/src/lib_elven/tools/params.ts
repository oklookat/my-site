import { getRecordLength, searchParamsByObject, stringToNormal } from '$lib/tools';
import type { Article } from '$lib_elven/types/articles';
import {
	By as FilesBy,
	Start,
	type File,
	type Params as FileParamsT
} from '$lib_elven/types/files';
import { type Params as ArticleParamsT, By as ArticlesBy } from '$lib_elven/types/articles';
import type { Items } from '$lib_elven/types';
import { goto } from '$app/navigation';
import { browser } from '$app/env';

export type ParamTypeKey<T> = keyof ParamType<T>;

export type ParamTypeValue<T> = ParamType<T>[ParamTypeKey<T>];

export type ParamType<T extends {}> = T extends File
	? FileParamsT
	: T extends Article
	? ArticleParamsT
	: never;

/** RouteParamHandler / param change event */
export type RPH_Event<T> = {
	name: ParamTypeKey<T>;
	val: ParamTypeValue<T>;
};

/** RouteParamHandler / data */
export type RPH_Data<T> = {
	items: Items<T>;
	params: Params<T>;
};

/** RouteParamHandler / get data from API */
export interface RPH_DataGetter<T> {
	getAll(params: ParamType<T>): Promise<Response>;
}

/** params for what? */
type ParamsFor = 'file' | 'article';

/** manage request params (file, article) */
export class Params<T> {
	private paramsFor: ParamsFor;

	/** request params object */
	private self: ParamType<T>;

	/** like request params object, but URLSearchParams */
	private searchparams: URLSearchParams;

	/** create new request params */
	constructor(paramsFor: ParamsFor, init?: URLSearchParams) {
		if (paramsFor !== 'article' && paramsFor !== 'file') {
			throw Error("wrong 'paramsFor' param");
		}

		this.paramsFor = paramsFor;
		this.self = this.getDefault();
		this.searchparams = searchParamsByObject(this.self);
		this.import(this.searchparams);
		this.import(init);
	}

	/** import values from searchparams */
	public import(params?: URLSearchParams) {
		if (!params || !(params instanceof URLSearchParams)) {
			return;
		}

		params.forEach((value, key) => {
			// @ts-ignore
			this.setParam(key, value);
		});
	}

	/** set request param */
	public setParam(name: ParamTypeKey<T>, val: ParamTypeValue<T>) {
		let normalized = stringToNormal(val);

		// @ts-ignore
		name = String(name);

		// @ts-ignore
		name = name.toLowerCase();
		let isNeedToDelete = false;
		if (typeof normalized === 'string') {
			const normalized_BIG = normalized.toUpperCase();
			isNeedToDelete =
				normalized_BIG === '' || normalized_BIG === 'NULL' || normalized_BIG === 'UNDEFINED';
		} else {
			isNeedToDelete = normalized === undefined || normalized === null;
		}

		/** is param exists in defaults? */
		const isExists = name in this.getDefault();

		if (isNeedToDelete || !isExists) {
			// @ts-ignore
			this.searchparams.delete(name);
			delete this.self[name];
			return;
		}

		// @ts-ignore
		if (name === 'page') {
			if (normalized < 1 || isNaN(Number(normalized))) {
				// @ts-ignore
				normalized = 1;
			}
		}

		// @ts-ignore
		this.searchparams.set(name, normalized);

		// @ts-ignore
		this.self[name] = normalized;
	}

	/** get request param */
	public getParam(name: ParamTypeKey<T>): any {
		return this.self[name];
	}

	/** get request params searchparams copy */
	public toSearchparams(): URLSearchParams {
		return new URLSearchParams(this.searchparams);
	}

	/** get request params copy */
	public toObject(): ParamType<T> {
		return { ...this.self };
	}

	/** get default request params */
	private getDefault(): ParamType<T> {
		if (this.paramsFor === 'article') {
			return {
				page: 1,
				drafts: false,
				newest: true,
				by: ArticlesBy.published,
				title: undefined
			} as ParamType<T>;
		}
		return {
			page: 1,
			start: Start.newest,
			by: FilesBy.created,
			extensions: undefined,
			filename: undefined
		} as ParamType<T>;
	}
}

/** fetch data / set searchparams when you change params */
export async function HandleRouteParam<T>(event: RPH_Event<T>, data: RPH_Data<T>) {
	let windowSearch = ''
	let searchparams: URLSearchParams;
	if (browser) {
		windowSearch = window.location.search;
		searchparams = new URLSearchParams(windowSearch);
		data.params.import(searchparams);
	} else {
		searchparams = data.params.toSearchparams();
	}

	let trashParam = '';

	const isPageParamChanged = typeof event.name === 'string' && event.name.toUpperCase() === 'PAGE';
	if (!isPageParamChanged) {
		// @ts-ignore
		data.params.setParam('page', 1);
	} else {
		// if page > total_pages = set max page
		if (event.val > data.items.meta.total_pages) {
			// @ts-ignore
			event.val = data.items.meta.total_pages;
		} else if (event.val < 1) {
			// if page < 1 = set page to 1
			// @ts-ignore
			event.val = 1;
		}
		const prevPage = data.params.getParam('page');

		// if current page and event page the same
		if (prevPage === event.val) {
			// goto will not work because params not changed
			// we need to add trash param at the end like 'something changed'
			trashParam = '&hack=yes';

			// remove if exists
			if (windowSearch && windowSearch.endsWith(trashParam)) {
				trashParam = '';
			}
		}
	}

	data.params.setParam(event.name, event.val);

	const toGoto = `?${data.params.toSearchparams().toString()}${trashParam}`;
	await goto(toGoto, {
		keepfocus: true
	});
}

/** check is items exists, and if not - refresh */
export async function Refresh<T>(getData: () => Promise<RPH_Data<T>>, useGoto = true) {
	let data = await getData();

	const refeshData = async (page: number) => {
		// @ts-ignore
		data.params.setParam('page', page);

		if (useGoto) {
			// @ts-ignore
			await HandleRouteParam<T>({ name: 'page', val: page }, data);
			data = await getData();
			return;
		}
		data = await getData();
	};

	let newPage = data.params.getParam('page');
	let totalPages = data.items.meta.total_pages;

	// if page invalid, set page to 1 and get data
	let isPageInvalid = newPage < 1 || newPage > totalPages;
	if (isPageInvalid || newPage === 1) {
		if (newPage > totalPages) {
			newPage = totalPages;
		} else {
			newPage = 1;
		}
		await refeshData(newPage);
		return data;
	}

	let dataLength = getRecordLength(data.items.data);

	if (
		dataLength > 0 &&
		dataLength < data.items.meta.per_page &&
		newPage < data.items.meta.total_pages
	) {
		await refeshData(newPage);
		dataLength = getRecordLength(data.items.data);
		if (dataLength >= data.items.meta.per_page || newPage === data.items.meta.total_pages) {
			return data;
		}
	}

	while (newPage > 1 && dataLength < 1) {
		newPage--;
		try {
			await refeshData(newPage);
			dataLength = getRecordLength(data.items.data);
		} catch (err) {
			break;
		}
	}

	return data;
}
