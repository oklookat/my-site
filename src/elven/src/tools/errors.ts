import type { HookOutput } from "@oklookat/duck"

/** handle request / response errors */
export class AdapterError {

    public static handle(err: HookOutput.onError): string {
        const message = this.sort(err)
        if (window.$notify) {
            window.$notify.add({message})
        }
        return message
    }

    private static sort(output: HookOutput.onError): string {
        // server send response with error
        const err = output.data
        if(err.type === "cancel") {
            if(!err.message) {
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
                case 400:
                    return 'Bad request.'
                case 401:
                    return 'Wrong credentials.'
                case 403:
                    return 'Access denied.'
                case 404:
                    return 'Not found.'
                case 413:
                    return 'Upload size too big.'
                default:
                    const str = statusCode.toString()
                    // 5**
                    if (str.startsWith('5')) {
                        return 'Server error.'
                    }
                    return 'Unknown server error.'
            }
        }
        // server not sent response
        if (err.type === "timeout") {
            return `Server not responding.`
        }
        if(err.type === "network") {
            return 'Network error.'
        }
        console.error(err)
        return 'Very unknown error.'
    }
}