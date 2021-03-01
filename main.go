package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/akhettar/rec-engine/app"
)


func main() {
	
	log.Info("starting the server on 3000")
	app.InitialiseApp("redis://34.66.203.46:6379").Run(":3000")
}

	

// 	arguments, _ := docopt.Parse(usage, nil, true, "redis-recommend", false)

// 	rr, err = redrec.New("redis://localhost:6379")
// 	chekErrorAndExit(err)

// 	if arguments["rate"].(bool) {
// 		user := arguments["<user>"].(string)
// 		item := arguments["<item>"].(string)
// 		score, err := strconv.ParseFloat(arguments["<score>"].(string), 64)
// 		chekErrorAndExit(err)
// 		rate(user, item, score)
// 	}

// 	if arguments["get-probability"].(bool) {
// 		user := arguments["<user>"].(string)
// 		item := arguments["<item>"].(string)
// 		getProbability(user, item)
// 	}

// 	if arguments["suggest"].(bool) {
// 		user := arguments["<user>"].(string)
// 		results, err := strconv.Atoi(arguments["--results"].(string))
// 		chekErrorAndExit(err)
// 		suggest(user, results)
// 	}

// 	if arguments["batch-update"].(bool) {
// 		results, err := strconv.Atoi(arguments["--results"].(string))
// 		chekErrorAndExit(err)
// 		update(results)
// 	}
// }

// func chekErrorAndExit(err error) {
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err.Error())
// 		rr.CloseConn()
// 		os.Exit(1)
// 	}
// }

// func rate(user string, item string, score float64) {
// 	fmt.Printf("User %s ranked item %s with %.2f\n", user, item, score)
// 	err := rr.Rate(item, user, score)
// 	chekErrorAndExit(err)
// }

// func getProbability(user string, item string) {
// 	score, err := rr.CalcItemProbability(item, user)
// 	chekErrorAndExit(err)
// 	fmt.Printf("%s %s %.2f\n", user, item, score)
// }

// func suggest(user string, max int) {
// 	fmt.Printf("Getting %d results for user %s\n", max, user)
// 	rr.UpdateSuggestedItems(user, max)
// 	s, err := rr.GetUserSuggestions(user, max)
// 	chekErrorAndExit(err)
// 	fmt.Println("results:")
// 	fmt.Println(s)
// }

// func update(max int) {
// 	fmt.Printf("Updating DB\n")
// 	err := rr.BatchUpdateSimilarUsers(max)
// 	chekErrorAndExit(err)
// }
