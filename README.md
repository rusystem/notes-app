# Rest API for create notes.

### Installing:
```
git clone https://github.com/rusystem/notes-app.git
```
```
make build && make run && make swag
```

### If this is the first launch then:
```
make migrate
```

### This Rest API contains the following methods:
[post] /auth/sign-up - to create new user.<br />
[post] /auth/sign-in - user authentication.<br />
[get] /api/note - get all notes.<br />
[post] /api/note - create new note.<br />
[get] /api/note/{id} - get note by id.<br />
[put] /api/note/{id} - update note by id.<br />
[delete] /api/note/{id} - delete note by id.<br />

#### Or after launching the application visit the page localhost:8080/swagger/index.html where all available methods are described.

