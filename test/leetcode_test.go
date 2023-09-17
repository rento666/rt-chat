package test

import (
	"fmt"
	"testing"
)

func TestLiCo(t *testing.T) {
	fmt.Printf("%v", numJewelsInStones("aA", "aAAAABBBb"))
}
func maxProfit(prices []int) int {
	minPrice := prices[0]
	maxPrice := 0
	res := 0
	for k, v := range prices {
		if minPrice >= v {
			minPrice = v
		}
		// 不是最后一个数 且 前 <= 后
		if k+1 < len(prices) && minPrice < prices[k+1] {
			maxPrice = prices[k+1]
		}
		// 如果 前利润 < 后利润
		if res < maxPrice-minPrice {
			res = maxPrice - minPrice
			maxPrice = 0
		}
		fmt.Printf("k: %v minprice: %v maxprice: %v res: %v \n", k, minPrice, maxPrice, res)
	}
	return res
}

func numJewelsInStones(jewels string, stones string) int {
	count := 0
	for _, j := range jewels {
		for _, s := range stones {
			if j == s {
				count++
			}
		}
	}
	return count
}
