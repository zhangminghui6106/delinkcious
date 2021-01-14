package service

import (
	_ "errors"
	//http传输包对于使用HTTP传输的服务是必须的
	httptransport "github.com/go-kit/kit/transport/http"
	//提供路由功能
	"github.com/gorilla/mux"
	sgm "github.com/zhangminghui6106/delinkcious/pkg/social_graph_manager"
	//记录日志
	"log"
	//提供http服务
	"net/http"
	"os"
	"strconv"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Run() {
	log.Println("Service started...")
	dbHost := os.Getenv("SOCIAL_GRAPH_DB_SERVICE_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPortStr := os.Getenv("SOCIAL_GRAPH_DB_SERVICE_PORT")
	if dbPortStr == "" {
		dbPortStr = "5432"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	log.Println("DB host:", dbHost, "DB port:", dbPortStr)

	dbPort, err := strconv.Atoi(dbPortStr)
	check(err)

	store, err := sgm.NewDbSocialGraphStore(dbHost, dbPort, "postgres", "postgres")
	check(err)

	svc, err := sgm.NewSocialGraphManager(store)
	check(err)

	followHandler := httptransport.NewServer(
		makeFollowEndpoint(svc),
		decodeFollowRequest,
		encodeResponse,
	)

	unfollowHandler := httptransport.NewServer(
		makeUnfollowEndpoint(svc),
		decodeUnfollowRequest,
		encodeResponse,
	)

	getFollowingHandler := httptransport.NewServer(
		makeGetFollowingEndpoint(svc),
		decodeGetFollowingRequest,
		encodeResponse,
	)

	getFollowersHandler := httptransport.NewServer(
		makeGetFollowersEndpoint(svc),
		decodeGetFollowersRequest,
		encodeResponse,
	)

	r := mux.NewRouter()
	r.Methods("POST").Path("/follow").Handler(followHandler)
	r.Methods("POST").Path("/unfollow").Handler(unfollowHandler)
	r.Methods("POST").Path("/following").Handler(getFollowingHandler)
	r.Methods("POST").Path("/followers").Handler(getFollowersHandler)

	log.Printf("Listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
