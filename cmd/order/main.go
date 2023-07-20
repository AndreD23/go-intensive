package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rabbitmq/amqp091-go"
	"githut.com/AndreD23/go-intensive/internal/entity"
	"githut.com/AndreD23/go-intensive/internal/infra/database"
	"githut.com/AndreD23/go-intensive/internal/usecase"
	"githut.com/AndreD23/go-intensive/pkg/rabbitmq"
)

type Car struct {
	Model string
	Color string
}

// metodo
func (c Car) Start() {
	fmt.Println(c.Model + " has been started")
}

// Esta eh uma funcao errada, pois nao esta trabalhando com ponteiros
func (c Car) ChangeColor(color string) {
	c.Color = color // duplicano o valor de c.Color na memoria
	fmt.Println("New color: " + c.Color)
}

// Esta eh a melhor forma de se trabalhar
func (c *Car) ChangeColorP(color string) {
	c.Color = color
	fmt.Println("New color: " + c.Color)
}

// funcao
func soma(x int, y int) int {
	return x + y
}

func main() {
	car := Car{
		Model: "Ferrari",
		Color: "Red",
	}

	car.Model = "FIAT"

	fmt.Println(car.Model)

	fmt.Println(soma(2, 3))

	car.Start()

	car.ChangeColor("Blue")
	fmt.Println(car.Color)

	car.ChangeColorP("Blue")
	fmt.Println(car.Color)

	a := 10
	b := a
	b = 20

	println(a)
	println(b)
	println(&b)

	// Alterando diretamente no endereco da memoria
	c := &a
	*c = 30

	println(c)
	println(&a)
	println(a)

	order, err := entity.NewOrder("1", 20, 3)
	if err != nil {
		fmt.Println(err)
	}

	err = order.CalculateFinalPrice()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(order.FinalPrice)

	db, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	orderRepository := database.NewOrderRepository(db)

	uc := usecase.NewCalculateFinalPrice(orderRepository)

	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	msgRabbitmqChannel := make(chan amqp091.Delivery)
	go rabbitmq.Consume(ch, msgRabbitmqChannel) // escutando a fila // trava // T2

	rabbitmqWorder(msgRabbitmqChannel, uc)
}

func rabbitmqWorder(msgChan chan amqp091.Delivery, uc *usecase.CalculateFinalPrice) {
	fmt.Println("Starting worker Rabbitmq")

	for msg := range msgChan {
		var input usecase.OrderInput
		err := json.Unmarshal(msg.Body, &input)
		if err != nil {
			panic(err)
		}
		output, err := uc.Execute(input)
		if err != nil {
			panic(err)
		}
		msg.Ack(false) // avisa ao rabbitmq que a mensagem foi consumida
		fmt.Println("Mensagem processada e salva no BD:", output)
	}
}
