## Install

```bash
go get -u github.com/gagliardetto/eta
```

## Why I wrote it

I needed time estimates for when running tasks will be completed.

## Usage

```golang
package main

import (
	"fmt"
	"time"

	"github.com/gagliardetto/eta"
	"github.com/hako/durafmt"
)

func main() {
	totalTasks := int64(60)
	etac := eta.New(totalTasks)

	// Execute tasks:
	go func() {
		for {
			func() {
				// Mark one task done:
				defer etac.Done(1)
				err := runTask()
				if err != nil {
					panic(err)
				}
			}()
		}
	}()

	// Print stats:
	for {
		time.Sleep(time.Second)
		averagedETA := etac.GetETA()
		thisETA := durafmt.Parse(averagedETA.Round(time.Second)).String()

		percentDone := etac.GetFormattedPercentDone()

		fmt.Println(thisETA, percentDone)
	}
}
func runTask() error {
	time.Sleep(time.Second * 2)
	return nil
}

```