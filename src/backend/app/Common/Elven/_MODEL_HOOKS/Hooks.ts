import User from "App/Models/Elven/User"
import Article from "App/Models/Elven/Article"
import UserValidator from "App/Common/Elven/_VALIDATORS/UserValidator"
import SlugValidator from "App/Common/Elven/_VALIDATORS/SlugValidator"
import EL_Slug from "App/Common/Elven/_TOOLS/EL_Slug"
import EL_Random from "App/Common/Elven/_TOOLS/EL_Random"

const bcrypt = require('bcrypt')

class Hooks {

  public static async hashPassword(user: User) {
    if (user.$dirty.password) {
      const saltRounds = 14
      user.password = await bcrypt.hash(user.password, saltRounds)
    }
    return Promise.resolve(user.password)
  }

  public static async userValidate(user: User) {
    try {
      UserValidator.validateCredentials(user.username, user.password)
      return Promise.resolve()
    } catch (error) {
      return Promise.reject(error)
    }
  }

  public static autoSlug(article: Article): string {
    if(!article.id){
      return unknownSlug()
    }
    let slug = EL_Slug.make(`${article.title}-${article.id}`)
    const isValid = SlugValidator.validate(slug)
    if (!isValid) {
      return unknownSlug()
    }
    return slug
    function unknownSlug(){
      let slug = `unknown-${EL_Random.randString(EL_Random.randInt(4, 12), "alphabetic")}`
      slug = EL_Slug.make(slug)
      return slug
    }
  }

  public static setPublishedDate(article: Article) {
    if (article.is_published && !article.published_at) {
      article.published_at = new Date()
    }
    return article.published_at
  }

}

export default Hooks
