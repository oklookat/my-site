/* eslint-disable no-console */

import { register } from 'register-service-worker'

if (process.env.NODE_ENV === 'production') {
  register(`${process.env.BASE_URL}service-worker.js`, {
    ready () {
      console.log(
        'Приложение обслуживается через service worker.\n' +
        'Для получения информации посетите https://goo.gl/AFskqB'
      )
    },
    registered () {
      console.log('Service worker зарегестрирован.')
    },
    cached () {
      console.log('Контент кеширован для использования в оффлайне.')
    },
    updatefound () {
      console.log('Контент загружается...')
    },
    updated () {
      console.log('В приложении появился новый контент, обновите его.')
    },
    offline () {
      console.log('Нет интернета. Приложение работает в оффлайне.')
    },
    error (error) {
      console.error('Ошибка при регистрации service worker:', error)
    }
  })
}
