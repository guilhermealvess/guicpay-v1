-- name: SaveWallet :exec
insert into wallets (id,wallet_type,customer_name,document_number,email,encoded_password,phone_number,created_at,updated_at) values($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: SaveTransaction :exec
insert into transactions(id,reference_id,wallet_id,"type",entries_type,currency,amount,"timestamp") values($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetTransactionByWalletID :many
select * from transactions where wallet_id = $1 and snapshot is null;

-- name: GetWalletByID :many
select 
    w.id as wallet_id,
    w.wallet_type,
    w.customer_name,
    w.document_number,
    w.email,
    w.encoded_password,
    w.phone_number,
    w.created_at,
    w.updated_at,
    t.id as transaction_id,
    t.reference_id,
    t.type as transaction_type,
    t.entries_type,
    t.currency,
    t.amount,
    t.timestamp
from wallets w
inner join transactions t
on w.id = t.wallet_id and t.snapshot is null
where w.id = $1;