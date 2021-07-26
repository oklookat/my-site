import User from "App/Models/Elven/User"
import Article from "App/Models/Elven/Article"
import UserValidator from "App/Common/Elven/_VALIDATORS/UserValidator"
import SlugValidator from "App/Common/Elven/_VALIDATORS/SlugValidator"
import Slug from "App/Common/Elven/_TOOLS/Slug"

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
    return await UserValidator.validateReg(user)
      .then(() =>{
        return Promise.resolve(true)
      })
      .catch(error =>{
        return Promise.reject(error)
      })
  }

  public static async autoSlug(article: Article){
    let slug = await Slug.make(article.title)
    const isValid = await SlugValidator.validate(slug)
    if(!isValid){
      slug = 'unknown'
      slug = await Slug.make(slug)
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
