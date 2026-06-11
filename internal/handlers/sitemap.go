package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type sitemapRow struct {
	ID        int
	Slug      string
	UpdatedAt time.Time
}

const baseURL = "https://cdlytics.com"

func (h *Handler) GetSitemap(c *gin.Context) {
	if h.db == nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	ctx, cancel := getContext(10)
	defer cancel()
	db := h.db.WithContext(ctx)

	var players []sitemapRow
	db.Raw("SELECT id, updated_at FROM players ORDER BY id").Scan(&players)

	var teams []sitemapRow
	db.Raw("SELECT id, updated_at FROM teams WHERE is_cdl_franchise = true ORDER BY id").Scan(&teams)

	var events []sitemapRow
	db.Raw(`
		SELECT id, slug, updated_at FROM tournaments
		WHERE tournament_type NOT IN ('season_summary','unknown')
		ORDER BY start_date DESC
	`).Scan(&events)

	now := time.Now().UTC().Format("2006-01-02")

	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	b.WriteString("\n")
	b.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	b.WriteString("\n")

	staticRoutes := []struct {
		path     string
		priority string
		freq     string
	}{
		{"/", "1.0", "weekly"},
		{"/players", "0.9", "weekly"},
		{"/teams", "0.9", "weekly"},
		{"/events", "0.9", "weekly"},
		{"/stats", "0.8", "weekly"},
		{"/transfers", "0.7", "weekly"},
	}

	for _, r := range staticRoutes {
		fmt.Fprintf(&b, "  <url>\n    <loc>%s%s</loc>\n    <lastmod>%s</lastmod>\n    <changefreq>%s</changefreq>\n    <priority>%s</priority>\n  </url>\n",
			baseURL, r.path, now, r.freq, r.priority)
	}

	for _, p := range players {
		mod := p.UpdatedAt.UTC().Format("2006-01-02")
		fmt.Fprintf(&b, "  <url>\n    <loc>%s/players/%d</loc>\n    <lastmod>%s</lastmod>\n    <changefreq>monthly</changefreq>\n    <priority>0.6</priority>\n  </url>\n",
			baseURL, p.ID, mod)
	}

	for _, t := range teams {
		mod := t.UpdatedAt.UTC().Format("2006-01-02")
		fmt.Fprintf(&b, "  <url>\n    <loc>%s/teams/%d</loc>\n    <lastmod>%s</lastmod>\n    <changefreq>monthly</changefreq>\n    <priority>0.6</priority>\n  </url>\n",
			baseURL, t.ID, mod)
	}

	for _, e := range events {
		mod := e.UpdatedAt.UTC().Format("2006-01-02")
		fmt.Fprintf(&b, "  <url>\n    <loc>%s/events/%s</loc>\n    <lastmod>%s</lastmod>\n    <changefreq>monthly</changefreq>\n    <priority>0.7</priority>\n  </url>\n",
			baseURL, e.Slug, mod)
	}

	b.WriteString("</urlset>\n")

	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.Header("Cache-Control", "public, max-age=86400")
	c.String(http.StatusOK, b.String())
}
