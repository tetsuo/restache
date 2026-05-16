import { FlatNav, MultiNav, AccordionNav, NavVariant } from '.'
import * as NavItem from './NavItem'
import * as Route from './Route'

export const Home = f => FlatNav([NavItem.Home(f)])

export const Settings = x =>
  FlatNav([NavItem.Home, NavItem.AccountSettings, NavItem.NotificationSettings, NavItem.PrivacySettings].map(f => f(x)))

// TODO: fix recursive components
// export const Settings = x =>
//   MultiNav([
//     NavVariant(FlatNav([NavItem.Home()])),
//     NavVariant(
//       undefined,
//       AccordionNav(
//         'Settings',
//         'Update account settings',
//         [NavItem.AccountSettings, NavItem.NotificationSettings, NavItem.PrivacySettings].map(f => f(x)),
//       ),
//     ),
//   ])

export const NotificationSettings = Settings

export const PrivacySettings = Settings

export const Changelog = Home

export const System = x => FlatNav([NavItem.Changelog, NavItem.Settings, NavItem.Logout].map(f => f(x)))

const constantFalse = () => false

export function fromRoute(route) {
  switch (route.key) {
    case Route.Home.key:
      return Home
    case Route.Settings.key:
      return Settings
    case Route.NotificationSettings.key:
      return NotificationSettings
    case Route.PrivacySettings.key:
      return PrivacySettings
    case Route.Changelog.key:
      return Changelog
  }
  return Home(constantFalse)
}
