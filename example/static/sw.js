/* eslint-disable no-undef */
const ESMSH_URL = 'https://esm.sh/'
const CACHE_NAME = 'esm-sh-cache-10022026'

self.addEventListener('install', () => {
  self.skipWaiting()
})

self.addEventListener('fetch', event => {
  const url = event.request.url

  if (url.startsWith(ESMSH_URL)) {
    event.respondWith(
      caches.open(CACHE_NAME).then(cache => {
        return cache.match(event.request).then(cachedResponse => {
          if (cachedResponse) {
            return cachedResponse
          }
          console.log('[Service Worker] Fetching new resource:', url)
          return fetch(event.request).then(response => {
            if (response.ok) {
              console.log('[Service Worker] Cached:', url)
              cache.put(event.request, response.clone())
            } else {
              console.warn('[Service Worker] Could not cache:', url)
            }
            return response
          })
        })
      }),
    )
  }
})

self.addEventListener('activate', event => {
  event.waitUntil(
    caches.keys().then(cacheNames => {
      return Promise.all(
        cacheNames.map(cacheName => {
          if (cacheName !== CACHE_NAME) {
            console.log('[Service Worker] Deleting old cache:', cacheName)
            return caches.delete(cacheName)
          }
        }),
      )
    }),
  )
})
