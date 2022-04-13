import type { File } from "$lib_elven/types/files";
import { writable, type Writable } from "svelte/store";

type store_file = {
    withSelectOption: Writable<boolean>
    selected: Writable<File | null>
}

/** app store */
export default class Store {
    public static onUploadedFileExists: Writable<File> = writable(undefined)

    public static file: store_file = {
        withSelectOption: writable(false),
        selected: writable(null)
    }
}