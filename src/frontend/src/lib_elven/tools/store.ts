import type { File } from "$lib_elven/types/files";
import { writable, type Writable } from "svelte/store";

type store_file = {
    withSelectOption: Writable<boolean>
    selected: Writable<File | null>
    uploadedExists: Writable<File | null>
}

/** app store */
export default class Store {
    
    /** files store */
    public static files: store_file = {
        /** add 'select' option */
        withSelectOption: writable(false),

        /** when 'select' option clicked */
        selected: writable(null),

        /** when uploaded file exists */
        uploadedExists: writable(null)
    }
}