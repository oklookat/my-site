/** handle request / response errors */
export class NetworkError {

    private static logStrange(tag: string | number, msg: string | number) {
        console.log('------------')
        console.error(`[${tag}] ${msg}`)
        console.log('------------')
    }

    public static handle(output: Response | Error): string {
        const isWindowExists = typeof window !== 'undefined' && window.$notify
        if (!isWindowExists) {
            return
        }
        const message = this.sort(output)
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
            case 401:
            case 403:
                return 'Wrong request or credentials.'
            case 404:
                return 'Not found.'
            case 409:
                return 'Already exists.'
            case 413:
                return 'Upload size too big.'
            default:
                const str = statusCode.toString()
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