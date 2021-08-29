import Route from '@ioc:Adonis/Core/Route'

Route.get('/', async () => {
  return { hello: 'world' }
})

Route.group(() => {
  // auth
  Route.post('/auth/login', 'Elven/AuthController.login')
  Route.post('/auth/logout', 'Elven/AuthController.logout')
  Route.post('/auth/check', 'Elven/AuthController.check').middleware(['elvenPerm:adminOnly'])

  // content
  Route.resource('/articles', 'Elven/ArticlesController')
    .except(['create', 'edit'])
    .middleware({
      '*': ['elvenPerm:readOnly'],
    })
  Route.resource('/files', 'Elven/FilesController')
    .except(['create', 'edit', 'update', 'show'])
    .middleware({
      '*': ['elvenPerm:adminOnly'],
    })
}).prefix('/elven')
