DELETE FROM jenkins_slave 
WHERE name in ('maven',
               'gradle',
               'dotnet');