package sync

import (
	"fmt"
	"math/rand"
	"time"

	"../config"
	"../network/localip"
)

func Sync(id int, syncCh config.SyncChns, esmChns config.EsmChns) {
	const numPeers = config.NumElevs - 1
	masterID := id
	var (
		elev               config.Elevator
		onlineIPs          []int
		receivedReceipt    []int
		currentMsgID       int
		numTimeouts        int
		updatedLocalOrders [config.NumElevs][config.NumFloors][config.NumButtons]bool
		currentAllOrders   [config.NumElevs][config.NumFloors][config.NumButtons]bool
		online             bool
		allElevs           [config.NumElevs]config.Elevator
	)
	go func() {
		for {
			select {
			case b := <-syncCh.Online:
				if b {
					online = true
					fmt.Println("Yaho, we are online!")
				} else {
					online = false
					fmt.Println("Boo, we are offline.")
				}
			case elev = <-esmChns.Elev:
				if updatedLocalOrders[id] != elev.Orders {
					updatedLocalOrders[id] = elev.Orders
				}
				allElevs[id] = elev
			}
		}
	}()

	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}

	go func() {
		for {
			if currentAllOrders != updatedLocalOrders {
				if !online {
					updatedLocalOrders = mergeAllOrders(id, updatedLocalOrders)
					esmChns.CurrentAllOrders <- updatedLocalOrders
					currentAllOrders = updatedLocalOrders
				} else {
					if newCabOrdersOnly(id, &currentAllOrders, &updatedLocalOrders) {
						esmChns.CurrentAllOrders <- updatedLocalOrders
						currentAllOrders = updatedLocalOrders
					}
				}
			}
			time.Sleep(20 * time.Millisecond)
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
			//esmChns.CurrentAllOrders <- currentAllOrders
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		select {
		case incomming := <-syncCh.RecChn:
			recID := incomming.LocalID
			if id != recID {
				if !contains(onlineIPs, recID) {
					onlineIPs = append(onlineIPs, recID)
					if len(onlineIPs) == numPeers {
						syncCh.Online <- true
						for i := 0; i < numPeers; i++ {
							theID := onlineIPs[i]
							if theID < masterID {
								masterID = theID
							}
						}
					}
				}
				if online {
					allElevs[recID] = incomming.Elev
					allElevs[recID].Orders = incomming.AllOrders[recID]
					if masterID == id {
						updatedLocalOrders = CostFunction(allElevs)
					} else if masterID == recID {
						if currentAllOrders == updatedLocalOrders {
							// Ikke noe nytt lokalt - ta inn det vi får fra master
							updatedLocalOrders = incomming.AllOrders
						} else {
							// nytt lokalt - merge med det nye
							//fmt.Println("Lokale endringer, merger med masterbeskjed")
							updatedLocalOrders = mergeLocalOrders(id, &elev.Orders, incomming.AllOrders)
						}
					}
					if currentAllOrders != updatedLocalOrders {
						esmChns.CurrentAllOrders <- updatedLocalOrders
						currentAllOrders = updatedLocalOrders
					}
				}
				if !incomming.Receipt {
					msg := config.Message{elev, updatedLocalOrders, incomming.MsgId, true, localIP, id}
					for i := 0; i < 5; i++ {
						syncCh.SendChn <- msg
						time.Sleep(10 * time.Millisecond)
					}
				} else {
					if incomming.MsgId == currentMsgID {
						if !contains(receivedReceipt, recID) {
							receivedReceipt = append(receivedReceipt, recID)
							if len(receivedReceipt) == numPeers {
								numTimeouts = 0
								msgTimer.Stop()
								receivedReceipt = receivedReceipt[:0]
							}
						}
					}
				}
			}
		case <-msgTimer.C:
			numTimeouts++
			if numTimeouts > 2 {
				syncCh.Online <- false
				fmt.Println("Three timeouts in a row")
				numTimeouts = 0
				onlineIPs = onlineIPs[:0]
				masterID = id
			}
		}
	}
}

func contains(elevs []int, new int) bool {
	for _, a := range elevs {
		if a == new {
			return true
		}
	}
	return false
}

func mergeAllOrders(id int, all [config.NumElevs][config.NumFloors][config.NumButtons]bool) [config.NumElevs][config.NumFloors][config.NumButtons]bool {
	var merged [config.NumElevs][config.NumFloors][config.NumButtons]bool
	merged[id] = all[id]
	for elev := 0; elev < config.NumElevs; elev++ {
		if elev == id {
			continue
		}
		for floor := 0; floor < config.NumFloors; floor++ {
			for btn := 0; btn < config.NumButtons; btn++ {
				if all[elev][floor][btn] && btn != config.NumButtons-1 {
					merged[id][floor][btn] = true
					merged[elev][floor][btn] = false
				}
			}
		}
	}
	return merged
}

func newCabOrdersOnly(id int, current *[config.NumElevs][config.NumFloors][config.NumButtons]bool, updated *[config.NumElevs][config.NumFloors][config.NumButtons]bool) bool {
	var newCab bool
	for floor := 0; floor < config.NumFloors; floor++ {
		for btn := 0; btn < config.NumButtons-1; btn++ {
			if current[id][floor][btn] != updated[id][floor][btn] {
				return false
			}
		}
		if current[id][floor][2] != updated[id][floor][2] {
			newCab = true
		}
	}
	return newCab
}

func mergeLocalOrders(id int, local *[config.NumFloors][config.NumButtons]bool, incomming [config.NumElevs][config.NumFloors][config.NumButtons]bool) [config.NumElevs][config.NumFloors][config.NumButtons]bool {
	var merged [config.NumElevs][config.NumFloors][config.NumButtons]bool
	merged = incomming
	for floor := 0; floor < config.NumFloors; floor++ {
		for btn := 0; btn < config.NumButtons-1; btn++ {
			if local[floor][btn] {
				merged[id][floor][btn] = true
			}
		}
	}
	return merged
}
