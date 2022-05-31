package job

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"encore.dev/storage/sqldb"
)

type Job struct {
	ID               string
	Title            string
	OrganizationName string
	ContractType     string
	Location         string
	Link             string
}

type StoreJobParams struct {
	Title            string
	OrganizationName string
	ContractType     string
	Location         string
	Link             string
}

type ListResponse struct {
	Jobs []*Job
}

// Store a job.
//encore:api public method=POST path=/job
func Store(ctx context.Context, p *StoreJobParams) error {
	id, err := generateID()
	if err != nil {
		return err
	}

	if err := insert(ctx, id, p); err != nil {
		return err
	}

	return nil
}

// Insert a job into the database
func insert(ctx context.Context, id string, p *StoreJobParams) error {
	_, err := sqldb.Exec(ctx, `
		INSERT INTO jobs (id, title, organizationName, contractType, location, link)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, id, p.Title, p.OrganizationName, p.ContractType, p.Location, p.Link)

	return err
}

// generateID generates a random short ID.
func generateID() (string, error) {
	var data [6]byte // 6 bytes of entropy
	if _, err := rand.Read(data[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data[:]), nil
}

// List jobs.
//encore:api public method=GET path=/job
func List(ctx context.Context) (*ListResponse, error) {
	rows, err := sqldb.Query(ctx, "SELECT * FROM jobs")
	if err != nil {
		return nil, fmt.Errorf("could not list jobs")
	}
	defer rows.Close()

	var jobs []*Job

	for rows.Next() {
		j := &Job{}
		err := rows.Scan(&j.ID, &j.Title, &j.OrganizationName, &j.ContractType, &j.Location, &j.Link)
		if err != nil {
			return nil, fmt.Errorf("could not scan jobs")
		}
		jobs = append(jobs, j)
	}

	return &ListResponse{Jobs: jobs}, nil
}
