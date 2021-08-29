/**
 * Config source: https://git.io/JfefC
 *
 * Feel free to let us know via PR, if you find something broken in this config
 * file.
 */

import { CorsConfig } from '@ioc:Adonis/Core/Cors'

const corsConfig: CorsConfig = {

  enabled: true,

  origin: ["https://elven.oklookat.ru", "https://static.oklookat.ru", "https://oklookat.ru"],

  methods: ['GET', 'HEAD', 'POST', 'PUT', 'DELETE'],

  /*
  |--------------------------------------------------------------------------
  | Headers
  |--------------------------------------------------------------------------
  |
  | List of headers to be allowed for `Access-Control-Allow-Headers` header.
  | The value can be one of the following:
  |
  | https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Headers
  |
  | Boolean(true)     - Allow all headers mentioned in `Access-Control-Request-Headers`.
  | Boolean(false)    - Disallow all headers.
  | String            - Comma separated list of allowed headers.
  | Array             - An array of allowed headers.
  | Function          - Receives the current header and should return one of the above values.
  |
  */
  headers: true,

  exposeHeaders: [
    'cache-control',
    'content-language',
    'content-type',
    'expires',
    'last-modified',
    'pragma',
  ],

  credentials: true,

  maxAge: 90,
}

export default corsConfig
