package main

import (
	"math/rand"
	"time"
)

const (
MAXBET float64 = 100.0
)

type winnings struct {
	bank float64
	player float64
}


func makeRandCABoard() []float64 {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	cboard := []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}

	for i := 0 ; i < 6 ; i++ {

		if r.Uint32()&0x00000001 == 1 {
			cboard[i] = r.Float64() * MAXBET
		} else {
			cboard[i] = 0.0
		}
	}
	return cboard
}


func rollOne() int32 {

	now := time.Now()
	r := rand.New(rand.NewSource(now.UnixNano()))

	return r.Int31n(6) // 0 - 5
}

func caWorker(resultCh chan winnings, trialsCh chan int64, runtime int) {

	result := winnings{0.0, 0.0}
	now := time.Now()
	totalTrials := int64(0)

	for time.Since(now) < (time.Duration(runtime) * time.Minute) {
		totalTrials += 1
		board := makeRandCABoard()
		hits := []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}

		//бросок трех костей
		for i := 0; i < 3; i++ {
			hits[rollOne()] += 1.0
		}

		for i := 0; i < 6; i++ {
			//если выпала нужная кость - игрок получает деньги, а банк их теряет
			if hits[i] > 0.99 {
				won := hits[i] * board[i]
				result.player += won + board[i]
				result.bank -= won
			} else {
				//тогда наоборот
				result.bank += board[i]
			}
		}

	}
	resultCh <- result
	trialsCh <- totalTrials
}


func MonteCA(cores int, runtime int) (pctPlayer float64, pctBank float64, totalTrials int64) {


	totalCh := make(chan winnings, cores)
	trialsCh := make(chan int64, cores)

// launch the workers
	for i := 0 ; i < cores ; i++ {
		go caWorker(totalCh, trialsCh, runtime)
	}

// drain the channels
	totalWinnings := winnings{0.0, 0.0}
	totalTrials = 0
	var total winnings

	for i := 0 ; i < cores ; i++ {
		total = <- totalCh
		totalWinnings.bank += total.bank
		totalWinnings.player += total.player
		totalTrials +=  <- trialsCh
	}

	playerPlusBank := totalWinnings.bank + totalWinnings.player
	pctPlayer = (totalWinnings.player / playerPlusBank) * 100
	pctBank = 100.0 - pctPlayer
	return
}
