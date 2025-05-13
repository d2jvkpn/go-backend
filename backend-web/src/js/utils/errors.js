class ApiError extends Error {
  constructor(code, msg, details = {}) {
    super(msg);

    this.name = 'ApiError';
    this.code = code;
    this.msg = msg;
    this.details = details; // {statusCode: 404, requestId: xxx}

    Object.setPrototypeOf(this, ApiError.prototype);

    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, ApiError);
    }
  }
}

function biz() {
  throw new ApiError('not_exists', "account not found", { statusCode: 404, requestId: "xxxx-xxxx" });
}

try {
  let ans = biz();
} catch (err) {
  // console.log((err instanceof Error) && (err instanceof ApiError));
  // console.log(`${err}`);
  console.log(`--> ApiError: code=${err.code}, statusCode=${err.msg}, details=${JSON.stringify(err.details)}`);
}

/*

1. NetworkError
if (err instanceof TypeError && err.message.startsWith("NetworkError")) {
  return;
}

2. SyntaxError
if (err instanceof SyntaxError) {
  return;
}

3. ServerError
if (response.statusCode >= 500) { // response.statusCode >= 600
  return;
}

4. BadRequest(400, http.StatusBadRequest)

5. Unauthorized(401, http.StatusUnauthorized)
if (code == "unauthorized") {
  return;
}

6. Forbidden(403, http.StatusForbidden)
if (code == "forbidden") {
  return;
}

7. ApiError
res = { "requestId": "xxxx-xxxx", "code": "not_found", "msg": "item not found" }

let err = new ApiError(
  res.code, res.msg,
  { requestId: response.requestId, statusCode: response.statusCode },
);

callback.error(err);

8. OK
res = { "requestId": "xxxx-xxxx", "code": "ok", "data": {} }
callback.ok(res.data);


function newCallback(ok=null, error=null) {
  let callback = {ok, error};

  if (!callback.ok) {
    callback.ok = (data) => {};
  }

  if (!callback.error) {
    callback.error = (err) => {};
  }

  return callback;
}
*/
