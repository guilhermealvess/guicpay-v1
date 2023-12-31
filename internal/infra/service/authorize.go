package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/guilhermealvess/guicpay/internal/domain/service"
)

type authorize struct {
	url string
}

func NewAuthorizeService(url string) service.AuthorizeService {
	return &authorize{url}
}

func (a authorize) RegisterUser(ctx context.Context, id, email, password string) error {
	payload := map[string]any{
		"id":       id,
		"email":    email,
		"password": password,
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, a.url, bytes.NewReader(raw))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		return errors.New("TODO")
	}

	return nil
}

func (a authorize) Authorize(ctx context.Context, email, password string) error {
	payload := map[string]any{
		"email":    email,
		"password": password,
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, a.url, bytes.NewReader(raw))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		return errors.New("TODO")
	}

	return nil
}
