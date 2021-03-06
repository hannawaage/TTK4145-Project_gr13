package sync

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"../config"
)

const (
	NumElevs   = config.NumElevs
	NumFloors  = config.NumFloors
	NumButtons = config.NumButtons
)

func Sync(id int, syncCh config.SyncChns, esmChns config.EsmChns) {
	masterID := id

	var (
		numPeers         int
		currentMsgID     int
		elev             config.Elevator
		onlineIDs        []int
		receivedReceipt  []int
		updatedAllOrders [NumElevs][NumFloors][NumButtons]int
		currentAllOrders [NumElevs][NumFloors][NumButtons]int
		allElevs         [NumElevs]config.Elevator
		orderTimeStamps  [NumFloors]int
		online           bool
		faultyElev       int = -1
	)

	go func() {
		for {
			select {
			case elev = <-esmChns.Elev:
				allElevs[id] = elev
				if !online {
					if currentAllOrders[id] != elev.Orders {
						updatedAllOrders[id] = elev.Orders
						esmChns.CurrentAllOrders <- updatedAllOrders
						currentAllOrders = updatedAllOrders
					}
				}
			}
		}
	}()

	msgTimer := time.NewTimer(10 * time.Second)
	msgTimer.Stop()

	go func() {
		for {
			UpdateTimeStamp(&orderTimeStamps, &currentAllOrders, &allElevs)
			if OrderTimeout(&orderTimeStamps) {
				go func() { syncCh.OrderTimeout <- true }()
			}
			currentMsgID = rand.Intn(256)
			msg := config.Message{elev, updatedAllOrders, currentMsgID, false, id}
			syncCh.SendChn <- msg
			msgTimer.Reset(200 * time.Millisecond)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	for {
		select {
		case incomming := <-syncCh.RecChn:
			recID := incomming.LocalID
			if id != recID {
				if !Contains(onlineIDs, recID) {
					onlineIDs = append(onlineIDs, recID)
					numPeers = len(onlineIDs)
					online = true
					for i := 0; i < numPeers; i++ {
						theID := onlineIDs[i]
						if theID < masterID {
							masterID = theID
						}
					}
				}
				if incomming.IsReceipt {
					if incomming.MsgId == currentMsgID {
						if !Contains(receivedReceipt, recID) {
							receivedReceipt = append(receivedReceipt, recID)
							if len(receivedReceipt) == numPeers {
								msgTimer.Stop()
								receivedReceipt = receivedReceipt[:0]
							}
						}
					}
				} else {
					allElevs[recID] = incomming.Elev
					for elevator := 0; elevator < NumElevs; elevator++ {
						if !Contains(onlineIDs, allElevs[elevator].Id) && (elevator != id) {
							allElevs[elevator].Orders = [NumFloors][NumButtons]int{}
						}
					}
					if id == masterID {
						updatedAllOrders = CostFunction(id, allElevs, onlineIDs)
					} else if recID == masterID {
						updatedAllOrders = incomming.AllOrders
					}
					if currentAllOrders != updatedAllOrders {
						esmChns.CurrentAllOrders <- updatedAllOrders
						currentAllOrders = updatedAllOrders
					}
					msg := config.Message{elev, updatedAllOrders, incomming.MsgId, true, id}
					for i := 0; i < 5; i++ {
						syncCh.SendChn <- msg
						time.Sleep(10 * time.Millisecond)
					}
				}
			}
		case <-msgTimer.C:
			numPeers = 0
			onlineIDs = onlineIDs[:0]
			receivedReceipt = receivedReceipt[:0]
			masterID = id
			online = false
			updatedAllOrders = MergeAllOrders(id, updatedAllOrders)
			esmChns.CurrentAllOrders <- updatedAllOrders
			currentAllOrders = updatedAllOrders
			
		case <-syncCh.OrderTimeout:
			faultyElev = FindFaultyElev(&currentAllOrders, &orderTimeStamps)
			if id == faultyElev {
				fmt.Println("Exiting for reinitialization")
				os.Exit(1)
			} else {
				updatedAllOrders = MergeAllOrders(id, updatedAllOrders)
			}
			esmChns.CurrentAllOrders <- updatedAllOrders
			currentAllOrders = updatedAllOrders
			orderTimeStamps = [NumFloors]int{}
			numPeers = 0
			onlineIDs = onlineIDs[:0]
			receivedReceipt = receivedReceipt[:0]
			masterID = id
			online = false
		}
	}
}
