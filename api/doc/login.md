# Login

Used to collect a Token for a registered User. The token will be returned in the header with the 'Authorization' key.

**URL** : `/login/`

**Method** : `POST`

**Auth required** : NO

**Data constraints**

```json
{
    "username": "[valid username]",
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

**Code** : `200 OK`

**Header example**

```yaml
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTAwNzkwNzAsInVzZXJuYW1lIjoidXNlcm5hbWUifQ.XZvmtNBSFejAoMfZB4xrcUx_mjtIJ_FLKQjNVcqsYxA
Content-Type: application/json
```

**Content example**

```json
{
  "status": 0,
  "message": "ok"
}
```

## Error Response

**Condition** : If 'password' does not match with the one stored in the database.

**Code** : `401 UNAUTHORIZED`

**Content** :

```json
{
  "status": 1,
  "message": "bad password"
}
```

### OR

**Condition** : If 'username' is not a known user.

**Code** : `401 UNAUTHORIZED`

**Content** :

```json
{
  "status": 2,
  "message": "bad credentials"
}
```
