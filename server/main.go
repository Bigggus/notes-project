package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net"
	repo "notes_server/db"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	host := "localhost"      
	port := "5432"        
	user := "appuser"       
	password := "secret" 
	dbname := "notes"       
	

	databaseDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
	println(databaseDSN)

	db, err := sql.Open("postgres", databaseDSN)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("БД не доступна:", err)
	}

	log.Println("Подключение к БД успешно")

	err = StartTCPServer(db, ":9000")
	if err != nil {
		log.Fatal("Ошибка TCP сервера:", err)
	}
}



func StartTCPServer(db *sql.DB, address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	fmt.Println("TCP сервер слушает:", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn, db)
	}
}

func handleClient(conn net.Conn, db *sql.DB) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		line = strings.TrimSpace(line)
		parts := strings.Split(line, "|")
		command := strings.ToUpper(parts[0])

		switch command {

		case "CREATE":
			if len(parts) < 3 {
				writeError(writer, "CREATE|title|content")
				continue
			}
			err := repo.CreateNote(db, parts[1], parts[2])
			if err != nil {
				writeError(writer, err.Error())
			} else {
				writeOK(writer, "Заметка добавлена")
			}

		case "LIST":
			notes, err := repo.ListNotes(db)
			if err != nil {
				writeError(writer, err.Error())
				continue
			}
			for _, note := range notes {
				line := fmt.Sprintf("%d|%s|%s|%s",
					note.ID,
					note.Title,
					note.Content,
					note.CreatedAt.Format("2006-01-02 15:04:05"))
				writer.WriteString(line + "\n")
			}
			writer.WriteString("END\n")
			writer.Flush()

		case "GET":
			if len(parts) < 2 {
				writeError(writer, "GET|id")
				continue
			}
			id, err := strconv.Atoi(parts[1])
			if err != nil {
				writeError(writer, "Неверный ID")
				continue
			}
			note, err := repo.GetNoteByID(db, id)
			if err != nil {
				writeError(writer, "Заметка не найдена")
				continue
			}
			line := fmt.Sprintf("%d|%s|%s|%s",
				note.ID, note.Title, note.Content,
				note.CreatedAt.Format("2006-01-02 15:04:05"))
			writer.WriteString(line + "\n")
			writer.Flush()

		case "DELETE":
			if len(parts) < 2 {
				writeError(writer, "DELETE|id")
				continue
			}
			id, err := strconv.Atoi(parts[1])
			if err != nil {
				writeError(writer, "Неверный ID")
				continue
			}
			err = repo.DeleteNote(db, id)
			if err != nil {
				writeError(writer, err.Error())
			} else {
				writeOK(writer, "Заметка удалена")
			}

		default:
			writeError(writer, "Неизвестная команда")
		}
	}
}

func writeOK(w *bufio.Writer, message string) {
	w.WriteString("OK: " + message + "\n")
	w.Flush()
}

func writeError(w *bufio.Writer, message string) {
	w.WriteString("ERROR: " + message + "\n")
	w.Flush()
}
