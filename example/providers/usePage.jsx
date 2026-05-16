import * as React from 'react'
import format from 'date-fns/format'
import parseISO from 'date-fns/parseISO'
import { Status } from '../model'
import * as Route from '../model/Route'
import * as Nav from '../model/Nav'
import { useLocation } from './Location'

export default function usePage() {
  const loc = useLocation()

  const [page, setPage] = React.useState(() => ({
    route: Route.fromLocation(loc),
    data: null,
    status: Status.Idle,
    error: undefined,
  }))

  const [revokeDialogOpen, setRevokeSessionDialogOpen] = React.useState(false)

  React.useEffect(() => {
    const pendingRoute = Route.fromLocation(loc)
    const aborter = new AbortController()

    setPage(prev => ({
      ...prev,
      status: Status.Pending,
      error: undefined,
    }))

    const pathname = pendingRoute.path === '/' ? '/timeline' : pendingRoute.path

    fetch(`mockapi${pathname}.json`, { signal: aborter.signal })
      .then(r => {
        if (!r.ok) throw new Error(`HTTP ${r.status}`)
        return r.json()
      })
      .then(json => {
        setPage({
          route: pendingRoute,
          data: json,
          status: Status.Ready,
          error: undefined,
        })
      })
      .catch(err => {
        if (err.name !== 'AbortError') {
          setPage({
            route: pendingRoute,
            data: null,
            status: Status.Error,
            error: err,
          })
        }
      })

    return () => aborter.abort()
  }, [loc.pathname, loc.key])

  const nav = React.useMemo(
    () => ({
      primary: Nav.fromRoute(page.route)(Route.keyEquals(page.route)),
      secondary: Nav.System(Route.nsEquals(page.route)),
    }),
    [page.route.key],
  )

  const showRevokeDialog = () => setRevokeSessionDialogOpen(true)
  const closeRevokeDialog = () => setRevokeSessionDialogOpen(false)

  const pageProps = {}

  if (page.data || page.status === Status.Error) {
    switch (page.route.key) {
      case Route.Home.key:
        pageProps.timeline = {
          pagination: {
            prev: {
              label: format(parseISO(page.data.overview.start), page.data.overview.tooltipFormat),
              disabled: false,
            },
            next: {
              label: format(parseISO(page.data.overview.end), page.data.overview.tooltipFormat),
              disabled: false,
            },
          },
          ruler: {
            start: parseISO(page.data.focus.start),
            end: parseISO(page.data.focus.end),
          },
          table: page.data.overview,
        }
        break
      case Route.Settings.key:
        pageProps.settings = {
          account: {
            basicInfo: page.data.basicInfo,
            sessions: page.data.sessions.map(x => ({ ...x, onClick: showRevokeDialog })),
            isRevokeDialogOpen: revokeDialogOpen,
            closeDialog: closeRevokeDialog,
          },
        }
        break
      case Route.NotificationSettings.key:
        pageProps.settings = { notifications: page.data }
        break
      case Route.PrivacySettings.key:
        pageProps.settings = { privacy: {} }
        break
      case Route.Changelog.key:
        pageProps.changelog = { items: page.data }
        break
    }
  }

  return {
    nav,
    loading: page.status === Status.Pending,
    pageProps,
  }
}
