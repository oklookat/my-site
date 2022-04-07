export default class Validator {
    /** is device with touchscreen & its default input? */
    public static isTouchDevice(): boolean {
        return matchMedia('(hover: none)').matches;
    }
}