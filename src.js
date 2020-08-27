if (
  typeof ServiceWorkerGlobalScope !== 'undefined' &&
  self instanceof ServiceWorkerGlobalScope
) {
  self.addEventListener('install', (event) => {
    event.waitUntil(self.skipWaiting());
  });
  self.addEventListener('activate', (event) => {
    const p = (async () => {
      await self.clients.claim();

      const existingClients = await clients.matchAll({
        includeUncontrolled: true,
        type: 'window',
      });

      // Safari might crash here.
      try {
        existingClients.forEach((client) => client.navigate(client.url));
      } catch (e) {}
    })();

    event.waitUntil(p);
  });
}
