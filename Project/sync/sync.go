package sync

import (
	"fmt"
	"math/rand"
	"time"

	"../config"
	"../network/localip"
)

func Sync(id int, syncCh config.SyncChns, esmChns config.EsmChns) {
	masterID := id
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}

	var (
		numPeers           int
		currentMsgID       int
		numTimeouts        int
		elev               config.Elevator
		onlineIDs          []int
		receivedReceipt    []int
		updatedLocalOrders [config.NumElevs][config.NumFloors][config.NumButtons]bool
		currentAllOrders   [config.NumElevs][config.NumFloors][config.NumButtons]bool
		allElevs           [config.NumElevs]config.Elevator
		masterAck          bool
		online             bool
	)

	go func() {
		for {
			select {
			case elev = <-esmChns.Elev:
				if updatedLocalOrders[id] != elev.Orders {
					if online && !masterAck {
						updatedLocalOrders = mergeLocalOrders(id, &elev.Orders, updatedLocalOrders)
					} else {
						updatedLocalOrders[id] = elev.Orders
					}
				}
				allElevs[id].Orders = updatedLocalOrders[id]
				masterAck = false
			}
		}
	}()

	go func() {
		for {
			if currentAllOrders != updatedLocalOrders {
				if !online {
					esmChns.CurrentAllOrders <- updatedLocalOrders
					currentAllOrders = updatedLocalOrders
					/*
						if newCabOrdersOnly(id, &currentAllOrders, &updatedLocalOrders) {
							esmChns.CurrentAllOrders <- updatedLocalOrders
							currentAllOrders = updatedLocalOrders
						}*/
				} /*else {
					esmChns.CurrentAllOrders <- updatedLocalOrders
					currentAllOrders = updatedLocalOrders
				}*/
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()

	msgTimer := time.NewTimer(5 * time.Second)
	msgTimer.Stop()

	go func() {
		for {
			currentMsgID = rand.Intn(256)
			msg := config.Message{elev, updatedLocalOrders, currentMsgID, false, localIP, id}
			syncCh.SendChn <- msg
			msgTimer.Reset(800 * time.Millisecond)
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		select {
		case incomming := <-syncCh.RecChn:
			recID := incomming.LocalID
			if id != recID {
				if !contains(onlineIDs, recID) {
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
				if id == masterID {
					allElevs[recID] = incomming.Elev
					allElevs[recID].Orders = incomming.AllOrders[recID]
					updatedLocalOrders = CostFunction(id, allElevs, onlineIDs)
					if currentAllOrders != updatedLocalOrders {
						esmChns.CurrentAllOrders <- updatedLocalOrders
						currentAllOrders = updatedLocalOrders
					}
				} else if recID == masterID {
					if (currentAllOrders != updatedLocalOrders) && !masterAck {
						updatedLocalOrders = mergeLocalOrders(id, &elev.Orders, updatedLocalOrders)
					} else {
						updatedLocalOrders = incomming.AllOrders
						if currentAllOrders != updatedLocalOrders {
							esmChns.CurrentAllOrders <- updatedLocalOrders
							currentAllOrders = updatedLocalOrders
						}
					}
				}
				if incomming.IsReceipt {
					if incomming.MsgId == currentMsgID {
						if recID == masterID {
							masterAck = true
						}
						if !contains(receivedReceipt, recID) {
							receivedReceipt = append(receivedReceipt, recID)
							if len(receivedReceipt) == numPeers {
								numTimeouts = 0
								msgTimer.Stop()
								receivedReceipt = receivedReceipt[:0]
							}
						}
					}
				} else {
					msg := config.Message{elev, updatedLocalOrders, incomming.MsgId, true, localIP, id}
					for i := 0; i < 5; i++ {
						syncCh.SendChn <- msg
						time.Sleep(10 * time.Millisecond)
					}
				}
			}
		case <-msgTimer.C:
			numTimeouts++
			if numTimeouts > 2 {
				fmt.Println("Three timeouts in a row")
				numTimeouts = 0
				numPeers = 0
				onlineIDs = onlineIDs[:0]
				receivedReceipt = receivedReceipt[:0]
				masterID = id
				online = false
				updatedLocalOrders = mergeAllOrders(id, updatedLocalOrders)
				esmChns.CurrentAllOrders <- updatedLocalOrders
				currentAllOrders = updatedLocalOrders
			}
		}
	}
}
