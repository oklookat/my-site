import type { File } from "$lib_elven/types/files";
import { writable, type Writable } from "svelte/store";

/** app store */
export default class Store {
    public static onUploadedFileExists: Writable<File> = writable(undefined)
}