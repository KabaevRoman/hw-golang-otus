package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	scanner := bufio.NewScanner(r)
	i := 0
	for scanner.Scan() {
		var user User
		if err = json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return
		}
		result[i] = user
		i++
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	re, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}
	for _, user := range u {
		matched := re.Match([]byte(user.Email))
		if matched {
			currDom := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			num := result[currDom]
			num++
			result[currDom] = num
		}
	}
	return result, nil
}
