import { invalidate } from '$app/navigation';
import type { Data } from '$lib_elven/types';
import * as cookie from 'cookie';

export default class Utils {

  /** if string > 31 length = cut it and add '...' */
  public static cutString(str: string): string {
    const maxLength = 31;
    if (str.length < maxLength) {
      return str;
    }
    return str.substring(0, maxLength) + '...';
  }

  /** get Record keys length */
  public static getRecordLength<T>(record: Record<any, T>): number {
    const isObject = typeof record === 'object' && record !== null
    if (!isObject) {
      throw Error("not a object")
    }
    return Object.keys(record).length
  }

  public static debounce(f: Function, ms: number) {
    let isCooldown = false
    return function () {
      if (isCooldown) {
        return
      }

      f.apply(this, arguments)

      isCooldown = true

      setTimeout(() => {
        isCooldown = false
      }, ms)
    }
  }

  /** refresh data. Used for pagination
   * 
   * @param fetchItems func to fetch items. if **isInitial = true** you must return items what you have now.
   */
  public static async refresh<T>(getPage: () => Promise<number>,
    setPage: (newPage: number) => Promise<void>, fetchItems: (isInitial: boolean) => Promise<Data<T>>) {
    //
    let page = await getPage()

    // if page === 1 - just refresh
    if (page < 2) {
      await fetchItems(false)
      return
    }

    let items = await fetchItems(true)
    let dataLength = 0

    try {
      dataLength = this.getRecordLength(items.data)
    } catch (err) {
      return
    }

    if (dataLength > 0) {
      return
    }

    while (page > 1 && dataLength < 1) {
      page = await getPage()
      page--
      await setPage(page)

      try {
        items = await fetchItems(false)
      } catch (err) {
        break
      }
      console.log(items)
      try {
        dataLength = this.getRecordLength(items.data)
      } catch (err) {
        break
      }
    }
  }

  /** convert URLSearchParams to object */
  public static searchParamsToObject(params: URLSearchParams): Object {
    if (!(params instanceof URLSearchParams)) {
      return
    }

    const result = {}

    params.forEach((value, key) => {
      result[key] = value;
    })

    return result;
  }

  /** add params to URLSearchParams by object */
  public static searchParamsByObject(params: URLSearchParams, data: Record<string, string>) {
    if (!(params instanceof URLSearchParams)) {
      return
    }

    for (const key in data) {
      // no add param if it exists in params
      if (params.has(key)) {
        params.delete(key)
      }

      const value = data[key]
      if (!value) {
        continue
      }
      params.append(key, value)
    }

    const result = {}
    params.forEach((value, key) => {
      result[key] = value;
    })

    return result;
  }

  /** get token from request headers (cookie) */
  public static getTokenFromRequestHeaders(headers: Headers): string | null {
    const isValid = !!(headers) && headers instanceof Headers && headers.has('cookie')
    if (!isValid) {
      return null
    }
    const cookiesStr = headers.get('cookie')
    let token = ''
    try {
      const parsed = cookie.parse(cookiesStr)
      if (!parsed || !parsed.token) {
        return null
      }
      token = parsed.token
    } catch (err) {
      return null
    }
    return token
  }

  public static addTokenToHeaders(token: string, headers: Headers) {
    if (!token || !(headers instanceof Headers)) {
      return
    }
    headers.append('Authorization', `Elven ${token}`)
  }

  /** check is element not out of screen, and if it is, correct position */
  public static correctElementOverflow(el: HTMLElement, evt: MouseEvent): { x: number; y: number } {
    let x = evt.clientX;
    let y = evt.clientY;
    const moveOffset = 10;

    // left-right (X)
    const popupWidth = el.offsetWidth;
    const overflowDifferenceX = x + popupWidth - window.screen.width;
    if (overflowDifferenceX > 0) {
      x = x - overflowDifferenceX - moveOffset;
    }

    // top-bottom (Y)
    const popupHeight = el.offsetHeight;
    const overflowDifferenceY = y + popupHeight - window.screen.height;
    if (overflowDifferenceY > 0) {
      y = y - overflowDifferenceY - moveOffset;
    }

    return { x, y };
  }

}