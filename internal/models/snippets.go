package models

import (
	"database/sql"
	"strconv"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := "INSERT INTO snippets (title, content, created, expires) VALUES($1, $2, current_timestamp, current_date + interval '" + strconv.Itoa(expires) + " days') returning id"

	var pk int
	err := m.DB.QueryRow(stmt, title, content).Scan(&pk)
	if err != nil {
		return 0, err
	}

	return pk, nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, (created::timestamp(0)), (expires::timestamp(0)) FROM snippets
				WHERE expires > current_timestamp AND id=$1`

	s := &Snippet{}

	row := m.DB.QueryRow(stmt, id)

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, (created::timestamp(0)), (expires::timestamp(0)) FROM snippets
				WHERE expires > current_timestamp ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []*Snippet
	for rows.Next() {
		s := &Snippet{}
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
