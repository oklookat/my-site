// from https://github.com/Starcounter-Jack/JSON-Patch
export type Operation =
	| AddOperation<any>
	| RemoveOperation
	| ReplaceOperation<any>
	| MoveOperation
	| CopyOperation
	| TestOperation<any>
	| GetOperation<any>;
export interface BaseOperation {
	path: string;
}
export interface AddOperation<T> extends BaseOperation {
	op: 'add';
	value: T;
}
export interface RemoveOperation extends BaseOperation {
	op: 'remove';
}
export interface ReplaceOperation<T> extends BaseOperation {
	op: 'replace';
	value: T;
}
export interface MoveOperation extends BaseOperation {
	op: 'move';
	from: string;
}
export interface CopyOperation extends BaseOperation {
	op: 'copy';
	from: string;
}
export interface TestOperation<T> extends BaseOperation {
	op: 'test';
	value: T;
}
export interface GetOperation<T> extends BaseOperation {
	op: '_get';
	value: T;
}
