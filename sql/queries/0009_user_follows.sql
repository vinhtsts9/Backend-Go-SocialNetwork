-- name: GetFollowersByUserId :many
select 
    ub.user_id,
    ub.user_logout_time,
    case 
        when ub.user_state = 2 and ub.user_logout_time is not null and TIMESTAMPDIFF(MINUTE, ub.user_logout_time, NOW()) >= 10 then 3
        when ub.user_state = 2 and (ub.user_logout_time is null or TIMESTAMPDIFF(MINUTE, ub.user_logout_time, NOW()) < 10) then 2
        when ub.user_state = 1 then 1
        else 3
    end as calculated_user_state
from user_base ub
inner join user_follows uf 
    on ub.user_id = uf.follower_id
where uf.follower_id = ?;

-- name: UpdateUserState :exec 
update user_base
set user_state = ?
where user_id = ?;