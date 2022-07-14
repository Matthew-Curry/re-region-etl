CREATE TABLE federal_deductions (
    single_deduction SMALLINT NOT NULL, 
    married_deduction SMALLINT NOT NULL, 
    head_deduction SMALLINT NOT NULL,
    CONSTRAINT ux_deductions UNIQUE (single_deduction, married_deduction, head_deduction)
);