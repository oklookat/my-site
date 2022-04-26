import { goto } from "$app/navigation";
import { getRecordLength } from "$lib_elven/tools";
import type { Params } from "$lib_elven/tools/params";

import type { Items } from "$lib_elven/types";

interface RPH_DataGetter {
    getAll(params: Record<string | number, any>): Promise<Response>
}

/** RouteParamHandler param change event */
export type RPH_Event = {
    name: string,
    val: string | boolean | number
}

export type RPH_Data<T> = {
    items: Items<T>,
    params: Params<T>
}

/** fetch data / set searchparams when you change params */
export async function HandleRouteParam<T>(fetchData: RPH_DataGetter, event: RPH_Event, data: RPH_Data<T>): Promise<RPH_Data<T>> {
    const isPageParamChanged = event.name.toUpperCase() === 'PAGE';

    const setParam = (name: string, val: any) => {
        // @ts-ignore
        data.params.setParam(name, val)
    };

    const setPage = (newPage: number) => {
        if (newPage > data.items.meta.total_pages) {
            newPage = 1;
        }
        if (isPageParamChanged) {
            event.val = newPage;
        }
        setParam('page', newPage);
    };

    // correct page if invalid
    if (data.params.getParam('page') > data.items.meta.total_pages) {
        setPage(data.items.meta.total_pages || 1)
    }

    if (isPageParamChanged) {
        let newPage = event.val;
        const isPageInvalid = newPage < 1 || newPage > data.items.meta.total_pages;
        if (isPageInvalid) {
            newPage = 1;
        }

        // @ts-ignore
        setPage(newPage);
    } else {
        // other param changed,
        // reset page
        setPage(1);
        setParam(event.name, event.val);
    }

    const resp = await fetchData.getAll(data.params.toObject());
    if (!resp.ok) {
        throw Error(resp.statusText);
    }

    const result = await resp.json();
    data.items = result

    await goto(`?${data.params.toSearchparams().toString()}`, { replaceState: true })

    return data
}

/** refresh page when you change something */
export async function Refresh<T>(fetchData: RPH_DataGetter, data: RPH_Data<T>, useGoto = true): Promise<RPH_Data<T>> {
    const refreshByNetwork = async (page: number) => {
        data.params.setParam('page', page)

        const resp = await fetchData.getAll(data.params.toObject());
        if (!resp.ok) {
            throw Error(resp.statusText);
        }

        const result = await resp.json();
        data.items = result
        if (useGoto) {
            await goto(`?${data.params.toSearchparams().toString()}`, { replaceState: true, keepfocus: true });
        }
    };

    let newPage = data.params.getParam('page');
    let totalPages = data.items.meta.total_pages

    // if page invalid, set page to 1 and get data
    let isPageInvalid = newPage < 1 || newPage > totalPages;
    if (isPageInvalid || newPage === 1) {
        if (newPage > totalPages) {
            newPage = totalPages
        } else {
            newPage = 1;
        }
        await refreshByNetwork(newPage);
        return data;
    }

    let dataLength = getRecordLength(data.items.data)

    if (dataLength > 0 && dataLength < data.items.meta.per_page && newPage < data.items.meta.total_pages) {
        await refreshByNetwork(newPage)
        dataLength = getRecordLength(data.items.data)
        if (dataLength >= data.items.meta.per_page || newPage === data.items.meta.total_pages) {
            return data
        }
    }

    while (newPage > 1 && dataLength < 1) {
        newPage--;
        try {
            await refreshByNetwork(newPage);
            dataLength = getRecordLength(data.items.data);
        } catch (err) {
            break;
        }
    }

    return data
}