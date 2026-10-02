package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type Notification struct {
	id        int64 //генерируется через инкремент id
	msg       string
	isRead    bool      //по дефолту проставляется false
	createdAt time.Time //по дефолту проставляется текущая дата-время
}

func (n Notification) String() string {
	return fmt.Sprintf("id: %d; msg: %s; isRead: %t, createdAt: %s", n.id, n.msg, n.isRead, n.createdAt)
}

var idSeq int64 = 0
var notifications = make(map[int64]*Notification)

func addNotification() error {
	notification := Notification{}

	fmt.Print("Введите сообщение: ")
	_, err := fmt.Scanln(&notification.msg)
	if err != nil {
		return err
	}
	idSeq += 1
	notification.id = idSeq
	notification.isRead = false
	notification.createdAt = time.Now()

	notifications[notification.id] = &notification

	return nil
}

func deleteNotification() error {
	var id int64

	fmt.Print("Введите id уведомления, которое хотите изменить: ")
	fmt.Scanln(&id)

	_, ok := notifications[id]
	if !ok {
		return errors.New("notification not found")
	}

	delete(notifications, id)

	return nil
}

func printAllNotification() {
	if len(notifications) == 0 {
		fmt.Println("Уведомлений нет")
	}

	for _, notification := range notifications {
		fmt.Println(notification)
	}
}

func updateNotification() error {

	var id int64

	fmt.Print("Введите id уведомления, которое хотите изменить: ")
	fmt.Scanln(&id)

	notification, ok := notifications[id]
	if !ok {
		return errors.New("notification not found")
	}

	fmt.Print("Введите новое сообщение: ")
	_, err := fmt.Scanln(&notification.msg)
	if err != nil {
		return err
	}
	return nil
}

func exit() error {
	os.Exit(0)
	return nil
}

func main() {
	operations := map[int8]func() error{
		1: addNotification,
		2: updateNotification,
		3: deleteNotification,
		4: exit,
	}

	startText := `Введите номер операции:
	1 - добавить уведомление
	2 - изменить сообщение уведомления
	3 - удалить уведомление
	4 - отключить приложение
	`

	var indexOperation int8
	for {
		time.Sleep(1 * time.Second) //останавливаем процесс на 2 секунды, чтобы пользователь успел увидеть результат прошлой команды
		clearScreen()               //отчищаем консоль
		printAllNotification()
		fmt.Println(startText)

		_, err := fmt.Scanln(&indexOperation)
		if err != nil {
			fmt.Println("Неверный формат ввода")
		}

		op, ok := operations[indexOperation]
		if !ok {
			fmt.Printf("Операции с индексом %d не существует", indexOperation)
			continue
		}
		err = op()
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Print("Операция прошла успешно")
	}

}

func clearScreen() {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		// Linux и macOS
		cmd = exec.Command("clear")
	}

	// Подключаем стандартный вывод нашей программы к команде
	cmd.Stdout = os.Stdout
	// Запускаем команду и ждём завершения
	_ = cmd.Run()
}
