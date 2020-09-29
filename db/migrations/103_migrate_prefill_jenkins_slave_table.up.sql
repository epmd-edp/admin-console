insert into jenkins_slave(name)
VALUES ('npm')
on conflict (name) do nothing;