import i18n from 'sveltekit-i18n';
import type { Config } from 'sveltekit-i18n'


function generateOneLocale(loc: string) {
  // https://github.com/rollup/plugins/tree/master/packages/dynamic-import-vars#imports-must-start-with--or-
  const elvenPath = `elven/${loc}`
  return [
    //////// elven
    {
      locale: loc,
      key: 'elven.components',
      loader: async () => (
        await import(`./${elvenPath}/components.json`)
      ).default,
    },
    {
      locale: loc,
      key: 'elven.plugins',
      loader: async () => (
        await import(`./${elvenPath}/plugins.json`)
      ).default,
    },
    //
    {
      locale: loc,
      key: 'elven.index',
      routes: ["/elven"],
      loader: async () => (
        await import(`./${elvenPath}/index.json`)
      ).default,
    },
    {
      locale: loc,
      key: 'elven.articles',
      routes: ["/elven/articles", "/elven/articles/create", "/blog"],
      loader: async () => (
        await import(`./${elvenPath}/articles.json`)
      ).default,
    },
    {
      locale: loc,
      key: 'elven.files',
      routes: ["/elven/files", "/elven/articles/create"],
      loader: async () => (
        await import(`./${elvenPath}/files.json`)
      ).default,
    },
    {
      locale: loc,
      key: 'elven.auth',
      routes: ["/elven/login", "/elven/logout"],
      loader: async () => (
        await import(`./${elvenPath}/auth.json`)
      ).default,
    },
    {
      locale: loc,
      key: 'elven.settings',
      routes: ["/elven/settings"],
      loader: async () => (
        await import(`./${elvenPath}/settings.json`)
      ).default,
    }
  ]
}

function generateRouteLoaders(locales: string[]): Array<any> {
  const localed = []
  for (const locName of locales) {
    const comps = generateOneLocale(locName)
    localed.push(...comps)
  }
  return localed
}

const config: Config = ({
  fallbackLocale: 'en',
  loaders: generateRouteLoaders(["en"])
});

export const { t, locale, locales, loading, loadTranslations } = new i18n(config);