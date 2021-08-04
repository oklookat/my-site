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

    /*
    |--------------------------------------------------------------------------
    | Files to be processed manually
    |--------------------------------------------------------------------------
    |
    | You can turn off `autoProcess` for certain routes by defining
    | routes inside the following array.
    |
    | NOTE: Make sure the route pattern starts with a leading slash.
    |
    | Correct
    | ```js
    | /projects/:id/file
    | ```
    |
    | Incorrect
    | ```js
    | projects/:id/file
    | ```
    */

    processManually: [], // '/api/elven/files'

    encoding: 'utf-8',

    convertEmptyStringsToNull: true,

    maxFields: 1000,

    limit: '242mb',

    types: [
      'multipart/form-data',
    ],
  },
}

export default bodyParserConfig
