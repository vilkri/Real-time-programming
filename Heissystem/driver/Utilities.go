package driver

import ("fmt"
		"time"
)

func Utilities_bubble_sort(array []int) {
	
	for i := 0; i < len(array); i++ {
		for j := 0; j < len(array)-1; j++ {
			if array[j] > array[j+1] {
				temp := array[j+1]
				array[j+1] = array[j]
				array[j] = temp			
			}
		}	
	}
}

func Utilities_find_column_in_state_matrix(value int, array []int) (int) {

	for i := 0; i < len(array); i ++ {
		if value == array[i] {
			return i
		}	
	}
	fmt.Println("Utilities: no column to use in state matrix")
	return -1
}

func Utilities_i_am_alive(n_elevators int, m_floors int, port string) {
	
	var msg Message
	terminate_ch 		:= make(chan int, 1)
	
	msg.ID 				= 1 // for <I am alive>

	go UDP_broadcast("129.241.187.255:" + port, msg, 2000, terminate_ch)
}

func Utilities_broadcast_state(n_elevators int, m_floors int, port string) {
	
	var msg Message
	msg.ID 				= 2 // for <STATE>
	terminate_ch 		:= make(chan int, 1)
	
	msg.Latest_floor 	= Elev_get_latest_floor()
	msg.Direction 		= Elev_get_direction()
	

	go UDP_broadcast("129.241.187.255:" + port, msg, 5000, terminate_ch)
	
}

func Utilities_listen(port string, IP_list []int, active_elevator_list_ch chan []int){
	
	listen_ch 				:= make(chan Message, 1024)
	auxilary_list			:= make([]int, len(IP_list))
	active_elevator_list	:= make([]int, len(IP_list))
	go UDP_receive(port, listen_ch, 1000)
	
	for i := range active_elevator_list {
		active_elevator_list[i] = 1
	}
	
	for {
		message			:= <- listen_ch
		if message.ID == 1 {
			for i := 0; i < len(IP_list); i++ {
				if message.Trunc_IP == IP_list[i] {
					auxilary_list[i] = 0
				} else {
					if active_elevator_list[i] != 0 {
						auxilary_list[i] = auxilary_list[i] + 1
						fmt.Println("Incremented elevator value:", auxilary_list[i])
					}
				}
				if auxilary_list[i] > 5 {
					fmt.Println("Elevator", IP_list[i], "is d-e-a-d DEAD")
					active_elevator_list[i] = 0
					auxilary_list[i] = 0
				}			
			}
		}
		time.Sleep(10*time.Millisecond)
		active_elevator_list_ch <- active_elevator_list
		fmt.Println("Active elevator list:", active_elevator_list)
	}
	
}




