/** you can cancel request by passing CancelToken instance in request config */
export default class CancelToken {

    /** cancel request */
    public cancel(message?: string) {
        console.warn('ducksios: you called empty CancelToken. Pass CancelToken instance in request config')
    }

}