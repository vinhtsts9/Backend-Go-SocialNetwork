-- name: GetFollowersByUserId :many
select 
ub.user_id,
ub.user_logout_time,
ub.user_state
from user_base ub
inner join user_follows uf 
on ub.user_id = uf.follower_id
where uf.follower_id = ?;

-- name: UpdateUserState :exec 
update user_base
set user_state = ?
where user_id = ?;