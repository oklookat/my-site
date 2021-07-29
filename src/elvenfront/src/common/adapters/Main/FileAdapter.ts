import Axios from '@/common/adapters/Axios.js'

export default class FileAdapter {

    public static async getFiles(page = '1', sortBy = 'created', start = 'newest'){
        return await Axios.get('files/', {params: {page: page, by: sortBy, start: start}})
            .then(response => {
                if (response.data) {
                    return Promise.resolve(response.data)
                } else {
                    return Promise.reject('Нет данных.')
                }
            })
            .catch(error => {
                return Promise.reject(error)
            })
    }

    public static async upload(files: FileList){
        for(const file in files){
            if(typeof files[file] !== "object"){
                continue
            }
            const formData = new FormData()
            formData.append("file", files[file])
            await Axios.post('files/', formData, {
                headers: {
                    'Content-Type': 'multipart/form-data'
                }
            })
        }
    }
}