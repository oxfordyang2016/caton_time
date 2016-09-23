/*
package main

import (
	"fmt"
	"strconv"
	"time"
)

var i int

func makeCakeAndSend(cs chan string) {
	i = i + 1
	cakeName := "Strawberry Cake " + strconv.Itoa(i)
	fmt.Println("now is make time ", time.Now())
	fmt.Println("Making a cake and sending ...", cakeName)

	cs <- cakeName //send a strawberry cake
}

func receiveCakeAndPack(cs chan string) {
	s := <-cs //get whatever cake is on the channel
	fmt.Println("now is pack time ", time.Now())
	fmt.Println("Packing received cake: ", s)
}

func main() {
	cs := make(chan string)
	for i := 0; i < 3; i++ {
		go makeCakeAndSend(cs)

		go receiveCakeAndPack(cs)

		//sleep for a while so that the program doesn’t exit immediately and output is clear for illustration
		time.Sleep(1 * 1e9)
	}
}

*/
/*
package main

import (
	"fmt"
	"time"
)

func say(s string) {
	for i := 0; i < 5; i++ {
		//	time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

func main() {
	go say("world")
	say("hello")
	fmt.Println(time.Now())
	//time.Sleep(100 * time.Millisecond)
}
*/
/*
package main

import (
	"fmt"
	"time"
)

func main() {
	my_sleep_func()
	fmt.Println("Control doesnt reach here till my_sleep_func finishes executing")

}

func my_sleep_func() {
	//sleeps for 5 seconds
	time.Sleep(5 * time.Second)
	fmt.Println("My func out of sleep")
}
*/
package main

import (
	"fmt"
	"time"
)

/*
you must know that





*/
func main() {
	//by using the go construct the function execution can be made concurrent
	go my_sleep_func1()

	go my_sleep_func2()
	//sleeping in main to make sure the main doesn't end
	time.Sleep(3 * time.Second) //if there is no sleep ,goroutine  may will no time to  excute
	fmt.Println("Control reach here even though my_sleep_func has not finished executing")
	fmt.Println("the intution is  similar to being asynchrinous without callbacks :D")
}

func my_sleep_func1() { //it means that
	//fmt.Println("i ma try")
	//sleeps for 5 seconds
	fmt.Println("fun1 ---->i want to know  it is when to excute")
	fmt.Println("i am func 1")
	time.Sleep(2 * time.Second) //if there sleep time over 3,go routine will may not excute!,but test <3,it will alaways excute
	fmt.Println("fun2 time ---->", time.Now())
	fmt.Println("func1--->My func out of sleep, but its executed concurrently")
}
func my_sleep_func2() {
	//fmt.Println("i ma try")
	//sleeps for 5 seconds
	fmt.Println("i am func 2")
	time.Sleep(1 * time.Second) //if there sleep time over 3,go routine will may not excute!,but test <3,it will alaways excute
	fmt.Println("fun2---->time ", time.Now())
	fmt.Println(" fun2 --->My func out of sleep, but its executed concurrently")
}
