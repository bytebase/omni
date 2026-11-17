CREATE SCHEMA [dump];
GO

CREATE TABLE [dump].[parents] (
	[id] INT IDENTITY(1, 1) NOT NULL,
	[name] NVARCHAR(100) NOT NULL,
	[status] VARCHAR(20) NOT NULL CONSTRAINT [DF_parents_status] DEFAULT ('open'),
	[created_at] DATETIME2(3) NOT NULL CONSTRAINT [DF_parents_created] DEFAULT (SYSUTCDATETIME()),
	[balance] DECIMAL(12, 2) NOT NULL CONSTRAINT [DF_parents_balance] DEFAULT (0),
	[computed] AS ([balance] + 1) PERSISTED,
	CONSTRAINT [PK_parents] PRIMARY KEY CLUSTERED ([id] ASC),
	CONSTRAINT [UQ_parents_name] UNIQUE NONCLUSTERED ([name] ASC),
	CONSTRAINT [CK_parents_balance] CHECK ([balance] >= 0)
);
GO

EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'parent rows', @level0type = N'SCHEMA', @level0name = N'dump', @level1type = N'TABLE', @level1name = N'parents';
GO

EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'the display name', @level0type = N'SCHEMA', @level0name = N'dump', @level1type = N'TABLE', @level1name = N'parents', @level2type = N'COLUMN', @level2name = N'name';
GO

CREATE TABLE [dump].[children] (
	[id] INT NOT NULL,
	[parent_id] INT NOT NULL,
	[status] VARCHAR(20) NOT NULL,
	[note] NVARCHAR(50) NULL,
	CONSTRAINT [PK_children] PRIMARY KEY NONCLUSTERED ([id] ASC),
	CONSTRAINT [FK_children_parents] FOREIGN KEY ([parent_id]) REFERENCES [dump].[parents] ([id]) ON DELETE CASCADE
);
GO

CREATE NONCLUSTERED INDEX [IX_children_status] ON [dump].[children] ([status] DESC, [id] ASC);
GO

CREATE VIEW [dump].[v_parents] AS
	SELECT [id], [name] FROM [dump].[parents];
GO

CREATE FUNCTION [dump].[fn_name] (@id INT) RETURNS NVARCHAR(100) AS
BEGIN
	RETURN (SELECT [name] FROM [dump].[parents] WHERE [id] = @id);
END;
GO

CREATE PROCEDURE [dump].[sp_parents] AS
BEGIN
	SELECT [id], [name] FROM [dump].[parents];
END;
GO

CREATE TRIGGER [dump].[trg_parents] ON [dump].[parents] AFTER INSERT AS
BEGIN
	SET NOCOUNT ON;
END
GO

CREATE SEQUENCE [dump].[s] AS INT;
GO
