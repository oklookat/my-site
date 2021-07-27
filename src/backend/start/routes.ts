import Route from '@ioc:Adonis/Core/Route'

Route.get('/', async () => {
  return { hello: 'world' }
})

Route.group(() => {
  Route.post('/auth/login', 'Elven/AuthController.login')
  Route.post('/auth/logout', 'Elven/AuthController.logout')
  Route.resource('/articles', 'Elven/ArticlesController')
    .except(['create', 'edit'])
    .middleware({
      '*': ['elvenPerm:readOnly'],
    })
  Route.resource('/files', 'Elven/FilesController')
    .except(['create', 'edit'])
    .middleware({
      '*': ['elvenPerm:readOnly'],
    })
}).prefix('/api/elven')
