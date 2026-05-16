import * as React from 'react'
import { Layout } from '../../model'
import Sidebar from './sidebar.html'

export default function SidebarWrapper(props) {
  const [layout, setLayout] = React.useState(Layout.Expanded)
  const onToggle = evt => {
    evt.preventDefault()
    setLayout(prev => (prev === Layout.Compact ? Layout.Expanded : Layout.Compact))
  }
  let nextProps
  if (layout !== Layout.Compact) {
    nextProps = props
  } else {
    nextProps = {
      ...props,
      ...{
        nav: {
          ...props.nav,
          ...{
            primary: clearNavSectionLabels(props.nav.primary),
            secondary: clearNavSectionLabels(props.nav.secondary),
          },
        },
      },
    }
  }
  return <Sidebar onToggle={onToggle} className={layout.description} {...nextProps} />
}

function clearNavSectionLabels(nav) {
  return {
    ...nav,
    ...{
      items: nav.items.map(items => ({
        ...items,
        ...{
          label: '',
        },
      })),
    },
  }
}
