package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type payload struct {
	Id        int    `json:"id,omitempty"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

const apiURL = "http://arch.homework/user"

func put(ctx context.Context, id int, p payload) (int, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}
	data, err := json.Marshal(p)
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%d", apiURL, id), strings.NewReader(string(data)))
	if err != nil {
		return 0, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	return res.StatusCode, nil
}

func post(ctx context.Context, p payload) (int, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}
	data, err := json.Marshal(p)
	if err != nil {
		return 0, err
	}
	res, err := http.Post(apiURL, "application/json", strings.NewReader(string(data)))
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	return res.StatusCode, nil
}

func get(ctx context.Context, id int) (int, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}
	res, err := http.Get(fmt.Sprintf("%s/%d", apiURL, id))
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	return res.StatusCode, nil
}

func delete(ctx context.Context, id int) (int, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%d", apiURL, id), nil)
	if err != nil {
		return 0, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	return res.StatusCode, nil
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	counter := 1
	for {
		select {
		case <-ctx.Done():
			fmt.Println("done")
			return
		default:
			go func(counter int) {
				p := payload{
					Username:  fmt.Sprintf("user_%d", counter),
					FirstName: "John",
					LastName:  "Doe",
					Email:     fmt.Sprintf("user_%d@gmail.com", counter),
					Phone:     fmt.Sprintf("+3806612345%d", counter),
				}
				code, err := post(ctx, p)
				if err != nil {
					slog.Error("err", err.Error(), "counter", counter)
					return
				}
				if code != http.StatusCreated {
					slog.Error("post error", "counter", counter, "code", code)
					return
				}

				code, err = get(ctx, counter)
				if err != nil {
					slog.Error("err", err.Error(), "counter", counter)
					return
				}
				if code != http.StatusOK {
					slog.Error("get error", "counter", counter, "code", code)
					return
				}

				p.Username = fmt.Sprintf("user_%d_updated", counter)

				code, err = put(ctx, counter, p)
				if err != nil {
					slog.Error("err", err.Error(), "counter", counter)
					return
				}
				if code != http.StatusOK {
					slog.Error("put error", "counter", counter, "code", code)
					return
				}

				code, err = get(ctx, counter)
				if err != nil {
					slog.Error("err", err.Error(), "counter", counter)
					return
				}
				if code != http.StatusOK {
					slog.Error("get error", "counter", counter, "code", code)
					return
				}

				if counter%5 == 0 {
					code, err = delete(ctx, counter)
					if err != nil {
						slog.Error("err", err.Error(), "counter", counter)
						return
					}
					if code != http.StatusOK {
						slog.Error("delete error", "counter", counter, "code", code)
						return
					}
				}

			}(counter)
			if counter%60 == 0 {
				time.Sleep(time.Duration(rand.Intn(8000)+500) * time.Millisecond)
			}
		}
		counter++
	}
}
