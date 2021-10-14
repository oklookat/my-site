interface IE_ERROR {
    statusCode: number
    errorCode: string
    issuers: string[]
    message: string
}


export default class ErrorHandler {
    public static sortError(error) { // axios error
        if (error.response) {
            window.$elvenNotify.error(`Неизвестная ошибка. Код: ${error.response.status}`)
        } else if (error.request) {
            window.$elvenNotify.error(`Нет ответа от сервера.`)
        } else {
            window.$elvenNotify.error(`${error.message}`)
        }
    }
}