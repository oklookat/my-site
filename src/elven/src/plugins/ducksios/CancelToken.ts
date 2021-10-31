export default class CancelToken {

    public cancel(message?: string) {
        console.warn('ducksios: you called empty CancelToken. Pass CancelToken instance in request settings.')
    }

}