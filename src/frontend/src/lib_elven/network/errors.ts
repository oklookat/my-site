/** handle request / response errors */
export class NetworkError {

    private static logStrange(tag: string | number, msg: string | number) {
        console.log('------------')
        console.error(`[${tag}] ${msg}`)
        console.log('------------')
    }

    public static handle(output: Response | Error): string {
        const message = this.sort(output)
        const isWindowExists = typeof window !== 'undefined' && window.$notify
        if (!isWindowExists) {
            console.error(message)
            return message
        }
        window.$notify.add({ message })
        return message
    }

    private static sort(output: Response | Error): string {
        if (output instanceof Error) {
            return 'Network error.'
        }

        if (!(output instanceof Response)) {
            return 'Unknown error.'
        }

        if (output.ok) {
            this.logStrange(output.status, output.statusText)
            return 'Not error. But why i displayed this message?'
        }

        let statusCode = output.status
        switch (statusCode) {
            case 400:
                return 'Bad request.'
            case 401:
                return 'Authorization failed.'
            case 403:
                return 'Authentication failed.'
            case 404:
                return 'Not found.'
            case 409:
                return 'Already exists.'
            case 413:
                return 'Body size too big.'
            default:
                // 5**
                const is5xx = statusCode > 499 && statusCode < 600
                if (is5xx) {
                    return 'Server error. Try later.'
                }
        }

        // ***
        this.logStrange(output.status, output.statusText)
        return 'Very unknown error.'
    }
}