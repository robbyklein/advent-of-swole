-- name: GetMuscleStats :many
SELECT 
    UNNEST(c.muscle_groups) AS muscle_group,
    COUNT(*) AS count,
    ROUND(
        (COUNT(*) * 100.0) / SUM(COUNT(*)) OVER (), 
        2
    )::FLOAT AS percentage,
    SUM(COUNT(*)) OVER ()::FLOAT AS total_count
FROM 
    challenges c
JOIN 
    day_challenges dc 
    ON c.id = dc.challenge_id
JOIN 
    days d 
    ON dc.day_id = d.id
JOIN 
    user_challenge_completions ucc 
    ON c.id = ucc.challenge_id 
   AND d.id = ucc.day_id
WHERE 
    d.challenge_month_id = $1
    AND ucc.user_id = $2
GROUP BY 
    muscle_group
ORDER BY 
    percentage DESC;

-- name: GetCategoryStats :many
SELECT 
    c.category,
    COUNT(*) AS count,
    ROUND(
        (COUNT(*) * 100.0) / SUM(COUNT(*)) OVER (), 
        2
    )::FLOAT AS percentage
FROM 
    challenges c
JOIN 
    day_challenges dc 
    ON c.id = dc.challenge_id
JOIN 
    days d 
    ON dc.day_id = d.id
JOIN 
    user_challenge_completions ucc 
    ON c.id = ucc.challenge_id 
   AND d.id = ucc.day_id
WHERE 
    d.challenge_month_id = $1
    AND ucc.user_id = $2
GROUP BY 
    c.category
ORDER BY 
    percentage DESC;

-- name: GetCaloriesStats :many
SELECT 
    d.day_number,
    COALESCE(SUM(c.calories_burned_estimate)::INT, 0) AS total_calories
FROM 
    days d
LEFT JOIN 
    day_challenges dc ON d.id = dc.day_id
LEFT JOIN 
    challenges c ON dc.challenge_id = c.id
LEFT JOIN 
    user_challenge_completions ucc ON c.id = ucc.challenge_id AND d.id = ucc.day_id
WHERE 
    d.challenge_month_id = $1
    AND ucc.user_id = $2
GROUP BY 
    d.day_number
ORDER BY 
    d.day_number;

-- name: GetTotalParticipantsForMonth :one
SELECT 
    COUNT(DISTINCT ucc.user_id) AS total_participants
FROM 
    user_challenge_completions ucc
JOIN 
    days d ON ucc.day_id = d.id
WHERE 
    d.challenge_month_id = $1;

-- name: GetUserRankForMonth :one
WITH leaderboard AS (
    SELECT 
        u.id AS user_id,
        COALESCE(SUM(c.points), 0) AS total_points
    FROM 
        users u
    LEFT JOIN 
        user_challenge_completions ucc ON u.id = ucc.user_id
    LEFT JOIN 
        challenges c ON ucc.challenge_id = c.id
    LEFT JOIN 
        days d ON ucc.day_id = d.id
    WHERE 
        d.challenge_month_id = $1
    GROUP BY 
        u.id
    ORDER BY 
        total_points DESC
)
SELECT 
    RANK() OVER (ORDER BY total_points DESC) AS rank
FROM 
    leaderboard
WHERE 
    user_id = $2;

-- name: GetTotalChallengesCompletedForMonth :one
SELECT 
    COUNT(*) AS total_challenges_completed
FROM 
    user_challenge_completions ucc
JOIN 
    days d ON ucc.day_id = d.id
WHERE 
    d.challenge_month_id = $1
    AND ucc.user_id = $2;