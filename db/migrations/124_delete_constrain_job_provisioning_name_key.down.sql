alter table job_provisioning
  add constraint job_provisioning_name_key
    unique (name);