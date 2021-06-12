import Server from '@ioc:Adonis/Core/Server'


Server.middleware.register([
  () => import('@ioc:Adonis/Core/BodyParser'),
])


Server.middleware.registerNamed({
  'elvenPerm': () => import('App/Middleware/Elven/Perm')
})
