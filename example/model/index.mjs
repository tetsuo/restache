import history from 'history/hash'

export const Route = (key, path, namespace) => ({
  key,
  path,
  namespace: namespace || key,
})

export const ExternalLink = href => ({ href })

export const NavItemTarget = (route, external) => ({ route, external })

export const NavItem = (target, label, icon, routeActive) => {
  const o = {
    target,
    label,
    icon,
    active: false,
  }
  if (target.route) {
    if (routeActive?.(target.route)) {
      o.active = true
    }
    o.href = target.route.path
    o.key = target.route.key.description
    o.triggerClick = triggerRouteClick
  } else {
    o.href = target.external.href
    o.key = target.external.href
    o.triggerClick = () => void 0
  }
  return o
}

export const FlatNav = items => ({ items })

export const AccordionNav = (title, subtitle, items, open) => ({
  title,
  subtitle,
  items,
  open,
})

export const NavVariant = (flat, accordion, multi) => ({
  flat,
  accordion,
  multi,
})

// TODO: fix recursives in restache
export const MultiNav = variants => ({ variants })

export const Layout = {
  Compact: Symbol('compact'),
  Expanded: Symbol('expanded'),
}

export const Status = {
  Idle: Symbol('idle'),
  Pending: Symbol('pending'),
  Ready: Symbol('ready'),
  Error: Symbol('error'),
}

export const Page = (route, title, navVariant, layout, status) => ({
  route,
  title,
  nav: navVariant,
  layout: layout || Layout.Expanded,
  status: status || Status.Idle,
})

function triggerRouteClick(event) {
  event.preventDefault()

  const href = event.currentTarget.getAttribute('href')
  if (!href) {
    return
  }

  const url = new URL(href, window.location.origin)
  const path = url.pathname

  if (path !== history.location.pathname) {
    history.push(path)
  } else {
    history.replace(path)
  }
}
