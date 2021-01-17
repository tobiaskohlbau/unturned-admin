var CACHE_NAME = 'unturned-admin-cache-v1';
var urlsToCache = [
  '/',
];

self.addEventListener('install', function(event) {
  // Perform install steps
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then(function(cache) {
        return cache.addAll(urlsToCache);
      })
  );
});

self.addEventListener('fetch',() => {});
