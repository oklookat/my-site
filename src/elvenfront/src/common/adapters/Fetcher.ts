export default class Fetcher{
    private static api = import.meta.env.VITE_API_URL

    public static async check(){
        return await fetch(`${this.api}auth/check`, {credentials: 'include', method: 'POST'})
            .then(response => {
                if (!response.ok) {
                    return Promise.resolve(false)
                }
                return Promise.resolve(true)
            })
            .catch(error =>{
                return Promise.reject(error)
            })
    }
}