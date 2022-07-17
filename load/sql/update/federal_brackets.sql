-- update rates if they change for a set of brackets
ON CONFLICT (single_bracket, married_bracket, head_bracket) DO UPDATE SET 
    rate = EXCLUDED.rate,
    single_bracket = EXCLUDED.single_bracket,
    married_bracket = EXCLUDED.married_bracket,
    head_bracket = EXCLUDED.head_bracket;
