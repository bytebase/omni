CREATE SCHEMA app;
GO

CREATE TABLE app.accounts (
	id BIGINT IDENTITY(1, 2) NOT NULL,
	email NVARCHAR(190) NOT NULL,
	bio NVARCHAR(MAX) NULL,
	status VARCHAR(10) NOT NULL CONSTRAINT df_accounts_status DEFAULT 'open',
	balance DECIMAL(12, 4) NOT NULL CONSTRAINT df_accounts_balance DEFAULT 0.0000,
	flags VARBINARY(16) NULL,
	token UNIQUEIDENTIFIER NULL,
	created_at DATETIME2(3) NOT NULL CONSTRAINT df_accounts_created DEFAULT SYSUTCDATETIME(),
	updated_at DATETIME2(3) NULL,
	email_lower AS LOWER(email) PERSISTED,
	PRIMARY KEY (id),
	CONSTRAINT uk_accounts_email UNIQUE (email),
	CONSTRAINT chk_accounts_balance CHECK (balance >= 0)
);
GO

EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'customer accounts', @level0type = N'SCHEMA', @level0name = N'app', @level1type = N'TABLE', @level1name = N'accounts';
GO

EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'owner''s note', @level0type = N'SCHEMA', @level0name = N'app', @level1type = N'TABLE', @level1name = N'accounts', @level2type = N'COLUMN', @level2name = N'bio';
GO

CREATE INDEX idx_accounts_status ON app.accounts (status) INCLUDE (email) WHERE status = 'open';
GO

CREATE TABLE app.orders (
	id BIGINT IDENTITY NOT NULL PRIMARY KEY,
	account_id BIGINT NOT NULL CONSTRAINT fk_orders_account REFERENCES app.accounts (id) ON DELETE CASCADE,
	total INT NOT NULL,
	note NVARCHAR(100) NULL
);
GO

CREATE TABLE app.metrics (
	ts INT NOT NULL,
	val BIGINT NOT NULL,
	CONSTRAINT pk_metrics PRIMARY KEY (ts)
);
GO

CREATE SEQUENCE app.seq_orders AS BIGINT START WITH 10 INCREMENT BY 2 MINVALUE 1 MAXVALUE 10000 CYCLE CACHE 50;
GO

EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'order sequence', @level0type = N'SCHEMA', @level0name = N'app', @level1type = N'SEQUENCE', @level1name = N'seq_orders';
GO

CREATE VIEW app.v_open_accounts AS
	SELECT id AS acct_id, email AS acct_email FROM app.accounts WHERE status = 'open';
GO

CREATE FUNCTION app.fn_format_tag (@prefix NVARCHAR(16), @id BIGINT) RETURNS NVARCHAR(64) AS
BEGIN
	RETURN @prefix + N'-' + CAST(@id AS NVARCHAR(20));
END;
GO

CREATE FUNCTION app.fn_orders_of (@account BIGINT) RETURNS TABLE AS
	RETURN (SELECT id, total FROM app.orders WHERE account_id = @account);
GO

CREATE PROCEDURE app.sp_order_total @id BIGINT, @total INT OUTPUT AS
BEGIN
	SELECT @total = total FROM app.orders WHERE id = @id;
END;
GO

CREATE PROCEDURE app.sp_account_orders @account BIGINT AS
BEGIN
	SELECT id, total FROM app.orders WHERE account_id = @account;
	SELECT COUNT(*) AS n FROM app.orders WHERE account_id = @account;
END;
GO

CREATE TRIGGER app.trg_orders_bi ON app.orders AFTER INSERT AS
BEGIN
	UPDATE app.orders SET total = 0 WHERE total < 0 AND id IN (SELECT id FROM inserted);
END;
GO

CREATE SYNONYM app.syn_accounts FOR app.accounts;
GO
