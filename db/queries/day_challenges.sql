-- name: LinkChallengeToDay :exec
INSERT INTO day_challenges (
  day_id,
  challenge_id
) VALUES (
  $1,
  $2
);

-- name: UnlinkChallengeFromDay :exec
DELETE FROM day_challenges
WHERE day_id = $1 AND challenge_id = $2;

-- name: ListChallengesForDay :many
SELECT c.*
FROM challenges c
JOIN day_challenges dc ON c.id = dc.challenge_id
WHERE dc.day_id = $1
ORDER BY c.id;

-- name: ListDaysForChallenge :many
SELECT d.*
FROM days d
JOIN day_challenges dc ON d.id = dc.day_id
WHERE dc.challenge_id = $1
ORDER BY d.day_number;
