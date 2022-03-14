export default class Utils {

  /** get object keys length */
  public static getObjectLength(o: object): number {
    if (!o || !(o instanceof Object)) {
      return 0
    }
    return Object.keys(o).length
  }

  public static debounce(f: Function, ms: number) {
    let isCooldown = false;
    return function () {
      if (isCooldown) return;

      f.apply(this, arguments);

      isCooldown = true;

      setTimeout(() => isCooldown = false, ms);
    };
  }

  /** is device with touchscreen & its default input? */
  public static isTouchDevice(): boolean {
    return matchMedia('(hover: none)').matches;
  }
  
}