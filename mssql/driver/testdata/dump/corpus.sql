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

CREATE UNIQUE INDEX [UX_children_note] ON [dump].[children] ([note]) WHERE [note] IS NOT NULL;
GO

CREATE INDEX [IX_children_status_note] ON [dump].[children] ([status]) INCLUDE ([note]);
GO

-- A partitioned rowstore index reports the partitioning column it does not
-- key at key_ordinal 0 with is_included_column 0. The dump carries no
-- partition scheme to restore it with, so the index must be dumped with its
-- own key alone: taking the implicit column as the leading key would make the
-- source and the restored index disagree.
CREATE PARTITION FUNCTION [dump_pf] (INT) AS RANGE RIGHT FOR VALUES (100, 200);
GO

CREATE PARTITION SCHEME [dump_ps] AS PARTITION [dump_pf] TO ([PRIMARY], [PRIMARY], [PRIMARY]);
GO

CREATE TABLE [dump].[events] (
	[id] INT NOT NULL,
	[tenant] INT NOT NULL,
	[payload] NVARCHAR(50) NULL
) ON [dump_ps] ([tenant]);
GO

CREATE CLUSTERED INDEX [IX_events_id] ON [dump].[events] ([id] ASC) ON [dump_ps] ([tenant]);
GO

CREATE VIEW [dump].[v_parents] AS
	SELECT [id], [name] FROM [dump].[parents];
GO

EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'view note', @level0type = N'SCHEMA', @level0name = N'dump', @level1type = N'VIEW', @level1name = N'v_parents';
GO

CREATE FUNCTION [dump].[fn_name] (@id INT) RETURNS NVARCHAR(100) AS
BEGIN
	RETURN (SELECT [name] FROM [dump].[parents] WHERE [id] = @id);
END;
GO

EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'function note', @level0type = N'SCHEMA', @level0name = N'dump', @level1type = N'FUNCTION', @level1name = N'fn_name';
GO

CREATE PROCEDURE [dump].[sp_parents] AS
BEGIN
	SELECT [id], [name] FROM [dump].[parents];
END;
GO

EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'procedure note', @level0type = N'SCHEMA', @level0name = N'dump', @level1type = N'PROCEDURE', @level1name = N'sp_parents';
GO

CREATE TRIGGER [dump].[trg_parents] ON [dump].[parents] AFTER INSERT AS
BEGIN
	SET NOCOUNT ON;
END
GO

CREATE SEQUENCE [dump].[s] AS INT START WITH 10 INCREMENT BY 5 MINVALUE 1 MAXVALUE 1000 CACHE 20;
GO

CREATE SEQUENCE [dump].[decimal_s] AS DECIMAL(38, 0) START WITH 1 INCREMENT BY 1;
GO

EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'sequence note', @level0type = N'SCHEMA', @level0name = N'dump', @level1type = N'SEQUENCE', @level1name = N's';
GO
