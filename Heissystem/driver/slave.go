package driver

import ("fmt"
		"strconv"
		"time")
		

func Slave_main(Alive_port string, Order_port string, state_port string, row int, IP_list []int) {

	time.Sleep(2000*time.Millisecond)
	fmt.Print("I am slave\n")

	terminate_ch				:= make(chan int, 1)
	active_elevator_list_ch 	:= make(chan []int, 100)
	receive_ch					:= make(chan Message, 500)
	
	state_matrix 				:= Orders_make_state_matrix()
	local_order_queue		:= make([]int, len(state_matrix))
	
	go UDP_receive(Order_port, receive_ch, 10)
	go Utilities_send_i_am_alive(Alive_port)
	go Lamps_main(row, state_matrix, terminate_ch)
	go Sensors_main(state_matrix, row, receive_ch)
	//go Motor_main()
	
	go Utilities_whos_alive(Alive_port, IP_list, active_elevator_list_ch)
	go Slave_send_state_to_master(state_port)
	
	for {
		select {
		case i := <- receive_ch:
			if i.ID == ORDER_ASSIGN {
				fmt.Println("I will execute order", i.Order_type)
				Utilities_ack_order(Order_port, ORDER_ASSIGN_ACK)
				Orders_update_elevator_queue(local_order_queue, i, 1)
				time.Sleep(100*time.Millisecond)
				Utilities_send_order_done(i, Order_port, receive_ch)
			} else if i.ID == STATE_MATRIX_UPDATE {
				state_matrix = i.State_matrix
			}	
		default:
			Orders_execute_orders(local_order_queue)
		}
	}
	
}

func Slave_send_state_to_master(port string) {

	var latest_floor int 		= Sensors_get_latest_floor()
	var current_direction int 	= Motor_get_direction()

	for { //Sjekker om endring i tilstand
		if (Sensors_get_latest_floor() != latest_floor || Motor_get_direction() != current_direction) {
			Utilities_broadcast_state(port)
		}
		latest_floor 		= Sensors_get_latest_floor()
		current_direction 	= Motor_get_direction()
		time.Sleep(10*time.Millisecond)
	}
}
