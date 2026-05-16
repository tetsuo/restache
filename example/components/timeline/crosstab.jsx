import * as React from 'react'
import parseISO from 'date-fns/parseISO'
import format from 'date-fns/format'
import eachHourOfInterval from 'date-fns/eachHourOfInterval'
import eachDayOfInterval from 'date-fns/eachDayOfInterval'
import eachMinuteOfInterval from 'date-fns/eachMinuteOfInterval'
import { fromZonedTime } from 'date-fns-tz/fromZonedTime'
import Crosstab from './crosstab.html'

const EACH_INTERVAL = {
  minute: eachMinuteOfInterval,
  hour: eachHourOfInterval,
  day: eachDayOfInterval,
}

export default function CrosstabWrapper({ $ }) {
  const { start, end, resolution, labelFormat, headers, tooltipFormat = '', step = 1, timeZone } = $

  if (step <= 0) {
    throw new Error('step must be > 0')
  }

  const startDate = parseISO(start)
  const endDate = parseISO(end)
  if (startDate > endDate) {
    return []
  }

  const eachFn = EACH_INTERVAL[resolution]
  if (!eachFn) {
    throw new Error(`Unsupported resolution: ${resolution}`)
  }

  const dates = eachFn({ start: startDate, end: endDate }, { step })

  const rows = dates.map((d, i) => {
    const dateForLabel = timeZone ? fromZonedTime(d, timeZone) : d
    return {
      key: `d-${i}`,
      label: format(dateForLabel, labelFormat),
      tooltipText: tooltipFormat ? format(dateForLabel, tooltipFormat) : '',
      items: Array.from({ length: 4 }, (_, j) => ({
        key: j,
        label: '…',
      })),
    }
  })

  return <Crosstab stickyHeader={true} headers={headers} rows={rows} />
}
