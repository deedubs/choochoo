-- name: CreateWebhookEvent :one
INSERT INTO webhook_events (
    delivery_id,
    event_type,
    repository_name,
    sender_login,
    action,
    payload
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetWebhookEventByDeliveryID :one
SELECT * FROM webhook_events 
WHERE delivery_id = $1;

-- name: ListWebhookEventsByType :many
SELECT * FROM webhook_events 
WHERE event_type = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListWebhookEventsByRepository :many
SELECT * FROM webhook_events 
WHERE repository_name = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountWebhookEventsByType :one
SELECT COUNT(*) FROM webhook_events 
WHERE event_type = $1;

-- name: DeleteOldWebhookEvents :exec
DELETE FROM webhook_events 
WHERE created_at < $1;