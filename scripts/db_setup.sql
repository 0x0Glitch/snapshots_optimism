-- Moonwell Database Setup Script

-- Create the moonwell database if it doesn't exist
-- Note: This command needs to be run as a PostgreSQL superuser
-- CREATE DATABASE moonwell;

-- Connect to the database
-- \c moonwell

-- Drop table if it exists to start fresh
DROP TABLE IF EXISTS public.user_balances_optimism;

-- Create the main table for user balances
CREATE TABLE public.user_balances_optimism (
    user_addr VARCHAR(42) PRIMARY KEY,
    update_time TIMESTAMP WITH TIME ZONE
);

-- Add token columns for each token in data/tokens.json
-- Each token gets a column with an array type that stores [mToken, borrow]
ALTER TABLE public.user_balances_optimism ADD COLUMN "DAI" NUMERIC[] DEFAULT ARRAY[0, 0];
ALTER TABLE public.user_balances_optimism ADD COLUMN "USDC" NUMERIC[] DEFAULT ARRAY[0, 0];
ALTER TABLE public.user_balances_optimism ADD COLUMN "WETH" NUMERIC[] DEFAULT ARRAY[0, 0];
ALTER TABLE public.user_balances_optimism ADD COLUMN "cbETH" NUMERIC[] DEFAULT ARRAY[0, 0];
ALTER TABLE public.user_balances_optimism ADD COLUMN "wstETH" NUMERIC[] DEFAULT ARRAY[0, 0];
ALTER TABLE public.user_balances_optimism ADD COLUMN "rETH" NUMERIC[] DEFAULT ARRAY[0, 0];
ALTER TABLE public.user_balances_optimism ADD COLUMN "weETH" NUMERIC[] DEFAULT ARRAY[0, 0];
ALTER TABLE public.user_balances_optimism ADD COLUMN "wrsETH" NUMERIC[] DEFAULT ARRAY[0, 0];
ALTER TABLE public.user_balances_optimism ADD COLUMN "WBTC" NUMERIC[] DEFAULT ARRAY[0, 0];
ALTER TABLE public.user_balances_optimism ADD COLUMN "USDT" NUMERIC[] DEFAULT ARRAY[0, 0];
ALTER TABLE public.user_balances_optimism ADD COLUMN "OP" NUMERIC[] DEFAULT ARRAY[0, 0];
ALTER TABLE public.user_balances_optimism ADD COLUMN "VELO" NUMERIC[] DEFAULT ARRAY[0, 0];
ALTER TABLE public.user_balances_optimism ADD COLUMN "USDT0" NUMERIC[] DEFAULT ARRAY[0, 0];

-- Create index for faster queries
CREATE INDEX idx_user_addr ON public.user_balances_optimism(user_addr);
CREATE INDEX idx_update_time ON public.user_balances_optimism(update_time);

-- Grant permissions (if needed)
-- GRANT ALL PRIVILEGES ON TABLE public.user_balances_optimism TO your_user;

-- Display the table structure
\d public.user_balances_optimism; 