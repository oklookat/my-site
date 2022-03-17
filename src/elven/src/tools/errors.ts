import { AuthStorage } from "@/tools/storage"
import type { DuckHook } from "@oklookat/duck"


/** handle request / response errors */
export class AdapterError {

    public static handle(output: DuckHook.Output.onError): string {
        const message = this.sort(output)
        if (window.$notify) {
            window.$notify.add({ message })
        }
        return message
    }

    private static sort(output: DuckHook.Output.onError): string {
        // server send response with error
        const err = output.data
        if (err.type === "cancel") {
            if (!err.message) {
                return 'Request cancelled.'
            }
            return `Request cancelled: ${err.message}`
        }
        if (err.type === "response") {
            let statusCode: number = err.statusCode
            let apiStatusCode: number = err.body.statusCode
            // if status code exists in response body
            if (apiStatusCode) {
                statusCode = apiStatusCode
            }
            switch (statusCode) {
                /**
                 * Logout user if:
                 * 1. Bad request (can means p.2 or p.3)
                 * 2. Not authorized
                 * 3. Authorized, but not admin
                 */
                case 400:
                case 401:
                case 403:
                    // no logout if user maek mistake while changing credentials
                    const requestURL = output.config.url
                    const isChangeCredentialsRoute = requestURL.includes("users/me/change")
                    if (isChangeCredentialsRoute) {
                        return 'Wrong password.'
                    }
                    // no logout if user try to login
                    const isLoginRoute = requestURL.includes("auth/login")
                    if(isLoginRoute) {
                        return 'Wrong username or password.'
                    }
                    AuthStorage.remove()
                    return 'Something goes wrong.'
                case 404:
                    return 'Not found.'
                case 409:
                    return 'Already exists.'
                case 413:
                    return 'Upload size too big.'
                default:
                    const str = statusCode.toString()
                    // 5**
                    const is5xx = str.startsWith('5')
                    if (is5xx) {
                        return 'Server error. Try later.'
                    }
                    // ***
                    return 'Unknown server error.'
            }
        }
        // server not sent response
        if (err.type === "timeout") {
            return `Server not responding.`
        }
        if (err.type === "network") {
            return 'Network error.'
        }
        console.error(err)
        return 'Very unknown error.'
    }
}