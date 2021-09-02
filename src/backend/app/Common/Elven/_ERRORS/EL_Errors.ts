interface IE_ERROR {
  statusCode: number
  errorCode: string
  issuers: string[]
  message: string
}

interface IE_UNKNOWN extends IE_ERROR {

}

interface IE_CUSTOM extends IE_ERROR {
  data: object
}

interface IE_AUTH_INCORRECT extends IE_ERROR {

}

interface IE_AUTH_FORBIDDEN extends IE_ERROR {

}

interface IE_NOTFOUND extends IE_ERROR {

}

interface IE_VALIDATION_ALLOWED extends IE_ERROR {
  allowed: string[]
}

interface IE_VALIDATION_MINMAX extends IE_ERROR {
  min: number
  max: number
}

interface IE_VALIDATION_EMPTY extends IE_ERROR {

}

interface IE_VALIDATION_INVALID extends IE_ERROR {

}


// like 500 error
export class E_UNKNOWN {
  constructor(issuers: string [], message: string) {
    const bErr: IE_UNKNOWN = {
      statusCode: 500,
      errorCode: 'E_UNKNOWN',
      issuers: issuers,
      message: message
    }
    return bErr
  }
}

// like 'i need show message for users'
export class E_CUSTOM {
  constructor(issuers: string [], data: object = {}, message: string) {
    const bErr: IE_CUSTOM = {
      statusCode: 500,
      errorCode: 'E_CUSTOM',
      issuers: issuers,
      data: data,
      message: message
    }
    return bErr
  }
}

// like 'wrong username or password'
export class E_AUTH_INCORRECT {
  constructor(issuers: string []) {
    const err: IE_AUTH_INCORRECT = {
      statusCode: 403,
      errorCode: 'E_AUTH_INCORRECT',
      issuers: issuers,
      message: 'Wrong credentials.'
    }
    return err
  }
}

// like 'only admins'
export class E_AUTH_FORBIDDEN {
  constructor(issuers: string []) {
    const err: IE_AUTH_FORBIDDEN = {
      statusCode: 403,
      errorCode: 'E_AUTH_FORBIDDEN',
      issuers: issuers,
      message: 'Access denied.'
    }
    return err
  }
}

// like 'article not found'
export class E_NOTFOUND {
  constructor(issuers: string []) {
    const err: IE_NOTFOUND = {
      statusCode: 404,
      errorCode: 'E_NOTFOUND',
      issuers: issuers,
      message: 'Not found.'
    }
    return err
  }
}

// like 'allowed only numbers' or 'is_published must be true or false'
export class E_VALIDATION_ALLOWED {
  constructor(issuers: string [], allowed: string []) {
    const err: IE_VALIDATION_ALLOWED = {
      statusCode: 400,
      errorCode: 'E_VALIDATION_ALLOWED',
      issuers: issuers,
      allowed: allowed,
      message: 'These things not allowed.',
    }
    return err
  }
}

// like 'min length for username is 4 symbols'
export class E_VALIDATION_MINMAX {
  constructor(issuers: string [], min: number = 0, max: number = 0) {
    const err: IE_VALIDATION_MINMAX = {
      statusCode: 400,
      errorCode: 'E_VALIDATION_MINMAX',
      issuers: issuers,
      min: min,
      max: max,
      message: 'Too many or not enough characters.',
    }
    return err
  }
}

// like 'title cannot be empty'
export class E_VALIDATION_EMPTY {
  constructor(issuers: string []) {
    const err: IE_VALIDATION_EMPTY = {
      statusCode: 400,
      errorCode: 'E_VALIDATION_EMPTY',
      issuers: issuers,
      message: 'These things cannot be empty.',
    }
    return err
  }
}

// like 'request contains file, but file broken'
export class E_VALIDATION_INVALID {
  constructor(issuers: string [], message: string) {
    const err: IE_VALIDATION_INVALID = {
      statusCode: 400,
      errorCode: 'E_VALIDATION_INVALID',
      issuers: issuers,
      message: message,
    }
    return err
  }
}
