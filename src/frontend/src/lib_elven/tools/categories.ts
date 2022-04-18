import NetworkCategory from "$lib_elven/network/network_category";
import type { Items } from "$lib_elven/types";
import type { Category } from "$lib_elven/types/articles/categories";

export class ToolsCategories {

    /** get category counter by name */
    public static async getCounterByName(token: string, name: string): Promise<string> {
        const networkCategory = new NetworkCategory(token);
        const resp = await networkCategory.getAll();
        if (resp.status !== 200) {
            throw Error(resp.statusText)
        }

        const cats: Items<Category> = await resp.json();
        for (const counter in cats.data) {
            const cat = cats.data[counter];
            if (cat.name === name) {
                return counter
            }
        }
        throw Error("not found")
    }
}