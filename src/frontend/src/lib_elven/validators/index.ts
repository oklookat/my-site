export default class Validator {

    /** is device with touchscreen & its default input? */
    public static isTouchDevice(): boolean {
        return matchMedia('(hover: none)').matches;
    }

    public static isAdminPanelPage(url: URL): boolean {
        if (!(url instanceof URL)) {
            return false
        }
        return url.pathname.startsWith("/elven")
    }

    public static isAdminPanelLoginPage(url: URL): boolean {
        if (!(url instanceof URL)) {
            return false
        }
        return url.pathname.startsWith("/elven/login")
    }

    public static isAdminPanelLogoutPage(url: URL): boolean {
        if (!(url instanceof URL)) {
            return false
        }
        return url.pathname.startsWith("/elven/logout")
    }

}