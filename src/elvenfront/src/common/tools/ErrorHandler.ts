interface IE_ERROR {
    statusCode: number
    errorCode: string
    issuers: string[]
    message: string
}

interface IE_UNKNOWN extends IE_ERROR {

}

interface IE_CUSTOM extends IE_ERROR {
    data: object
}

interface IE_AUTH_INCORRECT extends IE_ERROR {

}

interface IE_AUTH_FORBIDDEN extends IE_ERROR {

}

interface IE_NOTFOUND extends IE_ERROR {

}

interface IE_VALIDATION_ALLOWED extends IE_ERROR {
    allowed: string[]
}

interface IE_VALIDATION_MINMAX extends IE_ERROR {
    min: number
    max: number
}

interface IE_VALIDATION_EMPTY extends IE_ERROR {

}

interface IE_VALIDATION_INVALID extends IE_ERROR {

}


const errorTypes = {
    unknown: 'E_UNKNOWN',
    custom: 'E_CUSTOM',
    authIncorrect: 'E_AUTH_INCORRECT',
    authForbidden: 'E_AUTH_FORBIDDEN',
    notFound: 'E_NOTFOUND',
    validationAllowed: 'E_VALIDATION_ALLOWED',
    validationMinMax: 'E_VALIDATION_MINMAX',
    validationEmpty: 'E_VALIDATION_EMPTY',
    validationInvalid: 'E_VALIDATION_INVALID'
}

export default class ErrorHandler {
    public static sortError(error) { // axios error
        if (error.response) {
            console.log(error.response)
            if (error.response.data && error.response.data.errors) {
                this.errorCollectorHandler(error.response.data.errors)
            } else {
                window.app.$elvenNotify.error(`Неизвестная ошибка. Код: ${error.response.status}`)
            }
        } else if (error.request) {
            window.app.$elvenNotify.error(`Нет ответа от сервера.`)
        } else {
            window.app.$elvenNotify.error(`${error.message}`)
        }
    }

    private static errorCollectorHandler(errors) {
        for (let error of errors) {
            switch (error.errorCode) {
                case errorTypes.unknown:
                    window.app.$elvenNotify.error('Ошибка сервера. Попробуйте позже.')
                    break
                case errorTypes.custom:
                    const err: IE_CUSTOM = error
                    err.message
                    window.app.$elvenNotify.info(`Сообщение от сервера: ${err.message}`)
                    break
                case errorTypes.authIncorrect:
                    window.app.$elvenNotify.error('Ошибка авторизации.')
                    break
                case errorTypes.authForbidden:
                    window.app.$elvenNotify.error('Доступ запрещен.')
                    break
                case errorTypes.notFound:
                    window.app.$elvenNotify.error('Не найдено.')
                    break
                case errorTypes.validationAllowed:
                    const v_err: IE_VALIDATION_ALLOWED = error
                    let allowedString = ''
                    for(let allowed of v_err.allowed){
                        allowedString = allowedString + allowed + ';'
                    }
                    window.app.$elvenNotify.error(`${v_err.issuers[0]}: разрешены только ${allowedString}`)
                    break
                case errorTypes.validationMinMax:
                    const vm_err: IE_VALIDATION_MINMAX = error
                    window.app.$elvenNotify.error(`${vm_err.issuers[0]}: мин ${vm_err.min} / макс ${vm_err.max}`)
                    break
                case errorTypes.validationEmpty:
                    const em_err: IE_VALIDATION_EMPTY = error
                    window.app.$elvenNotify.error(`${em_err.issuers[0]}: не может быть пустым`)
                    break
                case errorTypes.validationInvalid:
                    const vin_err: IE_VALIDATION_INVALID = error
                    window.app.$elvenNotify.error(`${vin_err.issuers[0]}: сервер не может обработать это. Проверьте правильность.`)
                    break
                default:
                    window.app.$elvenNotify.error(`Неизвестная ошибка.`)
                    break
            }
        }
    }
}