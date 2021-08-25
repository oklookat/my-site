const randomstring = require("randomstring")

export default class EL_Random {

  public static randInt(min, max): number {
    min = Math.ceil(min)
    max = Math.floor(max)
    return Math.floor(Math.random() * (max - min + 1)) + min
  }

  public static randString(length: number, charset: string): string {
    // charsets:
    // alphanumeric, alphabetic, numeric, hex, binary, octal, custom
    // https://www.npmjs.com/package/randomstring
    let random
    try {
      random = randomstring.generate({
        length: length,
        charset: charset
      })
    } catch (error) {
      throw error
    }
    return random
  }

}
