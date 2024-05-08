package database

const CreateQuery string = "insert into employee values(?,?,?)"
const GetQuery string = "select * from employee where id = ?"
const DeleteQuery string = "delete from employee where id = ?"
const GetAllQuery string = "SELECT id, name, position FROM employee LIMIT ? OFFSET ?"
