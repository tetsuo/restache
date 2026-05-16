import { Route } from '.'

const notFound = Symbol('NotFound')
const internalError = Symbol('InternalError')

const home = Symbol('Home')
const accountSettings = Symbol('AccountSettings')
const notificationSettings = Symbol('NotificationSettings')
const privacySettings = Symbol('PrivacySettings')
const changelog = Symbol('Changelog')

const settingsNS = Symbol('Settings')

export const NotFound = path => Route(notFound, path)

export const InternalError = path => Route(internalError, path)

export const Home = Route(home, '/')

export const Settings = Route(accountSettings, '/settings', settingsNS)

export const NotificationSettings = Route(notificationSettings, '/settings/notifications', settingsNS)

export const PrivacySettings = Route(privacySettings, '/settings/privacy', settingsNS)

export const Changelog = Route(changelog, '/changelog')

export function fromLocation(loc) {
  switch (loc.pathname) {
    case Home.path:
      return Home
    case Settings.path:
      return Settings
    case NotificationSettings.path:
      return NotificationSettings
    case PrivacySettings.path:
      return PrivacySettings
    case Changelog.path:
      return Changelog
  }
  return NotFound(loc.pathname)
}

export const keyEquals = a => b => a.key == b.key

export const nsEquals = a => b => a.namespace == b.namespace
