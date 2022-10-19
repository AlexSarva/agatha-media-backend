package handlers

import (
	"AlexSarva/media/admin"
	"AlexSarva/media/internal/app"
	"AlexSarva/media/models"
	"AlexSarva/media/storage/storagepg"
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// messageResponse additional respond generator
// useful in case of error handling in outputting results to respond
func messageResponse(w http.ResponseWriter, message, ContentType string, httpStatusCode int) {
	w.Header().Set("Content-Type", ContentType)
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, jsonRespErr := json.Marshal(resp)
	if jsonRespErr != nil {
		log.Println(jsonRespErr)
	}
	w.Write(jsonResp)
}

// readBodyBytes compressed request processing function
func readBodyBytes(r *http.Request) (io.ReadCloser, error) {
	// GZIP decode
	if len(r.Header["Content-Encoding"]) > 0 && r.Header["Content-Encoding"][0] == "gzip" {
		// Read body
		bodyBytes, readErr := io.ReadAll(r.Body)
		if readErr != nil {
			return nil, readErr
		}
		defer r.Body.Close()

		newR, gzErr := gzip.NewReader(io.NopCloser(bytes.NewBuffer(bodyBytes)))
		if gzErr != nil {
			log.Println(gzErr)
			return nil, gzErr
		}
		defer newR.Close()

		return newR, nil
	} else {
		return r.Body, nil
	}
}

// gzipContentTypes request types that support data compression
var gzipContentTypes = "application/x-gzip, application/javascript, application/json, text/css, text/html, text/plain, text/xml"

func GetSearch(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%+v\n", r.Header)
		headerContentType := r.Header.Get("Content-Type")
		if !strings.Contains("application/json, application/x-gzip", headerContentType) {
			messageResponse(w, "Content Type is not application/json or application/x-gzip", "application/json", http.StatusBadRequest)
			return
		}

		var query models.SearchQuery
		var unmarshalErr *json.UnmarshalTypeError

		b, err := readBodyBytes(r)
		if err != nil {
			messageResponse(w, "Problem in body", "application/json", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(b)
		decoder.DisallowUnknownFields()
		errDecode := decoder.Decode(&query)

		if errDecode != nil {
			if errors.As(errDecode, &unmarshalErr) {
				messageResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, "application/json", http.StatusBadRequest)
			} else {
				messageResponse(w, "Bad Request. "+errDecode.Error(), "application/json", http.StatusBadRequest)
			}
			return
		}

		searchRes, searchResErr := database.Repo.GetSearch(query.Text)
		if searchResErr != nil {
			if searchResErr == admin.ErrNoValues {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNoContent)
				return
			}
			messageResponse(w, "Internal Server Error: "+searchResErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		res, resErr := json.Marshal(searchRes)
		if resErr != nil {
			panic(resErr)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}

func GetGraph(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%+v\n", r.Header)
		headerContentType := r.Header.Get("Content-Type")
		if !strings.Contains("application/json, application/x-gzip", headerContentType) {
			messageResponse(w, "Content Type is not application/json or application/x-gzip", "application/json", http.StatusBadRequest)
			return
		}

		var query models.GraphQuery
		var unmarshalErr *json.UnmarshalTypeError

		b, err := readBodyBytes(r)
		if err != nil {
			messageResponse(w, "Problem in body", "application/json", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(b)
		decoder.DisallowUnknownFields()
		errDecode := decoder.Decode(&query)

		if errDecode != nil {
			if errors.As(errDecode, &unmarshalErr) {
				messageResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, "application/json", http.StatusBadRequest)
			} else {
				messageResponse(w, "Bad Request. "+errDecode.Error(), "application/json", http.StatusBadRequest)
			}
			return
		}

		graphInfo, graphInfoErr := database.Repo.GetGraphByURL(query.Query)
		if graphInfoErr != nil {
			if graphInfoErr == admin.ErrNoValues {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNoContent)
				return
			}
			messageResponse(w, "Internal Server Error: "+graphInfoErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		graphRes, graphResErr := json.Marshal(graphInfo)
		if graphResErr != nil {
			panic(graphResErr)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(graphRes)
	}
}

func GetGraphByID(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%+v\n", r.Header)
		headerContentType := r.Header.Get("Content-Type")
		if !strings.Contains("application/json, application/x-gzip", headerContentType) {
			messageResponse(w, "Content Type is not application/json or application/x-gzip", "application/json", http.StatusBadRequest)
			return
		}

		var query models.GraphQueryID
		var unmarshalErr *json.UnmarshalTypeError

		b, err := readBodyBytes(r)
		if err != nil {
			messageResponse(w, "Problem in body", "application/json", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(b)
		decoder.DisallowUnknownFields()
		errDecode := decoder.Decode(&query)

		if errDecode != nil {
			if errors.As(errDecode, &unmarshalErr) {
				messageResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, "application/json", http.StatusBadRequest)
			} else {
				messageResponse(w, "Bad Request. "+errDecode.Error(), "application/json", http.StatusBadRequest)
			}
			return
		}

		graphInfo, graphInfoErr := database.Repo.GetGraphByID(query.ID)
		if graphInfoErr != nil {
			if graphInfoErr == admin.ErrNoValues {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNoContent)
				return
			}
			messageResponse(w, "Internal Server Error: "+graphInfoErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		graphRes, graphResErr := json.Marshal(graphInfo)
		if graphResErr != nil {
			panic(graphResErr)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(graphRes)
	}
}

func GetSourceByURL(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%+v\n", r.Header)
		headerContentType := r.Header.Get("Content-Type")
		if !strings.Contains("application/json, application/x-gzip", headerContentType) {
			messageResponse(w, "Content Type is not application/json or application/x-gzip", "application/json", http.StatusBadRequest)
			return
		}

		var query models.GraphQuery
		var unmarshalErr *json.UnmarshalTypeError

		b, err := readBodyBytes(r)
		if err != nil {
			messageResponse(w, "Problem in body", "application/json", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(b)
		decoder.DisallowUnknownFields()
		errDecode := decoder.Decode(&query)

		if errDecode != nil {
			if errors.As(errDecode, &unmarshalErr) {
				messageResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, "application/json", http.StatusBadRequest)
			} else {
				messageResponse(w, "Bad Request. "+errDecode.Error(), "application/json", http.StatusBadRequest)
			}
			return
		}

		sourceInfo, sourceInfoErr := database.Repo.GetSourceInfoByURL(query.Query)
		if sourceInfoErr != nil {
			if sourceInfoErr == admin.ErrNoValues {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNoContent)
				return
			}
			messageResponse(w, "Internal Server Error: "+sourceInfoErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		graphRes, graphResErr := json.Marshal(sourceInfo)
		if graphResErr != nil {
			panic(graphResErr)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(graphRes)
	}
}

func GetSourceByID(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%+v\n", r.Header)
		headerContentType := r.Header.Get("Content-Type")
		if !strings.Contains("application/json, application/x-gzip", headerContentType) {
			messageResponse(w, "Content Type is not application/json or application/x-gzip", "application/json", http.StatusBadRequest)
			return
		}

		var query models.GraphQueryID
		var unmarshalErr *json.UnmarshalTypeError

		b, err := readBodyBytes(r)
		if err != nil {
			messageResponse(w, "Problem in body", "application/json", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(b)
		decoder.DisallowUnknownFields()
		errDecode := decoder.Decode(&query)

		if errDecode != nil {
			if errors.As(errDecode, &unmarshalErr) {
				messageResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, "application/json", http.StatusBadRequest)
			} else {
				messageResponse(w, "Bad Request. "+errDecode.Error(), "application/json", http.StatusBadRequest)
			}
			return
		}

		sourceInfo, sourceInfoErr := database.Repo.GetSourceInfoByID(query.ID)
		if sourceInfoErr != nil {
			if sourceInfoErr == admin.ErrNoValues {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNoContent)
				return
			}
			messageResponse(w, "Internal Server Error: "+sourceInfoErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		graphRes, graphResErr := json.Marshal(sourceInfo)
		if graphResErr != nil {
			panic(graphResErr)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(graphRes)
	}
}

func GetFullGraph(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerContentType := r.Header.Get("Content-Length")
		if len(headerContentType) != 0 {
			messageResponse(w, "Content-Length is not equal 0", "application/json", http.StatusBadRequest)
			return
		}

		graph, graphErr := database.Repo.GetFullGraph()
		if graphErr != nil {
			if errors.Is(graphErr, sql.ErrNoRows) {
				messageResponse(w, "user doesnt exist", "application/json", http.StatusUnauthorized)
				return
			}
			messageResponse(w, "Internal Server Error: "+graphErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		jsonResp, _ := json.Marshal(graph)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
	}
}

func AddNewGraph(database *app.Database, adminDB *admin.PostgresDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%+v\n", r.Header)
		headerContentType := r.Header.Get("Content-Type")
		if !strings.Contains("application/json, application/x-gzip", headerContentType) {
			messageResponse(w, "Content Type is not application/json or application/x-gzip", "application/json", http.StatusBadRequest)
			return
		}

		// Проверка авторизации по токену
		userID, tokenErr := GetToken(r)
		if tokenErr != nil {
			messageResponse(w, "User unauthorized: "+tokenErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		userInfo, userInfoErr := adminDB.GetUserInfo(userID)
		if userInfoErr != nil {
			if errors.Is(userInfoErr, sql.ErrNoRows) {
				messageResponse(w, "user doesnt exist", "application/json", http.StatusUnauthorized)
				return
			}
			messageResponse(w, "Internal Server Error: "+userInfoErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		userInfo.Type = "Bearer"
		log.Printf("USER: %+v\n", userInfo)

		//defer r.Body.Close()
		//bodyBytes, err := io.ReadAll(r.Body)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//bodyString := string(bodyBytes)
		//log.Println(bodyString)

		var newGraph models.NewGraph
		var unmarshalErr *json.UnmarshalTypeError

		b, err := readBodyBytes(r)
		if err != nil {
			messageResponse(w, "Problem in body", "application/json", http.StatusBadRequest)
			return
		}

		//log.Println(string(b))

		decoder := json.NewDecoder(b)
		decoder.DisallowUnknownFields()
		errDecode := decoder.Decode(&newGraph)

		if errDecode != nil {
			if errors.As(errDecode, &unmarshalErr) {
				messageResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, "application/json", http.StatusBadRequest)
			} else {
				messageResponse(w, "Bad Request. "+errDecode.Error(), "application/json", http.StatusBadRequest)
			}
			return
		}

		newGraph.UserID = userID
		newGraph.Cnt = len(newGraph.Sources)

		resp, respErr := database.Repo.AddNewGraph(newGraph)
		if respErr != nil {
			if errors.As(storagepg.ErrDuplicatePK, &respErr) {
				messageResponse(w, "GraphID already exists", "application/json", http.StatusConflict)
				return
			} else {
				messageResponse(w, "Internal Server Error: "+respErr.Error(), "application/json", http.StatusInternalServerError)
				return
				log.Println(respErr)
			}

		}

		log.Printf("NEWDATA: %+v\n", newGraph)

		ordersList, ordersListErr := json.Marshal(resp)
		if ordersListErr != nil {
			panic(ordersListErr)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(ordersList)
	}
}

func GetGraphCards(database *app.Database, adminDB *admin.PostgresDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerContentType := r.Header.Get("Content-Length")
		if len(headerContentType) != 0 {
			messageResponse(w, "Content-Length is not equal 0", "application/json", http.StatusBadRequest)
			return
		}

		// Проверка авторизации по токену
		userID, tokenErr := GetToken(r)
		if tokenErr != nil {
			messageResponse(w, "User unauthorized: "+tokenErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		userInfo, userInfoErr := adminDB.GetUserInfo(userID)
		if userInfoErr != nil {
			if errors.Is(userInfoErr, sql.ErrNoRows) {
				messageResponse(w, "user doesnt exist", "application/json", http.StatusUnauthorized)
				return
			}
			messageResponse(w, "Internal Server Error: "+userInfoErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		userInfo.Type = "Bearer"

		graphCards, graphCardsErr := database.Repo.GetGraphCards(userID)
		if graphCardsErr != nil {
			if errors.Is(graphCardsErr, sql.ErrNoRows) {
				messageResponse(w, "no graphs exist", "application/json", http.StatusUnauthorized)
				return
			}
			messageResponse(w, "Internal Server Error: "+graphCardsErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		jsonResp, _ := json.Marshal(graphCards)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
	}
}

func DeleteGraphCard(database *app.Database, adminDB *admin.PostgresDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%+v\n", r.Header)
		headerContentType := r.Header.Get("Content-Type")
		if !strings.Contains("application/json, application/x-gzip", headerContentType) {
			messageResponse(w, "Content Type is not application/json or application/x-gzip", "application/json", http.StatusBadRequest)
			return
		}

		// Проверка авторизации по токену
		userID, tokenErr := GetToken(r)
		if tokenErr != nil {
			messageResponse(w, "User unauthorized: "+tokenErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		userInfo, userInfoErr := adminDB.GetUserInfo(userID)
		if userInfoErr != nil {
			if errors.Is(userInfoErr, sql.ErrNoRows) {
				messageResponse(w, "user doesnt exist", "application/json", http.StatusUnauthorized)
				return
			}
			messageResponse(w, "Internal Server Error: "+userInfoErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		userInfo.Type = "Bearer"
		log.Printf("USER: %+v\n", userInfo)

		var graphDel models.GraphDel
		var unmarshalErr *json.UnmarshalTypeError

		b, err := readBodyBytes(r)
		if err != nil {
			messageResponse(w, "Problem in body", "application/json", http.StatusBadRequest)
			return
		}

		//log.Println(string(b))

		decoder := json.NewDecoder(b)
		decoder.DisallowUnknownFields()
		errDecode := decoder.Decode(&graphDel)

		if errDecode != nil {
			if errors.As(errDecode, &unmarshalErr) {
				messageResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, "application/json", http.StatusBadRequest)
			} else {
				messageResponse(w, "Bad Request. "+errDecode.Error(), "application/json", http.StatusBadRequest)
			}
			return
		}

		resp, respErr := database.Repo.DeleteGraphCard(userID, graphDel.GraphID)
		if respErr != nil {
			if errors.As(storagepg.ErrNoData, &respErr) {
				messageResponse(w, storagepg.ErrNoData.Error(), "application/json", http.StatusConflict)
				return
			} else {
				messageResponse(w, "Internal Server Error: "+respErr.Error(), "application/json", http.StatusInternalServerError)
				return
				log.Println(respErr)
			}

		}

		graphCardList, graphCardListErr := json.Marshal(resp)
		if graphCardListErr != nil {
			panic(graphCardListErr)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		w.Write(graphCardList)
	}
}

func GetGraphByUUID(database *app.Database, adminDB *admin.PostgresDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%+v\n", r.Header)
		headerContentType := r.Header.Get("Content-Type")
		if !strings.Contains("application/json, application/x-gzip", headerContentType) {
			messageResponse(w, "Content Type is not application/json or application/x-gzip", "application/json", http.StatusBadRequest)
			return
		}

		// Проверка авторизации по токену
		userID, tokenErr := GetToken(r)
		if tokenErr != nil {
			messageResponse(w, "User unauthorized: "+tokenErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		userInfo, userInfoErr := adminDB.GetUserInfo(userID)
		if userInfoErr != nil {
			if errors.Is(userInfoErr, sql.ErrNoRows) {
				messageResponse(w, "user doesnt exist", "application/json", http.StatusUnauthorized)
				return
			}
			messageResponse(w, "Internal Server Error: "+userInfoErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		userInfo.Type = "Bearer"
		log.Printf("USER: %+v\n", userInfo)

		var query models.GraphUUID
		var unmarshalErr *json.UnmarshalTypeError

		b, err := readBodyBytes(r)
		if err != nil {
			messageResponse(w, "Problem in body", "application/json", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(b)
		decoder.DisallowUnknownFields()
		errDecode := decoder.Decode(&query)

		if errDecode != nil {
			if errors.As(errDecode, &unmarshalErr) {
				messageResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, "application/json", http.StatusBadRequest)
			} else {
				messageResponse(w, "Bad Request. "+errDecode.Error(), "application/json", http.StatusBadRequest)
			}
			return
		}

		graphInfo, graphInfoErr := database.Repo.GetGraphByUUID(query.GraphID)
		if graphInfoErr != nil {
			if graphInfoErr == admin.ErrNoValues {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNoContent)
				return
			}
			messageResponse(w, "Internal Server Error: "+graphInfoErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		graphRes, graphResErr := json.Marshal(graphInfo)
		if graphResErr != nil {
			panic(graphResErr)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(graphRes)
	}
}

func MyAllowOriginFunc(r *http.Request, origin string) bool {
	if origin == "http://localhost:3000" || origin == "http://10.2.3.197:3000" {
		return true
	}
	return false
}

// MyHandler - the main handler of the server
// contains middlewares and all routes
func MyHandler(database *app.Database, adminDatabase *admin.PostgresDB) *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowOriginFunc:  MyAllowOriginFunc,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentEncoding("gzip"))
	r.Use(middleware.AllowContentType("application/json", "text/plain", "application/x-gzip"))
	r.Use(middleware.Compress(5, gzipContentTypes))
	r.Mount("/debug", middleware.Profiler())
	r.Get("/api/graph", GetFullGraph(database))
	r.Post("/api/graph/url", GetGraph(database))
	r.Post("/api/graph/id", GetGraphByID(database))
	r.Post("/api/graph/new", AddNewGraph(database, adminDatabase))
	r.Get("/api/graph/all", GetGraphCards(database, adminDatabase))
	r.Post("/api/graph/uuid", GetGraphByUUID(database, adminDatabase))
	r.Delete("/api/graph/del", DeleteGraphCard(database, adminDatabase))
	r.Post("/api/source/url", GetSourceByURL(database))
	r.Post("/api/source/id", GetSourceByID(database))
	//
	r.Post("/api/user/register", UserRegistration(adminDatabase))
	r.Post("/api/user/login", UserAuthentication(adminDatabase))
	r.Get("/api/users/me", GetUserInfo(adminDatabase))
	r.Post("/api/search", GetSearch(database))
	//r.Get("/api/user/orders", GetOrders(database))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, nfErr := w.Write([]byte("route does not exist"))
		if nfErr != nil {
			log.Println(nfErr)
		}
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, naErr := w.Write([]byte("sorry, only GET and POST methods are supported."))
		if naErr != nil {
			log.Println(naErr)
		}
	})
	return r
}
