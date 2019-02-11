# RESTAPIDocs v1

API documentation for HIBK

The documentation contains only relative urls, the root url is 'http://domain.tld/api/v1'.

## Contents
- [Open Endpoints](#open-endpoints)
- [Endpoints that require Authentication](#endpoints-that-require-authentication)
  - [Current User related](#curent-user-related)
  - [Music related](#music-related)
- [Global errors](#global-errors)

## Open Endpoints

Open endpoints require no Authentication.

* [Login](login.md) : `POST /login/`

## Endpoints that require Authentication

Closed endpoints require a valid Token to be included in the header of the
request. A Token can be acquired from the Login view above.

### Current User related

Each endpoint manipulates or displays information related to the User whose
Token is provided with the request:

* [Show info](user/get.md) : `GET /user/`
* [Update password](user/put.md) : `PUT /password/`

### Music related

Endpoints for viewing and playing musics from the server.
Only authenticated users has permissions to access to these endpoints.

* [Show Available Musics](musics/get.md) : `GET /tracks/`
* [Create Account](accounts/post.md) : `POST /api/accounts/`
* [Show An Account](accounts/pk/get.md) : `GET /api/accounts/:pk/`
* [Update An Account](accounts/pk/put.md) : `PUT /api/accounts/:pk/`
* [Delete An Account](accounts/pk/delete.md) : `DELETE /api/accounts/:pk/`

## Global errors

### Error Response

**Condition** : If the request body does not correspond to the expected structure.

**Code** : `400 BAD REQUEST`

**Content** :

```json
{
  "status": 20,
  "message": "bad request"
}
```
