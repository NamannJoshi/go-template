package router

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"template/router/handlers"
	"template/services"
	"template/utils"
	"template/utils/auth"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

func InitRouter(svc services.Services, client *redis.Client) {
	router := mux.NewRouter()

	router.Use(RateLimit(client))

	handlers := handlers.RouteBundler(svc)
	authRouterHandler(router, handlers)
	userRouterHandler(router, handlers)

	http.ListenAndServe("127.0.0.1:8080", router)
}

func authRouterHandler(r *mux.Router, handlers *handlers.RouteHandlers) {
	r.HandleFunc("/login", handlers.AuthRoutes.Login)
	r.HandleFunc("/register", handlers.AuthRoutes.Register)
	r.HandleFunc("/refresh", handlers.AuthRoutes.RefreshAccessToken)
}

func userRouterHandler(r *mux.Router, handlers *handlers.RouteHandlers) {

	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.Use(AuthMiddleware)
	userRouter.Use(RequireRole("User"))

	userRouter.HandleFunc("/{id}", handlers.UserRoutes.UserHandlerById).
		Methods("GET", "PUT")

	adminRouter := r.PathPrefix("/admin/user").Subrouter()
	adminRouter.Use(AuthMiddleware)
	adminRouter.Use(RequireRole("Admin"))
	adminRouter.HandleFunc("", handlers.UserRoutes.UserHandler).Methods("GET", "POST")
	adminRouter.HandleFunc("/{id}", handlers.UserRoutes.UserHandlerById).Methods("PUT", "DELETE")
}

func protectedHandler(f http.HandlerFunc, requiredRole string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("access_token")
		if accessToken == "" {
			utils.WriteJSON(w, http.StatusBadRequest, fmt.Errorf("missing access token header."))
			return
		}

		claims, err := auth.VerfiyAccessToken(accessToken)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, fmt.Errorf("access token is not valid.."))
			return
		}

		roles := strings.Split(claims["roles"].(string), ",")

		for _, role := range roles {
			if role == "Admin" || role == requiredRole {
				f(w, r)
				return
			}
		}

		utils.WriteJSON(w, http.StatusForbidden, fmt.Errorf(""))
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("access_token")
		if accessToken == "" {
			utils.WriteJSON(w, http.StatusUnauthorized, fmt.Errorf("missing access token"))
			return
		}

		claims, err := auth.VerfiyAccessToken(accessToken)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, fmt.Errorf("invalid access token"))
			return
		}

		// attach claims to context
		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRole(requiredRole string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value("claims").(jwt.MapClaims)

			rawRoles := claims["roles"].(string)
			roles := strings.Split(rawRoles, ",")

			for _, role := range roles {
				if role == requiredRole {
					next.ServeHTTP(w, r)
					return
				}
			}
			fmt.Println(rawRoles, requiredRole)

			utils.WriteJSON(w, http.StatusForbidden, fmt.Errorf("forbidden"))
		})
	}
}

func RateLimit(client *redis.Client) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getClientIP(r)
			isRequestAllowed, err := allowRequest(r.Context(), client, ip)
			if err != nil {
				utils.WriteJSON(w, http.StatusInternalServerError, err)
				return
			}

			if !isRequestAllowed {
				utils.WriteJSON(w, http.StatusTooManyRequests, nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

}

func getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	return r.RemoteAddr
}

func allowRequest(context context.Context, client *redis.Client, ip string) (bool, error) {
	fmt.Println(ip)
	limit := int64(1)
	window := 5 * time.Second
	count, err := client.Incr(context, ip).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		client.Expire(context, ip, window)
	}

	if count > limit {
		return false, nil
	}

	return true, nil
}
