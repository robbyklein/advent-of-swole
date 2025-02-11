// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Challenge struct {
	ID                     int64
	Description            string
	DescriptionMetric      string
	Category               string
	MuscleGroups           []string
	Difficulty             int32
	CaloriesBurnedEstimate int32
	CreatedAt              pgtype.Timestamptz
	UpdatedAt              pgtype.Timestamptz
	Points                 int32
}

type ChallengeMonth struct {
	ID        int64
	Month     int32
	Year      int32
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

type Day struct {
	ID               int64
	ChallengeMonthID int64
	DayNumber        int32
	CreatedAt        pgtype.Timestamptz
	UpdatedAt        pgtype.Timestamptz
}

type DayChallenge struct {
	DayID       int64
	ChallengeID int64
}

type User struct {
	ID                int64
	OauthProvider     string
	OauthProviderID   string
	Email             string
	Timezone          string
	DisplayName       string
	MeasurementSystem string
	CreatedAt         pgtype.Timestamptz
	UpdatedAt         pgtype.Timestamptz
}

type UserChallengeCompletion struct {
	ID          int64
	UserID      int64
	ChallengeID int64
	DayID       int64
	CompletedAt pgtype.Timestamptz
}

type UserPoint struct {
	UserID      int64
	TotalPoints interface{}
}
