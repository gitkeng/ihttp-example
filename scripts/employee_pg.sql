-- note
-- numeric(5,2) range from -999.99 to 999.99

create or replace procedure initial_employee_table()
    language plpgsql
AS
$$
begin
    if exists(select from pg_catalog.pg_tables where schemaname = 'public' and tablename = 'employee') THEN
        RAISE NOTICE 'Table employee already exist.';
    else
        create table employee
        (
            employee_code varchar(50)    NOT NULL,
            first_name    varchar(50)    NOT NULL,
            last_name     varchar(50)    NOT NULL,
            email         varchar(100)   NULL,
            age           int            NOT NULL default 0,
            department    varchar(100)   NOT NULL,
            salary        numeric(13, 2) NOT NULL,
            update_time   bigint         NOT NULL,
            PRIMARY KEY (employee_code)
        );
        create index idx_employee_first_name on employee (first_name) ;
        create index idx_employee_last_name on employee (last_name) ;
        create index idx_employee_department on employee (department) ;
        create index idx_employee_salary on employee (salary) ;
        create index idx_employee_update_time on employee (update_time) ;
    end if;
end
$$;

CALL initial_employee_table();
DROP procedure if exists initial_employee_table();