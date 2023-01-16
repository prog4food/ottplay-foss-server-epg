package internal

import (
	"net/http"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

const (
	HeaderAccessControlAllowOrigin = "Access-Control-Allow-Origin"
	HeaderToken = "X-Token"
	SymWildcard = "*"
)


func CheckApiKey(keystore []string, w http.ResponseWriter, r *http.Request) bool {
	client_token := r.Header.Get(HeaderToken)
	if client_token != "" {
		for i := 0; i < len(keystore); i++ {
			if client_token == keystore[i] {
				return true
			}
		}
	}
	w.WriteHeader(http.StatusUnauthorized)
	return false
}

func addCors(w http.ResponseWriter) {
	w.Header().Add(HeaderAccessControlAllowOrigin, SymWildcard) // CORS: *
}

func ServeStatic(w http.ResponseWriter, r *http.Request) {
	addCors(w)
	
	log.Printf("%s", filepath.Join(".", r.URL.Path))
	http.ServeFile(w, r, filepath.Join(".", r.URL.Path))
}

func RouterPublic(w http.ResponseWriter, r *http.Request) {
	addCors(w)

	switch r.URL.Path {
		case "/m3u/match-logos":
			IcoMatch(w, r)
		case "/m3u/match-channels":
			EpgMatch(w, r)
		case
			"/m3u/gelist.php",
			"/m3u/geicons.php":
			//TODO: backward compability (maybe)
			http.Error(w, "Not Implemented", http.StatusNotImplemented)
		default:
			return
	}
}

func PrivateApi(w http.ResponseWriter, r *http.Request) {
	addCors(w)

	switch r.URL.Path {
		case "/api/epg/update-providers":
			if CheckApiKey(Config.AdminTokens, w, r) {
				log.Info().Msg("api: call UpdateProviders")
				SchedulerCall(P)
			}
		case "/api/reload":
			if CheckApiKey(Config.AdminTokens, w, r) {
				log.Info().Msg("api.config: reload")
				ReLoadConfig()
			}
		case "/api/provider":
			if CheckApiKey(Config.AdminTokens, w, r) {
				// WIP
			}
		// TODO: todo
		case
			"/api/provider-reload",
			"/proxy",
			"/proxy302",
			"/api/update-provider":
			http.Error(w, "Not Implemented", http.StatusNotImplemented)
		default:
			return
	}
}