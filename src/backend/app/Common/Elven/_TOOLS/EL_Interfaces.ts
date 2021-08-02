export interface EL_IError_AUTH_INCORRECT {
  readonly errorCode: 'E_AUTH_INCORRECT',
  message: string,
}

export interface EL_IError_VALIDATION_MUSTBE {
  readonly errorCode: 'E_VALIDATION_MUSTBE',
  issuer: string | string [],
  available: string [],
}

export interface EL_IError_VALIDATION_FORBIDDEN {
  readonly errorCode: 'E_VALIDATION_FORBIDDEN',
  issuer: string | string [],
}

export interface EL_IError_VALIDATION_MINMAX {
  readonly errorCode: 'E_VALIDATION_MINMAX',
  issuer: string | string [],
  min?: string,
  max?: string,
}

export interface EL_IError_UNKNOWN {
  readonly errorCode: 'E_UNKNOWN',
  message: string,
}

export interface EL_IError_CUSTOM {
  readonly errorCode: 'E_CUSTOM',
  data: object,
}
