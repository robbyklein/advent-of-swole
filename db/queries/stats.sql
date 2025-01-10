-- name: GetMuscleStats :many
SELECT 
    UNNEST(muscle_groups) AS muscle_group,
    COUNT(*) AS count,
    ROUND((COUNT(*) * 100.0) / SUM(COUNT(*)) OVER (), 2) AS percentage
FROM 
    challenges c
JOIN 
    day_challenges dc ON c.id = dc.challenge_id
JOIN 
    days d ON dc.day_id = d.id
WHERE 
    d.challenge_month_id = $1
GROUP BY 
    muscle_group
ORDER BY 
    percentage DESC;

-- name: GetCategoryStats :many
SELECT 
    category,
    COUNT(*) AS count,
    ROUND((COUNT(*) * 100.0) / SUM(COUNT(*)) OVER (), 2) AS percentage
FROM 
    challenges c
JOIN 
    day_challenges dc ON c.id = dc.challenge_id
JOIN 
    days d ON dc.day_id = d.id
WHERE 
    d.challenge_month_id = $1
GROUP BY 
    category
ORDER BY 
    percentage DESC;

-- name: GetCaloriesStats :many
SELECT 
    d.day_number,
    COALESCE(SUM(c.calories_burned_estimate), 0) AS total_calories
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

