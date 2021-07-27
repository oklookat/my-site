import Axios from '@/common/adapters/Axios.js'

export default class FileAdapter {
    public static async upload(files: FileList){
        const formData = new FormData();
        for(const file in files){
            formData.append("files", files[file]);
        }
        await Axios.post('files/', formData, {
            headers: {
                'Content-Type': 'multipart/form-data'
            }
        })
    }
}