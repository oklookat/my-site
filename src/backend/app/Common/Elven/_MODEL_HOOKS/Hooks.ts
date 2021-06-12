import ElvenTools from "App/Common/Elven/_TOOLS/ElvenTools"
import User from "App/Models/Elven/User"
import Article from "App/Models/Elven/Article"
import validator from "validator"
import ElvenValidators from "App/Common/Elven/_VALIDATORS/ElvenValidators";

const bcrypt = require('bcrypt')

class Hooks{

  public static async hashPassword (user: User) {
    if (user.$dirty.password) {
      const saltRounds = 12
      user.password = await bcrypt.hash(user.password, saltRounds)
    }
    return Promise.resolve(user.password)
  }

  public static async userValidate(user: User){
    let isUsername = validator.isLength(user.username, {min: 4, max: 24})
    if (!isUsername) {
      const err = new Error('Имя пользователя должно быть больше 4, и меньше 24 символов.')
      return Promise.reject(err)
    }
    isUsername = validator.isAlphanumeric(user.username)
    if (!isUsername) {
      const err = new Error('Имя пользователя должно быть без странных символов, и только на английском языке.')
      return Promise.reject(err)
    }

    const isPass = validator.isLength(user.password, {min: 8, max: 64})
    if (!isPass) {
      const err = new Error('Пароль должен быть больше 8, и меньше 64 символов.')
      return Promise.reject(err)
    }
    return Promise.resolve(true)
  }

  public static async autoSlug(article: Article){
    let slug = await ElvenTools.makeSlug(article.title)
    const isValid = await ElvenValidators.slugValidate(slug)
    if(!isValid){
      slug = 'unknown'
      slug = await ElvenTools.makeSlug(slug)
    }
    // reference
    // https://github.com/adonisjs/adonis-lucid-slugify/blob/develop/src/Strategies/dbIncrement.js
    const articleFound = await Article.query()
      .whereRaw(`?? ~* ?`, ['slug', `^${slug}(-[0-9]*)?$`])
      .orderBy('id', 'desc')
      .limit(1)
    if (articleFound.length < 1) {
      article.slug = slug
      return Promise.resolve(article.slug)
    }
    const lastNum = Number(articleFound[0].slug.replace(`${slug}-`, ''))
    article.slug = !lastNum || isNaN(lastNum) ? `${slug}-1` : `${slug}-${lastNum + 1}`
    return Promise.resolve(article.slug)
  }

  public static async setPublishedDate(article: Article){
    if(article.is_published && !article.published_at){
      article.published_at = new Date()
    }
    return Promise.resolve(article.published_at)
  }

}

export default Hooks
