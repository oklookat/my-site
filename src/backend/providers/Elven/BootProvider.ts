import { ApplicationContract } from '@ioc:Adonis/Core/Application'


export default class BootProvider {
  public static needsApplication = true
  constructor (protected application: ApplicationContract) {
  }

  public register () {
    // Register your own bindings
  }

  public async boot () {

  }

  public async ready () {

  }

  public async shutdown () {

  }
}
