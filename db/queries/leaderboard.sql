-- name: GetLeaderboard :many
SELECT 
  u.id AS user_id,
  u.display_name::TEXT AS display_name,
  u.email::TEXT AS email,
  COALESCE(SUM(c.points), 0)::INTEGER AS total_points
FROM 
  users u
LEFT JOIN 
  user_challenge_completions ucc ON u.id = ucc.user_id
LEFT JOIN 
  challenges c ON ucc.challenge_id = c.id
GROUP BY 
  u.id, u.display_name, u.email
ORDER BY 
  total_points DESC
LIMIT $1;
