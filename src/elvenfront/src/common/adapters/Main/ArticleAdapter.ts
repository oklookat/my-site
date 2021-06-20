import Axios from '@/common/adapters/Axios'

class ArticleAdapter {

    public static async getArticles(page = '1') {
        return await Axios.get('articles/', {params: {page: page}})
            .then(response =>{
                if(response.data){
                    return Promise.resolve(response.data)
                } else{
                    return Promise.reject('Нет данных.')
                }
            })
            .catch(error =>{
                return Promise.reject(error)
            })
    }

    public static async deleteArticle(id){
        return await Axios.delete(`articles/${id}`)
            .then(() =>{
                return Promise.resolve(true)
            })
            .catch(error =>{
                return Promise.reject(error)
            })
    }

    public static async createArticle(article){
        return await Axios.post('articles/', article)
            .then(response =>{
                if(response.data){
                    return Promise.resolve(response.data)
                } else{
                    return Promise.reject('Нет данных.')
                }
            })
            .catch(error =>{
                return Promise.reject(error)
            })
    }

    public static async saveArticle(article){
        return await Axios.post(`articles/${article.id}`, article)
            .then(response =>{
                if(response.data){
                    return Promise.resolve(response.data)
                } else{
                    return Promise.reject('Нет данных.')
                }
            })
            .catch(error =>{
                return Promise.reject(error)
            })
    }
}

export default ArticleAdapter