package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
    conn, err := net.Dial("tcp", "localhost:9000")
    if err != nil {
        fmt.Println("Ошибка подключения к серверу:", err)
        return
    }
    defer conn.Close()

    fmt.Println("Подключение к серверу установлено.")

    reader := bufio.NewReader(conn)
    writer := bufio.NewWriter(conn)

    for {
        showMenu()
        fmt.Print("Выберите действие: ")
        input := readLineFromStdin()
        choice, err := strconv.Atoi(input)
        if err != nil {
            fmt.Println("Введите число от 1 до 5")
            continue
        }

        switch choice {
        case 1:
            fmt.Print("Введите заголовок заметки: ")
            title := readLineFromStdin()
            fmt.Print("Введите текст заметки: ")
            content := readLineFromStdin()
            command := fmt.Sprintf("CREATE|%s|%s\n", title, content)
            sendCommand(writer, command)
            printResponse(reader)

        case 2:
            sendCommand(writer, "LIST\n")
            printMultiLineResponse(reader)

        case 3:
            fmt.Print("Введите ID заметки: ")
            id := readLineFromStdin()
            command := fmt.Sprintf("GET|%s\n", id)
            sendCommand(writer, command)
            printResponse(reader)

        case 4:
            fmt.Print("Введите ID заметки для удаления: ")
            id := readLineFromStdin()
            command := fmt.Sprintf("DELETE|%s\n", id)
            sendCommand(writer, command)
            printResponse(reader)

        case 5:
            fmt.Println("Выход из клиента...")
            return

        default:
            fmt.Println("Некорректный выбор. Введите число от 1 до 5.")
        }
    }
}

func showMenu() {
    fmt.Println("\n===== Меню заметок =====")
    fmt.Println("1. Создать заметку")
    fmt.Println("2. Показать все заметки")
    fmt.Println("3. Получить заметку по ID")
    fmt.Println("4. Удалить заметку")
    fmt.Println("5. Выход")
}

func readLineFromStdin() string {
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Scan()
    return scanner.Text()
}

func sendCommand(writer *bufio.Writer, command string) {
    writer.WriteString(command)
    writer.Flush()
}

func printResponse(reader *bufio.Reader) {
    response, _ := reader.ReadString('\n')
    fmt.Print(response)
}

func printMultiLineResponse(reader *bufio.Reader) {
    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Ошибка чтения:", err)
            break
        }

        line = strings.TrimSpace(line)

        if line == "END" {
            break
        }

        fmt.Println(line)
    }
}