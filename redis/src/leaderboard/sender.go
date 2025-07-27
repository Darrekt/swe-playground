package leaderboard

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	pb "github.com/Darrekt/swe-playground/leaderboard/proto"
	"google.golang.org/protobuf/proto"
)

func getContestGenInput() (int, int, int) {
	var submissions, users, questions int

	fmt.Println("Set how many submissions you want to simulate")
	n, err := fmt.Scan(&submissions)
	if n != 1 || err != nil {
		log.Fatalf("Error during user input for submissions: %s", err)
	}

	fmt.Println("Set how many users you want to simulate")
	n, err = fmt.Scan(&users)
	if n != 1 || err != nil {
		log.Fatalf("Error during user input for users: %s", err)
	}

	fmt.Println("Set how many questions you want to simulate")
	n, err = fmt.Scan(&questions)
	if n != 1 || err != nil {
		log.Fatalf("Error during user input for questions: %s", err)
	}

	return submissions, users, questions
}

func generateContestEntries(submissions, users, questions int) []*pb.Submission {
	// Submit a userID, questionID, score and name as a protobuf
	subs := make([]*pb.Submission, submissions)
	for i := range submissions {
		subs[i] = &pb.Submission{
			Name:       "Darrick Lau",
			UserId:     int32(rand.Int() % users),
			Score:      int32(rand.Int() % 100),
			QuestionId: int32(rand.Int() % questions),
		}
	}

	return subs
}

func sendSubmissions(server string, port int, submissions []*pb.Submission) {
	for _, sub := range submissions {
		data, err := proto.Marshal(sub)
		if err != nil {
			log.Fatalf("Failed to marshal protobuf.")
		}

		addr := fmt.Sprintf("http://%s:%v/leaderboard", server, port)
		res, err := http.Post(addr, "application/x-protobuf", bytes.NewReader(data))
		if err != nil {
			log.Fatalf("Failed to send POST request to %s", addr)
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			log.Fatalf("Server responded a non-OK status: %v", res.StatusCode)
		}
	}
}

func StartClient(host string, port int) {
	submissions, users, questions := getContestGenInput()
	subs := generateContestEntries(submissions, users, questions)
	sendSubmissions(host, port, subs)
}
