import type { TError } from "@/plugins/ducksios/Types"

export class AdapterError {

    public static handle(err: TError): string {
        const message = this.sort(err)
        if (window.$elvenNotify) {
            window.$elvenNotify.add(message)
        }
        return message
    }

    private static sort(err: TError): string {
        // server send response with error
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
                    return 'Size too big.'
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
        console.log(err)
        return 'Very unknown error.'
    }
}