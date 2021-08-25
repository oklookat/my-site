import validator from "validator";

export default class SlugValidator {

  public static validate(slug: string) {
    return validator.isSlug(slug)
  }

}
