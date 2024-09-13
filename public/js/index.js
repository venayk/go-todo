const BASE_URL = "http://localhost:3000/api/todos"


window.onload = () => {
  document.querySelectorAll(".delete-list-item").forEach((el) => el.onclick = deleteTodo);
  document.querySelector(".new-todo-form").onsubmit = newTodoSubmit
}

function newTodoSubmit(e) {
  e.preventDefault();
  body = {text: e.srcElement.querySelector("input[type='text']").value}
  fetch(BASE_URL, {

    method: "POST",
    body: JSON.stringify(body),
    headers: {
      "content-type": "application/json"
    }
  })
  .then(interpret)
  .then(responseHandler);
  return false;
}

function deleteTodo(e) {
  id = e.srcElement.dataset.id;
  fetch(`${BASE_URL}/${id}`, {
    method: "DELETE"
  })
  .then(interpret)
  .then(responseHandler);
}

function interpret(response) {
  result = null
  try {
    result = response.json()
  } catch (error) {
    result = {}
  }
  result["ok"] = !(response.status >= 400)
  return result;
}

function responseHandler(result) {
  if (result["ok"]) {
    console.log(result);
    window.location.reload();
  } else {
    console.error(result);
  }
}