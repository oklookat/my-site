// like 500 error
export class E_UNKNOWN{
  constructor(issuers: string [], message?: string) {
    return {
      statusCode: 500,
      errorCode: 'E_UNKNOWN',
      issuers: issuers,
      message: message,
    }
  }
}

// like 'i need show message for users'
export class E_CUSTOM {
  constructor(issuers: string [], data?: object, message?: string) {
    return {
      statusCode: 500,
      errorCode: 'E_CUSTOM',
      issuers: issuers,
      data: data,
      message: message,
    }
  }
}

// like 'wrong username or password'
export class E_AUTH_INCORRECT {
  constructor(issuers: string []) {
    return {
      statusCode: 403,
      errorCode: 'E_AUTH_INCORRECT',
      issuers: issuers,
      message: 'Wrong credentials.'
    }
  }
}

// like 'only admins'
export class E_AUTH_FORBIDDEN {
  constructor(issuers: string []) {
    return {
      statusCode: 403,
      errorCode: 'E_AUTH_FORBIDDEN',
      issuers: issuers,
      message: 'Access denied.'
    }
  }
}

// like 'article not found'
export class E_NOTFOUND {
  constructor(issuers: string []) {
    return {
      statusCode: 404,
      errorCode: 'E_NOTFOUND',
      issuers: issuers,
      message: 'Not found.'
    }
  }
}

// like 'is_published must be true or false'
export class E_VALIDATION_MUSTBE {
  constructor(issuers: string [], available: string []) {
    return {
      statusCode: 400,
      errorCode: 'E_VALIDATION_MUSTBE',
      issuers: issuers,
      available: available,
      message: 'Value not in range.',
    }
  }
}

// like 'min length for username is 4 symbols'
export class E_VALIDATION_MINMAX {
  constructor(issuers: string [], min?: number, max?: number) {
    return {
      statusCode: 400,
      errorCode: 'E_VALIDATION_MINMAX',
      issuers: issuers,
      min: min,
      max: max,
      message: 'Too many or not enough characters.',
    }
  }
}

// like 'title cannot be empty'
export class E_VALIDATION_EMPTY {
  constructor(issuers: string []) {
    return {
      statusCode: 400,
      errorCode: 'E_VALIDATION_EMPTY',
      issuers: issuers,
      message: 'These things cannot be empty.',
    }
  }
}

// like 'allowed only numbers'
export class E_VALIDATION_ALLOWED {
  constructor(issuers: string [], allowed: string []) {
    return {
      statusCode: 400,
      errorCode: 'E_VALIDATION_ALLOWED',
      issuers: issuers,
      allowed: allowed,
      message: 'These things not allowed.',
    }
  }
}

// like 'request contains file, but file broken'
export class E_VALIDATION_INVALID {
  constructor(issuers: string [], message?: string) {
    return {
      statusCode: 400,
      errorCode: 'E_VALIDATION_INVALID',
      issuers: issuers,
      message: message,
    }
  }
}
