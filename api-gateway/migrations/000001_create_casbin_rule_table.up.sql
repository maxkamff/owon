CREATE TABLE IF NOT EXISTS casbin_rule(p_type character varying(100), 
    v0 character varying(100), 
    v1 character varying(100),
    v2 character varying(100),
    v3 character varying(100),
    v4 character varying(100),
    v5 character varying(100));

INSERT INTO 
casbin_rule
    (p_type, v0, v1, v2) 
VALUES
    ('p', 'unauthorized', '/v1/swagger/*', 'GET'),
    ('p', 'unauthorized', '/v1/swagger/index.html', 'GET'),
    ('p', 'unauthorized', '/v1/swagger/index.html', 'POST'),
    ('p', 'unauthorized', '/v1/', 'POST'),
    ('p', 'unauthorized', '/v1/register', 'POST'),
    ('p', 'unauthorized', '/v1/verify/{email}/{code}', 'GET'),
    ('p', 'unauthorized', '/v1/login/{email}/{password}', 'GET'),
    ('p', 'admin', '/v1/admin/add/policy', 'POST');
