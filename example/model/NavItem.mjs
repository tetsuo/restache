import * as React from 'react'
import CampaignIcon from '@mui/icons-material/Campaign'
import SettingsIcon from '@mui/icons-material/Settings'
import ExitToAppIcon from '@mui/icons-material/ExitToApp'
import HomeIcon from '@mui/icons-material/Home'
import AccountBoxIcon from '@mui/icons-material/AccountBox'
import NotificationsIcon from '@mui/icons-material/Notifications'
import LockIcon from '@mui/icons-material/Lock'
import * as Route from './Route'
import * as ExternalLink from './ExternalLink'
import { NavItem, NavItemTarget } from '.'

/* targets */

const HomeTarget = NavItemTarget(Route.Home)

const ChangelogTarget = NavItemTarget(Route.Changelog)

const SettingsTarget = NavItemTarget(Route.Settings)

const NotificationSettingsTarget = NavItemTarget(Route.NotificationSettings)

const PrivacySettingsTarget = NavItemTarget(Route.PrivacySettings)

const LogoutTarget = NavItemTarget(undefined, ExternalLink.Logout)

/* icon elements */

const HomeIconElement = React.createElement(HomeIcon)

const CampaignIconElement = React.createElement(CampaignIcon)

const SettingsIconElement = React.createElement(SettingsIcon)

const AccountBoxIconElement = React.createElement(AccountBoxIcon)

const NotificationsIconElement = React.createElement(NotificationsIcon)

const LockIconElement = React.createElement(LockIcon)

const ExitToAppIconElement = React.createElement(ExitToAppIcon, { style: { transform: 'scaleX(-1)' } })

/* navitems */

export const Home = f => NavItem(HomeTarget, 'Home', HomeIconElement, f)

export const Changelog = f => NavItem(ChangelogTarget, "What's new?", CampaignIconElement, f)

export const Settings = f => NavItem(SettingsTarget, 'Settings', SettingsIconElement, f)

export const AccountSettings = f => NavItem(SettingsTarget, 'Account', AccountBoxIconElement, f)

export const NotificationSettings = f =>
  NavItem(NotificationSettingsTarget, 'Notifications', NotificationsIconElement, f)

export const PrivacySettings = f => NavItem(PrivacySettingsTarget, 'Privacy', LockIconElement, f)

const constantFalse = () => false

export const Logout = () => NavItem(LogoutTarget, 'Log out', ExitToAppIconElement, constantFalse)
