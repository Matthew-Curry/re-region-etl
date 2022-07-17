CREATE TABLE states (
    state_id SMALLINT PRIMARY KEY,
    state_name VARCHAR ( 50 ) NOT NULL,
    -- all metrics are not null. Use zero value in load if not applicable.
    single_deduction INTEGER NOT NULL,
    married_deduction INTEGER NOT NULL,
    single_exemption SMALLINT NOT NULL,
    married_exemption SMALLINT NOT NULL,
    dependent_exemption SMALLINT NOT NULL
);