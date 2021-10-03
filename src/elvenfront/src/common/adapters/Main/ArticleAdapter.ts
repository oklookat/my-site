import Axios from '@/common/adapters/Axios.js'


class ArticleAdapter {

    public static async getArticles(cursor = '', show = 'published', sortBy = 'updated', start = 'newest') {
        const config = {
            params:
                {
                    cursor: cursor, show: show, by: sortBy, start: start
                }
        }
        return await Axios.get('articles', config)
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

    public static async getArticle(id) {
        return await Axios.get(`articles/${id}`)
            .then(response => {
                if (response.data) {
                    return Promise.resolve(response.data)
                } else {
                    return Promise.reject('Запись не найдена.')
                }
            })
            .catch(error => {
                return Promise.reject(error.response.status)
            })
    }

    public static async deleteArticle(id) {
        return await Axios.delete(`articles/${id}`)
            .then(() => {
                return Promise.resolve(true)
            })
            .catch(error => {
                return Promise.reject(error)
            })
    }

    public static async createArticle(article) {
        console.log(article)
        return await Axios.post('articles', article)
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

    public static async saveArticle(article) {
        return await Axios.put(`articles/${article.id}`, article)
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

    public static async publishArticle(article) {
        return await Axios.put(`articles/${article.id}`, {is_published: true})
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

    public static async makeDraftArticle(article) {
        return await Axios.put(`articles/${article.id}`, {is_published: false})
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
}

export default ArticleAdapter