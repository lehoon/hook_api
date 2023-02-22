package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/lehoon/hook_api/v2/api"
)

type RouteInfo struct {
	Method string
	Path   string
}

var routes = []RouteInfo{}

func PushRoute(method, path string) {
	routeInfo := RouteInfo{
		Method: method,
		Path:   path,
	}

	routes = append(routes, routeInfo)
}

func GetRoutes(w http.ResponseWriter, r *http.Request) {
	render.Respond(w, r, api.SuccessBizResultWithData(routes))
}

func Routes() http.Handler {
	r := chi.NewRouter()

	r.Route("/stream", func(r chi.Router) {
		r.Get("/", api.StreamList)
		r.Post("/not_found", api.StreamNotFound)
		r.Post("/change", api.StreamChanged)
		r.Post("/none_reader", api.StreamNoneReader)
		r.Get("/isonline/{streamId}", api.StreamIsOnline)
		r.Get("/open/{streamId}", api.StreamOpen)
		r.Get("/close/{streamId}", api.StreamClose)
		r.Get("/play/url/{streamId}", api.StreamPlayUrl)
	})

	r.Route("/record", func(r chi.Router) {
		r.Post("/ts_finish", api.RecordTsFinish)
		r.Post("/mp4_finish", api.RecordMP4Finish)
	})

	r.Route("/service", func(r chi.Router) {
		r.Post("/startup_report", api.ServerStartupReport)
		r.Post("/keepalive_report", api.KeepAliveReport)
		r.Post("/rtp_close_report", api.RtpCloseReport)
		r.Post("/rtp_timeout_report", api.RtpTimeoutReport)
		r.Post("/flow_report", api.FlowReport)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/http_file", api.AuthHttpFile)
		r.Post("/stream_play", api.AuthPlay)
		r.Post("/stream_publish", api.AuthPublish)
		r.Post("/rtsp_play", api.AuthRtspPlay)
		r.Post("/shell", api.AuthShell)
		r.Post("/rtsp_auth", api.IsRtspAuth)
	})

	r.Route("/device", func(r chi.Router) {
		r.Post("/", api.PublishDevice)
		r.Put("/", api.UpdateDevice)
		r.Get("/", api.DeviceList)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", api.QueryDeviceInfo)
			r.Delete("/", api.DeleteDevice)
		})
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/exception", api.QueryOperateCodeMessage)
		r.Get("/", GetRoutes)
		r.Post("/notify", api.ShowPostMessage)
	})

	return r
}
