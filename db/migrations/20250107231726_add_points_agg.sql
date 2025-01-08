-- +goose Up
-- +goose StatementBegin
CREATE VIEW user_points AS
SELECT 
  u.id AS user_id,
  COALESCE(SUM(c.points), 0) AS total_points
FROM 
  users u
LEFT JOIN 
  user_challenge_completions ucc ON u.id = ucc.user_id
LEFT JOIN 
  challenges c ON ucc.challenge_id = c.id
GROUP BY 
  u.id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS user_points;
-- +goose StatementEnd