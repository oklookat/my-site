package cmd

var helpPage = `
-- elven cmd --

usage:
--long-command arg1 arg2 arg3
--long-command=arg1
-short-command arg1 arg2 arg3
-short-command=arg1

commands:

[--elven-help / -eh] - show help

[--create-superuser / -csu] [username] [password] [delete_user_if_exists bool] - create superuser

[--create-user / -cu] [username] [password] [delete_user_if_exists bool] - create user

[--rollback / -rb] - delete all tables from database

[--migrate / -mg] [.sql path] - create tables in database
`
