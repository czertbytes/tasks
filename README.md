# Tasks
Tasks is simple restful webservice which manages personal tasks in tree hierarchy and allows to list, create and update tasks.

Webservice api is specified in `api-doc.md`

For more details please read the source code, it's quite nice documented and tests.

## Usage
You should have Make, Docker, Docker Compose and connection to the internet (to be able to download Docker containers)

### Make image
- `make all`

### Run service
- `make run`
- Webservice is then running on port 8091. It can be changed in `docker-compose.yml` file.
- To stop simple press Ctrl+C

### Example queries
- `curl -v -X POST -H "Content-Type: application/json" -d '{"label":"foo1"}' "http://localhost:8091/tasks"`
- `curl -v -X POST -H "Content-Type: application/json" -d '{"label":"foo2"}' "http://localhost:8091/tasks/1"`
- `curl -v -X POST -H "Content-Type: application/json" -d '{"label":"foo3"}' "http://localhost:8091/tasks/1/2"`
- `curl -v "http://localhost:8091/tasks"`
- `curl -v "http://localhost:8091/tasks/1"`
- `curl -v "http://localhost:8091/tasks/1/2"`
- `curl -v "http://localhost:8091/tasks/1/2/3"`
- `curl -v -X DELETE "http://localhost:8091/tasks/1/2/3"`
- `curl -v -X PUT -H "Content-Type: application/json" -d '{"label":"foo2_update","completed":true}' "http://localhost:8091/tasks/1/2"`

## License
The MIT License (MIT)

Copyright (c) 2016 Daniel Hodan

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

Copyright (c) 2016 Daniel Hodan
