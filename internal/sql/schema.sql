create type wallet_type AS ENUM ('seller', 'commom');

create table if not exists wallets(
    id uuid primary key,
    created_at timestamptz not null,
    updated_at timestamptz not null,
    customer_name varchar(255) not null,
    document_number varchar(14) not null,
    wallet_type wallet_type not null,
    email varchar(255) not null,
    encoded_password varchar(255) not null,
    phone_number varchar(20) not null
);

create type transaction_type AS ENUM ('deposit', 'transfer');
create type transaction_entry_type AS ENUM ('inbound', 'outbound');

create table if not exists "transactions" (
    id uuid primary key,
    reference_id uuid not null,
    wallet_id uuid not null,
    "type" transaction_type not null,
    entries_type transaction_entry_type not null,
    currency varchar(3) not null,
    amount integer not null,
    "timestamp" timestamptz not null,
    snapshot uuid,
    foreign key (id) references wallets(id)
);
