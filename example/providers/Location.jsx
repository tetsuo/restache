import * as React from 'react'
import history from 'history/hash'

export const LocationCtx = React.createContext(null)

export default function LocationProvider({ children }) {
  const [loc, setLoc] = React.useState(history.location) // { pathname, search, key }

  React.useEffect(() => {
    const unlisten = history.listen(({ location }) => setLoc(location))
    return unlisten
  }, [])

  return <LocationCtx.Provider value={loc}>{children}</LocationCtx.Provider>
}

export function useLocation() {
  const loc = React.useContext(LocationCtx)
  if (!loc) {
    throw new Error('useLocation must be inside <LocationProvider>')
  }
  return loc
}
