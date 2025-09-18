package api

import (
	"net/http"
	"time"
)

type TraceNameRequest struct {
	Name string `json:"name"`
	Page int    `json:"page"`
}

type TraceNameItem struct {
	Name  string `json:"name"`
	Mykad string `json:"mykad"`
	ID    int    `json:"id"`
}

type TraceNameResponse struct {
	Data    []TraceNameItem `json:"data"`
	Message string          `json:"message"`
	User    User            `json:"user"`
}

type Client struct {
	httpClient *http.Client
	apiToken   string
}

type TraceDetailID struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Mykad       string `json:"mykad"`
	Address     string `json:"address"`
	City        string `json:"city"`
	Postcode    string `json:"postcode"`
	State       string `json:"state"`
	Phone       string `json:"phone"`
	Gender      string `json:"gender"`
	Mobile      string `json:"mobile"`
	Race        string `json:"race"`
	Religion    string `json:"religion"`
	Income      string `json:"income"`
	Occupations string `json:"occupations"`
	Addresses   string `json:"addresses"`
}

type TraceDetailResponse struct {
	Message string          `json:"message"`
	Data    []TraceDetailID `json:"data"`
	User    User            `json:"user"`
}

type TraceDetailRequest struct {
	ID int `json:"id"`
}

type TraceNameMsg struct {
	Items []TraceNameItem
	Err   error
}

type TraceDetailsMsg struct {
	Details []TraceDetailID
	Err     error
}

type TraceRelationsRequest struct {
	ID     int `json:"id"`
	Offset int `json:"offset"`
}

type TraceRelationsItem struct {
	UserID        int    `json:"user_id"`
	RelatedUserID int    `json:"related_user_id"`
	Relation      string `json:"relation"`
}

type TraceRelationsResponse struct {
	Data          []TraceDetailID    `json:"data"`
	Relationships TraceRelationsItem `json:"relationships"`
	User          User               `json:"user"`
}

type TraceRelationsMsg struct {
	Details   []TraceDetailID
	Relations TraceRelationsItem
	Err       error
}

type User struct {
	ID            int            `json:"id"`
	Name          string         `json:"name"`
	Role          string         `json:"role"`
	APIToken      string         `json:"api_token"`
	Balance       float64        `json:"balance"`
	Subscriptions []Subscription `json:"subscriptions"`
}

type Subscription struct {
	Plan      string    `json:"plan"`
	Active    bool      `json:"active"`
	ExpiredAt time.Time `json:"expired_at"`
}

type TraceNRICRequest struct {
	NRIC string `json:"nric"`
}

type TraceNRICResponse struct {
	Data []TraceDetailID `json:"data"`
	User User            `json:"user"`
}

type HistoryRequest struct {
	Page int `json:"page"`
}

type HistoryItem struct {
	Type   string  `json:"type"`
	Email  string  `json:"email"`
	Query  string  `json:"query"`
	Result string  `json:"result_summary"`
	Cost   float64 `json:"cost"`
}

type HistoryResponse struct {
	Data    []HistoryItem `json:"data"`
	Message string        `json:"message"`
}

type HistoryMsg struct {
	Data []HistoryItem
	Err  error
}
