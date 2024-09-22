package server

import (
	"context"
	"log"
	"net"
	"net/http"

	"hshelby-tkcled-product/src/graph/directive"
	generated_admin "hshelby-tkcled-product/src/graph/generated/admin"
	resolver_admin "hshelby-tkcled-product/src/graph/resolver/admin"
	"hshelby-tkcled-product/src/middleware"
	service_product "hshelby-tkcled-product/src/service/product"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi/v5"
)

func ServeGraph(ctx context.Context, addr string) (err error) {
	defer log.Println("HTTP server stopped", err)

	r := chi.NewRouter()
	v1(r)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	srv := http.Server{
		Addr:    addr,
		Handler: r,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	errChan := make(chan error, 1)

	go func(ctx context.Context, errChan chan error) {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}(ctx, errChan)

	log.Printf("Listen and Serve Product-Graph-Service API at: %s\n", addr)

	select {
	case <-ctx.Done():
		return nil
	case err = <-errChan:
		return err
	}
}

func v1(r chi.Router) {
	configAdmin := generated_admin.Config{Resolvers: &resolver_admin.Resolver{}}
	configAdmin.Directives = directive.AdminDirective

	srvAdmin := handler.NewDefaultServer(generated_admin.NewExecutableSchema(configAdmin))

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.AllowAll().Handler)
		r.With(middleware.Middleware()).Route("/graphql", func(r chi.Router) {
			r.Handle("/tkcled", srvAdmin)
		})

		r.With(middleware.AuthMiddleware())
		r.Route("/upload-file", func(r chi.Router) {
			r.Post("/image", service_product.UploadImage)
			r.Post("/import-product", service_product.ImportProduct)
		})
	})
}
