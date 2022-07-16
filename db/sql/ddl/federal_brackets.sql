CREATE TABLE federal_brackets (
    rate DECIMAL(3, 2) NOT NULL,
    single_bracket INTEGER NOT NULL,
    married_bracket INTEGER NOT NULL,
    head_bracket INTEGER NOT NULL,
    CONSTRAINT ux_brackets UNIQUE (single_bracket, married_bracket, head_bracket)
);