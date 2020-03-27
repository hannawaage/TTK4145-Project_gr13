package sync

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"../config"
	"../network/localip"
)

func Sync(id string, syncCh config.SyncChns, esmChns config.EsmChns) {
	const numPeers = config.NumElevs - 1
	idDig, _ := strconv.Atoi(id)
	masterID := idDig
	var (
		elev            config.Elevator
		onlineIPs       []string
		receivedReceipt []string
		currentMsgID    int
		numTimeouts     int
		allOrders       [config.NumElevs][config.NumFloors][config.NumButtons]bool
	)

	/*
		go func() {
			for {
				select {
				case newElev := <-esmChns.Elev:
					elev = newElev
					if allOrders[idDig-1] != elev.Orders {
						allOrders[idDig-1] = elev.Orders
					}
				}
			}
		}()*/

	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}

	msgTimer := time.NewTimer(5 * time.Second)
	msgTimer.Stop()

	go func() {
		for {
			currentMsgID = rand.Intn(256)
			msg := config.Message{elev, allOrders, currentMsgID, false, localIP, id}
			syncCh.SendChn <- msg
			msgTimer.Reset(800 * time.Millisecond)
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		select {
		case incomming := <-syncCh.RecChn:
			recID := incomming.LocalID
			recIDDig, _ := strconv.Atoi(recID)

			if id != recID { //Hvis det ikke er fra oss selv, BYTTES TIL IP VED KJØRING PÅ FORSKJELLIGE MASKINER
				if !contains(onlineIPs, recID) {
					// Dersom heisen enda ikke er registrert, sjekker vi om vi nå er online og sjekker om vi er master
					onlineIPs = append(onlineIPs, recID)
					if len(onlineIPs) == numPeers {
						syncCh.Online <- true
						for i := 0; i < numPeers; i++ {
							theID, _ := strconv.Atoi(onlineIPs[i])
							if theID < masterID {
								masterID = theID
							}
						}
						if masterID == idDig {
							syncCh.IAmMaster <- true
						} else {
							syncCh.IAmMaster <- false
						}
						/*
							Dette er ved diff på IP:
							localDig, _ := strconv.Atoi(localIP[len(localIP)-3:])
							for i := 0; i <= numPeers; i++ {
								theIP := onlineIPs[i]
								lastDig, _ := strconv.Atoi(theIP[len(theIP)-3:])
								if localDig < lastDig {
									iAmMaster = false
									break
								}
							}
						*/
					}
				}

				if !incomming.Receipt {
					// Hvis det ikke er en kvittering, skal vi svare med kvittering
					msg := config.Message{elev, allOrders, incomming.MsgId, true, localIP, id}
					//sender ut fem kvitteringer på femti millisekunder
					for i := 0; i < 5; i++ {
						syncCh.SendChn <- msg
						time.Sleep(10 * time.Millisecond)
					}
					theID, _ := strconv.Atoi(recID)
					if allOrders[theID-1] != incomming.Elev.Orders {
						// Hvis vi mottar noe nytt
						if (masterID == theID) || (recIDDig == masterID) {
							syncCh.ReceivedOrders <- incomming.AllOrders
						}
					}
				} else { // Hvis det er en kvittering
					if incomming.MsgId == currentMsgID {
						if !contains(receivedReceipt, recID) {
							receivedReceipt = append(receivedReceipt, recID)
							if len(receivedReceipt) == numPeers {
								//Hvis vi har fått bekreftelse fra alle andre peers på meldingen
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
				masterID = idDig
				//iAmMaster = false
			}
		}
	}
}

func OrdersDistribute(id int, syncCh config.SyncChns, esmCh config.EsmChns) {
	var (
		online         bool //initiates to false
		iAmMaster      bool = true
		allOrders      [config.NumElevs][config.NumFloors][config.NumButtons]bool
		receivedOrders [config.NumElevs][config.NumFloors][config.NumButtons]bool
		newLocalOrders config.Elevator
		//updateOrders   bool
		//updateOffline  bool
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
			case b := <-syncCh.IAmMaster:
				if b {
					iAmMaster = true
				} else {
					iAmMaster = false
				}
			case receivedOrders = <-syncCh.ReceivedOrders:
				if iAmMaster {
					esmCh.CurrentAllOrders <- costfcn()
					fmt.Println(".. I am Master and I just updated my orders")
				} else {
					esmCh.CurrentAllOrders <- receivedOrders
					fmt.Println(".. and I am backup and I just updated my orders")
				}
			case newLocalOrders = <-esmCh.Elev:
				if allOrders[id] != newLocalOrders.Orders {
					allOrders[id] = newLocalOrders.Orders
					esmCh.CurrentAllOrders <- allOrders
				}
			}

		}
	}()
	for {
		/*
			if online {

			} else {
				//fmt.Println("Singlemode")
				if updateOffline {
					//fmt.Println("Just updated my currentAllOrders")

				}
			}
		}*/
	}
}

func contains(elevs []string, str string) bool {
	for _, a := range elevs {
		if a == str {
			return true
		}
	}
	return false
}

func costfcn() [config.NumElevs][config.NumFloors][config.NumButtons]bool {
	var allOrderMat [config.NumElevs][config.NumFloors][config.NumButtons]bool
	allOrderMat[1][2][1] = true
	allOrderMat[2][2][2] = true
	return allOrderMat
}
