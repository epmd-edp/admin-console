alter table if exists codebase add perf_server_id integer constraint perf_server_fk references perf_server;