import validator from "validator";

export default class SlugValidator {

  public static validate(slug: string) {
    if (!validator.isAscii(slug)) {
      return Promise.resolve(false)
    }
    if (validator.isEmpty(slug, {ignore_whitespace: true})) {
      return Promise.resolve(false)
    }
    return Promise.resolve(true)
  }

}
