const randomstring = require("randomstring")

export default class EL_Random {

  public static async randInt(min, max): Promise<number> {
    min = Math.ceil(min)
    max = Math.floor(max)
    return Promise.resolve(Math.floor(Math.random() * (max - min + 1)) + min)
  }

  public static async randString(length: number, charset: string): Promise<string> {
    // charsets:
    // alphanumeric, alphabetic, numeric, hex, binary, octal, custom
    // https://www.npmjs.com/package/randomstring
    let random
    try {
      random = await randomstring.generate({
        length: length,
        charset: charset
      })
    } catch (error) {
      throw error
    }
    return Promise.resolve(random)
  }

}
