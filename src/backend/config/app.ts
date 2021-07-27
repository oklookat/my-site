import proxyAddr from 'proxy-addr'
import Env from '@ioc:Adonis/Core/Env'
import { ServerConfig } from '@ioc:Adonis/Core/Server'
import { LoggerConfig } from '@ioc:Adonis/Core/Logger'
import { ProfilerConfig } from '@ioc:Adonis/Core/Profiler'
import { ValidatorConfig } from '@ioc:Adonis/Core/Validator'


export const appKey: string = Env.get('APP_KEY')


export const http: ServerConfig = {

  allowMethodSpoofing: false,


  subdomainOffset: 2,

  /*
  |--------------------------------------------------------------------------
  | Request Ids
  |--------------------------------------------------------------------------
  |
  | Setting this value to `true` will generate a unique request id for each
  | HTTP request and set it as `x-request-id` header.
  |
  */
  generateRequestId: false,

  /*
  |--------------------------------------------------------------------------
  | Trusting proxy servers
  |--------------------------------------------------------------------------
  |
  | Define the proxy servers that AdonisJs must trust for reading `X-Forwarded`
  | headers.
  |
  */
  trustProxy: proxyAddr.compile('loopback'),

  /*
  |--------------------------------------------------------------------------
  | Generating Etag
  |--------------------------------------------------------------------------
  |
  | Whether or not to generate an etag for every response.
  |
  */
  etag: false,


  jsonpCallbackName: 'callback',


  cookie: {
    domain: '',
    path: '/',
    maxAge: '2h',
    httpOnly: true,
    secure: false,
    sameSite: false,
  },


  forceContentNegotiationTo: 'application/json',
}

/*
|--------------------------------------------------------------------------
| Logger
|--------------------------------------------------------------------------
*/
export const logger: LoggerConfig = {

  name: Env.get('APP_NAME'),


  enabled: true,


  level: Env.get('LOG_LEVEL', 'info'),


  prettyPrint: Env.get('NODE_ENV') === 'development',
}


export const profiler: ProfilerConfig = {
  enabled: true,


  blacklist: [],

  whitelist: [],
}


export const validator: ValidatorConfig = {
}
