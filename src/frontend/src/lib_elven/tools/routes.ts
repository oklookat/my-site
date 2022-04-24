import { goto } from "$app/navigation";
import { getRecordLength, setSearchParam, stringToNormal, stringToNumber } from "$lib_elven/tools";
import type { Items } from "$lib_elven/types";

type Params = { page?: number }

interface RPH_DataGetter {
    getAll(params: any): Promise<Response>
}

/** RouteParamHandler param change event */
export type RPH_Event = {
    name: string,
    val: string | boolean | number
}

export type RPH_Data<T> = {
    searchparams: URLSearchParams,
    items: Items<T>,
    params: Params
}

/** fetch data / set searchparams when you change params */
export async function HandleRouteParam<T>(fetchData: RPH_DataGetter, event: RPH_Event, data: RPH_Data<T>): Promise<RPH_Data<T>> {
    // normalize value type
    event.val = stringToNormal(event.val);

    const isPageParamChanged = event.name.toUpperCase() === 'PAGE';

    const setParam = (name: string, val: any) => {
        data.params[name] = val;
        setSearchParam(data.searchparams, name, val);
    };

    const setPage = (newPage: number) => {
        if (newPage < 1 || newPage > data.items.meta.total_pages) {
            newPage = 1;
        }
        if (isPageParamChanged) {
            event.val = newPage;
        }
        setParam('page', newPage);
    };

    if (isPageParamChanged) {
        let newPage = event.val;
        const isPageInvalid = newPage < 1 || newPage > data.items.meta?.total_pages;
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

    const resp = await fetchData.getAll(data.params);
    if (!resp.ok) {
        throw Error(resp.statusText);
    }

    const result = await resp.json();
    data.items = result

    await goto(`?${data.searchparams.toString()}`, { replaceState: true })

    return data
}

/** refresh page when you change something */
export async function Refresh<T>(fetchData: RPH_DataGetter, data: RPH_Data<T>, useGoto = true): Promise<RPH_Data<T>> {
    const refreshByNetwork = async (page: number) => {
        data.params.page = page;

        if (useGoto) {
            setSearchParam(data.searchparams, 'page', page);
        }

        const resp = await fetchData.getAll(data.params);
        if (resp.ok) {
            const result = await resp.json();
            data.items = result
            if (useGoto) {
                await goto(`?${data.searchparams.toString()}`, { replaceState: true, keepfocus: true });
            }
            return;
        }

        throw Error(resp.statusText);
    };

    let newPage = data.params.page;

    // if page invalid, set page to 1 and get data with goto
    let isPageInvalid = newPage < 1 || newPage > data.items?.meta?.total_pages;
    if (isPageInvalid || newPage === 1) {
        newPage = 1;
        await refreshByNetwork(newPage);
        return data;
    }

    let dataLength = getRecordLength(data.items.data)

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