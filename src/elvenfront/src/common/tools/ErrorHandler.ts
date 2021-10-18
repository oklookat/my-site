interface IE_ERROR {
    statusCode: number
    errorCode: string
    issuers: string[]
    message: string
}


export default class ErrorHandler {
    public static sortError(error) { // axios error
        if (error.response) {
            window.$elvenNotify.add(`Unknown error. Code: ${error.response.status}`)
        } else if (error.request) {
            window.$elvenNotify.add(`No response from server.`)
        } else {
            window.$elvenNotify.add(`${error.message}`)
        }
    }
}