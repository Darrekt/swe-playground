package leaderboard

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	pb "github.com/Darrekt/swe-playground/leaderboard/proto"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
)

func (a *API) LeaderboardSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Encountered error while reading request body.", http.StatusInternalServerError)
	}
	defer r.Body.Close()

	p := &pb.Submission{}
	err = proto.Unmarshal(body, p)
	if err != nil {
		http.Error(w, "Unable to unmarshal the Submission protobuf.", http.StatusBadRequest)
	}

	z := redis.Z{
		Score:  float64(p.Score),
		Member: p.UserId,
	}
	err = a.RedisClient.ZAdd(ctx, strconv.Itoa(int(p.QuestionId)), z).Err()
	if err != nil {
		log.Printf("Failed to update Redis leaderboard: %v", err)
		http.Error(w, "Failed to update leaderboard", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("Successfully processed submission from user %v for question %v with score %v.\n",
		p.UserId,
		p.QuestionId,
		p.Score,
	)
}

func RunLeaderboard(rdb *redis.Client, port int) {
	api := NewAPI(rdb)
	http.HandleFunc("/leaderboard", api.LeaderboardSubmissionHandler)
	log.Printf("Starting leaderboard server on port %v.\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
