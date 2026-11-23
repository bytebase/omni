CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
GO
CREATE SCHEMA [s];
GO
CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
GO
CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
GO
CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
GO
CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
