package di

import (
	"context"
	"net/http"

	"go.uber.org/fx"
	"vk-search/internal/api"
	"vk-search/internal/api/handlers"
	"vk-search/internal/app/auth"
	"vk-search/internal/infrastructure/config"
	"vk-search/internal/infrastructure/mocks"
)

func BuildApp() *fx.App {
	return fx.New(
		fx.Provide(
			config.Load,
			mocks.NewUserMockRepository,
			fx.Annotate(
				func(cfg *config.Config) auth.TokenConfig { return cfg },
			),
			auth.NewAuthUseCase,
			handlers.NewAuthHandler,
			api.NewRouter,
		),
		fx.Invoke(func(lc fx.Lifecycle, handler http.Handler, cfg *config.Config) {
			srv := &http.Server{
				Addr:    ":" + cfg.GetHTTPPort(),
				Handler: handler,
			}
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go srv.ListenAndServe()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return srv.Shutdown(ctx)
				},
			})
		}),
	)
}