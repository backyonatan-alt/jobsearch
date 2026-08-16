package reminder

import (
	"context"
	"log/slog"
	"time"
	_ "time/tzdata"

	"github.com/backyonatan-alt/jobsearch/internal/mail"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Runner sends the pre-round reminder email: one per scheduled interview,
// fired when the round is ≤24h away. sent_emails rows are claimed before
// sending, so ticks are idempotent across restarts.
type Runner struct {
	Pool        *pgxpool.Pool
	Mail        *mail.Client
	Logger      *slog.Logger
	BaseURL     string
	UnsubSecret string
}

const tickEvery = 10 * time.Minute

func (r *Runner) Run(ctx context.Context) {
	if r.Mail == nil {
		return
	}
	r.tick(ctx)
	t := time.NewTicker(tickEvery)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			r.tick(ctx)
		}
	}
}

type dueRound struct {
	InterviewID   int64
	UserID        int64
	ApplicationID int64
	Email         string
	Summary       string
	Company       string
	Role          string
	StartsAt      time.Time
	AllDay        bool
	TZ            *string
}

func (r *Runner) tick(ctx context.Context) {
	rows, err := r.Pool.Query(ctx, `
		SELECT iv.id, iv.user_id, a.id, u.email, iv.summary, a.company, a.role,
		       iv.starts_at, iv.all_day, iv.tz
		FROM interviews iv
		JOIN applications a ON a.id = iv.application_id
		JOIN users u ON u.id = iv.user_id
		WHERE iv.scheduled
		  AND u.email_reminders
		  AND iv.starts_at > now()
		  AND iv.starts_at <= now() + interval '24 hours'
		  AND (a.notes IS NULL OR a.notes NOT LIKE '[demo] %')
		  AND NOT EXISTS (
		      SELECT 1 FROM sent_emails se
		      WHERE se.interview_id = iv.id AND se.kind = 'pre_round')
		ORDER BY iv.starts_at
		LIMIT 100`)
	if err != nil {
		r.Logger.Error("reminder query", "err", err)
		return
	}
	var due []dueRound
	for rows.Next() {
		var d dueRound
		if err := rows.Scan(&d.InterviewID, &d.UserID, &d.ApplicationID, &d.Email,
			&d.Summary, &d.Company, &d.Role, &d.StartsAt, &d.AllDay, &d.TZ); err != nil {
			r.Logger.Error("reminder scan", "err", err)
			rows.Close()
			return
		}
		due = append(due, d)
	}
	rows.Close()

	for _, d := range due {
		r.sendOne(ctx, d)
	}
}

func (r *Runner) sendOne(ctx context.Context, d dueRound) {
	// Claim before send: the unique (interview_id, kind) row is the lock.
	tag, err := r.Pool.Exec(ctx, `
		INSERT INTO sent_emails (user_id, interview_id, kind)
		VALUES ($1, $2, 'pre_round')
		ON CONFLICT (interview_id, kind) DO NOTHING`, d.UserID, d.InterviewID)
	if err != nil || tag.RowsAffected() == 0 {
		return
	}

	subject, html, text := composePreRound(r.BaseURL, r.UnsubSecret, d, time.Now())
	if err := r.Mail.Send(ctx, d.Email, subject, html, text); err != nil {
		r.Logger.Error("reminder send", "err", err, "interview_id", d.InterviewID)
		// Release the claim so the next tick retries.
		_, _ = r.Pool.Exec(ctx, `
			DELETE FROM sent_emails WHERE interview_id = $1 AND kind = 'pre_round'`, d.InterviewID)
		return
	}
	r.Logger.Info("reminder sent", "interview_id", d.InterviewID, "user_id", d.UserID)
	_, _ = r.Pool.Exec(ctx, `
		INSERT INTO events (user_id, name, props)
		VALUES ($1, 'reminder_email_sent', $2)`,
		d.UserID, map[string]any{"interview_id": d.InterviewID, "application_id": d.ApplicationID})
}
