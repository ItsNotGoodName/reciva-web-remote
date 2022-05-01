package app

import (
	"context"
	"log"

	"github.com/ItsNotGoodName/reciva-web-remote/core/dto"
	"github.com/ItsNotGoodName/reciva-web-remote/core/preset"
	"github.com/ItsNotGoodName/reciva-web-remote/core/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
)

type App struct {
	build        dto.Build
	hubService   radio.HubService
	presetStore  preset.PresetStore
	radioService radio.RadioService
}

func New(
	build dto.Build,
	hubService radio.HubService,
	presetStore preset.PresetStore,
	radioService radio.RadioService,
) *App {
	return &App{
		build:        build,
		hubService:   hubService,
		presetStore:  presetStore,
		radioService: radioService,
	}
}

func (a *App) Build() dto.Build {
	return a.build
}

func (a *App) PresetGet(ctx context.Context, req *dto.PresetGetRequest) (*dto.PresetGetResponse, error) {
	p, err := a.presetStore.Get(ctx, req.URL)
	if err != nil {
		return nil, err
	}

	return &dto.PresetGetResponse{Preset: newPreset(p)}, nil
}

func (a *App) PresetList(ctx context.Context) (*dto.PresetListResponse, error) {
	pts, err := a.presetStore.List(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.PresetListResponse{Presets: newPresets(pts)}, nil
}

func (a *App) PresetUpdate(ctx context.Context, req *dto.PresetUpdateRequest) error {
	p, err := preset.ParsePreset(req.Preset.URL, req.Preset.TitleNew, req.Preset.URLNew)
	if err != nil {
		return err
	}

	return a.presetStore.Update(ctx, p)
}

func (a *App) RadioDiscover(req *dto.RadioDiscoverRequest) (*dto.RadioDiscoverResponse, error) {
	count, err := a.hubService.Discover(req.Force)
	if err != nil {
		return nil, err
	}

	return &dto.RadioDiscoverResponse{Count: count}, nil
}

func (a *App) RadioGet(req *dto.RadioRequest) (*dto.RadioGetResponse, error) {
	r, err := a.hubService.Get(req.UUID)
	if err != nil {
		return nil, err
	}

	return &dto.RadioGetResponse{Radio: newRadio(&r)}, nil
}

func (a *App) RadioList() (*dto.RadioListResponse, error) {
	return &dto.RadioListResponse{Radios: newRadios(a.hubService.List())}, nil
}

func (a *App) RadioRefreshSubscription(ctx context.Context, req *dto.RadioRequest) error {
	rd, err := a.hubService.Get(req.UUID)
	if err != nil {
		return err
	}

	return a.radioService.RefreshSubscription(ctx, rd)
}
func (a *App) RadioRefreshVolume(ctx context.Context, req *dto.RadioRequest) error {
	rd, err := a.hubService.Get(req.UUID)
	if err != nil {
		return err
	}

	return a.radioService.RefreshVolume(ctx, rd)
}

func (a *App) StateGet(ctx context.Context, req *dto.StateRequest) (*dto.StateGetResponse, error) {
	rd, err := a.hubService.Get(req.UUID)
	if err != nil {
		return nil, err
	}

	s, err := a.radioService.GetState(ctx, rd)
	if err != nil {
		return nil, err
	}

	return &dto.StateGetResponse{State: *s}, nil
}

func (a *App) StateList(ctx context.Context) (*dto.StateListResponse, error) {
	pts := a.hubService.List()
	var states []state.State
	for _, p := range pts {
		s, err := a.radioService.GetState(ctx, p)
		if err != nil {
			log.Println("app.App.StateList:", err)
			continue
		}

		states = append(states, *s)
	}

	return &dto.StateListResponse{States: states}, nil
}

func (a *App) StatePatch(ctx context.Context, req *dto.StatePatchRequest) error {
	rd, err := a.hubService.Get(req.UUID)
	if err != nil {
		return err
	}

	if req.Power != nil {
		if err := a.radioService.SetPower(ctx, rd, *req.Power); err != nil {
			return err
		}
	}

	if req.AudioSource != nil {
		if err := a.radioService.SetAudioSource(ctx, rd, *req.AudioSource); err != nil {
			return err
		}
	}

	if req.Preset != nil {
		if err := a.radioService.PlayPreset(ctx, rd, *req.Preset); err != nil {
			return err
		}
	}

	if req.Volume != nil {
		if err := a.radioService.SetVolume(ctx, rd, *req.Volume); err != nil {
			return err
		}
	}

	return nil
}
