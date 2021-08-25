import slugify from "slugify"

const slugifyOptions = {
  replacement: '-',  // replace spaces with replacement character, defaults to `-`
  remove: /[*+~.()'"!:@]/g, // remove characters that match regex, defaults to `undefined`
  lower: true,      // convert to lower case, defaults to `false`
  strict: true,     // strip special characters except replacement, defaults to `false`
  locale: 'ru'       // language code of the locale to use
}


export default class EL_Slug {
  public static make(text: string) {
    let slug = slugify(text, slugifyOptions)
    let arr = slug.split("")
    if (arr[0] == "-") {
      arr[0] = ""
    }
    if (arr[arr.length - 1] == "-") {
      arr[arr.length - 1] = ""
    }
    slug = arr.join("")
    return slug
  }
}
