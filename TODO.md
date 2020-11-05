rabbitmq:
- //p.changeConnection(ctx, conn, ch)
  
documentation
tests
when we get a cancel we need to defer it. check where this is relevant

organize the utils
- functions in core.url should be moved to the encoding_decoding

Handle this:
return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		// Instance will look something like this "http://localhost:8080"

		breakermw := tlecb.NewDefaultHystrixCommandMiddleware("get_user_by_id")
		limitermw := tlelimiter.NewDefaultErrorLimitterMiddleware()
		fallbackmw := tlefallback.NewFallbackMiddleware(getDefaultUser())

		tgt, _ := url.Parse(instance) // e.g. parse http://localhost:8080"
		tgt.Path = path

		// TODO: need to create a client with defaults...
		endpoint := httptransport.NewClient(
			"GET",
			tgt,
			encodeGetUserByIDRequest,
			decodeGetUserByIDResponse,
			tlectxhttp.WriteBefore()).Endpoint()

		endpoint = breakermw(endpoint)
		endpoint = limitermw(endpoint)
		// fallback should run last. if circuit was opened, it will return the response from the fallback
		endpoint = fallbackmw(endpoint)

		return endpoint, nil, nil
	}