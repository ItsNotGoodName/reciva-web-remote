package http

import (
	"context"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"github.com/ItsNotGoodName/reciva-web-remote/http/ws"
	"github.com/ItsNotGoodName/reciva-web-remote/internal"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/build"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/hub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/model"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/store"
)

//	@title			Reciva Web Remote
//	@version		1.0
//	@description	Control your legacy Reciva based internet radios (Crane, Grace Digital, Tangent, etc.) via web browser or REST API.
//	@BasePath		/api

type HTTPError struct {
	Message string `json:"message" validate:"required"`
}

type API struct {
	Hub        *hub.Hub
	Discoverer *radio.Discoverer
	Store      store.Store
}

func NewAPI(hub *hub.Hub, Discoverer *radio.Discoverer, store store.Store) API {
	return API{
		Hub:        hub,
		Discoverer: Discoverer,
		Store:      store,
	}
}

func (a API) RadioMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		r, err := a.Hub.Get(c.Param("uuid"))
		if err != nil {
			if err == internal.ErrRadioNotFound {
				return echo.ErrNotFound.WithInternal(err)
			}
			return err
		}
		cc := &RadioContext{Radio: r, Context: c}
		return next(cc)
	}
}

//	@Summary	Get build
//	@Tags		build
//	@Produce	json
//	@Success	200	{object}	model.Build
//	@Router		/build [get]
func (a API) GetBuild(c echo.Context) error {
	return c.JSON(http.StatusOK, build.CurrentBuild)
}

//	@Summary	List presets
//	@Tags		presets
//	@Produce	json
//	@Success	200	{array}		model.Preset
//	@Failure	500	{object}	HTTPError
//	@Router		/presets [get]
func (a API) ListPresets(c echo.Context) error {
	presets, err := a.Store.ListPresets(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, presets)
}

//	@Summary	Get preset
//	@Tags		presets
//	@Param		url	path	string	true	"Preset URL"
//	@Produce	json
//	@Success	200	{object}	model.Preset
//	@Failure	404	{object}	HTTPError
//	@Failure	500	{object}	HTTPError
//	@Router		/presets/{url} [get]
func (a API) GetPreset(c echo.Context) error {
	url, err := url.QueryUnescape(c.Param("*"))
	if err != nil {
		return err
	}

	preset, err := a.Store.GetPreset(c.Request().Context(), url)
	if err != nil {
		if err == internal.ErrPresetNotFound {
			return echo.ErrNotFound.WithInternal(err)
		}
		return err
	}

	return c.JSON(http.StatusOK, preset)
}

//	@Summary	Update preset
//	@Tags		presets
//	@Param		preset	body	model.Preset	true	"Preset"
//	@Produce	json
//	@Success	200
//	@Failure	400	{object}	HTTPError
//	@Failure	404	{object}	HTTPError
//	@Failure	500	{object}	HTTPError
//	@Router		/presets [post]
func (a API) UpdatePreset(c echo.Context) error {
	var preset model.Preset
	if err := c.Bind(&preset); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := preset.Validate(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := a.Store.UpdatePreset(c.Request().Context(), &preset); err != nil {
		if err == internal.ErrPresetNotFound {
			return echo.ErrNotFound.WithInternal(err)
		}
		return err
	}

	return nil
}

//	@Summary	Discover radios
//	@Tags		radios
//	@Produce	json
//	@Success	200
//	@Failure	409	{object}	HTTPError	"Discovery already in progress"
//	@Failure	500	{object}	HTTPError
//	@Router		/radios [post]
func (a API) DiscoverRadios(c echo.Context) error {
	if err := a.Discoverer.Discover(true); err != nil {
		if err == internal.ErrDiscovering {
			return echo.ErrConflict.WithInternal(err)
		}
		return err
	}

	return nil
}

//	@Summary	List radios
//	@Tags		radios
//	@Produce	json
//	@Success	200	{array}	model.Radio
//	@Router		/radios [get]
func (a API) ListRadios(c echo.Context) error {
	return c.JSON(http.StatusOK, radio.ListRadios(a.Hub))
}

//	@Summary	Get radio
//	@Tags		radios
//	@Param		uuid	path	string	true	"Radio UUID"
//	@Produce	json
//	@Success	200	{object}	model.Radio
//	@Failure	404	{object}	HTTPError
//	@Failure	500	{object}	HTTPError
//	@Router		/radios/{uuid} [get]
func (a API) GetRadio(c echo.Context) error {
	cc := c.(*RadioContext)
	return c.JSON(http.StatusOK, model.NewRadio(cc.Radio))
}

//	@Summary	Delete radio
//	@Tags		radios
//	@Param		uuid	path	string	true	"Radio UUID"
//	@Produce	json
//	@Success	200
//	@Failure	404	{object}	HTTPError
//	@Failure	500	{object}	HTTPError
//	@Router		/radios/{uuid} [delete]
func (a API) DeleteRadio(c echo.Context) error {
	cc := c.(*RadioContext)

	if err := radio.DeleteRadio(a.Hub, cc.Radio); err != nil {
		if err == internal.ErrRadioNotFound {
			return echo.ErrNotFound.WithInternal(err)
		}
		return err
	}

	return nil
}

//	@Summary	Refresh radio volume
//	@Tags		radios
//	@Param		uuid	path	string	true	"Radio UUID"
//	@Produce	json
//	@Success	200
//	@Failure	404	{object}	HTTPError
//	@Failure	500	{object}	HTTPError
//	@Router		/radios/{uuid}/volume [post]
func (a API) RefreshRadioVolume(c echo.Context) error {
	cc := c.(*RadioContext)
	return radio.RefreshVolume(cc.Request().Context(), cc.Radio)
}

//	@Summary	Refresh radio subscription
//	@Tags		radios
//	@Param		uuid	path	string	true	"Radio UUID"
//	@Produce	json
//	@Success	200
//	@Failure	404	{object}	HTTPError
//	@Failure	500	{object}	HTTPError
//	@Router		/radios/{uuid}/subscription [post]
func (a API) RefreshRadioSubscription(c echo.Context) error {
	cc := c.(*RadioContext)
	return radio.RefreshSubscription(cc.Request().Context(), cc.Radio)
}

//	@Summary	List states
//	@Tags		states
//	@Produce	json
//	@Success	200	{array}	state.State
//	@Router		/states [get]
func (a API) ListStates(c echo.Context) error {
	return c.JSON(http.StatusOK, radio.ListStates(c.Request().Context(), a.Hub))
}

//	@Summary	Get state
//	@Tags		states
//	@Param		uuid	path	string	true	"Radio UUID"
//	@Produce	json
//	@Success	200	{object}	state.State
//	@Failure	404	{object}	HTTPError
//	@Failure	500	{object}	HTTPError
//	@Router		/states/{uuid} [get]
func (a API) GetState(c echo.Context) error {
	cc := c.(*RadioContext)
	state, err := cc.Radio.State(cc.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, state)
}

type PostState struct {
	Power       *bool   `json:"power,omitempty" validate:"optional"`
	AudioSource *string `json:"audio_source,omitempty" validate:"optional"`
	Preset      *int    `json:"preset,omitempty" validate:"optional"`
	Volume      *int    `json:"volume,omitempty" validate:"optional"`
	VolumeDelta *int    `json:"volume_delta,omitempty" validate:"optional"`
}

//	@Summary	Update state
//	@Tags		states
//	@Param		uuid	path	string		true	"Radio UUID"
//	@Param		state	body	PostState	true	"Patch state"
//	@Success	200
//	@Failure	404	{object}	HTTPError
//	@Failure	500	{object}	HTTPError
//	@Router		/states/{uuid} [post]
func (a API) PostState(c echo.Context) error {
	var req PostState
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	cc := c.(*RadioContext)
	ctx := cc.Request().Context()
	rd := cc.Radio

	if req.Power != nil {
		if err := radio.SetPower(ctx, rd, *req.Power); err != nil {
			return err
		}
	}

	if req.AudioSource != nil {
		if err := radio.SetAudioSource(ctx, rd, *req.AudioSource); err != nil {
			return err
		}
	}

	if req.Preset != nil {
		if err := radio.PlayPreset(ctx, rd, *req.Preset); err != nil {
			return err
		}
	}

	if req.Volume != nil || req.VolumeDelta != nil {
		if req.Volume != nil {
			if err := radio.SetVolume(ctx, rd, *req.Volume); err != nil {
				return err
			}
		} else if req.VolumeDelta != nil {
			if err := radio.ChangeVolume(ctx, rd, *req.VolumeDelta); err != nil {
				return err
			}
		}
	}

	return nil
}

//	@Summary	WebSocket
//	@Tags		websocket
//	@Param		Command	body		ws.Command	false	"Command"
//	@Param		Event	body		ws.Event	false	"Event"
//	@Success	200		{object}	pubsub.Stale
func (a API) WS(upgrader *websocket.Upgrader) echo.HandlerFunc {
	return func(c echo.Context) error {
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		go ws.Handle(context.Background(), conn, a.Hub, a.Discoverer)

		return nil
	}
}
