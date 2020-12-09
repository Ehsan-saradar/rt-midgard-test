// Package api provides the HTTP interface.
package api

import (
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/julienschmidt/httprouter"
	"github.com/pascaldekloe/metrics"

	"gitlab.com/thorchain/midgard/internal/graphql"
	"gitlab.com/thorchain/midgard/internal/graphql/generated"
	"gitlab.com/thorchain/midgard/internal/util/timer"
)

// Handler serves the entire API.
var Handler http.Handler

func addMeasuredFunc(router *httprouter.Router, url string, handler http.HandlerFunc) {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic("Bad constant url regex.")
	}
	simplifiedUrl := reg.ReplaceAllString(url, "_")
	t := timer.NewMilli("serving" + simplifiedUrl)

	router.HandlerFunc(http.MethodGet, url, func(w http.ResponseWriter, r *http.Request) {
		m := t.One()
		handler(w, r)
		m()
	})
}

func InitHandler(nodeURL string, proxiedWhitelistedEndpoints []string) {
	var router = httprouter.New()
	Handler = router

	// apply some navigation pointers
	router.HandleMethodNotAllowed = true
	router.HandleOPTIONS = true
	router.HandlerFunc(http.MethodGet, "/", serveRoot)

	router.HandlerFunc(http.MethodGet, "/v2/debug/metrics", metrics.ServeHTTP)
	router.HandlerFunc(http.MethodGet, "/v2/debug/timers", timer.ServeHTTP)

	for _, endpoint := range proxiedWhitelistedEndpoints {
		midgardPath := "/v2/thorchain/" + endpoint
		addMeasuredFunc(router, midgardPath, proxiedEndpointHandlerFunc(nodeURL))
	}

	router.HandlerFunc(http.MethodGet, "/v2/doc", serveDoc)

	// version 1
	addMeasuredFunc(router, "/v2/health", jsonHealth)
	addMeasuredFunc(router, "/v2/history/swaps", jsonSwapHistory)
	addMeasuredFunc(router, "/v2/history/depths/:pool", jsonDepths)
	addMeasuredFunc(router, "/v2/history/earnings", jsonEarningsHistory)
	addMeasuredFunc(router, "/v2/history/liquidity_changes", jsonLiquidityHistory)
	addMeasuredFunc(router, "/v2/network", jsonNetwork)
	addMeasuredFunc(router, "/v2/nodes", jsonNodes)
	addMeasuredFunc(router, "/v2/pools", jsonPools)
	addMeasuredFunc(router, "/v2/pools/:pool", jsonPool)
	addMeasuredFunc(router, "/v2/members", jsonMembers)
	addMeasuredFunc(router, "/v2/members/:addr", jsonMemberDetails)
	addMeasuredFunc(router, "/v2/stats", jsonStats)
	addMeasuredFunc(router, "/v2/swagger.json", jsonSwagger)
	addMeasuredFunc(router, "/v2/tx", jsonTx)

	// version 2 with GraphQL
	router.HandlerFunc(http.MethodGet, "/v2/graphql", playground.Handler("Midgard Playground", "/v2"))
	router.Handle(http.MethodPost, "/v2", serverV2())
}

func serveDoc(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./openapi/generated/doc.html")
}

func serverV2() httprouter.Handle {
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graphql.Resolver{}}))
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		h.ServeHTTP(w, req)
	}
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain;charset=UTF-8")

	// Discarding errors
	_, _ = io.WriteString(w, `# THORChain Midgard

Welcome to the HTTP interface.
`)
}

func proxiedEndpointHandlerFunc(nodeURL string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// NOTE(elfedy): url may come with or without leading slash, so make sure we handle this
		// regardless
		// Path is the same without leading v1 (or /v2)
		targetPath := strings.ReplaceAll(r.URL.Path, "v2/thorchain", "")
		targetPath = strings.ReplaceAll(targetPath, "//", "/")
		targetPath = strings.TrimPrefix(targetPath, "/")
		url, err := url.Parse(nodeURL + "/" + targetPath)
		if err != nil {
			http.NotFound(w, r)
		}

		proxy := httputil.NewSingleHostReverseProxy(url)
		proxy.Director = func(req *http.Request) {
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Add("X-Origin-Host", url.Host)
			req.URL.Scheme = url.Scheme
			req.URL.Host = url.Host
			req.URL.Path = url.Path
		}
		proxy.ServeHTTP(w, r)
	}
}

// CORS returns a Handler which applies CORS on h.
func CORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/v2/thorchain") {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		h.ServeHTTP(w, r)
	})
}
