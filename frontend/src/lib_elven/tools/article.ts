import NetworkArticle from '$elven/network/article';
import type { Operation } from '$elven/tools/jsonpatch';
import type { RAW } from '$elven/types/article';

/** editable article */
export class Editable implements RAW {
	public id = '';
	public user_id = '';
	private _cover_id: string | undefined = '';
	private _is_published = false;
	private _title = '';
	private _content = '';
	public published_at = '';
	public updated_at = '';
	public cover_path: string | undefined = '';
	public cover_extension: string | undefined = '';

	///////////
	/** not apply PATCH/json-patch? */
	private avoidPatch = false;
	/** avoid creating new article */
	private avoidCreating = false;
	/** avoid too often saves */
	private patchThrottle: NodeJS.Timeout;
	// @ts-ignore
	private lastPatch: Record<keyof RAW, Operation> = {};
	public lastSaved: Date | undefined;
	/** on article created/saved hook */
	public onSaved: (resp: Response) => void

	/** is new article? */
	private get isNew() {
		return !this.id;
	}

	constructor(raw?: RAW) {
		if (!raw) {
			return;
		}
		this.fromRAW(raw);
	}

	public get is_published() {
		return this._is_published;
	}

	public set is_published(val: boolean) {
		let op: Operation = { op: 'replace', path: '/is_published', value: val };
		this._is_published = val;
		this.makePatchOrCreateNew('is_published', op);
	}

	public get cover_id() {
		return this._cover_id;
	}

	public set cover_id(val: string | undefined) {
		let op: Operation;
		if (!val) {
			op = { op: 'remove', path: '/cover_id' };
		} else {
			op = { op: 'replace', path: '/cover_id', value: val };
		}
		this._cover_id = val;
		this.makePatchOrCreateNew('cover_id', op);
	}

	public get title() {
		return this._title;
	}

	public set title(val: string) {
		let op: Operation;
		const prevValue = this.title;

		if (typeof val !== 'string' && typeof val !== 'undefined') {
			val = '';
			op = { op: 'remove', path: '/title' };
		} else if (val && val.length > 124) {
			throw Error('too big title');
		} else if (!prevValue) {
			op = { op: 'add', path: '/title', value: val };
		} else {
			op = { op: 'replace', path: '/title', value: val };
		}

		this._title = val;
		this.makePatchOrCreateNew('title', op);
	}

	public get content() {
		return this._content;
	}

	public set content(val: string) {
		const isEmpty = typeof val !== 'string' || val.length < 1;
		if (isEmpty) {
			val = '';
		} else {
			const isValid = val && val.length < 816000;
			if (!isValid) {
				return;
			}
		}
		const op: Operation = { op: 'replace', path: '/content', value: val };
		this._content = val;
		this.makePatchOrCreateNew('content', op);
	}

	/** apply patch to existing article (PATCH) / or create new article */
	private async makePatchOrCreateNew(toKey: keyof RAW, op: Operation): Promise<void> {
		if (this.isNew) {
			this.createNew();
			return;
		}

		if (this.avoidPatch) {
			return;
		}

		if (this.patchThrottle) {
			clearTimeout(this.patchThrottle);
		}

		this.lastPatch[toKey] = op;

		const handler = async () => {
			window.$progress?.startBasic();
			const patchesArray: Operation[] = [];
			for (const prop in this.lastPatch) {
				// @ts-ignore
				patchesArray.push(this.lastPatch[prop]);
			}
			const resp = await NetworkArticle.update(this.id, patchesArray);
			delete this.lastPatch[toKey];
			if (resp.ok) {
				this.lastSaved = new Date();
			}
			window.$progress?.finishBasic();
			if(this.onSaved) {
				this.onSaved(resp)
			}
		};

		this.patchThrottle = setTimeout(handler, 1000);
	}

	public toRAW(): RAW {
		return {
			id: this.id,
			user_id: this.user_id,
			cover_id: this.cover_id,
			is_published: this.is_published,
			title: this.title,
			content: this.content,
			published_at: this.published_at,
			updated_at: this.updated_at,
			cover_path: this.cover_path,
			cover_extension: this.cover_extension
		};
	}

	public fromRAW(raw: RAW) {
		this.avoidPatch = true;
		Object.assign(this, raw);
		this.avoidPatch = false;
	}

	/** create new article */
	private async createNew() {
		if (!this.isNew) {
			return;
		}

		if (this.avoidCreating) {
			return;
		}

		this.avoidCreating = true;
		let resp: Response;
		try {
			resp = await NetworkArticle.create(this.toRAW());
			if (!resp.ok) {
				this.avoidCreating = false;
				return;
			}
		} catch (err) {
			this.avoidCreating = false;
			return;
		}
		const jsond = (await resp.json()) as RAW;
		this.fromRAW(jsond);
		if(this.onSaved) {
			this.onSaved(resp)
		}
	}
}
