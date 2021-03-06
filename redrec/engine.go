package redrec

import (
	"fmt"
	"math"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/gomodule/redigo/redis"
)

// Redrec struct
type Redrec struct {
	rconn redis.Conn
}

// New returns a new Redrec
func New(url string) (*Redrec, error) {
	rconn, err := redis.DialURL(url, redis.DialUsername(os.Getenv("REDIS_USER")), redis.DialPassword(os.Getenv("REDIS_PWD")))
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	rr := &Redrec{
		rconn: rconn,
	}

	return rr, nil
}

// CloseConn closes the Redis connection
func (rr *Redrec) CloseConn() {
	rr.rconn.Close()
}

// Rate adds user->score to a given item
func (rr *Redrec) Rate(item string, user string, score float64) error {
	if _, err := rr.rconn.Do("ZADD", fmt.Sprintf("user:%s:items", user), score, item); err != nil {
		return err
	}

	if _, err := rr.rconn.Do("ZADD", fmt.Sprintf("item:%s:scores", item), score, user); err != nil {
		return err
	}

	if _, err := rr.rconn.Do("ZINCRBY", fmt.Sprintf("items"), 1, item); err != nil {
		log.Warnf("failed to store item: %s in the items list", item)
	}

	if err := rr.BatchUpdateSimilarUsers(100000); err != nil {
		log.Warnf("failed to update DB with error %v", err)
	}

	if _, err := rr.rconn.Do("SADD", "users", user); err != nil {
		return err
	}
	return nil
}

// GetUserSuggestions return the existing user
//suggestions range for a given user as a []string
func (rr *Redrec) GetUserSuggestions(user string, max int) ([]string, error) {
	items, err := redis.Strings(rr.rconn.Do("ZREVRANGE", fmt.Sprintf("user:%s:suggestions", user), 0, max, "WITHSCORES"))
	if err != nil {
		return nil, err
	}

	return items, nil
}

// BatchUpdateSimilarUsers runs on all the users,
// getting the similarity candidates for each user and storing the similar
// users and scores in a sorted set
func (rr *Redrec) BatchUpdateSimilarUsers(max int) error {
	users, err := redis.Strings(rr.rconn.Do("SMEMBERS", "users"))
	if err != nil {
		return err
	}
	for _, user := range users {
		candidates, err := rr.getSimilarityCandidates(user, max)
		args := []interface{}{}
		args = append(args, fmt.Sprintf("user:%s:similars", user))
		for _, candidate := range candidates {
			if candidate != user {
				score, _ := rr.calcSimilarity(user, candidate)
				args = append(args, score)
				args = append(args, candidate)
			}
		}

		_, err = rr.rconn.Do("ZADD", args...)
		if err != nil {
			log.Warnf("ZADD ERR: %v", err)
			return err
		}
	}

	return nil
}

// UpdateSuggestedItems gets the candidate suggest items for a given user and stores
// the calculated probability for each item in a sorted set
func (rr *Redrec) UpdateSuggestedItems(user string, max int) error {
	items, err := rr.getSuggestCandidates(user, max)

	// there are no suggestion candidate, abort updates
	if len(items) == 0 {
		return nil
	}

	// reset the maxumim number of suggestions to the possible value
	if max > len(items) {
		max = len(items)
	}

	args := []interface{}{}
	args = append(args, fmt.Sprintf("user:%s:suggestions", user))
	for _, item := range items {
		probability, _ := rr.CalcItemProbability(user, item)
		args = append(args, probability)
		args = append(args, item)
	}

	_, err = rr.rconn.Do("ZADD", args...)
	if err != nil {
		log.Println("ZADD ERR: ", err)
		return err
	}

	return nil
}

// CalcItemProbability takes all the user`s similars that rated the input item
// and calculates the average score.
func (rr *Redrec) CalcItemProbability(user string, item string) (float64, error) {
	_, err := rr.rconn.Do("ZINTERSTORE",
		"ztmp", 2, fmt.Sprintf("user:%s:similars", user), fmt.Sprintf("item:%s:scores", item), "WEIGHTS", 0, 1)
	if err != nil {
		return 0, err
	}

	scores, err := redis.Strings(rr.rconn.Do("ZRANGE", "ztmp", 0, -1, "WITHSCORES"))
	rr.rconn.Do("DEL", "ztmp")
	if err != nil {
		return 0, err
	}

	if len(scores) == 0 {
		return 0, nil
	}

	var score float64
	for i := 1; i < len(scores); i += 2 {
		val, _ := strconv.ParseFloat(scores[i], 64)
		score += val
	}
	score /= float64(len(scores) / 2)

	return score, nil
}

// GetUserItems gets all items for a given user
func (rr *Redrec) GetUserItems(user string, max int) ([]string, error) {
	items, err := redis.Strings(rr.rconn.Do("ZREVRANGE", fmt.Sprintf("user:%s:items", user), 0, max))
	if err != nil {
		return nil, err
	}

	return items, nil
}

// GetPopularItems return the existing user
//suggestions range for a given user as a []string
func (rr *Redrec) GetPopularItems(max int) ([]string, error) {
	items, err := redis.Strings(rr.rconn.Do("ZREVRANGE", "items", 0, max, "WITHSCORES"))
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (rr *Redrec) getItemScores(item string, max int) (map[string]string, error) {
	scores, err := redis.StringMap(rr.rconn.Do("ZREVRANGE", fmt.Sprintf("item:%s:scores", item), 0, max))
	if err != nil {
		return nil, err
	}

	return scores, nil
}

func (rr *Redrec) getSimilarityCandidates(user string, max int) ([]string, error) {
	items, err := rr.GetUserItems(user, max)
	if max > len(items) {
		max = len(items)
	}

	args := []interface{}{}
	args = append(args, "ztmp", float64(max))
	for i := 0; i < max; i++ {
		args = append(args, fmt.Sprintf("item:%s:scores", items[i]))
	}

	_, err = rr.rconn.Do("ZUNIONSTORE", args...)
	if err != nil {
		return nil, err
	}

	users, err := redis.Strings(rr.rconn.Do("ZRANGE", "ztmp", 0, -1))
	if err != nil {
		return nil, err
	}

	_, err = rr.rconn.Do("DEL", "ztmp")
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (rr *Redrec) getSuggestCandidates(user string, max int) ([]string, error) {
	similarUsers, err := redis.Strings(rr.rconn.Do("ZRANGE", fmt.Sprintf("user:%s:similars", user), 0, max))
	if err != nil {
		return nil, err
	}

	max = len(similarUsers)
	args := []interface{}{}
	args = append(args, "ztmp", float64(max+1), fmt.Sprintf("user:%s:items", user))
	weights := []interface{}{}
	weights = append(weights, "WEIGHTS", -1.0)
	for _, simuser := range similarUsers {
		args = append(args, fmt.Sprintf("user:%s:items", simuser))
		weights = append(weights, 1.0)
	}

	args = append(args, weights...)
	args = append(args, "AGGREGATE", "MIN")
	_, err = rr.rconn.Do("ZUNIONSTORE", args...)
	if err != nil {
		return nil, err
	}

	candidates, err := redis.Strings(rr.rconn.Do("ZRANGEBYSCORE", "ztmp", 0, "inf"))
	if err != nil {
		return nil, err
	}

	_, err = rr.rconn.Do("DEL", "ztmp")
	if err != nil {
		return nil, err
	}

	return candidates, nil
}

func (rr *Redrec) calcSimilarity(user string, simuser string) (float64, error) {
	_, err := rr.rconn.Do("ZINTERSTORE",
		"ztmp", 2, fmt.Sprintf("user:%s:items", user), fmt.Sprintf("user:%s:items", simuser), "WEIGHTS", 1, -1)
	if err != nil {
		return 0, err
	}

	userDiffs, err := redis.Strings(rr.rconn.Do("ZRANGE", "ztmp", 0, -1, "WITHSCORES"))
	rr.rconn.Do("DEL", "ztmp")
	if err != nil {
		return 0, err
	}

	if len(userDiffs) == 0 {
		return 0, nil
	}

	var score float64
	for i := 1; i < len(userDiffs); i += 2 {
		diffVal, _ := strconv.ParseFloat(userDiffs[i], 64)
		score += diffVal * diffVal
	}
	score /= float64(len(userDiffs) / 2)
	score = math.Sqrt(score)

	return score, nil
}
