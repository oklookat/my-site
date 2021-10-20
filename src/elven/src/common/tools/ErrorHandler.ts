import type { AxiosError } from "axios"

export default function ErrorHandler(axiosError): string {
    const message = sort(axiosError)
    if(window.$elvenNotify){
        window.$elvenNotify.add(message)
    }
    return message
}

function sort(error: AxiosError<any>): string { // axios error
    if (error.response) {
        let statusCode: number = error.response.status
        if(error.response.data.statusCode){
            statusCode = error.response.data.statusCode
        }
        switch (statusCode) {
            case 400:
                return 'Bad request.'
            case 401:
                return 'Wrong username or password.'
            case 403:
                return 'Access denied.'
            case 404:
                return 'Not found.'
            case 413:
                return 'Size too big.'
            default:
                const str = statusCode.toString()
                if(str.startsWith('5')){
                    return 'Server error.'
                }
                return 'Unknown error.'
        }
    } else if (error.request) {
        return 'Server not responding.'
    } else {
        console.error(error.message)
        return `Unknown error or request setup failed.`
    }
}