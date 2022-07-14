-- update rates if they change for a set of brackets for a given state
ON CONFLICT (state_id, single_bracket, married_bracket) DO UPDATE SET
    single_rate = EXCLUDED.single_rate,
    single_bracket = EXCLUDED.single_bracket,
    married_rate = EXCLUDED.married_rate,
    married_bracket = EXCLUDED.married_bracket;
