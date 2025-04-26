-- Create users table
CREATE TABLE users (
    userId INT PRIMARY KEY AUTO_INCREMENT,              -- Unique user identifier
    name VARCHAR(255) NOT NULL,                         -- Full name
    email VARCHAR(255) UNIQUE,                          -- Optional: for notifications or login
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP        -- Timestamp of user creation
);

-- Create trips table
CREATE TABLE trips (
    tripId INT PRIMARY KEY AUTO_INCREMENT,              -- Unique trip ID
    tripName VARCHAR(255) NOT NULL,                     -- Name or description of the trip
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP        -- When the trip was added
);

-- Create payup_transactions table
CREATE TABLE payup_transactions (
    transactionId INT PRIMARY KEY AUTO_INCREMENT,       -- Unique transaction ID
    tripId INT NOT NULL,                                -- Foreign key to trips table
    payerId INT NOT NULL,                               -- Foreign key to users table
    itemName VARCHAR(255),                              -- Description of item/service
    establishment VARCHAR(255),                         -- Free-text name of establishment
    share VARCHAR(255),                                 -- Cost split metadata
    currency VARCHAR(10) DEFAULT 'USD',                 -- ISO currency code
    totalCost DOUBLE NOT NULL,                          -- Actual cost of item
    grossAmount DOUBLE NOT NULL,                        -- Amount paid by the user
    notes TEXT,                                         -- Optional notes
    status ENUM('active', 'cancelled', 'refunded') DEFAULT 'active',  -- Transaction status
    deleted BOOLEAN DEFAULT FALSE,                      -- Soft delete flag
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,       -- Transaction time

    -- Foreign key constraints
    FOREIGN KEY (tripId) REFERENCES trips(tripId),
    FOREIGN KEY (payerId) REFERENCES users(userId)
);

