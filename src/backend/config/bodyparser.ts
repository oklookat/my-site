/**
 * Config source: https://git.io/Jfefn
 *
 * Feel free to let us know via PR, if you find something broken in this config
 * file.
 */

import { BodyParserConfig } from '@ioc:Adonis/Core/BodyParser'

const bodyParserConfig: BodyParserConfig = {

  whitelistedMethods: ['POST', 'PUT', 'PATCH', 'DELETE'],

  json: {
    encoding: 'utf-8',
    limit: '1mb',
    strict: true,
    types: [
      'application/json',
      'application/json-patch+json',
      'application/vnd.api+json',
      'application/csp-report',
    ],
  },

  form: {
    encoding: 'utf-8',
    limit: '1mb',
    queryString: {},

    convertEmptyStringsToNull: true,

    types: [
      'application/x-www-form-urlencoded',
    ],
  },

  raw: {
    encoding: 'utf-8',
    limit: '1mb',
    queryString: {},
    types: [
      'text/*',
    ],
  },


  multipart: {

    autoProcess: true,

    processManually: [], // '/api/elven/files'

    encoding: 'utf-8',

    convertEmptyStringsToNull: true,

    maxFields: 1000,

    limit: '256mb',

    types: [
      'multipart/form-data',
    ],
  },
}

export default bodyParserConfig
