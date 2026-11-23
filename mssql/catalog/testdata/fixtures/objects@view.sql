CREATE TABLE [dbo].[t] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_t] PRIMARY KEY CLUSTERED ([b]));
GO
CREATE VIEW [dbo].[a] AS SELECT [b], [a], [d] FROM [dbo].[t];
