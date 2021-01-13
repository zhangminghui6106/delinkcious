package service

import (
	"errors"
	//http传输包对于使用HTTP传输的服务是必须的
	httptransport "github.com/go-kit/kit/transport/http"
	//提供路由功能
	"github.com/gorilla/mux"
	sgm "github.com/zhangminghui6106/delinkcious/pkg/social_graph_manager"
	//记录日志
	"log"
	//提供http服务
	"net/http"
)

var (
	// return when an expected path variable is missing.
	BadRoutingError = errors.New("inconsistent mapping between route and handler")
)

func Run() {
	store, err := sgm.NewDbSocialGraphStore("localhost", 5432, "postgres", "postgres")
	if err != nil {
		log.Fatal(err)
	}
	svc, err := sgm.NewSocialGraphManager(store)
	if err != nil {
		log.Fatal(err)
	}

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
	r.Methods("GET").Path("/following/{username}").Handler(getFollowingHandler)
	r.Methods("GET").Path("/followers/{username}").Handler(getFollowersHandler)

	log.Println("Listening on port 9090...")
	log.Fatal(http.ListenAndServe(":9090", r))
}
