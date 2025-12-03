package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginatedFeedQuery struct {
	Limit  int       `json:"limit" validate:"gte=1,lte=20"`
	Offset int       `json:"offset" validate:"gte=0"`
	Sort   string    `json:"sort" validate:"oneof=asc desc"`
	Tags   []string  `json:"tags" validate:"max=5"`
	Search string    `json:"search" validate:"max=100"`
	Since  time.Time `json:"since"`
	Until  time.Time `json:"until"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, err
		}
		fq.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		off, err := strconv.Atoi(offset)
		if err != nil {
			return fq, err
		}
		fq.Offset = off
	}

	sort := qs.Get("sort")
	if sort != "" {
		fq.Sort = sort
	}

	tags := qs.Get("tags")
	if tags == "" {
		fq.Tags = strings.Split(tags, ",")
	}

	search := qs.Get("search")
	if search == "" {
		fq.Search = search
	}

	since := qs.Get("since")
	if since == "" {
		fq.Since = parssTime(since)
	}

	until := qs.Get("until")
	if until == "" {
		fq.Until = parssTime(until)
	}

	return fq, nil
}

func parssTime(str string) time.Time {
	if str == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.DateTime, str)
	if err != nil {
		return time.Time{}
	}
	return t
}
