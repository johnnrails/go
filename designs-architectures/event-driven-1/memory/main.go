package main

import (
	"fmt"
	eventbus "github.com/asaskevich/EventBus"
)

func calculator(a int, b int) {
	fmt.Printf("%d\n", a+b)
}

func main() {
	bus := eventbus.New()
	bus.Subscribe("main:calculator", calculator)
	bus.Publish("main:calculator", 20, 40)
	bus.Unsubscribe("main:calculator", calculator)
}
