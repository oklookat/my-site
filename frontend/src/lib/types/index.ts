import type { ElvenPlayer } from '$lib/plugins/elvenPlayer/types';
import type { ElvenNotify } from '$lib/plugins/elvenNotify/types';
import type { ElvenProgress } from '$lib/plugins/elvenProgress/types';
import type { ElvenConfirm } from '../plugins/elvenChoose/types';

export type Page = number;
export type Counter = number | string;

/** response with multiple entities */
export type Items<T> = {
	meta: Meta;
	data: Record<Counter, T>;
};

/** information about requested data */
export type Meta = {
	per_page: Page;
	total_pages: Page;
	current_page: Page;
};

declare global {
	interface Window {
		// plugins
		$player: ElvenPlayer | undefined;
		$notify: ElvenNotify | undefined;
		$progress: ElvenProgress | undefined;
		$confirm: ElvenConfirm | undefined;
	}
}
