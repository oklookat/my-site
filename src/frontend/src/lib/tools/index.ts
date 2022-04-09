export default class Utils {

  /** get Record keys length */
  public static getRecordLength<T>(record: Record<any, T>): number {
    const isObject = typeof record === 'object' && record !== null
    if (!isObject) {
      return 0
    }
    return Object.keys(record).length
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

  /** refresh data. Used for pagination */
  public static async refresh<T>(page: number,
    setPage: (val: number) => void, fetchData: () => Promise<Record<number, T>>) {
    let data = await fetchData()
    if (page < 2) {
      return
    }
    while (true) {
      const dataLength = this.getRecordLength(data)
      const isNoData = dataLength > 0 || page < 2
      if (isNoData) {
        break
      }
      page--
      setPage(page)
      try {
        data = await fetchData()
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

}