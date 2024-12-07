-- Create User Service Database
CREATE DATABASE user_DB;
USE user_DB;

-- Membership Tiers
CREATE TABLE MembershipTiers (
    TierID VARCHAR(5) NOT NULL PRIMARY KEY,
    TierName VARCHAR(50) NOT NULL UNIQUE,
    HourlyRate DECIMAL(10, 2) NOT NULL,
    BookingLimit INT NOT NULL
);

-- Users
CREATE TABLE Users (
    UserID VARCHAR(5) NOT NULL PRIMARY KEY,
    Name VARCHAR(100) NOT NULL,
    Email VARCHAR(100) NOT NULL UNIQUE,
    PhoneNumber VARCHAR(15) NOT NULL UNIQUE,
    PasswordHash VARCHAR(255) NOT NULL,
    TierID VARCHAR(5) NOT NULL,
    FOREIGN KEY (TierID) REFERENCES MembershipTiers(TierID)
);


-- Insert Sample Data
INSERT INTO MembershipTiers (TierID, TierName, HourlyRate, BookingLimit)
VALUES 
('T1', 'Basic', 5.00, 5), 
('T2', 'Premium', 10.00, 10), 
('T3', 'VIP', 20.00, 15);


INSERT INTO Users (UserID, Name, Email, PhoneNumber, PasswordHash, TierID)
VALUES 
('U1', 'Alice Johnson', 'alice@example.com', '12345678', 'password_1', 'T1'),
('U2', 'Bob Smith', 'bob@example.com', '98765432', 'password_2', 'T2'),
('U3', 'Catherine Lee', 'catherine@example.com', '11122233', 'password_3', 'T3'),
('U4', 'David Brown', 'david@example.com', '44455566', 'password_4', 'T1'),
('U5', 'Ella Williams', 'ella@example.com', '77788899', 'password_5', 'T2');

-- Create Vehicle Service Database
CREATE DATABASE vehicle_DB;
USE vehicle_DB;

-- Vehicles Table
CREATE TABLE Vehicles (
    VehicleID VARCHAR(5) NOT NULL PRIMARY KEY,
    Model VARCHAR(100) NOT NULL,
    VehicleNumber VARCHAR(10) NOT NULL UNIQUE,
    ChargeLevel INT NOT NULL CHECK (ChargeLevel BETWEEN 0 AND 100),
    Location VARCHAR(100) NOT NULL,
    Status ENUM('Available', 'Reserved', 'Maintenance') NOT NULL DEFAULT 'Available',
    RatePerHour DECIMAL(10, 2) NOT NULL DEFAULT 0.00
);


-- Insert Sample Data with RatePerHour
INSERT INTO Vehicles (VehicleID, Model, VehicleNumber, ChargeLevel, Location, Status, RatePerHour)
VALUES 
('V1', 'Tesla Model 3', 'SGP6101P', 90, 'Station A', 'Available', 25.00),
('V2', 'Nissan Leaf', 'SGP6102P', 70, 'Station B', 'Reserved', 20.00),
('V3', 'BMW i3', 'SGP6103P', 50, 'Station C', 'Maintenance', 30.00),
('V4', 'Chevrolet Bolt', 'SGP6104P', 80, 'Station D', 'Available', 22.00),
('V5', 'Hyundai Kona Electric', 'SGP6105P', 60, 'Station E', 'Reserved', 18.00);

-- Create Reservation Service Database
CREATE DATABASE reservation_DB;
USE reservation_DB;

-- Reservations
CREATE TABLE Reservations (
    ReservationID VARCHAR(5) NOT NULL PRIMARY KEY,
    UserID VARCHAR(5) NOT NULL, -- References UserServiceDB.Users.UserID
    VehicleID VARCHAR(5) NOT NULL, -- References VehicleServiceDB.Vehicles.VehicleID
    StartTime DATETIME NOT NULL,
    EndTime DATETIME NOT NULL,
    Status ENUM('Active', 'Cancelled', 'Completed') NOT NULL DEFAULT 'Active'
);

-- Insert Sample Data
INSERT INTO Reservations (ReservationID, UserID, VehicleID, StartTime, EndTime, Status)
VALUES 
('R1', 'U1', 'V1', '2024-12-01 10:00:00', '2024-12-01 12:00:00', 'Active'),
('R2', 'U2', 'V2', '2024-12-02 14:00:00', '2024-12-02 16:00:00', 'Completed'),
('R3', 'U3', 'V3', '2024-12-03 09:00:00', '2024-12-03 11:00:00', 'Cancelled'),
('R4', 'U4', 'V4', '2024-12-04 08:00:00', '2024-12-04 10:00:00', 'Active'),
('R5', 'U5', 'V5', '2024-12-05 15:00:00', '2024-12-05 17:00:00', 'Completed');

-- Create Billing Service Database
CREATE DATABASE billing_DB;
USE billing_DB;

-- Promotions
CREATE TABLE Promotions (
    PromoID VARCHAR(5) NOT NULL PRIMARY KEY,
    PromoCode VARCHAR(20) NOT NULL UNIQUE,
    DiscountPercent DECIMAL(5, 2) NOT NULL,
    StartDate DATE NOT NULL,
    EndDate DATE NOT NULL
);

-- Billing
CREATE TABLE Billing (
    BillingID VARCHAR(5) NOT NULL PRIMARY KEY,
    ReservationID VARCHAR(5) NOT NULL, -- References ReservationServiceDB.Reservations.ReservationID
    UserID VARCHAR(5) NOT NULL, -- References UserServiceDB.Users.UserID
    TotalAmount DECIMAL(10, 2) NOT NULL,
    PromoID VARCHAR(5), -- References BillingServiceDB.Promotions.PromoID
    FOREIGN KEY (PromoID) REFERENCES Promotions(PromoID)
);

-- Insert Sample Data
INSERT INTO Promotions (PromoID, PromoCode, DiscountPercent, StartDate, EndDate)
VALUES 
('P1', 'WELCOME10', 10.00, '2024-11-01', '2024-12-31'),
('P2', 'HOLIDAY20', 20.00, '2024-12-01', '2025-01-01'),
('P3', 'XMAS25', 25.00, '2024-12-01', '2024-12-31');


INSERT INTO Billing (BillingID, ReservationID, UserID, TotalAmount, PromoID)
VALUES 
('B1', 'R1', 'U1', 40.00, 'P1'),
('B2', 'R2', 'U2', 30.00, 'P2'),
('B3', 'R3', 'U3', 50.00, 'P3'),
('B4', 'R4', 'U4', 20.00, 'P3'),
('B5', 'R5', 'U5', 15.00, 'P3');


