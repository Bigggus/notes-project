package repo

import "database/sql"

func CreateNote(db *sql.DB, title string, content string) error {
    query := `INSERT INTO notes (title, content) VALUES ($1, $2)`
    _, err := db.Exec(query, title, content)
    return err
}

func DeleteNote(db *sql.DB, noteID int) error {
    query := `DELETE FROM notes WHERE id = $1`
    _, err := db.Exec(query, noteID)
    return err
}

func ListNotes(db *sql.DB) ([]Note, error) {
    query := `SELECT id, title, content, created_at FROM notes ORDER BY id`
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var notes []Note
    for rows.Next() {
        var note Note
        err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt)
        if err != nil {
            return nil, err
        }
        notes = append(notes, note)
    }
    return notes, nil
}

func GetNoteByID(db *sql.DB, noteID int) (*Note, error) {
    query := `SELECT id, title, content, created_at FROM notes WHERE id = $1`
    row := db.QueryRow(query, noteID)

    var note Note
    err := row.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt)
    if err != nil {
        return nil, err
    }
    return &note, nil
}