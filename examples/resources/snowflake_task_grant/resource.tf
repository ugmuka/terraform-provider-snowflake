resource snowflake_task_grant grant {
  database_name = "db"
  schema_name   = "schema"
  task_name = "task"

  privilege = "operate"
  roles = [
    "role1",
    "role2",
  ]

  on_future         = false
  with_grant_option = false
}
