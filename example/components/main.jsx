import * as React from 'react'
import Main from './main.html'
import usePage from '../providers/usePage'

export default function MainWrapper() {
  const { nav, loading, pageProps } = usePage()
  return <Main nav={nav} loading={loading} {...pageProps} />
}
