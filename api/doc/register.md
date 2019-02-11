# Register

Used to register an user in the app.

**URL** : `/register/`

**Method** : `POST`

**Auth required** : NO

**Data constraints**

```json
{
    "username": "[new username]",
    "password": "[password in plain text]"
}
```

**Data example**

```json
{
    "username": "maximalfr",
    "password": "mypassword"
}
```

## Success Response

**Code** : `201 CREATED`

**Content example**

```json
{
  "status": 0,
  "message": "created"
}
```

## Error Response

**Condition** : If 'username' already exists.

**Code** : `401 UNAUTHORIZED`

**Content** :

```json
{
  "status": 3,
  "message": "username already exists"
}
```
