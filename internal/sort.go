package internal

import (
	"sort"
	"strconv"

	"gitlab.com/slon/shad-go/gitfame/internal/util/constants"
)

func (stats *Statistics) SortResults(orderKey constants.OrderKey) {
	users := filterUsers(stats.CommitCountPerUser)

	sortUsers(users, stats.AggregatedData, orderKey)

	stats.OrderedData = formatSortedData(users, stats.AggregatedData)
}

func filterUsers(commitCountPerUser map[string]int) []string {
	var users []string
	for user := range commitCountPerUser {
		if user != "Not Committed Yet" {
			users = append(users, user)
		}
	}
	return users
}

func sortUsers(users []string, aggregatedData map[string][3]int, orderKey constants.OrderKey) {
	sort.SliceStable(users, func(i, j int) bool {
		return compareUsers(users[i], users[j], aggregatedData, orderKey)
	})
}

func compareUsers(user1, user2 string, aggregatedData map[string][3]int, orderKey constants.OrderKey) bool {
	switch orderKey {
	case constants.Lines:
		return compareByLines(user1, user2, aggregatedData)
	case constants.Commits:
		return compareByCommits(user1, user2, aggregatedData)
	case constants.Files:
		return compareByFiles(user1, user2, aggregatedData)
	default:
		return user1 < user2
	}
}

func compareByLines(user1, user2 string, aggregatedData map[string][3]int) bool {
	if aggregatedData[user1][0] == aggregatedData[user2][0] {
		if aggregatedData[user1][1] == aggregatedData[user2][1] {
			if aggregatedData[user1][2] == aggregatedData[user2][2] {
				return user1 < user2
			}
			return aggregatedData[user1][2] > aggregatedData[user2][2]
		}
		return aggregatedData[user1][1] > aggregatedData[user2][1]
	}
	return aggregatedData[user1][0] > aggregatedData[user2][0]
}

func compareByCommits(user1, user2 string, aggregatedData map[string][3]int) bool {
	if aggregatedData[user1][1] == aggregatedData[user2][1] {
		if aggregatedData[user1][0] == aggregatedData[user2][0] {
			if aggregatedData[user1][2] == aggregatedData[user2][2] {
				return user1 < user2
			}
			return aggregatedData[user1][2] > aggregatedData[user2][2]
		}
		return aggregatedData[user1][0] > aggregatedData[user2][0]
	}
	return aggregatedData[user1][1] > aggregatedData[user2][1]
}

func compareByFiles(user1, user2 string, aggregatedData map[string][3]int) bool {
	if aggregatedData[user1][2] == aggregatedData[user2][2] {
		if aggregatedData[user1][0] == aggregatedData[user2][0] {
			if aggregatedData[user1][1] == aggregatedData[user2][1] {
				return user1 < user2
			}
			return aggregatedData[user1][1] > aggregatedData[user2][1]
		}
		return aggregatedData[user1][0] > aggregatedData[user2][0]
	}
	return aggregatedData[user1][2] > aggregatedData[user2][2]
}

func formatSortedData(users []string, aggregatedData map[string][3]int) [][4]string {
	var sortedStats [][4]string
	for _, user := range users {
		sortedStats = append(sortedStats, [4]string{
			user,
			strconv.Itoa(aggregatedData[user][0]),
			strconv.Itoa(aggregatedData[user][1]),
			strconv.Itoa(aggregatedData[user][2]),
		})
	}
	return sortedStats
}
