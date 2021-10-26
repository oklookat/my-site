import type { AxiosError } from "axios"


export class AdapterError {

    public static handle(err: AxiosError<any>): string {
        const message = this.sort(err)
        if (window.$elvenNotify) {
            window.$elvenNotify.add(message)
        }
        return message
    }

    private static sort(err: AxiosError<any>): string {
        // server send response with error
        if (err.response) {
            let statusCode: number = err.response.status
            let apiStatusCode: number = err.response.data.statusCode
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
        if (err.request) {
            if (err.code && err.code === 'ECONNABORTED') {
                return `Server not responding.`
            }
            console.error(err.message)
            return `Unknown error while request.`
        }
        console.log(err)
        return 'Very unknown error.'
    }
}