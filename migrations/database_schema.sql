create table customers (
    id serial primary key,
    customer_id varchar(50) unique not null,
    name text,
    email text,
    address text
);

create table products (
    id serial primary key,
    product_id varchar(50) unique not null,
    name text,
    category text
);


create table orders (
    id serial primary key,
    order_id varchar(50) unique not null,
    customer_id integer references customers(id),
    region text,
    date_of_sale date,
    payment_method text,
    shipping_cost numeric
);


create table order_items (
    id serial primary key,
    order_id integer references orders(id),
    product_id integer references products(id),
    quantity_sold integer,
    unit_price numeric,
    discount numeric
);
