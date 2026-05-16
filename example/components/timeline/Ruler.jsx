import * as React from 'react'
import addMinutes from 'date-fns/addMinutes'
import getMinutes from 'date-fns/getMinutes'
import getHours from 'date-fns/getHours'

const shortH = 4
const longH = 7
const intervalDuration = 30 // in minutes

export default function TimeRuler({ startDate, endDate }) {
  if (!(startDate instanceof Date)) {
    startDate = new Date(startDate)
  }
  if (!(endDate instanceof Date)) {
    endDate = new Date(endDate)
  }

  const start = alignUp(startDate, intervalDuration)
  const end = alignDown(endDate, intervalDuration)

  const ticks = []
  const hourTicks = []

  for (let t = start, i = 0; t <= end; t = addMinutes(t, intervalDuration), i++) {
    ticks.push(t)
    if (getMinutes(t) === 0) {
      hourTicks.push({ date: t, idx: i })
    }
  }

  if (!ticks.length) {
    return null
  }

  const slotCount = ticks.length + 1

  return (
    <div
      style={{
        position: 'relative',
        width: '100%',
        height: longH + 18,
      }}
    >
      <div
        style={{
          position: 'relative',
          fontSize: 10,
          color: '#000000de',
          paddingBottom: 8,
        }}
      >
        {hourTicks.map(({ date, idx }) => {
          const pct = ((idx + 1) / slotCount) * 100
          return (
            <span
              key={idx}
              style={{
                position: 'absolute',
                left: `${pct}%`,
                transform: 'translateX(-50%)',
              }}
            >
              {pad(getHours(date))}
            </span>
          )
        })}
      </div>

      <svg
        width="100%"
        height={longH}
        viewBox={`0 0 ${slotCount} ${longH}`}
        preserveAspectRatio="none"
        shapeRendering="crispEdges"
      >
        {ticks.map((d, i) => {
          const x = i + 1
          const isHour = getMinutes(d) === 0
          const y1 = longH - (isHour ? longH : shortH)

          return (
            <line
              key={i}
              x1={x}
              x2={x}
              y1={y1}
              y2={longH}
              stroke="#dadada"
              strokeWidth={1}
              vectorEffect="non-scaling-stroke"
            />
          )
        })}
      </svg>
    </div>
  )
}

function alignUp(d, offset = 30) {
  const m = d.getMinutes()
  const delta = m === 0 ? 0 : m < offset ? offset - m : offset * 2 - m
  const out = addMinutes(d, delta)
  out.setSeconds(0, 0)
  return out
}

function alignDown(d, offset = 30) {
  const m = d.getMinutes()
  const delta = m >= offset ? -(m - offset) : -m
  const out = addMinutes(d, delta)
  out.setSeconds(0, 0)
  return out
}

function pad(n) {
  return String(n).padStart(2, '0')
}
