let tileDomain = 'tile.dchr.host';

/**
 * Return 503 if query to upstream tileserver fails
 * @return {Response} 503 response
 */
function unableToResolve() {
	return new Response('', {status: 503, statusText: 'service unavailable'});
}

/**
 * Check if cached tile data is still valid
 * @param  {Object}  response The response object
 * @return {Boolean} If true, cached data is valid
 */
let isValid = function (response) {
	if (!response) return false;
	let fetched = response.headers.get('sw-fetched-on');
	// cache tiles for three days before they expire and are re-fetched
	return !!(fetched &&
		(parseFloat(fetched) + (1000 * 60 * 60 * 24 * 3))
		> new Date().getTime());
};

/**
 * Intercept and cache requests for tiles. Expire cached tiles after three days.
 * Serve old content from cache if fetch fails for new or re-cached tiles.
 */
self.addEventListener('fetch', function (event) {
	let request = event.request;

	// if the cached response does not exist perform a fetch for that resource
	// and return the response, otherwise return response from cache
	let queriedCache = function (response) {
		// if valid cache and has not expired
		if (isValid(response)) return response;

		return fetch(request)
			.then(function (response) {
				let copy = response.clone();
				// cache response if request was successful and from the tileserver
				if (response.status === 200) {
					event.waitUntil(caches.open(tileDomain).then(function (cache) {
						// append fetch date header to cached response to check expiry later
						let headers = new Headers(copy.headers);
						headers.append('sw-fetched-on',
							new Date().getTime().toString());

						// store response in cache keyed by original request
						return copy.blob().then(function (body) {
							return cache.put(request, new Response(body, {
								status: copy.status,
								statusText: copy.statusText,
								headers: headers
							}));
						});
					}));
				}
				return response;
			}, unableToResolve)
			.catch(function (error) {
				console.error(error);
				return caches.match(request)
					.then(function (response) {
						return response || unableToResolve();
					});
			});
	};

	// only cache requests from tile server
	if (request.url.indexOf(tileDomain) === -1) {
		event.respondWith(fetch(request));
	} else {
		event.respondWith(caches.match(request)
			.then(queriedCache));
	}
});
