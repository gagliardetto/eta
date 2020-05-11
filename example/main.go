package main

import (
	"time"

	"github.com/gagliardetto/eta"
	. "github.com/gagliardetto/utils"
	"github.com/hako/durafmt"
)

func main() {
	etac := eta.New(
		60,
		time.Second*1,
	)

	go func() {
		for {
			time.Sleep(time.Second * 2)
			etac.Done(1)
		}
	}()

	go func() {

		for {
			time.Sleep(time.Second)
			averagedETA := etac.GetETA()
			thisETA := durafmt.Parse(averagedETA.Round(time.Second)).String()

			percentDone := GetFormattedPercent(etac.GetDone(), etac.GetTotal())

			Ln(thisETA, percentDone)
		}

	}()

	time.Sleep(time.Hour)
}
