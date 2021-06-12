const randomstring = require("randomstring")
import slugify from "slugify"
const slugifyOptions = {
  replacement: '-',  // replace spaces with replacement character, defaults to `-`
  remove: /[*+~.()'"!:@]/g, // remove characters that match regex, defaults to `undefined`
  lower: true,      // convert to lower case, defaults to `false`
  strict: true,     // strip special characters except replacement, defaults to `false`
  locale: 'ru'       // language code of the locale to use
}

class ElvenTools {

  public static async errorConstructor(type: string, message: string): Promise<{ type: string, message: string }> {
    return Promise.resolve({type: type, message: message})
  }

  public static async publicErrorConstructor(message: string): Promise<{ error: string }>{
    return Promise.resolve({error: message})
  }

  public static async publicSuccessConstructor(message: string): Promise<{ success: string }>{
    return Promise.resolve({success: message})
  }

  public static async getRandomInt(min, max): Promise<number> {
    min = Math.ceil(min)
    max = Math.floor(max)
    return Promise.resolve(Math.floor(Math.random() * (max - min + 1)) + min)
  }

  public static async getRandomString(length: number, charset: string): Promise<string> {
    // charsets
    // alphanumeric, alphabetic, numeric, hex, binary, octal, custom
    // https://www.npmjs.com/package/randomstring
    let random
    try {
      random = await randomstring.generate({
        length: length,
        charset: charset
      })
    } catch (error) {
      console.log(error)
      random = await randomstring.generate({length: length})
    }
    return Promise.resolve(random)
  }

  public static async makeSlug(text: string){
    return Promise.resolve(slugify(text, slugifyOptions))
  }
}

export default ElvenTools
