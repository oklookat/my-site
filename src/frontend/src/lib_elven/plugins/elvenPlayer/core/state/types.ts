/** available hooks */
export enum HookName {
	onChange = 'onChange'
}

/** remove hook */
export interface HookRemove {
	(): void;
}

/** user hook */
export interface Hook<T> {
	(value?: T): void;
}

/** hooks list in specific hook type (internal) */
export type HooksList<T> = {
	/** counter of hooks */
	count: number;
	/** hooks */
	items: {
		/** one hook */
		[id: number]: Hook<T>;
	};
};

/** all hooks (internal) */
export type Hooks<T> = {
	[name in HookName]: HooksList<T>;
};
