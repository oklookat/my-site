import Axios from '@/common/adapters/Axios.js'

export default class FileAdapter {

    public static async getFiles(cursor = '', start = 'newest'): Promise<IFilesData> {
        const config = {
            params:
                {
                    cursor: cursor, start: start
                }
        }
        return Axios.get('files', config)
            .then(response => {
                if (response.data) {
                    let data: IFilesData = response.data
                    return Promise.resolve(data)
                } else {
                    return Promise.reject('Нет данных.')
                }
            })
            .catch(error => {
                return Promise.reject(error)
            })
    }

    public static async upload(files: FileList) {
        for (const file in files) {
            if (typeof files[file] !== "object") {
                continue
            }
            const formData = new FormData()
            formData.append("file", files[file])
            await Axios.post('files', formData, {
                headers: {
                    'Content-Type': 'multipart/form-data'
                }
            })
        }
    }

    public static async delete(id) {
        return await Axios.delete(`files/${id}`)
            .then(() => {
                return Promise.resolve(true)
            })
            .catch(error => {
                return Promise.reject(error)
            })
    }
}