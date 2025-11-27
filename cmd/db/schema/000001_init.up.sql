CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username varchar(30) unique not null,
    name varchar(20) not null,
    password varchar(255) not null,
    phone varchar(15) not null,
    email varchar(255) unique not null,
    confirmed_email boolean default false,
    role varchar(20) not null check(role in ('landlord', 'buyer'))
);

CREATE TABLE properties (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    address TEXT NOT NULL,
    price INTEGER NOT NULL CHECK (price > 0),
    area INTEGER NOT NULL CHECK (area > 0),
    rooms_count INTEGER NOT NULL CHECK (rooms_count > 0),
    bathrooms_count INTEGER NOT NULL CHECK (bathrooms_count >= 0),
    property_type VARCHAR(20) NOT NULL CHECK (property_type IN ('house', 'apartment')),
    deal_type VARCHAR(20) NOT NULL CHECK (deal_type IN ('daily_rent', 'long_rent', 'purchase')),
    material_type VARCHAR(20) NOT NULL CHECK (material_type IN ('brick', 'concrete', 'wood', 'sip', 'pannel')),
    gas BOOLEAN DEFAULT FALSE,
    electricity BOOLEAN DEFAULT FALSE,
    internet BOOLEAN DEFAULT FALSE,
    sewerage BOOLEAN DEFAULT FALSE,
    plumbing BOOLEAN DEFAULT FALSE,
    renovation VARCHAR(20) CHECK (renovation IN ('clean', 'without_repairs', 'cosmetic')),
    floor INTEGER CHECK (floor > 0),
    land_area INTEGER CHECK (land_area > 0),
    floors INTEGER CHECK (floors > 0)
);