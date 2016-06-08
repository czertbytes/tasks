## Webservice endpoints:
Please notice that starting from second level endpoint is recursive. That means you can create nested tasks in task.

### `GET /tasks`

Returns a list of tasks.

```
> GET /tasks

< 200 OK
{
  tasks: Task[] = [
    { id: number, label: string, completed: boolean, sub_tasks: Task[] }
  ]
}
```

### `POST /tasks`

Creates a new task.

```
> POST /tasks
{ label: string }

< 201 Created
{
  task: { id: number, label: string, completed: boolean }
}
```

### `POST /tasks/:id`

Creates a new task.

```
> POST /tasks
{ label: string }

< 201 Created
{
  task: { id: number, label: string, completed: boolean }
}
```

```
> GET /tasks/:id

< 200 OK
{
  tasks: Task[] = [
    { id: number, label: string, completed: boolean, sub_tasks: Task[] }
  ]
}
```

### `PUT /tasks/:id`

Updates the task of the given ID.

```
> POST /tasks/:id
{ label: string } |
{ completed: boolean } |
{ label: string, completed: boolean }

< 200 OK
{
  task: Task = { id: number, label: string, completed: boolean, sub_tasks: Task[] }
}

< 404 Not Found
{ error: string }
```

### `DELETE /tasks/:id`

Deletes the task of the given ID.

```
> DELETE /tasks/:id

< 200 OK
{
  task: Task = { id: number, label: string, completed: boolean, sub_tasks: Task[] }
}

< 404 Not Found
{ error: string }
```
