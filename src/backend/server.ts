import 'reflect-metadata'
import sourceMapSupport from 'source-map-support'
import { Ignitor } from '@adonisjs/core/build/standalone'

sourceMapSupport.install({ handleUncaughtExceptions: false })

import {createServer} from "https";
import {readFileSync} from "fs";
import { join } from 'path';
const privateKey = readFileSync(join(__dirname + '/start/ssl/key.pem'), 'utf8')
const certificate = readFileSync(join(__dirname + '/start/ssl/cert.pem'), 'utf8')
const credentials = {key: privateKey, cert: certificate}

new Ignitor(__dirname).httpServer().start((handle) => {
  return createServer(
    credentials,
    handle
  );
});
