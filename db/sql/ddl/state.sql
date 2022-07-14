CREATE TABLE state (
    state_id SMALLINT PRIMARY KEY,
    state_name VARCHAR ( 50 ),
    single_deduction INTEGER,
    married_deduction INTEGER,
    single_exemption SMALLINT,
    married_exemption SMALLINT,
    dependent_exemption SMALLINT
);