import slugify from "slugify"

const slugifyOptions = {
  replacement: '-',  // replace spaces with replacement character, defaults to `-`
  remove: /[*+~.()'"!:@]/g, // remove characters that match regex, defaults to `undefined`
  lower: true,      // convert to lower case, defaults to `false`
  strict: true,     // strip special characters except replacement, defaults to `false`
  locale: 'ru'       // language code of the locale to use
}


export default class Slug {
  public static async make(text: string) {
    return Promise.resolve(slugify(text, slugifyOptions))
  }
}
